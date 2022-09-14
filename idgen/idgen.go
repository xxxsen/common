package idgen

import (
	"github.com/yitter/idgenerator-go/idgen"
	gen "github.com/yitter/idgenerator-go/idgen"
)

var defaultGenner IDGenerator

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
	opt := gen.NewIdGeneratorOptions(wrkid)
	idg := idgen.NewDefaultIdGenerator(opt)
	return &idgimpl{
		gen: idg,
	}
}
