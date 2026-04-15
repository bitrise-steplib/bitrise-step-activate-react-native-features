# React Native Build Cache - Setup Guide

## Introduction

Bitrise Build Cache for React Native speeds up your CI builds by caching native compilation artifacts across builds. It covers all three native build systems used in a React Native project:

- **Gradle** — Android build outputs (compiled classes, resources, dex files) are shared via Bitrise's remote build cache
- **Xcode** — iOS compilation results are cached via Bitrise's remote build cache, the same backend that is already available for [Gradle](https://devcenter.bitrise.io/en/dependencies-and-caching/remote-build-caching/remote-build-cache-for-gradle.html) and [Bazel](https://devcenter.bitrise.io/en/dependencies-and-caching/remote-build-caching/remote-build-cache-for-bazel.html)
- **C++ native modules** — Compiled native bridge code and third-party native modules are cached via ccache, stored in Bitrise's remote build cache, sharing cache entries across builds

## What are the benefits?

By caching native compilation outputs, subsequent CI builds can skip recompiling unchanged native code. This is especially effective for React Native projects where the native layer changes infrequently compared to the JS layer.

Depending on how much of your native code remains unchanged between builds, you can expect significant reductions in build times for both Android and iOS.

## Requirements

- **Active Bitrise Build Cache Trial or Subscription**
- A React Native project connected to Bitrise

If you try to use the Bitrise Build Cache without an active trial or subscription, the Activate step will report an error.

## How to set up - Workflow configuration

### 0. Ensure that you have an active Bitrise Build Cache subscription or trial

