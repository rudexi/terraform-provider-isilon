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

var testAccExportConfigClients = `
resource "isilon_volume_v1" "mydir2" {
    path = "mydir2"
}

resource "isilon_export_v2" "myexport2" {
    paths   = [ isilon_volume_v1.mydir2.absolute_path ]
    clients = ["192.168.1.10", "192.168.1.11"]
}
`

func TestAccExportCreateClient(t *testing.T) {
    resource.Test(t, resource.TestCase{
        Providers: testAccProviders,
        Steps: []resource.TestStep{
            {
                Config: testAccExportConfigClients,
                Check: resource.ComposeTestCheckFunc(
                    testAccCheckResourceExists("isilon_export_v2.myexport2"),
                ),
            },
        },
    })
}

var testAccExportUpdate1 = `
resource "isilon_volume_v1" "mydir3" {
    path = "mydir3"
}

resource "isilon_export_v2" "myexport3" {
    paths   = [ isilon_volume_v1.mydir3.absolute_path ]
}
`

var testAccExportUpdate2 = `
resource "isilon_volume_v1" "mydir3" {
    path = "mydir3"
}

resource "isilon_export_v2" "myexport3" {
    paths   = [ isilon_volume_v1.mydir3.absolute_path ]
    clients = ["192.168.1.10", "192.168.1.11"]
    root_clients = ["192.168.1.10", "192.168.1.11"]
}
`

func TestAccExportUpdate(t *testing.T) {
    resource.Test(t, resource.TestCase{
        Providers: testAccProviders,
        Steps: []resource.TestStep{
            {
                Config: testAccExportUpdate1,
                Check: resource.ComposeTestCheckFunc(
                    testAccCheckResourceExists("isilon_export_v2.myexport3"),
                ),
            },
            {
                Config: testAccExportUpdate2,
                Check: resource.ComposeTestCheckFunc(
                    testAccCheckResourceExists("isilon_export_v2.myexport3"),
                    resource.TestCheckResourceAttr("isilon_export_v2.myexport3", "clients.0", "192.168.1.10"),
                    resource.TestCheckResourceAttr("isilon_export_v2.myexport3", "clients.1", "192.168.1.11"),
                    resource.TestCheckResourceAttr("isilon_export_v2.myexport3", "root_clients.0", "192.168.1.10"),
                    resource.TestCheckResourceAttr("isilon_export_v2.myexport3", "root_clients.1", "192.168.1.11"),
                ),
            },
        },
    })
}
