# Build Cache for React Native

A Bitrise step that activates Bitrise Build Cache for React Native projects. It downloads the CLI binary and delegates to the `activate react-native` command, which handles Gradle, Xcode, and C++ (ccache) activation including background service startup.

## Architecture

The step is a thin wrapper around `bitrise-build-cache-cli activate react-native`. It parses two user-facing inputs (`xcode_cache_enabled`, `gradle_cache_enabled`) and maps them to CLI flags:

- `--gradle=false` / `--cpp=false` when Gradle is disabled (C++ follows Gradle since ccache handles native module compilation during Android builds)
- `--xcode=false` when Xcode is disabled

The CLI binary is embedded via Go module dependency — the step imports `cmd/reactnative` which registers the `activate react-native` Cobra subcommand, then calls `cli.RootCmd.SetArgs(...)` + `Execute()`.

## Development notes

### README generation

`README.md` is auto-generated from `step.yml` by the CI's `generate_readme` workflow (uses `bitrise-steplib/steps-readme-generator`). When changing step inputs or descriptions in `step.yml`, regenerate the README locally before pushing:

```sh
bitrise run generate_readme
```

The CI lint check diffs the committed README against a freshly generated one and fails if they differ.

### Go version

The `go` directive in `go.mod` is `1.24`. This must match the toolchain used by golangci-lint in CI.

### Mock generation

Mocks are generated with [mockery v3](https://github.com/vektra/mockery) from interfaces annotated with `//go:generate mockery`. Config is in `.mockery.yaml` at the repo root. To regenerate:

```sh
go generate ./...
```

Mocks live in `mocks/` subdirectories next to their source interfaces and are committed to the repo.

### Linting

```sh
bitrise run check
```

Runs `golangci-lint` via the `steps-check` step.

### E2E tests

E2E workflows are in `e2e/bitrise.yml`. They run the step with various input combinations and verify config files and environment variables via `e2e/scripts/check.sh`.

Key details:
- The `setup` bundle cleans config artifacts (`~/.bitrise/cache/ccache/config.json`, `~/.bitrise-xcelerate/`, `~/.gradle/init.d/bitrise-build-cache.init.gradle.kts`) between workflow runs to prevent state leaking between tests.
- `BITRISE_IO` is unset in the `run` bundle so the CLI does not detect a Bitrise CI provider — this prevents the benchmark phase API from overriding cache-enabled flags during E2E tests.
