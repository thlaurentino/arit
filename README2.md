# ARIT

ARIT is a simple tool for analyzing Clojure code and finding potential quality, style, and maintainability problems.

## What does it analyze?

ARIT looks for files with these extensions:

- `.clj`
- `.cljs`
- `.cljc`

It includes rules for detecting very long functions, excessive nesting, side effects, and non-idiomatic Clojure code.

## Installation

Go is required. From the project root, run:

```bash
go build -o arit .
```

This creates the `arit` executable. Check the installation with:

```bash
./arit --help
```

## Usage

Analyze a single file:

```bash
./arit src/core.clj
```

Analyze an entire directory:

```bash
./arit src/
```

ARIT recursively analyzes all `.clj`, `.cljs`, and `.cljc` files found in the directory.

## Output formats

By default, the result is displayed as a summary in the terminal. You can also choose another format:

```bash
./arit --format text src/
./arit --format json src/ > result.json
./arit --format html src/ > result.html
./arit --format markdown src/ > result.md
./arit --format csv src/
```

The available formats are `summary`, `text`, `json`, `html`, `markdown`, and `csv`.

## Useful options

```bash
./arit --verbose src/
./arit --timing src/
./arit --quiet src/
```

- `--verbose`: shows more details during the analysis.
- `--timing`: shows the execution time.
- `--quiet`: hides the banner and progress output.

Options can be combined:

```bash
./arit --format json --quiet src/ > result.json
```

## Available rules

To see all rules and their descriptions:

```bash
./arit info-rules
```

To see a shorter list:

```bash
./arit list-rules
```

## Configuration

The `.arit.yaml` file configures the analysis rules. It should be placed in the root of the project being analyzed.

Example:

```yaml
enabled-rules:
  long-function: true
  nested-forms: true
  verbose-checks: false

rule-config:
  long-function:
    max-lines: 30
```

A rule set to `true` is enabled. A rule set to `false` is disabled.

If no `.arit.yaml` file is found, ARIT uses its default settings.

## Testing the tool with simple examples

The `examples/smells/` directory contains small `.clj` files, each with a canonical code smell example.

Analyze all examples:

```bash
./arit --format text examples/smells/
```

The output should include, among others, `redundant-do-block`, `nested-forms`, `thread-ignorance`, `immutability-violation`, `production-doall`, and `unnecessary-into`.

Analyze only one smell:

```bash
./arit --format text examples/smells/redundant_do.clj
```

Save the result as JSON:

```bash
./arit --format json examples/smells/ > result.json
```

Available examples:

```text
examples/smells/redundant_do.clj
examples/smells/nested_forms.clj
examples/smells/thread_ignorance.clj
examples/smells/immutability_violation.clj
examples/smells/production_doall.clj
examples/smells/unnecessary_into.clj
```

The files in `internal/test/data/` are used by the framework's automated tests. Enabled rules are defined in `.arit.yaml`.

## Project statistics

The `stats` command collects information about project functions, such as line count, parameter count, and nesting depth:

```bash
./arit stats src/
```

Choose the output directory:

```bash
./arit stats --output-dir statistics src/
```

Generate non-aggregated data:

```bash
./arit stats --raw --output-dir statistics src/
```

## Development

Run the tests:

```bash
go test ./...
```

Build the project again:

```bash
go build -o arit .
```

## Example using the test framework

Rule tests are divided into two parts:

- `internal/test/data/`: Clojure files used as test input.
- `internal/test/suite/`: Go tests that verify the results.

For example, the `long-function` rule uses:

```text
internal/test/data/long_function.clj
internal/test/suite/long_function_test.go
```

The Clojure file contains long functions that should be detected. The Go test specifies:

```go
framework.RuleTestCase{
    FileToAnalyze: "long_function.clj",
    RuleID:         "long-function",
    ExpectedFindings: []framework.ExpectedFinding{
        {
            Message:   "is too long:",
            StartLine: 4,
        },
    },
}
```

In this example, the framework analyzes `long_function.clj`, enables only the `long-function` rule, and checks whether a finding exists on line 4.

Run this test with:

```bash
go test ./internal/test/suite -run TestLongFunction -v
```

Run all rule tests with:

```bash
go test ./internal/test/suite/...
```
