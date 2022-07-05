# Usage

{% include "woke.md" %}

## Config file

Configure your custom rules config in `.woke.yaml` or `.woke.yml`. `woke` uses the following precedence order. Each item takes precedence over the item below it:

- `current working directory`
- `$HOME`

This file will be picked up automatically up your customizations without needing to supply it with the `-c` flag.

See [example.yaml]({{config.repo_url}}blob/main/example.yaml) for an example of adding custom rules.
You can also supply your own rules with `-c path/to/rules.yaml` if you want to handle different rulesets.

### Remote config file

You can also use a remote config file by providing a publicly-accessible URL.

```bash
$ woke -c https://raw.githubusercontent.com/get-woke/woke/main/example.yaml
No findings found.
```

## Inputs

### File globs

By default, `woke` will run against all text files in your current directory.
To change this, supply a space-separated list of file glob patterns.
`woke` supports the following glob pattern:

```
pattern:
	{ term }
term:
	'*'         matches any sequence of non-Separator characters
	'?'         matches any single non-Separator character
	'[' [ '^' ] { character-range } ']'
	            character class (must be non-empty)
	c           matches character c (c != '*', '?', '\\', '[')
	'\\' c      matches character c

character-range:
	c           matches character c (c != '\\', '-', ']')
	'\\' c      matches character c
	lo '-' hi   matches character c for lo <= c <= hi
```

This can be something like `**/*.go`, or a space-separated list of filenames.

If `woke` is invoked from a shell, the invoking shell performs file glob pattern expansion according to the shell glob rules.

```bash
$ woke test.txt
test.txt:2:2-11: `Blacklist` may be insensitive, use `denylist`, `blocklist` instead (warning)
* Blacklist
  ^
test.txt:3:2-12: `White-list` may be insensitive, use `allowlist` instead (warning)
* White-list
  ^
test.txt:4:2-11: `whitelist` may be insensitive, use `allowlist` instead (warning)
* whitelist
  ^
test.txt:5:2-11: `blacklist` may be insensitive, use `denylist`, `blocklist` instead (warning)
* blacklist
  ^
```

### STDIN

You can also provide text to `woke` via STDIN (Standard Input)

```bash
$ echo "This has whitelist from stdin" | woke --stdin
/dev/stdin:1:9-18: `whitelist` may be insensitive, use `allowlist` instead (warning)
This has whitelist from stdin
         ^
```

