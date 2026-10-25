package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Options struct {
	Repos string
	Hub   string
	Dirty bool
	Clean bool
}

func defaults(home string) Options {
	config := os.Getenv("XDG_CONFIG_HOME")
	if config == "" {
		config = filepath.Join(home, ".config")
	}
	hub := os.Getenv("O_PROJ_HUB")
	if hub == "" {
		hub = filepath.Join(home, "hub", "o")
	}
	return Options{Repos: filepath.Join(config, "o-inventory", "repos.json"), Hub: hub}
}

func expand(path, home string) string {
	if path == "~" {
		return home
	}
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(home, path[2:])
	}
	return path
}

func discover(opts Options, dirs []string, home string) ([]Repo, error) {
	var repos []Repo
	if len(dirs) > 0 {
		for _, dir := range dirs {
			path := expand(dir, home)
			repos = append(repos, Repo{Name: filepath.Base(filepath.Clean(path)), Path: path})
		}
	} else {
		config := expand(opts.Repos, home)
		data, err := os.ReadFile(config)
		switch {
		case err == nil:
			if err := json.Unmarshal(data, &repos); err != nil {
				return nil, fmt.Errorf("%s: %w", config, err)
			}
			for i := range repos {
				if repos[i].Path == "" {
					return nil, fmt.Errorf("%s: repository %q has no path", config, repos[i].Name)
				}
				repos[i].Path = expand(repos[i].Path, home)
				if !filepath.IsAbs(repos[i].Path) {
					repos[i].Path = filepath.Join(filepath.Dir(config), repos[i].Path)
				}
			}
		case os.IsNotExist(err) && opts.Repos == defaults(home).Repos:
			for _, name := range []string{"oden", "ochazuke", "matomeru"} {
				repos = append(repos, Repo{Name: name, Path: filepath.Join(home, "proj", name)})
			}
		default:
			return nil, err
		}
		if opts.Hub != "" {
			hub := expand(opts.Hub, home)
			entries, err := os.ReadDir(hub)
			if err != nil && !(os.IsNotExist(err) && opts.Hub == defaults(home).Hub) {
				return nil, err
			}
			for _, entry := range entries {
				if entry.IsDir() || entry.Type()&os.ModeSymlink != 0 {
					repos = append(repos, Repo{Name: entry.Name(), Path: filepath.Join(hub, entry.Name())})
				}
			}
		}
	}
	seen := map[string]bool{}
	out := make([]Repo, 0, len(repos))
	for _, repo := range repos {
		path, err := filepath.Abs(repo.Path)
		if err != nil {
			return nil, err
		}
		if resolved, err := filepath.EvalSymlinks(path); err == nil {
			path = resolved
		}
		if seen[path] {
			continue
		}
		seen[path] = true
		repo.Path = path
		if repo.Name == "" {
			repo.Name = filepath.Base(path)
		}
		out = append(out, repo)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Name == out[j].Name {
			return out[i].Path < out[j].Path
		}
		return out[i].Name < out[j].Name
	})
	return out, nil
}
