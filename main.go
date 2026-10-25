// Command o-family moved into o proj: o family [status|repos] ARGS prints
// one deprecation line on stderr and runs o-proj status ARGS (repos adds
// --json), with o-proj's exit code. Discovery is o proj's manifest.
package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr, exec.LookPath))
}

// statusArgs maps o family's words onto o-proj status.
func statusArgs(args []string) []string {
	if len(args) > 0 {
		switch args[0] {
		case "status":
			args = args[1:]
		case "repos":
			args = append([]string{"--json"}, args[1:]...)
		}
	}
	return append([]string{"status"}, args...)
}

// catalog is the answer to o-family catalog: schema 1 with no nouns, so
// family tooling (o-inventory, nabe) still finds a well-formed catalog.
const catalog = `{"schema":1,"name":"o-family","version":"moved to o-proj","nouns":[]}`

func run(args []string, stdin io.Reader, stdout, stderr io.Writer, lookPath func(string) (string, error)) int {
	if len(args) == 1 && args[0] == "catalog" {
		fmt.Fprintln(stdout, catalog)
		return 0
	}
	argv := statusArgs(args)
	fmt.Fprintf(stderr, "o family moved into o proj: running o proj %s\n", joinArgs(argv))
	bin, err := lookPath("o-proj")
	if err != nil {
		fmt.Fprintln(stderr, "o-family: o-proj is not on PATH (github.com/unixsuperhero/o-proj)")
		return 1
	}
	cmd := exec.Command(bin, argv...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = stdin, stdout, stderr
	if err := cmd.Run(); err != nil {
		var exit *exec.ExitError
		if errors.As(err, &exit) {
			return exit.ExitCode()
		}
		fmt.Fprintln(stderr, "o-family:", err)
		return 1
	}
	return 0
}

func joinArgs(args []string) string {
	s := ""
	for i, a := range args {
		if i > 0 {
			s += " "
		}
		s += a
	}
	return s
}
