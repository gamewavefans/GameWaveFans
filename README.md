# GameWave Research & Emulation

[![Discord server invitation](https://discord.com/api/guilds/1149362963812188350/widget.png?style=shield)](https://discord.gg/Qrz8FM6CXQ)

This repository is our initial research and tooling repo for the GameWaveFans project.

Check out [the wiki](https://github.com/namgo/GameWaveFans/wiki) and our [Discord](https://discord.gg/Qrz8FM6CXQ)!

## Tools

Tis repository contains an array of tools, available for download at [https://github.com/gamewavefans/GameWaveFans/releases/latest](https://github.com/gamewavefans/GameWaveFans/releases/latest):

- zwf_unpack - can unpack .zwf audio files, and whole directories recursively
- zbm_unpack - can unpack .zbm image files, and whole directories recursively
- zbc_unpack - can unpack .zbc bytecode files, and whole directories recursively

# Building

To build tools you can run `make` to build tools just for your platform, or any of these make targets: `build_all`, `build_linux_32`, `build_linux_64`, `build_windows_32`, `build_windows_64`.
