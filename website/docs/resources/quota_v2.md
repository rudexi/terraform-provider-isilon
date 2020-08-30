# Quota (API v2) Resource

Manage quota for directories.
At the moment, only the hard threshold is supported.

## Example usage

```terraform
resource "isilon_volume_v1" "mydir" {
    path = "mydir"
}

resource "isilon_quota_v2" "myquota" {
    path       = isilon_volume_v1.mydir.path
    thresholds = { hard = 1 * 1024*1024*1024 }  # 1GB hard threshold
}
```

> It is recommended to use `isilon_volume_v1`'s `path` attribute in order
> to ensure the order

## Argument Reference
* `path` - (Required) The path of the directory which is concerned by quota. Must be relative
to the provider's `volume_path`
* `thresholds` - (Required) A key/value for quotas. At the moment, only the `hard` key is supported

## Import

Quotas can be imported using the path which they apply to:
```bash
terraform import isilon_quota_v2.myquota <path>
```
