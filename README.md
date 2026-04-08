# Activate React Native Features

[![Step changelog](https://shields.io/github/v/release/bitrise-steplib/bitrise-step-activate-react-native-features?include_prereleases&label=changelog&color=blueviolet)](https://github.com/bitrise-steplib/bitrise-step-activate-react-native-features/releases)

Activates Bitrise Build Cache features for React Native projects

<details>
<summary>Description</summary>

This Step activates Bitrise Build Cache for all build systems used in a React Native project.

After this Step executes,
- enabling C++ cache will result in: C++ native modules will be compiled using ccache, with a background storage helper started automatically.
- enabling Xcode cache will result in: iOS builds will use the Bitrise Build Cache for Xcode via Xcelerate, with a background proxy started automatically.
- enabling Gradle cache will result in: Android Gradle builds will automatically read from and push entries to the remote cache.

</details>

## 🧩 Get started

Add this step directly to your workflow in the [Bitrise Workflow Editor](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/steps/adding-steps-to-a-workflow.html).

You can also run this step directly with [Bitrise CLI](https://github.com/bitrise-io/bitrise).

## ⚙️ Configuration

<details>
<summary>Inputs</summary>

| Key | Description | Flags | Default |
| --- | --- | --- | --- |
| `xcode_cache_enabled` | Enables Bitrise Build Cache for Xcode. When enabled, activates Xcelerate and starts the Xcelerate proxy in the background. | required | `true` |
| `gradle_cache_enabled` | Enables Bitrise Build Cache for Gradle. When enabled, activates Gradle cache with analytics and remote cache plugins. | required | `true` |
| `verbose` | Enable logging additional information for troubleshooting | required | `false` |
</details>

<details>
<summary>Outputs</summary>
There are no outputs defined in this step
</details>

## 🙋 Contributing

We welcome [pull requests](https://github.com/bitrise-steplib/bitrise-step-activate-react-native-features/pulls) and [issues](https://github.com/bitrise-steplib/bitrise-step-activate-react-native-features/issues) against this repository.

For pull requests, work on your changes in a forked repository and use the Bitrise CLI to [run step tests locally](https://docs.bitrise.io/en/bitrise-ci/bitrise-cli/running-your-first-local-build-with-the-cli.html).

Learn more about developing steps:

- [Create your own step](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/developing-your-own-bitrise-step/developing-a-new-step.html)
