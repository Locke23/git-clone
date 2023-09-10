# Lit

## Requirements

- Go 1.21+

## Documentation

- `lit --help`
- `lit init`
  - initialize lit repository
- `lit add file1 file2`
  - add file to staging area
- `lit commit "commit message"`
  - create a commit pkg with files at staging area
- `lit status`
  - Displays paths that have differences between the index file and the current HEAD commit
- `lit log`
  - List commits that are reachable by following the parent links from the given commit hash
