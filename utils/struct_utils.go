package utils

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

func ConvStructJson(src interface{}, dst interface{}) error {
	c, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "json",
		Result:  dst,
	})
	if err != nil {
		return fmt.Errorf("create decoder failed, err:%w", err)
	}
	if err := c.Decode(src); err != nil {
		return fmt.Errorf("decode type failed, err:%w", err)
	}
	return nil
}
