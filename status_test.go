package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/unixsuperhero/oden"
)

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir, "-c", "core.hooksPath=/dev/null", "-c", "user.name=Family Test", "-c", "user.email=family@example.invalid", "-c", "commit.gpgSign=false"}, args...)...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

func put(t *testing.T, path, text string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
		t.Fatal(err)
	}
}

func realTempDir(t *testing.T) string {
	t.Helper()
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

func newRepo(t *testing.T) string {
	t.Helper()
	dir := realTempDir(t)
	runGit(t, dir, "init", "-q", "-b", "main")
	runGit(t, dir, "config", "status.renames", "true")
	for _, name := range []string{"tracked.txt", "rename.txt", "delete.txt"} {
		put(t, filepath.Join(dir, name), "base\n")
	}
	put(t, filepath.Join(dir, ".gitignore"), "ignored/\n")
	runGit(t, dir, "add", ".")
	runGit(t, dir, "commit", "-qm", "initial")
	return dir
}

func TestStatusCountsPathsRatherThanStatusColumns(t *testing.T) {
	dir := newRepo(t)
	clean := inspect(context.Background(), Repo{Name: "example", Path: dir})
	if clean.State != "clean" || clean.Tracked != 4 || clean.Changed != 0 || clean.Branch != "main" {
		t.Fatalf("clean: %+v", clean)
	}
	put(t, filepath.Join(dir, "tracked.txt"), "staged\n")
	runGit(t, dir, "add", "tracked.txt")
	put(t, filepath.Join(dir, "tracked.txt"), "worktree\n")
	runGit(t, dir, "mv", "rename.txt", "renamed\nname.txt")
	if err := os.Remove(filepath.Join(dir, "delete.txt")); err != nil {
		t.Fatal(err)
	}
	put(t, filepath.Join(dir, "new", "one.txt"), "one")
	put(t, filepath.Join(dir, "new", "two\nname.txt"), "two")
	put(t, filepath.Join(dir, "ignored", "cache"), "ignored")
	index := filepath.Join(dir, ".git", "index")
	before, err := os.ReadFile(index)
	if err != nil {
		t.Fatal(err)
	}
	s := inspect(context.Background(), Repo{Name: "example", Path: dir})
	if s.State != "dirty" || s.Tracked != 4 || s.Changed != 5 || s.Staged != 2 || s.Unstaged != 2 || s.Untracked != 2 || s.Conflicts != 0 {
		t.Fatalf("dirty: %+v", s)
	}
	after, err := os.ReadFile(index)
	if err != nil || !bytes.Equal(before, after) {
		t.Fatalf("inspection changed Git index: %v", err)
	}
}

func TestConflictsAndWorktrees(t *testing.T) {
	dir := newRepo(t)
	worktree := filepath.Join(t.TempDir(), "linked")
	runGit(t, dir, "worktree", "add", "-qb", "linked", worktree)
	s := inspect(context.Background(), Repo{Name: "linked", Path: worktree})
	if s.State != "clean" || s.Tracked != 4 || s.Branch != "linked" {
		t.Fatalf("linked worktree: %+v", s)
	}
	put(t, filepath.Join(worktree, "tracked.txt"), "linked change\n")
	runGit(t, worktree, "add", "tracked.txt")
	runGit(t, worktree, "commit", "-qm", "linked change")
	put(t, filepath.Join(dir, "tracked.txt"), "main change\n")
	runGit(t, dir, "add", "tracked.txt")
	runGit(t, dir, "commit", "-qm", "main change")
	merge := exec.Command("git", "-C", dir, "-c", "core.hooksPath=/dev/null", "-c", "user.name=Family Test", "-c", "user.email=family@example.invalid", "-c", "commit.gpgSign=false", "merge", "--no-ff", "linked")
	if err := merge.Run(); err == nil {
		t.Fatal("expected a merge conflict")
	}
	s = inspect(context.Background(), Repo{Name: "main", Path: dir})
	if s.State != "dirty" || s.Tracked != 4 || s.Changed != 1 || s.Conflicts != 1 || s.Staged != 0 || s.Unstaged != 0 {
		t.Fatalf("conflict: %+v", s)
	}
}

func TestDiscoveryDeduplicatesConfigAndHub(t *testing.T) {
	home := realTempDir(t)
	config := filepath.Join(home, "config", "repos.json")
	repo := filepath.Join(home, "project")
	if err := os.MkdirAll(repo, 0o755); err != nil {
		t.Fatal(err)
	}
	put(t, config, `[{"name":"named","path":"../project","roots":["tool"]}]`)
	hub := filepath.Join(home, "hub")
	if err := os.MkdirAll(hub, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(repo, filepath.Join(hub, "alias")); err != nil {
		t.Fatal(err)
	}
	repos, err := discover(Options{Repos: config, Hub: hub}, nil, home)
	if err != nil || len(repos) != 1 || repos[0].Name != "named" || repos[0].Path != repo {
		t.Fatalf("discovery: %+v, %v", repos, err)
	}
	explicit, err := discover(Options{Repos: "missing-config", Hub: "missing-hub"}, []string{repo}, home)
	if err != nil || len(explicit) != 1 || explicit[0].Path != repo {
		t.Fatalf("explicit directories must bypass config: %+v, %v", explicit, err)
	}
	put(t, config, `[{"name":"broken"}]`)
	if _, err := discover(Options{Repos: config}, nil, home); err == nil {
		t.Fatal("missing repository path must not silently resolve to the working directory")
	}
}

func TestFilteredReportRetainsErrorsAndFullSummary(t *testing.T) {
	clean, dirty, ordinary := newRepo(t), newRepo(t), realTempDir(t)
	put(t, filepath.Join(dirty, "new.txt"), "new")
	missing := filepath.Join(t.TempDir(), "missing")
	var out, errOut bytes.Buffer
	app := commands([]string{"repos", clean, "-d", dirty, ordinary, missing}, t.TempDir(), &out, &errOut)
	err := app.Run(context.Background())
	if oden.ExitCode(err) != 1 {
		t.Fatalf("unavailable paths must produce nonzero exit: %v, %s", err, errOut.String())
	}
	var report Report
	if err := json.Unmarshal(out.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if len(report.Repos) != 3 || report.Summary.Total != 4 || report.Summary.Dirty != 1 || report.Summary.Clean != 1 || report.Summary.Errors != 2 || report.Summary.Changed != 1 {
		t.Fatalf("filtered report: %+v", report)
	}
	states := map[string]string{}
	for _, repo := range report.Repos {
		states[repo.Path] = repo.State
	}
	if states[dirty] != "dirty" || states[ordinary] != "not-repo" || states[missing] != "missing" || states[clean] != "" {
		t.Fatalf("filter hid an error or included a clean repository: %+v", states)
	}
	out.Reset()
	app = commands([]string{"repos", "--", ordinary, missing}, t.TempDir(), &out, &errOut)
	if err := app.Run(context.Background()); oden.ExitCode(err) != 1 {
		t.Fatalf("missing/nonrepo exit: %v", err)
	}
	if err := json.Unmarshal(out.Bytes(), &report); err != nil || len(report.Repos) != 2 || report.Summary.Errors != 2 {
		t.Fatalf("unavailable report: %+v, %v", report, err)
	}
}

func TestGitEnvironmentAndUnbornBranch(t *testing.T) {
	dir := t.TempDir()
	runGit(t, dir, "init", "-q", "-b", "fresh")
	put(t, filepath.Join(dir, "new.txt"), "new")
	runGit(t, dir, "add", ".")
	other := newRepo(t)
	t.Setenv("GIT_DIR", filepath.Join(other, ".git"))
	t.Setenv("GIT_WORK_TREE", other)
	t.Setenv("GIT_INDEX_FILE", filepath.Join(other, ".git", "index"))
	s := inspect(context.Background(), Repo{Name: "fresh", Path: dir})
	if s.State != "dirty" || s.Branch != "fresh" || s.Tracked != 1 || s.Staged != 1 || s.Changed != 1 {
		t.Fatalf("unborn repository or Git environment isolation: %+v", s)
	}
}

func TestReportFlagAliasesAndDelimiter(t *testing.T) {
	dir := newRepo(t)
	put(t, filepath.Join(dir, "new.txt"), "new")
	for _, option := range []string{"-d", "--dirty"} {
		var out, errOut bytes.Buffer
		app := commands([]string{"repos", dir, option}, t.TempDir(), &out, &errOut)
		if err := app.Run(context.Background()); err != nil {
			t.Fatalf("%s: %v, %s", option, err, errOut.String())
		}
		var report Report
		if err := json.Unmarshal(out.Bytes(), &report); err != nil || len(report.Repos) != 1 || report.Repos[0].Changed != 1 {
			t.Fatalf("%s: %+v, %v", option, report, err)
		}
	}
	parent := realTempDir(t)
	for _, name := range []string{"-c", "-h"} {
		if err := os.Mkdir(filepath.Join(parent, name), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	t.Chdir(parent)
	var out, errOut bytes.Buffer
	app := commands([]string{"repos", "-d", "--", "-c", "-h"}, t.TempDir(), &out, &errOut)
	if err := app.Run(context.Background()); oden.ExitCode(err) != 1 {
		t.Fatalf("literal paths must not become options: %v, %s", err, errOut.String())
	}
	var report Report
	if err := json.Unmarshal(out.Bytes(), &report); err != nil || report.Summary.Total != 2 || report.Summary.Errors != 2 {
		t.Fatalf("literal directory report: %+v, %v", report, err)
	}
}
