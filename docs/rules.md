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
    #   categories: nil
```

A set of default rules is provided in [`pkg/rule/default.yaml`]({{config.repo_url}}/blob/main/pkg/rule/default.yaml).

!!! tip
    If you copy these rules into your config file, be sure to put them under the `rules:` key.

## Options

You can configure options for each rule. Add an `options` key to your rule definition to customize.

### `word_boundary`

:octicons-milestone-24: Default: `false`

* If `true`, terms will trigger findings when they are surrounded by ASCII word boundaries.
* If `false`, will trigger findings if the term if found anywhere in the line, regardless if it is within an ASCII word boundary.
* !!! warning "setting `word_boundary` to `true` will always win out over `word_boundary_start` and `word_boundary_end`"

### `word_boundary_start`

:octicons-milestone-24: Default: `false`

* If `true`, terms will trigger findings when they begin with an ASCII word boundaries.
* If `false`, will trigger findings if the term if found anywhere in the line, regardless if it begins with an ASCII word boundary.

### `word_boundary_end`

:octicons-milestone-24: Default: `false`

* If `true`, terms will trigger findings when they end with an ASCII word boundaries.
* If `false`, will trigger findings if the term if found anywhere in the line, regardless if it ends with an ASCII word boundary.

### `include_note`

:octicons-milestone-24: Default: `not set`

* If `true`, the rule note will be included in the output message explaining why this finding is not inclusive
* If `false`, the rule note will not be included in the output message
* If `not set`, `include_note` in your `woke` config file (ie `.woke.yml`) regulates if the note should be included in the output message (default: `false`).

### `categories`

:octicons-milestone-24: Default: `not set`

* A list of any number of string category names to associate with the rule
* These can be used as logical groupings for actions such as excluding certain categories of rules for example

## Disabling Default Rules

You can disable default rules by providing a rule in your `woke` config file (ie `.woke.yml`), with no terms or alternatives.

This will disable the default `whitelist` rule:

```yaml
rules:
  - name: whitelist
```

### Disable all Default Rules

There may be a case where you want full control over the rules you want to run with woke.

You can either disable each default rule via the instructions above.

Or you can run woke with `--disable-default-rules` to completely disable all default rules.

!!! note
    `woke` will fail to run if you use `--disable-default-rules` without providing your own rules
    because that would mean running `woke` without any rules, which is pointless.

## Excluding Categories of Rules

You can also specify any number of rule categories to be excluded, or filtered out, from within your `woke` configuration. If any rules in a configuration file have matching categories, they will be excluded and will not be run against the target files.

This example will exclude `rule1` and `rule2`, but will still use `rule3`:

```yaml
exclude_categories:
  - category1

rules:
  - name: rule1
    terms:
      - rule1
    alternatives:
      - alt-rule1
    severity: warning
    options:
      categories:
        - category1

  - name: rule2
    terms:
      - rule2
    alternatives:
      - alt-rule2
    severity: warning
    options:
      categories:
        - category1
        - category2

  - name: rule3
    terms:
      - rule3
    alternatives:
      - alt-rule3
    severity: warning
```
