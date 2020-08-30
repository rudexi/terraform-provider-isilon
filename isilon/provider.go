package isilon

import (
    "fmt"
    "context"
    "net/url"

    "github.com/thecodeteam/goisilon"
    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/hashicorp/go-cty/cty"
)

func Provider() *schema.Provider {
    return &schema.Provider{
        ResourcesMap: map[string]*schema.Resource{
            "isilon_volume_v1": volume_v1(),
            "isilon_quota_v2": quota_v2(),
            "isilon_export_v2": export_v2(),
        },
        Schema: map[string]*schema.Schema{
            "url": {
                Type:        schema.TypeString,
                Optional:    true,
                DefaultFunc: schema.EnvDefaultFunc("ISILON_ENDPOINT", nil),
            },
            "username": {
                Type: schema.TypeString,
                Optional: true,
                DefaultFunc: schema.EnvDefaultFunc("ISILON_USERNAME", nil),
            },
            "group": {
                Type: schema.TypeString,
                Optional: true,
                DefaultFunc: schema.EnvDefaultFunc("ISILON_GROUP", nil),
            },
            "password": {
                Type: schema.TypeString,
                Optional: true,
                Sensitive: true,
                StateFunc: hashSum,
                DefaultFunc: schema.EnvDefaultFunc("ISILON_PASSWORD", nil),
            },
            "skip_ssl_verify": {
                Type:        schema.TypeBool,
                Optional:    true,
                DefaultFunc: schema.EnvDefaultFunc("ISILON_SKIP_SSL_VERIFY", "0"),
            },
            "volume_path": {
                Type: schema.TypeString,
                Optional: true,
                DefaultFunc: schema.EnvDefaultFunc("ISILON_VOLUMEPATH", "/ifs"),
            },
        },
        ConfigureContextFunc: configure,
    }

}

func configure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

    var diags diag.Diagnostics

    url, err := url.Parse(d.Get("url").(string))
    if err != nil { return nil, diag.FromErr(err) }

    skip_ssl_verify := d.Get("skip_ssl_verify").(bool)
    username := d.Get("username").(string)
    group := d.Get("group").(string)
    password := d.Get("password").(string)
    volume_path := d.Get("volume_path").(string)

    diags = append(diags, diag.Diagnostic{
        Severity: diag.Warning,
        Summary: "Debugging provider",
        Detail: fmt.Sprintf("data: %s, %s, %s", username, group, volume_path),
        AttributePath: cty.Path{cty.GetAttrStep{Name: "d"}},
    })

    client, err := goisilon.NewClientWithArgs(
        ctx,
        url.String(),
        skip_ssl_verify,
        username,
        group,
        password,
        volume_path,
    )
    if err != nil {
        diags = append(diags, diag.Errorf("Error in provider connection: %s", err)...)
        return nil, diags
    }

    api_version := client.API.APIVersion()
    if api_version == 0 { return nil, diag.Errorf("Could not retrieve API version from Isilon") }

    return client, diags
}
