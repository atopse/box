//go:generate driver-gen

package exec

import (
	"bufio"
	"errors"
	"os/exec"
	"strings"

	"github.com/atopse/box/drivers"
)

// ExecDriver 命令执行驱动器
// @title Command
// @namespace exec.driver.atopse
// @tags #command #bash
// @desc 用于执行操作系统的各项命令
type ExecDriver struct {
	drivers.Driver
}

// ExecAction 执行Action
func (d *ExecDriver) ExecAction(action *drivers.Action) (output interface{}, err error) {
	switch action.Namespace {
	case "execute":
		return d.execute(action)
	default:
		return nil, errors.New("未知的Action:" + action.Namespace)
	}
}

// execute 执行OS的cmd命令
// @title execute
// @namespace execute
// @param command string required 执行命令内容
// @return map[string]string  命令执行结果信息
func (d *ExecDriver) execute(action *drivers.Action) (output interface{}, err error) {
	var command string
	err = action.Input.Bind("command", &command)
	if err != nil {
		return nil, err
	}

	c := strings.Split(command, " ")

	cmd := exec.Command(c[0], c[1:]...)

	// read and print stdout
	outReader, err := cmd.StdoutPipe()
	if err != nil {
		return nil, errors.New("创建StdoutPipe失败," + err.Error())
	}
	outBuffer := []string{}
	outScanner := bufio.NewScanner(outReader)
	go func() {
		for outScanner.Scan() {
			foo := outScanner.Text()
			outBuffer = append(outBuffer, foo)
		}
	}()

	// read and print stderr
	errReader, err := cmd.StderrPipe()
	if err != nil {
		return nil, errors.New("创建cmd.StderrPipe失败," + err.Error())
	}
	errBuffer := []string{}
	errScanner := bufio.NewScanner(errReader)
	go func() {
		for errScanner.Scan() {
			foo := errScanner.Text()
			// mod.Logln("Err: | ", foo)
			errBuffer = append(errBuffer, foo)
		}
	}()

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	err = cmd.Wait()
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"stdout": strings.Join(outBuffer, "\n"),
		"stderr": strings.Join(errBuffer, "\n"),
	}, nil
}
