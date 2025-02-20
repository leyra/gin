package gin

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

type Builder interface {
	Build() error
	Binary() string
	Errors() string
}

type builder struct {
	dir         string
	binary      string
	errors      string
	useGodep    bool
	customBuild string
}

func NewBuilder(dir string, bin string, useGodep bool, customBuild string) Builder {
	if len(bin) == 0 {
		bin = "bin"
	}

	// does not work on Windows without the ".exe" extension
	if runtime.GOOS == "windows" {
		if !strings.HasSuffix(bin, ".exe") { // check if it already has the .exe extension
			bin += ".exe"
		}
	}

	return &builder{dir: dir, binary: bin, useGodep: useGodep, customBuild: customBuild}
}

func (b *builder) Binary() string {
	return b.binary
}

func (b *builder) Errors() string {
	return b.errors
}

func (b *builder) Build() error {
	var command *exec.Cmd

	if b.customBuild != "" {
		parts := strings.Split(b.customBuild, " ")
		args := []string{}
		for k, v := range parts {
			if k > 0 {
				args = append(args, v)
			}
		}
		command = exec.Command(parts[0], args...)
	} else if b.useGodep {
		command = exec.Command("godep", "go", "build", "-o", b.binary)
	} else {
		command = exec.Command("go", "build", "-o", b.binary)
	}
	command.Dir = b.dir

	output, err := command.CombinedOutput()

	if command.ProcessState.Success() {
		b.errors = ""
	} else {
		b.errors = string(output)
	}

	if len(b.errors) > 0 {
		return fmt.Errorf(b.errors)
	}

	return err
}
