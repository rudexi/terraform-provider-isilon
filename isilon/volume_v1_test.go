package isilon

import (
    "fmt"
    "testing"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
    "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAccVolumeConfig = `
resource "isilon_volume_v1" "mydir" {
    path = "mydir"
}
`

func TestAccVolumeCreate(t *testing.T) {
    resource.Test(t, resource.TestCase{
            Providers: testAccProviders,
            Steps: []resource.TestStep{
                {
                    Config: testAccVolumeConfig,
                    Check: resource.ComposeTestCheckFunc(
                        testAccCheckVolumeExists("isilon_volume_v1.mydir"),
                        resource.TestCheckResourceAttrSet("isilon_volume_v1.mydir", "absolute_path"),
                    ),
                },
                {
                    ResourceName: "isilon_volume_v1.mydir",
                    ImportState:  true,
                },
            },
    })
}

func testAccCheckVolumeExists(resourceName string) resource.TestCheckFunc {
    return func(s *terraform.State) error {
        rs, ok := s.RootModule().Resources[resourceName]
        if !ok { return fmt.Errorf("Not found: %s", resourceName) }

        if rs.Primary.ID == "" {
            return fmt.Errorf("No ID set")
        }

        return nil
    }
}
