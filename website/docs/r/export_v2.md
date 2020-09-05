# Export v2 Resource

Manage NFS exports in Isilon

## Example usage

```terraform
resource "isilon_volume_v1" "mydir" {
    path = "mydir"
}

resource "isilon_export_v2" "myexport" {
    zone      = "myzone"
    paths     = [
        isilon_volume_v1.mydir.absolute_path
    ]
}
```

> It is recommended to use `isilon_volume_v1`'s `absolute_path` attribute
> to ensure the order between the directory and the export

## Argument Reference
* `paths` - (Required) An array of absolute paths of the directories to export
* `zone` - (Optional) The access zone of the NFS export. If ommited,
Isilon will select the default access zone
* `clients` (Optional) An array of IPs representing the allowed clients
* `clients` (Optional) An array of IPs representing the allowed clients which
will be mapped as root for the export

## Import

Exports can be imported by specifying the access zone and their ID:
```bash
terraform import isilon_export_v2.myexport <access_zone>/<id>
```
