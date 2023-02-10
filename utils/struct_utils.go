package utils

import (
	"github.com/mitchellh/mapstructure"
	"github.com/xxxsen/common/errs"
)

func ConvStructJson(src interface{}, dst interface{}) error {
	c, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "json",
		Result:  dst,
	})
	if err != nil {
		return errs.Wrap(errs.ErrParam, "create decoder fail", err)
	}
	if err := c.Decode(src); err != nil {
		return errs.Wrap(errs.ErrUnmarshal, "decode type fail", err)
	}
	return nil
}
