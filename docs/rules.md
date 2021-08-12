# Rules

## Schema

The syntax for rules is very basic. You just need a name, a list of terms to match that violate the rule,
and a list of alternative suggestions.

```yaml
rules:
  - name: whitelist
    terms:
      - whitelist
      - white-list
    alternatives:
      - allowlist
    note: An optional description why these terms are not inclusive. It can be optionally included in the output message.
    # options:
    #   word_boundary: false
    #   word_boundary_start: false
    #   word_boundary_end: false
    #   include_note: false
```

A set of default rules is provided in [`pkg/rule/default.yaml`]({{config.repo_url}}blob/main/pkg/rule/default.yaml).

!!! tip
    If you copy these rules into your config file, be sure to put them under the `rules:` key.

## Options

You can configure options for each rule. Add an `options` key to your rule definition to customize.

Current options include:

- `word_boundary` (default: `false`)
    * If `true`, terms will trigger findings when they are surrounded by ASCII word boundaries.
    * If `false`, will trigger findings if the term if found anywhere in the line, regardless if it is within an ASCII word boundary.
    * !!! note
          setting `word_boundary` to `true` will always win out over `word_boundary_start` and `word_boundary_end`

- `word_boundary_start` (default: `false`)
    * If `true`, terms will trigger findings when they begin with an ASCII word boundaries.
    * If `false`, will trigger findings if the term if found anywhere in the line, regardless if it begins with an ASCII word boundary.

- `word_boundary_end` (default: `false`)
    * If `true`, terms will trigger findings when they end with an ASCII word boundaries.
    * If `false`, will trigger findings if the term if found anywhere in the line, regardless if it ends with an ASCII word boundary.

- `include_note` (default: `not set`)
    * If `true`, the rule note will be included in the output message explaining why this finding is not inclusive
    * If `false`, the rule note will not be included in the output message
    * If `not set`, `include_note` in your `woke` config file (ie `.woke.yml`) regulates if the note should be included in the output message (default: `false`).

## Disabling Default Rules

You can disable default rules by providing a rule in your `woke` config file (ie `.woke.yml`), with no terms or alternatives.

This will disable the default `whitelist` rule:

```yaml
rules:
  - name: whitelist
```
