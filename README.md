# Version Parser

A CLI tool to parse version flags from `.txt` files. While technically we can use this pattern for any .txt file that includes the pattern of `versionFlag + versionNumber`, this tool is designed for `.txt` files exported from Microsoft Dynamics NAV.

## Installation

Compile your own binary:

```bash
go build
```

## Usage

### CLI:

```bash
./lastversionparser --path /path/to/file.txt --version-flag TX --direction max
```
