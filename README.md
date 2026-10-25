# o-family (moved)

The repository status report moved into **o proj** (github.com/unixsuperhero/o-proj, `~/proj/o-proj`): `o proj status [DIR...] [-d|--dirty] [-c|--clean] [-j|--json]`, over the repositories in o proj's manifest (`~/hub/repos.json`).

`o family` still works: it prints one deprecation line on stderr and runs `o proj status` with the same arguments (`o family repos` runs `o proj status --json`), with o-proj's exit code. The old `-r/--repos` and `-u/--hub` flags are gone; discovery is the manifest.

```sh
gofmt -l . && go vet ./... && go test -race ./...
go build -o ~/bin/o-family .
```
