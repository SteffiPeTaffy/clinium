package clinium

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	DefaultPromptDelimiter byte = ':'
)

type cli struct {
	t *testing.T
	path string
}

type UserInput struct {
	Prompt string
	Input  string
}

type result struct {
	output string
	err    error
	t      *testing.T
}

func NewCli(t *testing.T, path string) *cli {
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Errorf("could not find binary to cli: %s", err.Error())
	}

	cli := &cli{t, absPath}

	return cli
}

func (r *result) ToSucceed(msg string) *result {
	if r.err != nil {
		fmt.Printf(": ❌\nCommand was not suppossed to fail, but did:\n")
		fmt.Printf("Error: %s\n", r.err.Error())
		fmt.Printf("Output: \n")
		fmt.Printf(r.output)
	} else {
		fmt.Print(" ✅")
	}
	require.NoError(r.t, r.err, msg)
	return r
}

func (r *result) ToFail(msg string) *result {
	if r.err == nil {
		fmt.Printf(": ❌\nCommand was suppossed to fail, but didn't:\n")
		fmt.Printf(r.output)
	} else {
		fmt.Print(" ✅")
	}
	require.NotNil(r.t, r.err, msg)
	return r
}

func (r *result) ToContain(needle string) *result {
	require.Contains(r.t, r.output, needle)
	fmt.Print(" ✅")
	return r
}

func (r *result) ToNotContain(needle string) *result {
	require.NotContains(r.t, r.output, needle)
	fmt.Print(" ✅")
	return r
}

func (c *cli) Expect(args ...string) *result {
	output, err := c.executeInteractive([]UserInput{}, args...)
	return &result{output, err, c.t}
}

func (c *cli) ExpectInteractive(userInput []UserInput, args ...string) *result {
	output, err := c.executeInteractive(userInput, args...)
	return &result{output, err, c.t}
}

func (c *cli) executeInteractive(userInput []UserInput, args ...string) (string, error) {
	cmd := exec.Command(c.path, args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	defer stdin.Close()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	defer stdout.Close()
	stdoutReader := bufio.NewReader(stdout)

	result := make(chan string)
	go func(result chan<- string) {
		for len(userInput) > 0 {
			line, _ := stdoutReader.ReadString(DefaultPromptDelimiter)
			index, match, err := indexOf(line, userInput)
			if err == nil {
				_, _ = stdin.Write([]byte(match.Input + "\n"))
				userInput = remove(userInput, index)
			}
		}

		output, _ := ioutil.ReadAll(stdoutReader)
		result <- string(output)
	}(result)

	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	err = cmd.Wait()
	output := <-result
	return output, err
}

func indexOf(line string, userInput []UserInput) (int, *UserInput, error) {
	for i, input := range userInput {
		if strings.Contains(line, input.Prompt) {
			return i, &input, nil
		}
	}
	return -1, nil, fmt.Errorf("unexpected prompt %s", line)
}

func remove(slice []UserInput, s int) []UserInput {
	return append(slice[:s], slice[s+1:]...)
}

func (c *cli) execute(args ...string) (string, error) {
	return c.executeInteractive([]UserInput{}, args...)
}
