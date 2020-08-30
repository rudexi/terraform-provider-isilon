package isilon

import (
    "testing"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testAccExportConfig = `
resource "isilon_volume_v1" "mydir" {
    path = "mydir"
}

resource "isilon_export_v2" "myexport" {
    paths = [ isilon_volume_v1.mydir.absolute_path ]
}
`

func TestAccExportCreate(t *testing.T) {
    resource.Test(t, resource.TestCase{
        Providers: testAccProviders,
        Steps: []resource.TestStep{
            {
                Config: testAccExportConfig,
                Check: resource.ComposeTestCheckFunc(
                    testAccCheckResourceExists("isilon_export_v2.myexport"),
                ),
            },
            {
                ResourceName: "isilon_volume_v1.mydir",
                ImportState:  true,
            },
        },
    })
}
