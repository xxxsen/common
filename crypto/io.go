package crypto

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	defaultMaxFrameSizeToRead  = 64 * 1024
	defaultMaxFrameSizeToWrite = 63 * 1024
)

type reader struct {
	r     *bufio.Reader
	cc    ICodec
	key   []byte
	spare []byte
}

func NewReader(r io.Reader, cc ICodec, key []byte) io.Reader {
	bio, ok := r.(*bufio.Reader)
	if !ok {
		bio = bufio.NewReader(r)
	}
	return &reader{r: bio, cc: cc, key: key}
}

func (r *reader) Read(b []byte) (int, error) {
	if len(r.spare) > 0 {
		return r.readFromSpare(b)
	}
	if err := r.decodeFrame(); err != nil {
		return 0, err
	}
	return r.readFromSpare(b)
}

func (r *reader) readFromSpare(b []byte) (int, error) {
	cnt := copy(b, r.spare)
	r.spare = r.spare[cnt:]
	return cnt, nil
}

func (r *reader) decodeFrame() error {
	tmp, err := r.r.Peek(2)
	if err != nil {
		return err
	}
	length := binary.BigEndian.Uint16(tmp)
	if int(length) > defaultMaxFrameSizeToRead {
		return fmt.Errorf("frame too big, size:%d", length)
	}
	_, _ = r.r.Discard(2)
	frame := make([]byte, length)
	_, err = io.ReadAtLeast(r.r, frame, int(length))
	if err != nil {
		return err
	}
	buf, err := r.cc.Decrypt(frame, r.key)
	if err != nil {
		return err
	}
	r.spare = buf
	return nil
}

type writer struct {
	w   *bufio.Writer
	key []byte
	cc  ICodec
}

func NewWriter(w io.Writer, cc ICodec, key []byte) io.Writer {
	bwr, ok := w.(*bufio.Writer)
	if !ok {
		bwr = bufio.NewWriter(w)
	}
	return &writer{w: bwr, cc: cc, key: key}
}

func (w *writer) Write(data []byte) (int, error) {
	if err := w.circleWriteFrame(data); err != nil {
		return 0, err
	}
	if err := w.w.Flush(); err != nil {
		return 0, err
	}
	return len(data), nil
}

func (w *writer) circleWriteFrame(b []byte) error {
	length := make([]byte, 2)
	for i := 0; i < len(b); i += defaultMaxFrameSizeToWrite {
		end := i + defaultMaxFrameSizeToWrite
		if end > len(b) {
			end = len(b)
		}
		buf := b[i:end]
		data, err := w.cc.Encrypt(buf, w.key)
		if err != nil {
			return err
		}
		binary.BigEndian.PutUint16(length, uint16(len(data)))
		if _, err := w.w.Write(length); err != nil {
			return err
		}
		if _, err := w.w.Write(data); err != nil {
			return err
		}
	}
	return nil
}