You can check your subscription status on the [Bitrise Build Cache page](https://app.bitrise.io/build-cache/).

### 1. Ensure your build workflow works without Build Cache first

Before configuring Bitrise Build Cache, make sure your existing React Native workflow builds successfully. Run a build with your current setup to confirm everything is working, so you have a clear baseline to compare against.

### 2. Add the "Build Cache for React Native" step to the workflow

Add the **Build Cache for React Native** step to the workflow, **before** any step that triggers a native build (before your Script steps that run `yarn`, `npm`, `npx`, `expo`, `fastlane`, or direct Gradle/Xcode commands).

The step has two main inputs:

| Input | Default | Description |
|-------|---------|-------------|
| **Enable Xcode cache** | `true` | Activates Bitrise Build Cache for iOS builds via Xcelerate. A background proxy is started automatically. |
| **Enable Gradle cache** | `true` | Activates Bitrise Build Cache for Android builds. C++ native modules are also cached via ccache, with a background storage helper started automatically. |

If your workflow only builds for one platform, disable the other:

- **Android-only workflow**: Set *Enable Xcode cache* to `false`
- **iOS-only workflow**: Set *Enable Gradle cache* to `false`

### 3. Wrap your build commands with the CLI

After the Activate step runs, the `bitrise-build-cache` CLI is available on the machine. In your Script steps, prefix any command that triggers a native build with `bitrise-build-cache react-native run`.

This wrapper:

- Ensures the ccache storage helper is running during the build
- Tracks cache hit rates and build analytics
- Runs your command exactly as-is — all arguments, stdin, stdout, and exit codes are preserved

**Before:**

```bash
npx react-native run-android
```

**After:**

```bash
bitrise-build-cache react-native run npx react-native run-android
```

This works with any package manager or build tool:

```bash
# yarn
bitrise-build-cache react-native run yarn build:android

# npm
bitrise-build-cache react-native run npm run build:ios

# expo
bitrise-build-cache react-native run expo build:ios

# pnpm
bitrise-build-cache react-native run pnpm run build:android

# fastlane
bitrise-build-cache react-native run fastlane beta

# Direct Gradle invocation
bitrise-build-cache react-native run ./gradlew assembleRelease
```

#### Example workflow YAML

```yaml
workflows:
  build-react-native:
    steps:
    - activate-ssh-key@4: {}
    - git-clone@8: {}
    #
    # Add the Activate step **before** your build steps
    - activate-build-cache-for-react-native@0: {}
    #
    # Install JS dependencies (no wrapping needed here)
    - script@1:
        title: Install dependencies
        inputs:
        - content: yarn install
    #
    # Wrap commands that trigger native builds with the CLI
    - script@1:
        title: Build Android
        inputs:
        - content: bitrise-build-cache react-native run npx react-native run-android --mode=release
    - script@1:
        title: Build iOS
        inputs:
        - content: bitrise-build-cache react-native run npx react-native run-ios --configuration=Release
    #
    - deploy-to-bitrise-io@2: {}
```

### 4. Do a test build to validate the setup works

You should see the **Build Cache for React Native** step complete successfully.

After the build finishes, check the **Build Cache** tab on the build details page. You should see the commands that were executed with build cache enabled.

The first run will always have a **0% cache hit rate**, as the cache is still empty at this point. This is expected.

### 5. Perform additional builds to warm up the cache

After running a few additional builds (usually 1-3), the cache will warm up and you should start seeing cache hit rates greater than 0%.

You can monitor cache performance on the build's details page and on the [Build Cache list page](https://app.bitrise.io/build-cache/).

### Platform-specific workflows

If you have separate workflows for Android and iOS, add the Activate step to each one and disable the platform that doesn't apply:

**Android workflow:**

```yaml
- activate-build-cache-for-react-native@0:
    inputs:
    - xcode_cache_enabled: "false"
```

**iOS workflow:**

```yaml
- activate-build-cache-for-react-native@0:
    inputs:
    - gradle_cache_enabled: "false"
```

### Endpoint for Private Cloud customers

Customers on Bitrise Private Cloud need to set the following Environment Variable (we suggest setting it as a **Project Environment Variable** so you don't have to set it for every workflow):

- **Key:** `BITRISE_BUILD_CACHE_ENDPOINT`
- **Value:** `grpcs://remote-build-cache.services.bitrise.io`

Bitrise YML example:

```yaml
app:
  envs:
  - BITRISE_BUILD_CACHE_ENDPOINT: grpcs://remote-build-cache.services.bitrise.io
    opts:
      is_expand: false
```

This is only necessary for builds running on Bitrise Private Cloud; builds on Bitrise Public or Dedicated Cloud use the default endpoint.

## FAQ

### What exactly gets cached?

- **Android:** Gradle task outputs (compilation, resource processing, dex generation) via remote build cache
- **iOS:** Xcode compilation outputs (object files, module artifacts) via LLVM CAS cache
- **C++ native modules:** Compiled native bridge code and third-party native modules via ccache

### What does NOT get cached?

- **Metro JS bundling** — the JavaScript bundling step is not affected by this setup
- **node_modules** — package installation (yarn/npm) is not cached by this tool. Use Bitrise's standard caching steps for node_modules if needed.

### Do I need to wrap every command?

Only wrap commands that trigger **native builds**. You do not need to wrap:
- `yarn install` / `npm install` (dependency installation)
- `yarn test` / `npm test` (JS-only tests)
- Other commands that don't invoke Gradle or Xcode

Wrap commands like `npx react-native run-android`, `npx react-native run-ios`, `./gradlew assembleRelease`, `fastlane build`, or any script that ultimately calls `xcodebuild` or Gradle.

### Will this speed up my tests?

Build Cache reduces **compilation** times. If your test workflow includes a build step (e.g., `xcode-build-for-test`), that step will be faster. However, the actual test execution time is not affected.

### Can I use this alongside standalone Gradle or Xcode Build Cache steps?

The Build Cache for React Native step configures caching for all three platforms (Gradle, Xcode, C++) in one go. If you are already using a standalone Gradle or Xcode Build Cache step, replace it with this step to avoid conflicting configurations.

### How do I troubleshoot issues?

Enable verbose logging by setting the **Verbose logging** input to `true` on the Activate step. This logs additional details about cache configuration and background service startup.
