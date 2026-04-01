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

### Linting

```sh
bitrise run check
```

Runs `golangci-lint` via the `steps-check` step.
