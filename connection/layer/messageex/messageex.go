package messageex

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/xxxsen/common/connection/layer"

	"github.com/xxxsen/common/utils"
)

func init() {
	layer.Register(messageLayerName, createMessageLayer)
}

func parseAction(idx int, params string) (*action, error) {
	grp := strings.Split(params, "@")
	act := strings.ToLower(grp[0])
	var minLength uint64 = 0
	var maxLength uint64 = 0

	switch act {
	case actionTypeSend:
		if len(grp) < 2 {
			return nil, fmt.Errorf("send type should like: send@100-200, get nil params")
		}
		sendparams := strings.Split(grp[1], "-")
		minLength, _ = strconv.ParseUint(sendparams[0], 0, 64)
		maxLength = minLength
		if len(sendparams) > 1 {
			maxLength, _ = strconv.ParseUint(sendparams[1], 0, 64)
		}
		if maxLength == 0 {
			return nil, fmt.Errorf("send type, max length should not eq to 0, idx:%d", idx)
		}
	case actionTypeRecv:
	default:
		return nil, fmt.Errorf("action type:%s invalid", act)
	}
	return &action{
		Type:      act,
		MinLength: uint(minLength),
		MaxLength: uint(maxLength),
	}, nil
}

func createMessageLayer(params interface{}) (layer.ILayer, error) {
	c := &config{}
	if err := utils.ConvStructJson(params, c); err != nil {
		return nil, err
	}
	realact := make([]*action, 0, len(c.Actions))
	for idx, params := range c.Actions {
		act, err := parseAction(idx, params)
		if err != nil {
			return nil, err
		}
		realact = append(realact, act)
	}
	return &messageLayer{acts: realact}, nil
}

type messageLayer struct {
	acts []*action
}

func (d *messageLayer) Name() string {
	return messageLayerName
}

func (d *messageLayer) MakeLayerContext(ctx context.Context, conn net.Conn) (net.Conn, error) {
	for idx, act := range d.acts {
		if strings.EqualFold(act.Type, actionTypeSend) {
			if err := d.sendMessage(idx, conn, act.MinLength, act.MaxLength); err != nil {
				return nil, err
			}
		}
		if strings.EqualFold(act.Type, actionTypeRecv) {
			if err := d.discardMessage(idx, conn); err != nil {
				return nil, err
			}
		}
	}
	return conn, nil
}

func (d *messageLayer) sendMessage(idx int, w io.Writer, minlength uint, maxlength uint) error {
	rndData := utils.RandBytes(int(minlength), int(maxlength))
	buf := make([]byte, 2+len(rndData))
	binary.BigEndian.PutUint16(buf, uint16(len(rndData)))
	copy(buf[2:], rndData)
	if _, err := w.Write(buf); err != nil {
		return fmt.Errorf("send msg fail at idx:%d with length:%d, err:%w", idx, len(rndData), err)
	}
	return nil
}

func (d *messageLayer) discardMessage(idx int, r io.Reader) error {
	b := make([]byte, 2)
	_, err := io.ReadAtLeast(r, b, len(b))
	if err != nil {
		return fmt.Errorf("read packet length at idx:%d failed, err:%w", idx, err)
	}
	length := binary.BigEndian.Uint16(b)
	if _, err := io.CopyN(io.Discard, r, int64(length)); err != nil {
		return fmt.Errorf("discard packet at idx:%d with length:%d failed, err:%w", idx, length, err)
	}
	return nil
}
