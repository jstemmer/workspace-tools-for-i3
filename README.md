# Workspace tools for i3

A collection of small tools to make working with [i3](https://i3wm.org/)
workspaces a bit easier.

## `rename-i3-workspace`

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
