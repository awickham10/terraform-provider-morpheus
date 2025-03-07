---
page_title: "morpheus_ruby_script_task Resource - terraform-provider-morpheus"
subcategory: ""
description: |-
  Provides a Morpheus ruby script task resource
---

# morpheus_ruby_script_task

Provides a Morpheus ruby script task resource

## Example Usage

Creating the ruby script task with local script content:

```terraform
resource "morpheus_ruby_script_task" "tfexample_ruby_local" {
  name                = "tfexample_ruby_local"
  code                = "tfexample_ruby_local"
  source_type         = "local"
  script_content      = <<EOF
  puts "testing"
EOF
  retryable           = true
  retry_count         = 1
  retry_delay_seconds = 10
  allow_custom_config = true
}
```

Creating the ruby script task with the script fetched from a url:

```terraform
resource "morpheus_ruby_script_task" "tfexample_ruby_url" {
  name                = "tfexample_ruby_url"
  code                = "tfexample_ruby_url"
  source_type         = "url"
  result_type         = "json"
  script_path         = "https://example.com/example.rb"
  retryable           = true
  retry_count         = 1
  retry_delay_seconds = 10
  allow_custom_config = true
}
```

Creating the ruby script task with the script fetch via git:

```terraform
resource "morpheus_ruby_script_task" "tfexample_ruby_git" {
  name                = "tfexample_ruby_git"
  code                = "tfexample_ruby_git"
  source_type         = "repository"
  result_type         = "json"
  script_path         = "example.rb"
  version_ref         = "master"
  repository_id       = 1
  retryable           = true
  retry_count         = 1
  retry_delay_seconds = 10
  allow_custom_config = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the ruby script task
- `source_type` (String) The source of the ruby script (local, url or repository)

### Optional

- `allow_custom_config` (Boolean) Custom configuration data to pass during the execution of the ruby script
- `code` (String) The code of the ruby script task
- `repository_id` (Number) The ID of the git repository integration
- `result_type` (String) The expected result type (single value, key pairs, json)
- `retry_count` (Number) The number of times to retry the task if there is a failure
- `retry_delay_seconds` (Number) The number of seconds to wait between retry attempts
- `retryable` (Boolean) Whether to retry the task if there is a failure
- `script_content` (String) The content of the ruby script. Used when the local source type is specified
- `script_path` (String) The path of the ruby script, either the url or the path in the repository
- `version_ref` (String) The git reference of the repository to pull (main, master, etc.)

### Read-Only

- `id` (String) The ID of the ruby script task

## Import

Import is supported using the following syntax:

```shell
terraform import morpheus_ruby_script_task.tf_example_ruby_task_script 1
```
