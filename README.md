# o-family

Read-only status of local oden-family repositories. Runs as `o-family` or
`o family` when installed on PATH. Uses oden for dispatch and Git for status;
no other Go dependencies.

```sh
o family                         # all family directories
o family status --dirty           # dirty repositories and unavailable directories
o family status -c                # clean repositories and unavailable directories
o family status ~/proj/oden ~/proj/ochazuke
o family repos                    # machine-readable JSON report
```

## Discovery

With no explicit directories, combine:

- `$XDG_CONFIG_HOME/o-inventory/repos.json`, defaulting to
  `~/.config/o-inventory/repos.json`. Entries have `name` and `path`; other
  fields such as `roots` are accepted and ignored.
- Directories and symlinks in `$O_PROJ_HUB`, defaulting to `~/hub/o`.
- If the default repository config is absent, seed `~/proj/oden`,
  `~/proj/ochazuke`, and `~/proj/matomeru`, alongside the hub entries.

Nothing is written to the config. Hub projects created later by `o-proj`
appear automatically, including o-family itself. Existing paths are
canonicalized and deduplicated; a configured name wins over a hub alias.
`~/` expands to the current home directory. Relative paths inside a repos
file resolve against that file's directory. Explicit directory arguments
replace discovery entirely. Unregistered clones are not found by a broad
filesystem scan; pass their directories or add them to the existing config.

| Option | Meaning |
|---|---|
| `-r`, `--repos FILE` | Read a different repository config |
| `-u`, `--hub DIR` | Read a different hub; an empty value disables hub discovery |
| `-d`, `--dirty` | Show dirty repositories |
| `-c`, `--clean` | Show clean repositories |
| `-h`, `--help` | Command usage |

Options apply to `status`, the default report, and `repos`; they may appear
before or after directory arguments. `--` ends option parsing. Dirty and
clean filters are mutually exclusive. Unavailable directories remain visible
under either filter. The summary always counts all scanned directories,
not just the displayed subset.

## Counts and states

- `TRACKED`: unique paths in Git's index, including staged additions and
  unstaged deletions. A submodule gitlink counts as one path; this is not a
  recursive count of ignored files or submodule contents.
- `CHANGED`: changed paths reported by Git, including untracked files.
  A reported rename counts once. Untracked directories are expanded to
  individual files. Ignored files are excluded.
- `STAGED` and `UNSTAGED`: paths changed in the index and working tree.
  One path can appear in both, but counts only once under `CHANGED`.
- `NEW`: untracked, nonignored files. Staged new files count under `STAGED`.
- Conflicts count once per path and appear as `conflict(N)` in the human
  report, with a separate count in the summary and JSON. They do not also
  count as staged/unstaged.
- A repository is dirty if any changed path exists. Missing paths,
  non-repository roots, and inspection errors are never labeled clean.

Linked Git worktrees and unborn branches are supported. A directory nested
inside a repository is not treated as a separate repository. Bare repos
are reported as non-repository roots because there is no working tree to
inspect. Branch is reported from local Git metadata; there is no fetch or
network access.

Git runs with optional locks and fsmonitor disabled. Ambient repository,
worktree, object-directory, and index overrides are removed so the supplied
directory is inspected. The command does not stage, commit, restore, clean,
fetch, or refresh the index.

## Human and machine commands

`status` and the no-command invocation print a table plus totals. `repos`
is explicitly machine-facing and always emits a JSON object containing
`repos` and `summary`. Each repo has its canonical path, state, branch,
counts, and any error. Neither command changes format based on stdout's
terminal state, and no JSON or human-output selection flag is required.

Oden's `help`, `info`, and `which` are human-facing. Its `commands`,
`catalog`, `state`, and `command list/show/state` expose the command registry,
not repository status. Help separates commands by audience and category.

Exit codes: 0 for a successful report, even when repositories are dirty;
1 for discovery/inspection/output errors, including missing or non-repository
directories; 2 for invalid report options. An inspection error still emits
the available report before exiting 1. `repos` keeps stdout valid JSON in
that case; discovery errors occur before a report can be collected.

## Build and verify

The local module replacement expects oden at `../oden`.

```sh
gofmt -l .
go vet ./...
go test -race ./...
go build -o ~/bin/o-family .
```

Created and registered using `o-proj new o-family`; linked from
`~/hub/o/o-family`. No GitHub repository is created by local installation.
