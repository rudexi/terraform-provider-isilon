package isilon

import (
    "fmt"
    "context"
    "testing"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
    "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
    "github.com/thecodeteam/goisilon"
)

var testAccQuotaConfig = `
resource "isilon_volume_v1" "mydir" {
    path = "mydir"
}

resource "isilon_quota_v2" "mydir" {
    path           = isilon_volume_v1.mydir.path
    thresholds     = {
        hard = 1024
    }
    depends_on     = [ "isilon_volume_v1.mydir" ]
}
`

var testAccQuotaUpdateConfig = `
resource "isilon_volume_v1" "mydir" {
    path = "mydir"
}

resource "isilon_quota_v2" "mydir" {
    path           = isilon_volume_v1.mydir.path
    thresholds     = {
        hard = 2048
    }
}
`

func TestAccQuotaCreate(t *testing.T) {
    resource.Test(t, resource.TestCase{
        Providers: testAccProviders,
        Steps: []resource.TestStep{
            {
                Config: testAccQuotaConfig,
                Check: resource.ComposeTestCheckFunc(
                    testAccCheckResourceExists("isilon_volume_v1.mydir"),
                    testAccCheckResourceExists("isilon_quota_v2.mydir"),
                ),
            },
            {
                ResourceName: "isilon_volume_v1.mydir",
                ImportState:  true,
            },
        },
    })
}

func TestAccQuotaUpdate(t *testing.T) {
    resource.Test(t, resource.TestCase{
        Providers: testAccProviders,
        Steps: []resource.TestStep{
            {
                Config: testAccQuotaConfig,
                Check: resource.ComposeTestCheckFunc(
                    testAccCheckResourceExists("isilon_volume_v1.mydir"),
                    testAccCheckResourceExists("isilon_quota_v2.mydir"),
                ),
            },
            {
                Config: testAccQuotaUpdateConfig,
                Check: resource.ComposeTestCheckFunc(
                    testAccCheckResourceExists("isilon_volume_v1.mydir"),
                    testAccCheckResourceExists("isilon_quota_v2.mydir"),
                    resource.TestCheckResourceAttr("isilon_quota_v2.mydir", "thresholds.hard", "2048"),
                ),
            },
        },
    })
}

func testAccCheckResourceExists(resourceName string) resource.TestCheckFunc {
    return func(s *terraform.State) error {
        rs, ok := s.RootModule().Resources[resourceName]
        if !ok { return fmt.Errorf("Not found: %s", resourceName) }

        if rs.Primary.ID == "" {
            return fmt.Errorf("No ID set")
        }

        return nil
    }
}

func testAccCheckHardThreshold(resourceName string, hard_threshold int) resource.TestCheckFunc {
    return func(s *terraform.State) error {
        rs, ok := s.RootModule().Resources[resourceName]
        if !ok { return fmt.Errorf("Not found: %s", resourceName) }

        if rs.Primary.ID == "" {
            return fmt.Errorf("No ID set")
        }

        client := testAccProvider.Meta().(*goisilon.Client)
        quota, err := client.GetQuota(context.Background(), rs.Primary.ID)
        if err != nil { return fmt.Errorf("error: %s", err) }

        if quota.Thresholds.Hard != int64(hard_threshold) {
            return fmt.Errorf("Quota is different from updated one: %d bits", quota.Thresholds.Hard)
        }

        return nil
    }
}
