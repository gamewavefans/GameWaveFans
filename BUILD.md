# Build instructions

## Local builds

To build tools you can run `make` to build tools just for your platform, or any of these make targets: `build_all`, `build_linux_32`, `build_linux_64`, `build_windows_32`, `build_windows_64`.

## Continuous Integration

After each commit, binaries are build and stored in a Github Action artifacts.

## Releases

After a semver tag (`*.*.*`) is created, a release with that version is created automatically.

## contributions

All PRs are expected to follow [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) format to allow automatic release notes generation.
We're using [this list of allowed types](https://github.com/commitizen/conventional-commit-types/blob/master/index.json).
