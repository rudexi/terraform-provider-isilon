# Isilon Provider

Terraform provider plugin for Dell Isilon (OneFS).

## Example usage

```terraform
provider "isilon" {
    url         = "https://nas.example.com:8080"
    username    = "root"
    group       = "wheel"
    volume_path = "/ifs"
}
```

> The use of `ISILON_PASSWORD` environment variable to pass the
> password to the provider is recommended

## Arguments Reference

The following arguments are supported:
* `url` - (Optional) The URL of the Isilon endpoint. If ommited, the
`ISILON_ENDPOINT` environment variable will be used.
* `username` - (Optional) The username to use to connect to the API. If ommited,
the `ISILON_USERNAME` environment variable will be used.
* `group` - (Optional) The group to be identified as when making directories and
files. It's recommended to use `wheel` if you're using the `root` user. If ommited,
the `ISILON_GROUP` environment variable will be used.
* `password` - (Optional) The password to login with. If ommited, the `ISILON_PASSWORD`
environment variable will be used.
* `skip_ssl_verify` - (Optional) Turn on this option if you're using a self-signed certificate
and you wish to skip TLS certificate verification. If ommited, the `ISILON_SKIP_SSL_VERIFY`
environment variable will be used (`1` for `true`, `0` for `false`). If this variable is ommited
as well, will default to `false`
* `volume_path` - (Optional) The default volume path to start resources with. This is
mainly used for directories. All directories path will be relative to this path,
which makes it practical to dedicate a provider to a subdirectory.
If ommited, the `ISILON_VOLUMEPATH` environment variable will be used. If this variable
is ommited as well, `/ifs` will be used
