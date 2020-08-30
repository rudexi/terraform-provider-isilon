resource "isilon_volume_v1" "mydir" {
    path = "mydir"
}

resource "isilon_quota_v2" "mydir" {
    path = isilon_volume_v1.mydir.path
    thresholds = { hard = 1 * 1024*1024*1024 } # 1GB hard threshold
}

resource "isilon_export_v2" "myexport" {
    paths = [ isilon_volume_v1.mydir.absolute_path ]
}
