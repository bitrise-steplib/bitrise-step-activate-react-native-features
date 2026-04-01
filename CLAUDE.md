# bitrise-step-activate-react-native-features

A Bitrise step that activates Bitrise Build Cache for all build systems used in React Native projects: C++ via ccache, Xcode via Xcelerate, and Gradle.

Follows the same feature-plugin pattern as `bitrise-step-activate-gradle-features`.

## Development notes

### README generation

`README.md` is auto-generated from `step.yml` by the CI's `generate_readme` workflow (uses `bitrise-steplib/steps-readme-generator`). When changing step inputs or descriptions in `step.yml`, regenerate the README locally before pushing:

```sh
bitrise run generate_readme
```

The CI lint check diffs the committed README against a freshly generated one and fails if they differ.

### Go version

Keep the `go` directive in `go.mod` at `1.23` (with a `toolchain` line for the actual toolchain version). The CI runs `golangci-lint` which is built with go1.24 and will fail if the module targets a newer language version than the linter was compiled with.

**Warning:** `go mod tidy` will reset the `go` directive back to the current toolchain version. After running `go mod tidy`, manually restore `go 1.23` + `toolchain go1.24.0`.

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