This option may not be used at the same time as [File Globs](#file-globs)

## Outputs

Options for output include text (default), simple, json, github-actions, or sonarqube format.
The following fields are supported, depending on format:

| Field        | Description                                       |
| ------------ | ------------                                      |
| filepath     | Relative path to file including filename          |
| rulename     | Name of the rule from the config file             |
| termname     | Specific term that was found in the text          |
| alternative  | List of alternative terms to use instead          |
| note         | Note about reasoning for inclusion                |
| severity     | From config, one of "error", "warning", or "info" |
| optionbool   | Option value, true or false                       |
| linecontents | Contents of the line with finding                 |
| lineno       | Line number, 1 based                              |
| startcol     | Starting column number, 0 based                   |
| endcol       | Ending column number, 0 based                     |
| description  | Description of finding                            |

Output is sent to STDOUT (Standard Output), which may be redirected to a file to save the results of a scan.

### Text

!!! example ""
    `woke -o text`

`text` is the default output format for woke. Displays each result on two lines. Includes color formatting if the terminal supports it.

#### Structure

```text
<filepath>:<lineno>:<startcol>-<endcol>: <description> (<severity>)
<linecontents>
```

### Simple

!!! example ""
    `woke -o simple`

`simple` is a format similar to text, but without color support and with each result on a single line.

#### Structure

```text
<filepath>:<lineno>:<startcol>: <description>
```

### GitHub Actions

!!! example ""
    `woke -o github-actions`

The `github-actions` output type is intended for integration with [GitHub Actions](https://github.com/features/actions). See [woke-action](https://github.com/get-woke/woke-action) for more information on integration.

#### Structure

```text
::<severity> file=<filepath>,line=<lineno>,col=<startcol>::<description>
```

### JSON

!!! example ""
    `woke -o json`

Outputs the results as a series of [`json`](https://www.json.org/json-en.html) formatted structures, one per line. In order to parse as a JSON document, each line must be processed separately. This output type includes every field available in woke.

#### Structure

!!! info inline end
    Actual output from woke will be consolidated JSON. Pretty-JSON here is just for readability.

```json
{
  "Filename": "<filepath>",
  "Results": [
    {
      "Rule": {
        "Name": "<rulename>",
        "Terms": [
          "<termname>",
          ...
        ],
        "Alternatives": [
          "<alternative>",
          ...
        ],
        "Note": "<note>",
        "Severity": "<severity>",
        "Options": {
          "WordBoundary": <optionbool>,
          "WordBoundaryStart": <optionbool>,
          "WordBoundaryEnd": <optionbool>,
          "IncludeNote": <optionbool>
        }
      },
      "Finding": "<termname>",
      "Line": "<linecontents>",
      "StartPosition": {
        "Filename": "<filepath>",
        "Offset": 0,
        "Line": <lineno>,
        "Column": <startcol>
      },
      "EndPosition": {
        "Filename": "<filepath>",
        "Offset": 0,
        "Line": <lineno>,
        "Column": <endcol>
      },
      "Reason": "<description>"
    }
  ]
}
```

### SonarQube

!!! example ""
    `woke -o sonarqube`

Format used to populate results into the popular [SonarQube](https://www.sonarqube.org/) code quality tool. Note: `woke` is not executed as part of SonarQube itself, so must be run and the results file output prior to execution. Typically woke would be executed with an automation server such as Jenkins, Travis CI or Github Actions prior to creating the SonarQube report. For details on the file format, see [Generic Issue Input Format](https://docs.sonarqube.org/latest/analysis/generic-issue/). The [Analysis Parameter](https://docs.sonarqube.org/latest/analysis/analysis-parameters/) `sonar.externalIssuesReportPaths` is used to point to the path to the report file generated by `woke`.

#### Structure

!!! info inline end
    Actual output from woke will be consolidated JSON. Pretty-JSON here is just for readability.

```json
{
  "issues": [
    {
      "engineId": "woke",
      "ruleId": "<rulename>",
      "primaryLocation": {
        "message": "<description>",
        "filePath": "<filepath>",
        "textRange": {
          "startLine": <lineno>,
          "startColumn": <startcol>,
          "endColumn": <endcol>
        }
      },
      "type": "CODE_SMELL",
      "severity": "<sonarqubeseverity>"
    }
  ]
}
```

!!! note
    `<sonarqubeseverity>` is mapped from severity, such that an error in `woke` is translated to a `MAJOR`, warning to a `MINOR`, and info to `INFO`

## Exit Code

By default, `woke` will exit with a successful exit code when there are any rule failures.
The idea is, if you run `woke` on PRs, you may not want to block a merge, but you do
want to inform the author that they can make better word choices.

If you're using `woke` on PRs, you can choose to enforce these rules with a non-zero
exit code by running `woke --exit-1-on-failure`.

## Parallelism

!!! error "Advanced Configuration"
    It's unlikely that you will need to adjust parallelism. But in case you do, if you are running `woke`
    with limited resources, and/or against a very large directory, you may want to restrict the number of
    threads that `woke` uses.

By default, `woke` will parse files in parallel and will consume as many resources as it can, unbounded.
This means `woke` will be fast, but might run out of memory, depending on how large the files/lines are.

We can limit these allocations by bounding the number of files read in parallel. To accomplish this,
set the environment variable `WORKER_POOL_COUNT` to an integer value of the fixed number of goroutines
you would like to spawn for reading files.

Read more about go's concurrency patterns [here](https://blog.golang.org/pipelines).
