# Changelog

## 2026-10-26 07:49

### Commits
- o-family: 496c1ca docs: add the o family overview animation

## 2026-10-25 22:55:44

### Tools
- changed: `o-family` version `moved to o-proj`; was none
- changed: `o-family` catalog yes; was none

### Commits
- fixed: o-family answers catalog (schema 1, no nouns) so family tooling still reads it (o-family 57a0c2f)

## 2026-10-25 22:55

### Tools
- changed: `o-family` version none; was `v0.0.0-20261025222959-e20e56f787ff`
- changed: `o-family` commit `5850914`; was `e20e56f`
- changed: `o-family` catalog none; was yes
- changed: `o-family` state none; was `command.state`
- changed: `o-family` root none; was yes

### Nouns and verbs
- removed: noun `o-family command` and its verbs `list`, `show`, `state`

### Projects
- changed: project `o-family` status none; was `open`

## 2026-10-25 22:53:33

### Root and libraries
- removed: `o-family` → `oden`

### Commits
- o-family: 5850914 feat: o family moved into o proj status; it prints a deprecation line and delegates

## 2026-10-25 22:45

### Tools
- changed: `o-family` version `v0.0.0-20261025222959-e20e56f787ff`; was `v0.0.0-20261025181528-e89db52cecd2`
- changed: `o-family` commit `e20e56f`; was `e89db52`

## 2026-10-25 22:30

### Commits
- o-family: e20e56f docs: point AGENTS.md at ~/hub/o-plans for plans and todos

## 2026-10-25 18:47

### Tools
- changed: `o-family` version `v0.0.0-20261025181528-e89db52cecd2`; was `v0.0.0-20261025180612-b248112cc22d+dirty`
- changed: `o-family` commit `e89db52`; was `b248112`

## 2026-10-25 18:15

### Root and libraries
- added: `o-family` → `oden` by replace `~/proj/oden`

### Tools
- added: `o-family` at `~/bin/o-family` (o-family b248112), 1 noun, 3 verbs, 1 param

### Projects
- added: project `o-family` (o-proj) at `~/proj/o-family`, hub link

### Commits
- o-family: e89db52 'add o-family'
- o-family: b248112 chore: new project o-family

## Unreleased

- Add read-only local family status with clean/dirty classification, tracked and changed path counts, staged/unstaged/untracked counts, conflicts, and branch names.
- Discover repositories from the existing o-inventory config and o-proj family hub, with canonical-path deduplication and explicit directory overrides.
- Provide human `status` and machine-facing `repos` JSON, clean/dirty filters, short option aliases, and explicit unavailable-directory reporting.
