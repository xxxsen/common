package cmder

import (
	"context"
	"io"
	"os/exec"
	"syscall"
)

type Cmder struct {
	uid     *uint32
	gid     *uint32
	workdir string
	out     io.Writer
	err     io.Writer
}

// Deprecated: use Cmder
type CMDer = Cmder

// Deprecated: use New
func NewCMD(workdir string) *CMDer {
	return New(workdir)
}

func New(workdir string) *Cmder {
	return &Cmder{workdir: workdir, out: io.Discard, err: io.Discard}
}

func (c *Cmder) SetID(uid uint32, gid uint32) {
	c.uid = &uid
	c.gid = &gid
}

func (c *Cmder) SetOutput(out io.Writer, err io.Writer) *Cmder {
	c.out = out
	c.err = err
	return c
}

func (c *Cmder) Run(ctx context.Context, name string, args ...string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	cmd := exec.CommandContext(ctx, name, args...)
	if len(c.workdir) != 0 {
		cmd.Dir = c.workdir
	}
	cmd.Stdout = c.out
	cmd.Stderr = c.err

	if c.uid != nil && c.gid != nil {
		//need root
		cmd.SysProcAttr = &syscall.SysProcAttr{}
		cmd.SysProcAttr.Credential = &syscall.Credential{Uid: *c.uid, Gid: *c.gid}
	}

	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
