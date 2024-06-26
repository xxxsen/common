package iotool

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"sync"

	"github.com/xxxsen/common/utils"

	"github.com/xxxsen/common/errs"
)

var defaultPaddingTempBufSize = 2

var paddingTempBufPool = sync.Pool{New: func() interface{} {
	return make([]byte, defaultPaddingTempBufSize)
}}

type Padding struct {
	closer                io.Closer
	spareSz               int
	bioRw                 *bufio.ReadWriter
	paddingSizeIfLessThan uint
	maxBusiDataLength     uint
	min                   uint
	max                   uint
}

func WrapPadding(rwc io.ReadWriteCloser, min, max uint, paddingIfLessThan uint, maxBusiData uint) *Padding {
	rw := bufio.NewReadWriter(bufio.NewReader(rwc), bufio.NewWriter(rwc))
	padding := &Padding{
		closer:                rwc,
		bioRw:                 rw,
		min:                   min,
		max:                   max,
		paddingSizeIfLessThan: paddingIfLessThan,
		maxBusiDataLength:     maxBusiData,
	}
	return padding
}

func (p *Padding) readSpare(b []byte) (int, error) {
	if len(b) > p.spareSz {
		b = b[:p.spareSz]
	}
	sz, err := p.bioRw.Read(b)
	if err != nil {
		return sz, err
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
		return 0, err
	}
	_, _ = p.bioRw.Discard(len(szbuf))
	length := binary.BigEndian.Uint16(szbuf[:2])
	rndLength := binary.BigEndian.Uint16(szbuf[2:])
	if length == 0 {
		return 0, errs.New(errs.ErrParam, "data length == 0")
	}
	if rndLength > 0 {
		if _, err := p.bioRw.Discard(int(rndLength)); err != nil {
			return 0, errs.Wrap(errs.ErrIO, "skip pandding data fail", err)
		}
	}
	if len(b) > int(length) {
		b = b[:length]
	}
	sz, err := p.bioRw.Read(b)
	if err != nil {
		return sz, err
	}
	if sz < int(length) {
		p.spareSz = int(length) - sz
	}
	return sz, nil
}

func (p *Padding) acquireTempBuf() []byte {
	return paddingTempBufPool.Get().([]byte)
}

func (p *Padding) releaseTempBuf(res []byte) {
	paddingTempBufPool.Put(res)
}

func (p *Padding) writeUint16WithBuffer(w io.Writer, v uint16, b []byte) error {
	binary.BigEndian.PutUint16(b, v)
	if _, err := w.Write(b); err != nil {
		return err
	}
	return nil
}

func (p *Padding) circleWrite(data []byte, tempBuf []byte) error {
	if len(data) == 0 {
		return errs.New(errs.ErrParam, "write empty data")
	}
	//2字节长度+2字节填充长度+n字节填充+m字节数据流
	var rnd []byte
	if len(data) < int(p.paddingSizeIfLessThan) {
		rnd = utils.RandBytes(int(p.min), int(p.max))
	}
	if err := p.writeUint16WithBuffer(p.bioRw, uint16(len(data)), tempBuf); err != nil {
		return errs.Wrap(errs.ErrIO, "write total length fail", err)
	}
	if err := p.writeUint16WithBuffer(p.bioRw, uint16(len(rnd)), tempBuf); err != nil {
		return errs.Wrap(errs.ErrIO, "write padding length fail", err)
	}
	if len(rnd) != 0 {
		if _, err := p.bioRw.Write(rnd); err != nil {
			return errs.Wrap(errs.ErrIO, "write padding fail", err)
		}
	}
	if _, err := p.bioRw.Write(data); err != nil {
		return errs.Wrap(errs.ErrIO, "write raw data fail", err)
	}
	if err := p.bioRw.Flush(); err != nil {
		return errs.Wrap(errs.ErrIO, "flush write fail", err)
	}
	return nil
}

func (p *Padding) Write(b []byte) (int, error) {
	tempBuf := p.acquireTempBuf()
	defer p.releaseTempBuf(tempBuf)
	for i := 0; i < len(b); i += int(p.maxBusiDataLength) {
		l := i
		r := i + int(p.maxBusiDataLength)
		if r > len(b) {
			r = len(b)
		}
		subData := b[l:r]
		if err := p.circleWrite(subData, tempBuf); err != nil {
			return 0, errs.Wrap(errs.ErrIO, fmt.Sprintf("partial write fail, total:%d, write at:%d", len(b), l), err)
		}
	}
	return len(b), nil
}

func (p *Padding) Close() error {
	var err error
	if err = p.closer.Close(); err != nil {
		err = errs.Wrap(errs.ErrIO, "close err", err)
	}
	if p.bioRw.Writer.Buffered() > 0 {
		err = errs.Wrap(errs.ErrIO, fmt.Sprintf("write buffer not empty, sz:%d", p.bioRw.Writer.Buffered()), err)
	}
	if p.bioRw.Reader.Buffered() > 0 {
		err = errs.Wrap(errs.ErrIO, "read buffer not empty", err)
	}
	if p.spareSz > 0 {
		err = errs.Wrap(errs.ErrIO, fmt.Sprintf("spare data in buf, sz:%d", p.spareSz), err)
	}
	if err != nil {
		return err
	}
	return nil
}
