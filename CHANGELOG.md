# Changelog

## Unreleased

- Add read-only local family status with clean/dirty classification, tracked and changed path counts, staged/unstaged/untracked counts, conflicts, and branch names.
- Discover repositories from the existing o-inventory config and o-proj family hub, with canonical-path deduplication and explicit directory overrides.
- Provide human `status` and machine-facing `repos` JSON, clean/dirty filters, short option aliases, and explicit unavailable-directory reporting.
