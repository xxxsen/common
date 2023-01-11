package iotool

import "io"

type wrapReadWriter struct {
	r io.Reader
	w io.Writer
}

func (w *wrapReadWriter) Read(b []byte) (int, error) {
	return w.r.Read(b)
}

func (w *wrapReadWriter) Write(b []byte) (int, error) {
	return w.w.Write(b)
}

func WrapReadWriter(r io.Reader, w io.Writer) io.ReadWriter {
	return &wrapReadWriter{
		r: r,
		w: w,
	}
}
