# pango_utils

![CI/CD](https://github.com/mrichardson03/pango_utils/workflows/CI/CD/badge.svg)

Version: 0.0.1

Utility programs using the [pango](https://github.com/PaloAltoNetworks/pango) library.
Most of these are just built and repackaged from the source code on the 
[PAN-OS Terraform provider documentation](https://registry.terraform.io/providers/PaloAltoNetworks/panos/latest/docs)
page, written by Garfield Freeman ([@shinmog](https://github.com/shinmog)).

## Usage

Connection information for all utilities is expected in the following environment
variables:

- `PANOS_HOSTNAME`
- `PANOS_USERNAME`
- `PANOS_PASSWORD`

### panos_init

When deploying a VM-Series in AWS or GCP, the `admin` account will initially have no
password set, but authentication will be allowed using a SSH key.  This utility will
authenticate using the SSH key and set the password to be the contents of the
`PANOS_PASSWORD` environment variable.

### Example

```bash
$ panos_init my_ssh_key.pem
```

### make_api_key

API connections to PAN-OS require an API key. If you do not provide the API key to the
PAN-OS provider, the API key is generated before every single API call. Thus, some
slight speed gains can be realized in the panos provider by specifying the API key
instead of the username/password combo. The following may be used to generate the API
key, and will write to stdout:

```
$ make_api_key
```

### commit

Terraform does not provide native support for committing configuration to PAN-OS, so
this utility will issue the commit API command.

The first argument on the CLI will be used as a commit comment, but is optional.

```bash
$ commit 'My commit comment'
```

Commits can be issued in Terraform using a null resource:

```terraform
resource "null_resource" "commit_panorama" {
    provisioner "local-exec" {
        command = "commit 'Commit from Terraform'"
    }
    depends_on = [
        module.policies.security_rule_group,
        module.policies.nat_rule_group
    ]
}
```
