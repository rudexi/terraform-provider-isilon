# Volume (API v1) Resource

## Example usage

```terraform
resource "isilon_volume_v1" "mydir" {
    path = "mydir"
}
```

> Note that the `path` argument is relative to the provider's `volume_path` argument

## Argument reference
The following arguments are supported:
* `path` - (Required) The path of the directory relative to the provider's `volume_path`
* `mode` - (Optional) The permission mode of the directory. Defaults to `0777`. Note
that this might need to match the parent directory's mode
* `is_hidden` - (Optional) Whether the directory is hidden or visible. Defaults to `false`

## Attribute reference
* `path` - The path relative to `volume_path`
* `absolute_path` - The absolute path of the directory

## Import

Directory import can be done using its path:
```bash
terraform import isilon_volume_v1.mydir <path>
```
