# filesplit

CLI tool to split a file into smaller sub-files

## Build

```bash
$ go build
```

## Usage

```bash
filesplit [mode] [-F, --file] [-N, --number]
```

## Examples

```bash
# Split foo.txt into 5 sub-files
$ filesplit file -F foo.txt -N 5

# Split bar.txt into sub-files with 5 lines each
$ filesplit line -F bar.txt -N 5
```
