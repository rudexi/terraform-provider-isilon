# Terraform Isilon Provider

Terraform provider for Dell Isilon (OneFS) API.

## Requirements

* Terraform 0.13.x
* Go 1.13

## Documentation

* [Provider documentation](./website/docs/index.md)
* [Resources documentation](./website/docs/resources)

## Scope

The following is supported:
* Directories (`volume_v1`)
* Quotas
* NFS exports

There is the following limitations:
* Clients/user mapping for NFS exports is not supported yet (WIP)
* Only hard threshold for quota can be set, no other options (WIP)
* The mode for the directory might need to be adjusted to the one of
the parent directory, otherwise it will keep being changed

See provider and resources documentation for examples of usage.

## Building the provider

Using [gvm](https://github.com/moovweb/gvm) or equivalent:
```bash
gvm install go1.13
gvm use go1.13
go build -o terraform-provider-isilon
mkdir -p ~/.terraform.d/plugins/local/isilon/0.1/linux_amd64
mv terraform-provider-isilon ~/.terraform.d/plugins/local/isilon/0.1/linux_amd64
```


## Testing

### Automated testing

Basic unit testing:
```bash
go test -v ./isilon
```

Acceptance testing (require a working Isilon API):
```bash
TF_ACC=1 go test -v ./isilon
```

> Note that it is required to set the ISILON_* environment variables in order to
> run the acceptance testing. See the provider documentation for more info on
> what variables need to be defined.

### Manual testing

Create a file in `examples/provider.tf` containing your Isilon API information:
```hcl
terraform {
    required_providers {
        isilon = {
            source = "local/isilon"
            version = "0.1"
        }
    }
}

provider "isilon" {
    url         = "https://nas.example.com:8080"
    username    = "root"
    group       = "wheel"
    password    = "secret" # ISILON_PASSWORD env variable
    volume_path = "/ifs"
}
```

Then run the following
```
cd examples
terraform init
terraform plan
terraform apply
terraform destroy
```
