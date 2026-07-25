package iotool

import (
	"errors"
	"fmt"
	"io"
)

var (
	errDeadlineUnsupported = errors.New("deadline is not supported")
	errWriteUnsupported    = errors.New("write is not supported")
	errEmptyData           = errors.New("data is empty")
	errEmptyWrite          = errors.New("cannot write empty data")
	errPaddingRange        = errors.New("padding setting is out of range")
	errWriteBufferNotEmpty = errors.New("write buffer is not empty")
	errReadBufferNotEmpty  = errors.New("read buffer is not empty")
	errSpareData           = errors.New("spare data remains")
)

func wrapReadResult(count int, err error) (int, error) {
	if err == nil || errors.Is(err, io.EOF) {
		return count, err
	}
	return count, fmt.Errorf("read: %w", err)
}

func wrapWriteResult(count int, err error) (int, error) {
	if err == nil {
		return count, nil
	}
	return count, fmt.Errorf("write: %w", err)
}

func wrapCloseError(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("close: %w", err)
}
