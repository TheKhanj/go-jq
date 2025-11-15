package jq

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/sys/unix"
)

type keyValOption struct {
	key string
	val string
}

type option interface {
	string | keyValOption
}

type Jq struct {
	filter    io.Reader
	filePaths []string
	files     []io.Reader
	options   []interface{}
	ranOnce   bool
	tempDir   string
}

func (this *Jq) Exec() ([]byte, error) {
	if len(this.files)+len(this.filePaths) == 0 {
		return nil, errors.New("no input files provided")
	}

	if len(this.files)+len(this.filePaths) == 0 {
		return nil, errors.New("no input files provided")
	}

	if this.ranOnce {
		return nil, errors.New("cannot call function more than once")
	}
	this.ranOnce = true

	filterFifo, err := this.newFifo(this.filter)
	if err != nil {
		return nil, err
	}

	files := this.filePaths
	for _, file := range this.files {
		f, err := this.newFifo(file)
		if err != nil {
			return nil, err
		}

		files = append(files, f)
	}

	options := this.getOptions()

	args := make([]string, 0)
	args = append(args, options...)
	args = append(args, "-f", filterFifo)
	args = append(args, files...)

	cmd := exec.Command("jq", args...)

	return cmd.Output()
}

func (this *Jq) getOptions() []string {
	ret := make([]string, 0)

	for _, o := range this.options {
		switch v := o.(type) {
		case string:
			ret = append(ret, v)
		case keyValOption:
			ret = append(ret, v.key, v.val)
		default:
			panic(fmt.Errorf("invalid option: %v", v))
		}
	}

	return ret
}

func (this *Jq) newFifo(r io.Reader) (string, error) {
	err := assertDir(this.tempDir)
	if err != nil {
		return "", err
	}

	tmp, err := os.CreateTemp(this.tempDir, "jq-fifo-*")
	if err != nil {
		return "", err
	}

	path := tmp.Name()
	err = tmp.Close()
	if err != nil {
		return "", err
	}

	err = os.Remove(path)
	if err != nil {
		return "", err
	}

	err = unix.Mkfifo(path, 0600)
	if err != nil {
		return "", err
	}

	go func() {
		file, err := os.OpenFile(path, os.O_WRONLY, 0600)
		if err != nil {
			log.Printf("jq: opening fifo failed: %v", err)
			return
		}

		defer func() {
			err = file.Close()
			if err != nil {
				log.Printf("jq: closing fifo failed: %v", err)
			}

			err = os.Remove(path)
			if err != nil {
				log.Printf("jq: closing fifo failed: %v", err)
			}
		}()

		b, err := io.Copy(file, r)
		if err != nil {
			log.Printf(
				"jq: write failed: wrote %d bytes before error: %v",
				b, err,
			)
		}
	}()

	return path, nil
}

type JqOption = func(jq *Jq)

func WithFlag(flag string) JqOption {
	return func(jq *Jq) {
		jq.options = append(jq.options, flag)
	}
}

func WithOption(key, val string) JqOption {
	return func(jq *Jq) {
		jq.options = append(jq.options, keyValOption{key, val})
	}
}

func WithFilter(filter io.Reader) JqOption {
	return func(jq *Jq) {
		jq.filter = filter
	}
}

func WithFilterString(filter string) JqOption {
	r := strings.NewReader(filter)

	return WithFilter(r)
}

func WithFile(file io.Reader) JqOption {
	return func(jq *Jq) {
		jq.files = append(jq.files, file)
	}
}

func WithFilePath(path string) JqOption {
	return func(jq *Jq) {
		jq.filePaths = append(jq.filePaths, path)
	}
}

func WithFileData(data []byte) JqOption {
	r := bytes.NewReader(data)

	return WithFile(r)
}

func WithTempDir(dir string) JqOption {
	return func(jq *Jq) {
		jq.tempDir = dir
	}
}

func New(opts ...JqOption) *Jq {
	jq := Jq{
		filter:    nil,
		filePaths: make([]string, 0),
		files:     make([]io.Reader, 0),
		options:   make([]interface{}, 0),
		ranOnce:   false,
		tempDir:   os.TempDir(),
	}

	for _, opt := range opts {
		opt(&jq)
	}

	return &jq
}
