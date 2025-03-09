package idgen

import (
	"github.com/yitter/idgenerator-go/idgen"
)

func init() {
	_ = Init(1)
}

var defaultGenner IDGenerator

func Default() IDGenerator {
	return defaultGenner
}

func Init(wrkid uint16) error {
	defaultGenner = New(wrkid)
	return nil
}

func NextId() uint64 {
	return defaultGenner.NextId()
}

type IDGenerator interface {
	NextId() uint64
}

type idgimpl struct {
	gen *idgen.DefaultIdGenerator
}

func (p *idgimpl) NextId() uint64 {
	return uint64(p.gen.NewLong())
}

func New(wrkid uint16) IDGenerator {
	opt := idgen.NewIdGeneratorOptions(wrkid)
	idg := idgen.NewDefaultIdGenerator(opt)
	return &idgimpl{
		gen: idg,
	}
}
