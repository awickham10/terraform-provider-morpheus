---
page_title: "morpheus_aws_cloud Resource - terraform-provider-morpheus"
subcategory: ""
description: |-
  Provides a Morpheus AWS cloud resource.
---

# morpheus_aws_cloud

Provides a Morpheus AWS cloud resource.

## Example Usage

Creating the AWS cloud with local credentials:

```terraform
resource "morpheus_aws_cloud" "tf_example_aws_cloud" {
  name                       = "tf-aws-demo"
  code                       = "tf-aws-demo"
  location                   = "colorado"
  visibility                 = "public"
  tenant_id                  = 1
  enabled                    = true
  automatically_power_on_vms = true
  region                     = "us-east-1"
  access_key                 = "ADMEI422IMWIF824"
  secret_key                 = "34MPW23DQQFEWNGN112WEG"
  inventory                  = "full"
  vpc                        = "all"
  appliance_url              = "https://morpheus.local"
  time_zone                  = "America/Denver"
  ebs_encryption             = true
  datacenter_id              = "tfawsdemo"
  guidance                   = "manual"
  costing                    = "full"
  agent_install_mode         = "cloudInit"
}
```

Creating the AWS cloud with a credential store credential:

```terraform
data "morpheus_credential" "aws_credentials" {
  name = "awsdemo"
}

resource "morpheus_aws_cloud" "tf_example_aws_cloud" {
  name                       = "tf-aws-demo"
  code                       = "tf-aws-demo"
  location                   = "colorado"
  visibility                 = "public"
  tenant_id                  = 1
  enabled                    = true
  automatically_power_on_vms = true
  region                     = "us-east-1"
  credential_id              = data.morpheus_credential.aws_credentials.id
  inventory                  = "full"
  vpc                        = "all"
  appliance_url              = "https://morpheus.local"
  time_zone                  = "America/Denver"
  ebs_encryption             = true
  datacenter_id              = "tfawsdemo"
  guidance                   = "manual"
  costing                    = "full"
  agent_install_mode         = "cloudInit"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the cloud integration
- `region` (String) The AWS region associated with the cloud integration

### Optional

- `access_key` (String) The AWS access key used for authentication
- `agent_install_mode` (String) The method used to install the Morpheus agent on virtual machines provisioned in the cloud (ssh, cloudInit)
- `appliance_url` (String) The URL used by workloads provisioned in the cloud for interacting with the Morpheus server
- `automatically_power_on_vms` (Boolean) Determines whether to automatically power on cloud virtual machines
- `code` (String) Optional code for use with policies
- `costing` (String) Whether to enable costing on the cloud (off, costing, full)
- `credential_id` (Number) The ID of the credential store entry used for authentication
- `datacenter_id` (String) An arbitrary id used to reference the datacenter for the cloud
- `ebs_encryption` (Boolean) Determines whether to configure default EBS volume encryption or not
- `enabled` (Boolean) Determines whether the cloud is active or not
- `guidance` (String) Whether to enable guidance recommendations on the cloud (manual, off)
- `inventory` (String) Whether to import existing virtual machines (off, basic, full)
- `location` (String) Optional location for the cloud
- `secret_key` (String, Sensitive) The AWS secret key used for authentication
- `tenant_id` (String) The id of the morpheus tenant the cloud is assigned to
- `time_zone` (String) The time zone for the cloud
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- `visibility` (String) Determines whether the cloud is visible in sub-tenants or not
- `vpc` (String) The VPC ID for a specific VPC (all or the AWS VPC id (vpc-25e6dae))

### Read-Only

- `account_number` (String) The AWS account number associated with the cloud integration
- `id` (String) The ID of the cloud

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
- `read` (String)
- `update` (String)

## Import

Import is supported using the following syntax:

```shell
terraform import morpheus_vsphere_cloud.tf_example_aws_cloud 1
```
