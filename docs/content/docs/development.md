---
title: Develop & Deploy
linkTitle: Develop
weight: 1
---

## Development Environment

We provide a [Nix Flake](https://nixos.wiki/wiki/Flakes)
that provides a development environment (dev-shell)
with all requirements included.
To use it, run from the project's root:

``` sh
nix develop
```

For convenience, you may use [direnv](https://direnv.net/)
with the [nix-direnv](https://github.com/nix-community/nix-direnv) plugin
to automatically enter the devshell upon entering the project's root.
To do so, add a `.envrc` file at root with content:

```
use flake
```
