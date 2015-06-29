// command ev reads a newline-separated file with environment variables in
// `K=V` form and uses them to form an environment to run a program.
//
// It wires up the main file descriptors (stdin, stdout, and stderr) to the
// degree required to run most applications.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

var (
	file = flag.String("file", ".env", "File to read environment variables from.")
)

func main() {
	flag.Parse()

	f, err := os.Open(*file)
	if err != nil {
		panic(err)
	}

	r := bufio.NewReader(f)

	var env []string

	for {
		l, err := r.ReadString('\n')

		l = strings.TrimSpace(l)

		if len(l) > 0 {
			env = append(env, l)
		}

		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}

	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "Error: command not specified")
		os.Exit(1)
	}

	c := exec.Command(flag.Arg(0))
	c.Args = flag.Args()
	c.Env = env

	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Start(); err != nil {
		panic(err)
	}

	if err := c.Wait(); err != nil {
		os.Exit(1)
	}
}
