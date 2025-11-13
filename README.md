# Activate Gradle Features

[![Step changelog](https://shields.io/github/v/release/bitrise-steplib/bitrise-step-activate-gradle-features?include_prereleases&label=changelog&color=blueviolet)](https://github.com/bitrise-steplib/bitrise-step-activate-gradle-features/releases)

Activates Bitrise features for subsequent Gradle executions in the workflow

<details>
<summary>Description</summary>

This Step activates Bitrise's various features for subsequent Gradle executions in the workflow.

After this Step executes,
- enabling build-cache will result in: Gradle builds will automatically read from the remote cache and push new entries if it's enabled.
- enabling test distribution will result in: Tests will be distributed to Bitrise's remote worker pool.

</details>

## 🧩 Get started

Add this step directly to your workflow in the [Bitrise Workflow Editor](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/steps/adding-steps-to-a-workflow.html).

You can also run this step directly with [Bitrise CLI](https://github.com/bitrise-io/bitrise).

## ⚙️ Configuration

<details>
<summary>Inputs</summary>

| Key | Description | Flags | Default |
| --- | --- | --- | --- |
| `build_cache_enabled` | Enables Gradle build cache for subsequent Gradle executions. When enabled, Gradle builds will automatically read from the remote cache and push new entries if it's enabled. | required | `true` |
| `build_cache_push` | Whether the build can not only read, but write new entries to the remote cache | required | `true` |
| `build_cache_validation_level` | Level of cache entry validation for both uploads and downloads.  Levels: - `none`: no validation. - `warning`: print a warning about invalid cache entries, but don't interrupt the build - `error`: print an error about invalid cache entries and interrupt the build | required | `warning` |
| `test_distribution_enabled` | Enables Gradle Test Distribution for subsequent Gradle executions. When enabled, Gradle tests will automatically split their execution across multiple workers. | required | `false` |
| `test_distribution_shard_size` | Sets the number of tests per shard sent to the Bitrise remote worker pool. | required | `200` |
| `verbose` | Enable logging additional information for troubleshooting | required | `false` |
</details>

<details>
<summary>Outputs</summary>
There are no outputs defined in this step
</details>

## 🙋 Contributing

We welcome [pull requests](https://github.com/bitrise-steplib/bitrise-step-activate-gradle-features/pulls) and [issues](https://github.com/bitrise-steplib/bitrise-step-activate-gradle-features/issues) against this repository.

For pull requests, work on your changes in a forked repository and use the Bitrise CLI to [run step tests locally](https://docs.bitrise.io/en/bitrise-ci/bitrise-cli/running-your-first-local-build-with-the-cli.html).

Learn more about developing steps:

- [Create your own step](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/developing-your-own-bitrise-step/developing-a-new-step.html)
