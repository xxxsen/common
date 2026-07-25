package iotool

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"sync"

	"github.com/xxxsen/common/utils"
)

var defaultPaddingTempBufSize = 2

var paddingTempBufPool = sync.Pool{New: func() interface{} {
	buffer := make([]byte, defaultPaddingTempBufSize)
	return &buffer
}}

type Padding struct {
	closer                io.Closer
	spareSz               int
	bioRw                 *bufio.ReadWriter
	paddingSizeIfLessThan int
	maxBusiDataLength     int
	min                   int
	max                   int
}

func WrapPadding(
	rwc io.ReadWriteCloser,
	minimumPadding, maximumPadding uint,
	paddingIfLessThan uint,
	maxBusinessData uint,
) *Padding {
	if minimumPadding > math.MaxUint16 ||
		maximumPadding > math.MaxUint16 ||
		paddingIfLessThan > math.MaxUint16 ||
		maxBusinessData == 0 ||
		maxBusinessData > math.MaxUint16 {
		panic(errPaddingRange)
	}
	rw := bufio.NewReadWriter(bufio.NewReader(rwc), bufio.NewWriter(rwc))
	padding := &Padding{
		closer:                rwc,
		bioRw:                 rw,
		min:                   int(minimumPadding),
		max:                   int(maximumPadding),
		paddingSizeIfLessThan: int(paddingIfLessThan),
		maxBusiDataLength:     int(maxBusinessData),
	}
	return padding
}

func (p *Padding) readSpare(b []byte) (int, error) {
	if len(b) > p.spareSz {
		b = b[:p.spareSz]
	}
	sz, err := p.bioRw.Read(b)
	if err != nil {
		return wrapReadResult(sz, err)
	}
	p.spareSz -= sz
	return sz, nil
}

func (p *Padding) Read(b []byte) (int, error) {
	if p.spareSz > 0 {
		return p.readSpare(b)
	}
	szbuf, err := p.bioRw.Peek(4)
	if err != nil {
		return wrapReadResult(0, err)
	}
	_, _ = p.bioRw.Discard(len(szbuf))
	length := binary.BigEndian.Uint16(szbuf[:2])
	rndLength := binary.BigEndian.Uint16(szbuf[2:])
	if length == 0 {
		return 0, errEmptyData
	}
	if rndLength > 0 {
		if _, err := p.bioRw.Discard(int(rndLength)); err != nil {
			return 0, fmt.Errorf("skip pandding data failed, err:%w", err)
		}
	}
	if len(b) > int(length) {
		b = b[:length]
	}
	sz, err := p.bioRw.Read(b)
	if err != nil {
		return wrapReadResult(sz, err)
	}
	if sz < int(length) {
		p.spareSz = int(length) - sz
	}
	return sz, nil
}

func (p *Padding) acquireTempBuf() []byte {
	buffer, ok := paddingTempBufPool.Get().(*[]byte)
	if !ok {
		panic("iotool: invalid padding buffer pool value")
	}
	return *buffer
}

func (p *Padding) releaseTempBuf(res []byte) {
	paddingTempBufPool.Put(&res)
}

func (p *Padding) writeUint16WithBuffer(w io.Writer, v uint16, b []byte) error {
	binary.BigEndian.PutUint16(b, v)
	if _, err := w.Write(b); err != nil {
		return fmt.Errorf("write uint16: %w", err)
	}
	return nil
}

func (p *Padding) circleWrite(data []byte, tempBuf []byte) error {
	if len(data) == 0 {
		return errEmptyWrite
	}
	if len(data) > math.MaxUint16 {
		return errPaddingRange
	}
	// 2 字节长度 + 2 字节填充长度 + n 字节填充 + m 字节数据流。
	var rnd []byte
	if len(data) < p.paddingSizeIfLessThan {
		rnd = utils.RandBytes(p.min, p.max)
	}
	if len(rnd) > math.MaxUint16 {
		return errPaddingRange
	}
	dataLength, err := paddingLength(len(data))
	if err != nil {
		return err
	}
	if err := p.writeUint16WithBuffer(p.bioRw, dataLength, tempBuf); err != nil {
		return fmt.Errorf("write total length failed, err:%w", err)
	}
	randomLength, err := paddingLength(len(rnd))
	if err != nil {
		return err
	}
	if err := p.writeUint16WithBuffer(p.bioRw, randomLength, tempBuf); err != nil {
		return fmt.Errorf("write padding length failed, err:%w", err)
	}
	if len(rnd) != 0 {
		if _, err := p.bioRw.Write(rnd); err != nil {
			return fmt.Errorf("write padding failed, err:%w", err)
		}
	}
	if _, err := p.bioRw.Write(data); err != nil {
		return fmt.Errorf("write raw data failed, err:%w", err)
	}
	if err := p.bioRw.Flush(); err != nil {
		return fmt.Errorf("flush write failed, err:%w", err)
	}
	return nil
}

func (p *Padding) Write(b []byte) (int, error) {
	tempBuf := p.acquireTempBuf()
	defer p.releaseTempBuf(tempBuf)
	for i := 0; i < len(b); i += p.maxBusiDataLength {
		r := i + p.maxBusiDataLength
		if r > len(b) {
			r = len(b)
		}
		subData := b[i:r]
		if err := p.circleWrite(subData, tempBuf); err != nil {
			return 0, fmt.Errorf("partial write fail, total:%d, write at:%d, err:%w", len(b), i, err)
		}
	}
	return len(b), nil
}

func paddingLength(length int) (uint16, error) {
	if length < 0 || length > math.MaxUint16 {
		return 0, errPaddingRange
	}
	return uint16(length), nil
}

func (p *Padding) Close() error {
	var result error
	if err := p.closer.Close(); err != nil {
		result = errors.Join(result, fmt.Errorf("close padding stream: %w", err))
	}
	if p.bioRw.Writer.Buffered() > 0 {
		result = errors.Join(
			result,
			fmt.Errorf("%w: %d bytes", errWriteBufferNotEmpty, p.bioRw.Writer.Buffered()),
		)
	}
	if p.bioRw.Reader.Buffered() > 0 {
		result = errors.Join(result, errReadBufferNotEmpty)
	}
	if p.spareSz > 0 {
		result = errors.Join(result, fmt.Errorf("%w: %d bytes", errSpareData, p.spareSz))
	}
	return result
}
