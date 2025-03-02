package envflag

import (
	"flag"
	"os"
	"strings"
	"time"
)

// Deprecated: should not use this
type EnvFlag struct {
	fs *flag.FlagSet
}

func New(name string, handler flag.ErrorHandling) *EnvFlag {
	return &EnvFlag{
		fs: flag.NewFlagSet(name, handler),
	}
}

func (ev *EnvFlag) Parse(args ...string) error {
	// err := ev.fs.Parse(args)
	// if err != nil {
	// 	return err
	// }
	ev.resetAsEnv()
	return nil
}

func (ev *EnvFlag) rebuildName(name string) string {
	return strings.ToUpper(strings.NewReplacer("-", "_", ".", "_", "@", "_", "#", "_").Replace(name))
}

func (ev *EnvFlag) resetAsEnv() {
	//build env value
	ev.fs.VisitAll(func(f *flag.Flag) {
		value, ok := os.LookupEnv(ev.rebuildName(f.Name))
		if !ok {
			return
		}
		f.Value.Set(value)
	})
}

func (ev *EnvFlag) String(name string, value string, usage string) *string {
	return ev.fs.String(name, value, usage)
}

func (ev *EnvFlag) StringVar(p *string, name string, value string, usage string) {
	ev.fs.StringVar(p, name, value, usage)
}

func (ev *EnvFlag) Float64(name string, value float64, usage string) *float64 {
	return ev.fs.Float64(name, value, usage)
}

func (ev *EnvFlag) Float64Var(p *float64, name string, value float64, usage string) {
	ev.fs.Float64Var(p, name, value, usage)
}

func (ev *EnvFlag) Duration(name string, value time.Duration, usage string) *time.Duration {
	return ev.fs.Duration(name, value, usage)
}

func (ev *EnvFlag) DurationVar(p *time.Duration, name string, value time.Duration, usage string) {
	ev.fs.DurationVar(p, name, value, usage)
}

func (ev *EnvFlag) Uint64(name string, value uint64, usage string) *uint64 {
	return ev.fs.Uint64(name, value, usage)
}

func (ev *EnvFlag) Uint64Var(p *uint64, name string, value uint64, usage string) {
	ev.fs.Uint64Var(p, name, value, usage)
}

func (ev *EnvFlag) Uint(name string, value uint, usage string) *uint {
	return ev.fs.Uint(name, value, usage)
}

func (ev *EnvFlag) UintVar(p *uint, name string, value uint, usage string) {
	ev.fs.UintVar(p, name, value, usage)
}

func (ev *EnvFlag) Int64(name string, value int64, usage string) *int64 {
	return ev.fs.Int64(name, value, usage)
}

func (ev *EnvFlag) Int64Var(p *int64, name string, value int64, usage string) {
	ev.fs.Int64Var(p, name, value, usage)
}

func (ev *EnvFlag) Int(name string, value int, usage string) *int {
	return ev.fs.Int(name, value, usage)
}

func (ev *EnvFlag) IntVar(p *int, name string, value int, usage string) {
	ev.fs.IntVar(p, name, value, usage)
}

func (ev *EnvFlag) Bool(name string, value bool, usage string) *bool {
	return ev.fs.Bool(name, value, usage)
}

func (ev *EnvFlag) BoolVar(p *bool, name string, value bool, usage string) {
	ev.fs.BoolVar(p, name, value, usage)
}

func String(name string, value string, usage string) *string {
	return DefaultParser.String(name, value, usage)
}

func StringVar(p *string, name string, value string, usage string) {
	DefaultParser.StringVar(p, name, value, usage)
}

func Float64(name string, value float64, usage string) *float64 {
	return DefaultParser.Float64(name, value, usage)
}

func Float64Var(p *float64, name string, value float64, usage string) {
	DefaultParser.Float64Var(p, name, value, usage)
}

func Duration(name string, value time.Duration, usage string) *time.Duration {
	return DefaultParser.Duration(name, value, usage)
}

func DurationVar(p *time.Duration, name string, value time.Duration, usage string) {
	DefaultParser.DurationVar(p, name, value, usage)
}

func Uint64(name string, value uint64, usage string) *uint64 {
	return DefaultParser.Uint64(name, value, usage)
}

func Uint64Var(p *uint64, name string, value uint64, usage string) {
	DefaultParser.Uint64Var(p, name, value, usage)
}

func Uint(name string, value uint, usage string) *uint {
	return DefaultParser.Uint(name, value, usage)
}

func UintVar(p *uint, name string, value uint, usage string) {
	DefaultParser.UintVar(p, name, value, usage)
}

func Int64(name string, value int64, usage string) *int64 {
	return DefaultParser.Int64(name, value, usage)
}

func Int64Var(p *int64, name string, value int64, usage string) {
	DefaultParser.Int64Var(p, name, value, usage)
}

func Int(name string, value int, usage string) *int {
	return DefaultParser.Int(name, value, usage)
}

func IntVar(p *int, name string, value int, usage string) {
	DefaultParser.IntVar(p, name, value, usage)
}

func Bool(name string, value bool, usage string) *bool {
	return DefaultParser.Bool(name, value, usage)
}

func BoolVar(p *bool, name string, value bool, usage string) {
	DefaultParser.BoolVar(p, name, value, usage)
}

func (ev *EnvFlag) Parsed() bool {
	return ev.fs.Parsed()
}

var DefaultParser = New("default_env_parser", flag.ExitOnError)

func Parsed() bool {
	return DefaultParser.Parsed()
}

func Parse() error {
	return DefaultParser.Parse(os.Args[1:]...)
}
