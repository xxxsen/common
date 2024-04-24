package iotool

import "io"

type ioWrapper struct {
	r io.Reader
	w io.Writer
	c io.Closer
}

func (w *ioWrapper) Read(b []byte) (int, error) {
	return w.r.Read(b)
}

func (w *ioWrapper) Write(b []byte) (int, error) {
	return w.w.Write(b)
}

func (w *ioWrapper) Close() error {
	return w.c.Close()
}

func WrapReadWriter(r io.Reader, w io.Writer) io.ReadWriter {
	return &ioWrapper{
		r: r,
		w: w,
	}
}

func WrapReadWriteCloser(r io.Reader, w io.Writer, c io.Closer) io.ReadWriteCloser {
	return &ioWrapper{
		r: r,
		w: w,
		c: c,
	}
}
