package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/unixsuperhero/oden"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	home, err := os.UserHomeDir()
	if err == nil {
		err = commands(os.Args[1:], home, os.Stdout, os.Stderr).Run(ctx)
	}
	os.Exit(oden.Report(os.Stderr, "o-family", err))
}

func commands(args []string, home string, out, errOut io.Writer) *oden.App {
	app := oden.New("o-family", args, oden.WithBuiltins(), oden.WithStdio(os.Stdin, out, errOut))
	status := func(ctx context.Context, args []string) error { return runReport(ctx, args, home, out, errOut, false) }
	app.Default = status
	app.Add(
		oden.Command{Name: "status", Summary: "show clean and dirty family repositories with file counts", Category: "Repositories", Run: status},
		oden.Command{Name: "repos", Summary: "print the repository status report as JSON", Category: "Repositories", Audience: oden.AudienceMachine, Run: func(ctx context.Context, args []string) error {
			return runReport(ctx, args, home, out, errOut, true)
		}},
	)
	return app
}

func runReport(ctx context.Context, args []string, home string, out, errOut io.Writer, asJSON bool) error {
	opts := defaults(home)
	name := "status"
	if asJSON {
		name = "repos"
	}
	f := flag.NewFlagSet("o-family "+name, flag.ContinueOnError)
	f.SetOutput(errOut)
	f.StringVar(&opts.Repos, "repos", opts.Repos, "family repository JSON file")
	f.StringVar(&opts.Repos, "r", opts.Repos, "alias for --repos")
	f.StringVar(&opts.Hub, "hub", opts.Hub, "family hub directory")
	f.StringVar(&opts.Hub, "u", opts.Hub, "alias for --hub")
	f.BoolVar(&opts.Dirty, "dirty", false, "show only dirty repositories")
	f.BoolVar(&opts.Dirty, "d", false, "alias for --dirty")
	f.BoolVar(&opts.Clean, "clean", false, "show only clean repositories")
	f.BoolVar(&opts.Clean, "c", false, "alias for --clean")
	f.Usage = func() {
		fmt.Fprintf(errOut, "usage: o-family %s [OPTIONS] [DIR...]\n\nWith no directories, read %s and %s.\n\n  -r, --repos FILE  family repos.json\n  -u, --hub DIR     family hub; an empty value disables hub discovery\n  -d, --dirty       show only dirty repositories\n  -c, --clean       show only clean repositories\n  -h, --help        show this usage\n\nFlags may appear before or after directories. Explicit directories replace discovery.\n", name, opts.Repos, opts.Hub)
	}
	var dirs []string
	for len(args) > 0 {
		if args[0] == "--" {
			dirs = append(dirs, args[1:]...)
			break
		}
		if !strings.HasPrefix(args[0], "-") || args[0] == "-" {
			dirs = append(dirs, args[0])
			args = args[1:]
			continue
		}
		take := 1
		name, _, inline := strings.Cut(strings.TrimLeft(args[0], "-"), "=")
		if option := f.Lookup(name); option != nil && !inline {
			boolean, ok := option.Value.(interface{ IsBoolFlag() bool })
			if (!ok || !boolean.IsBoolFlag()) && len(args) > 1 {
				take = 2
			}
		}
		if err := f.Parse(args[:take]); err != nil {
			if errors.Is(err, flag.ErrHelp) {
				return nil
			}
			return &oden.ExitError{Code: 2}
		}
		args = args[take:]
	}
	if opts.Dirty && opts.Clean {
		return &oden.ExitError{Code: 2, Err: errors.New("--dirty and --clean are mutually exclusive")}
	}
	repos, err := discover(opts, dirs, home)
	if err != nil {
		return err
	}
	report := collect(ctx, repos)
	if opts.Dirty || opts.Clean {
		filtered := report.Repos[:0]
		for _, repo := range report.Repos {
			if repo.Error != "" || opts.Dirty && repo.State == "dirty" || opts.Clean && repo.State == "clean" {
				filtered = append(filtered, repo)
			}
		}
		report.Repos = filtered
	}
	if asJSON {
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		err = enc.Encode(report)
	} else {
		err = printReport(out, report, home)
	}
	if err != nil {
		return err
	}
	if report.Summary.Errors != 0 {
		return &oden.ExitError{Code: 1}
	}
	return nil
}

func printReport(out io.Writer, report Report, home string) error {
	w := tabwriter.NewWriter(out, 0, 4, 2, ' ', 0)
	fmt.Fprintln(w, "REPO\tSTATE\tTRACKED\tCHANGED\tSTAGED\tUNSTAGED\tNEW\tBRANCH\tDIRECTORY")
	for _, repo := range report.Repos {
		path := repo.Path
		if path == home {
			path = "~"
		} else if rel, ok := strings.CutPrefix(path, home+string(filepath.Separator)); ok {
			path = "~/" + rel
		}
		if repo.Error != "" {
			fmt.Fprintf(w, "%s\t%s\t-\t-\t-\t-\t-\t-\t%s\n", repo.Name, repo.State, path)
			fmt.Fprintf(w, "  %s\n", repo.Error)
			continue
		}
		state := repo.State
		if repo.Conflicts > 0 {
			state = fmt.Sprintf("conflict(%d)", repo.Conflicts)
		}
		fmt.Fprintf(w, "%s\t%s\t%d\t%d\t%d\t%d\t%d\t%s\t%s\n", repo.Name, state, repo.Tracked, repo.Changed, repo.Staged, repo.Unstaged, repo.Untracked, repo.Branch, path)
	}
	if err := w.Flush(); err != nil {
		return err
	}
	s := report.Summary
	_, err := fmt.Fprintf(out, "\n%d directories: %d dirty, %d clean, %d unavailable. Showing %d.\n%d tracked files; %d changed paths (%d staged, %d unstaged, %d new, %d conflicts).\n", s.Total, s.Dirty, s.Clean, s.Errors, len(report.Repos), s.Tracked, s.Changed, s.Staged, s.Unstaged, s.Untracked, s.Conflicts)
	return err
}
