package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStatusArgs(t *testing.T) {
	for in, want := range map[string]string{
		"":                 "status",
		"status -d":        "status -d",
		"repos ~/proj/x":   "status --json ~/proj/x",
		"-c":               "status -c",
		"~/proj/a --dirty": "status ~/proj/a --dirty",
	} {
		if got := strings.Join(statusArgs(strings.Fields(in)), " "); got != want {
			t.Errorf("%q: %q, want %q", in, got, want)
		}
	}
}

func TestRunDelegatesWithExitCode(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "o-proj")
	os.WriteFile(bin, []byte("#!/bin/sh\necho \"$*\"\nexit 1\n"), 0o755)
	var out, errOut bytes.Buffer
	code := run([]string{"repos", "-d"}, nil, &out, &errOut, func(string) (string, error) { return bin, nil })
	if code != 1 || out.String() != "status --json -d\n" || !strings.Contains(errOut.String(), "moved into o proj") {
		t.Fatalf("code %d out %q err %q", code, out.String(), errOut.String())
	}
	code = run(nil, nil, &out, &errOut, func(string) (string, error) { return "", errors.New("none") })
	if code != 1 || !strings.Contains(errOut.String(), "not on PATH") {
		t.Fatalf("no o-proj: %d %q", code, errOut.String())
	}
}

func TestCatalogIsEmpty(t *testing.T) {
	var out, errOut bytes.Buffer
	if code := run([]string{"catalog"}, nil, &out, &errOut, nil); code != 0 || !strings.Contains(out.String(), `"schema":1`) || errOut.Len() != 0 {
		t.Fatalf("catalog %d %q %q", code, out.String(), errOut.String())
	}
}
