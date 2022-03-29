# Workspace tools for i3

A collection of small tools to make working with [i3](https://i3wm.org/)
workspaces a bit easier.

- [move-i3-workspaces](#move-i3-workspaces)
- [new-i3-workspace](#new-i3-workspace)
- [rename-i3-workspace](#rename-i3-workspace)

## move-i3-workspaces

Quickly move i3 workspaces to a specified output.

### Usage

`move-i3-workspaces [flags] [workspaces] <output>`

Flags:
  * `--all`
  * `--except <workspaces>`
  * `-n`

Where `<workspaces>` is a comma-separated list of workspace numbers.

Use the dry-run flag `-n` to only print which workspaces will be moved without
actually moving them.

### Examples

Move all workspaces, except for workspaces `1` and `9`, to output `HDMI1`:

```bash
move-i3-workspaces --all --except 1,9 HDMI1
```

Move workspace `2`, `4` and `6` to output `DP1`:

```bash
move-i3-workspaces 2,4,6 DP1
```

## new-i3-workspace

Change to a new i3 workspace using the next available free number. Useful when
you quickly want a new empty workspace but don't want to check which number is
available.

### Usage

`new-i3-workspace [flags]`

Flags:
  * `max <number>` The maximum workspace number to consider (default 20).
  * `name <name>` Optional name to give to the new workspace.
  * `reserved <workspaces>` Comma-separated list of workspace numbers that
    should not be used. 

## rename-i3-workspace

Shows an i3-input prompt and renames the currently active workspace with the
given input while preserving the workspace number prefix.

i3 configuration example:

```
bindsym $mod+Shift+apostrophe exec rename-i3-workspace
```

## License

Apache 2.0; see [`LICENSE`](LICENSE) for details.

## Disclaimer

This is not an officially supported Google product.
