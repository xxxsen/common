package iotool

import "io"

func WrapCloser(f func() error) io.Closer {
	return &wrapCloser{
		do: f,
	}
}

type wrapCloser struct {
	do func() error
}

func (w *wrapCloser) Close() error {
	return w.do()
}
