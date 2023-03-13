package models

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"runtime"
)

type ServiceController struct {
	Title        string
	ExecArgs     []string
	Exec         *exec.Cmd
	StdoutReader io.ReadCloser
	StderrReader io.ReadCloser
	Scanner      *bufio.Scanner
	// ShouldStop   bool
}

func (s ServiceController) GetExecArgs() (res []string) {
	if len(s.ExecArgs) != 0 {
		return s.ExecArgs
	}
	if len(s.ExecArgs) == 0 {
		if runtime.GOOS == "windows" {
			res = append(res, fmt.Sprintf(".\\%s.exe", s.Title))
		} else {
			res = append(res, fmt.Sprintf("./%s", s.Title))
		}
	}
	return res
}

var (
	ErrCouldNotExec = errors.New("error could not exec")
)
