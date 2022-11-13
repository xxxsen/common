package cmder

import (
	"context"
	"io"
	"io/ioutil"
	"os/exec"
	"syscall"
)

type CMDer struct {
	uid     *uint32
	gid     *uint32
	workdir string
	out     io.Writer
	err     io.Writer
}

func NewCMD(workdir string) *CMDer {
	return &CMDer{workdir: workdir, out: ioutil.Discard, err: ioutil.Discard}
}

func (c *CMDer) SetID(uid uint32, gid uint32) {
	c.uid = &uid
	c.gid = &gid
}

func (c *CMDer) SetOutput(out io.Writer, err io.Writer) *CMDer {
	c.out = out
	c.err = err
	return c
}

func (c *CMDer) Run(ctx context.Context, name string, args ...string) error {
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
