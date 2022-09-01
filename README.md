# Terraform Provider OnFinality (Terraform Plugin Framework)

_This repository is built on the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework). 


These files contain boilerplate code that you will need to edit to create your own Terraform provider. Tutorials for creating Terraform providers can be found on the [HashiCorp Learn](https://learn.hashicorp.com/collections/terraform/providers) platform. _Terraform Plugin Framework specific guides are titled accordingly._

Please see the [GitHub template repository documentation](https://help.github.com/en/github/creating-cloning-and-archiving-repositories/creating-a-repository-from-a-template) for how to create a new repository from this template on GitHub.

Once you've written your provider, you'll want to [publish it on the Terraform Registry](https://www.terraform.io/docs/registry/providers/publishing.html) so that others can use it.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.18

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Using the provider

### OnFinality Platform Credentials
1. Generate your access_key and secret_key from app.onfinality.io Account Settings Page
2. Define variable and config the provider
```terraform
variable "onf_access_key" {}
variable "onf_secret_key" {}

provider "onfinality" {
  access_key = var.onf_access_key
  secret_key = var.onf_secret_key
}
```
3. Export them to env
```
export TF_VAR_onf_access_key=...
export TF_VAR_onf_secret_key=...
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```
