---
page_title: "morpheus_write_attributes_task Resource - terraform-provider-morpheus"
subcategory: ""
description: |-
  Provides a Morpheus write attributes task resource
---

# morpheus_write_attributes_task

Provides a Morpheus write attributes task resource

## Example Usage

```terraform
resource "morpheus_write_attributes_task" "tfexample_write_attributes" {
  name                = "tfexample_write_attributes"
  code                = "tfexample_write_attributes"
  attributes          = <<EOF
{"demo":"test"}
EOF
  retryable           = true
  retry_count         = 1
  retry_delay_seconds = 10
  allow_custom_config = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the write attributes task

### Optional

- `allow_custom_config` (Boolean) Custom configuration data to pass during the execution of the write attributes task
- `attributes` (String) The attributes payload
- `code` (String) The code of the write attributes task
- `retry_count` (Number) The number of times to retry the task if there is a failure
- `retry_delay_seconds` (Number) The number of seconds to wait between retry attempts
- `retryable` (Boolean) Whether to retry the task if there is a failure

### Read-Only

- `id` (String) The ID of the write attributes task

## Import

Import is supported using the following syntax:

```shell
terraform import morpheus_write_attributes.tfexample_write_attributes 1
```
