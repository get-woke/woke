# Ignoring

## Files

You can ignore files by adding to your config file. All ways of ignoring files below should follow the [gitignore](https://git-scm.com/docs/gitignore) convention.

```yaml
ignore_files:
  - other/files/in/repo
  - globs/too/*
```

!!! note "`.git`"
    Woke will always ignore the `.git` directory so there's no need to include it in any of the ignore configurations.

`woke` will also automatically ignore anything listed in `.gitignore`, `.ignore`, and `.git/info/exclude`.

## `.wokeignore`

You may also specify a `.wokeignore` file at the root of the directory to add additional ignore files.
This also follows the [gitignore](https://git-scm.com/docs/gitignore) convention.

See [.wokeignore.example]({{config.repo_url}}/blob/main/.wokeignore.example) for a collection of common files and directories that may contain generated [SHA](https://en.wikipedia.org/wiki/Secure_Hash_Algorithms) and [GUID](https://en.wikipedia.org/wiki/Universally_unique_identifier)s. Dependency directories are also shown in the example as the linter will parse dependency source code and possibly find errors.

## In-line and next-line ignoring

There may be times where you don't want to ignore an entire file.
You may ignore a specific line for one or more rules by creating an in-line or next-line comment.

This functionality is very rudimentary, it does a simple search for the phrase. Since
`woke` is just a text file analyzer, it has no concept of the comment syntax for every file
type it might encounter.

For in-line ignoring, simply add the following to the line you wish to ignore, using comment syntax that is supported for your file type.

!!! danger
    `woke` is not responsible for broken code due to in-line ignoring. Make sure you comment correctly!

Next-line ignoring works in a similar way. Instead of adding to the end of line you wish to ignore, you can create the ignore comment on its own line just before it. Any alphanumeric text to the left of the phrase will cause `woke` to treat it as an in-line ignore, but any text to the right of the phrase will not be considered.

!!! note
    Next-line ignore comments takes precedence over in-line ignores, so try to only use one for any given line!

```bash
This line has RULE_NAME but will be ignored # wokeignore:rule=RULE_NAME

# wokeignore:rule=RULENAME
Here is another line with RULE_NAME that will be ignored

# a couple of examples ignoring the following line for the whitelist rule
whitelist # wokeignore:rule=whitelist

# wokeignore:rule=whitelist
whitelist

# a couple of examples doing the same for multiple rules
# rule names must be comma-separated with no spaces
whitelist and blacklist # wokeignore:rule=whitelist,blacklist

# wokeignore:rule=whitelist,blacklist
whitelist and blacklist

# wokeignore:rule=whitelist text here won't be considered by woke even if it contains whitelist
this line with whitelist will still be ignored
```

Here's an example in go:

```go
func main() {
  fmt.Println("here is the whitelist") // wokeignore:rule=whitelist

  // wokeignore:rule=blacklist
  fmt.Println("and here is the blacklist")
}
```

## Nested Ignore Files

`woke` will apply ignore rules from nested ignore files to any child files/folders, similar to a nested `.gitignore` file. Nested ignore files work for any ignore file type listed above.

>Note: To disable nested ignore file functionality, run `woke` with the `--disable-nested-ignores` flag.

```txt
project
│   README.md
│   .wokeignore (applies to whole project)
│
└───folder1
│   │   file011.txt
│   │   file012.txt
│   │   .wokeignore (applies to file011.txt, file012.txt, and subfolder1)
│   │
│   └───subfolder1
│       │   file111.txt
│       │   file112.txt
│       │   ...
│
└───folder2
    │   file021.txt
    │   file022.txt
```
