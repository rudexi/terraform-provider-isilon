package isilon

import (
    "path"
    "context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/thecodeteam/goisilon"
    "github.com/thecodeteam/goisilon/api"
)

var (
    aclQS               = api.OrderedValues{{[]byte("acl")}}
    metadataQS          = api.OrderedValues{{[]byte("metadataQs")}}
    createVolumeHeaders = map[string]string{
            "x-isi-ifs-target-type":    "container",
            "x-isi-ifs-access-control": "public_read_write",
        }
)

func volume_v1() *schema.Resource {
    return &schema.Resource{
        CreateContext: createVolume,
        ReadContext:   readVolume,
        UpdateContext: updateVolume,
        DeleteContext: deleteVolume,
        Importer: &schema.ResourceImporter{ State: schema.ImportStatePassthrough },
        Schema: map[string]*schema.Schema{
            "path": &schema.Schema{
                Type:     schema.TypeString,
                Required: true,
                ForceNew: true,
            },
            "mode": &schema.Schema{
                Type:     schema.TypeString,
                Optional: true,
                Default:  "0777",
            },
            "is_hidden": &schema.Schema{
                Type:     schema.TypeBool,
                Optional: true,
                Default:  false,
            },
            "absolute_path": &schema.Schema{
                Type:     schema.TypeString,
                Optional: true,
                Computed: true,
            },
        },
    }
}

func createVolume (ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    client := meta.(*goisilon.Client)

    var diags diag.Diagnostics

    p := d.Get("path").(string)

    _, err := client.CreateVolume(ctx, p)
    if err != nil { return diag.FromErr(err) }

    d.SetId(p)
    d.Set("absolute_path", path.Join(client.API.VolumesPath(), p))
    return diags
}

func readVolume (ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    client := meta.(*goisilon.Client)

    var diags diag.Diagnostics

    p := d.Get("path").(string)

    volume, err := client.GetVolume(context.Background(), p, "")
    if err != nil { return diag.FromErr(err) }

    d.Set("path", volume.Name)
    for _, attribute := range volume.AttributeMap {
        name := attribute.Name
        value := attribute.Value
        switch name {
            case "mode": d.Set("mode", value.(string))
            case "is_hidden": d.Set("is_hidden", value.(bool))
        }
    }

    d.Set("absolute_path", path.Join(client.API.VolumesPath(), p))

    return diags
}

type attribute struct {
    Name  string      `json:"name"`
    Value interface{} `json:"value"`
}

type updateIsiVolumeAttributes struct {
    Action string `json:"action"`
    AttributeMap []attribute `json:"attrs"`
}

func updateVolume (ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    client := meta.(*goisilon.Client)
    path_name := d.Get("path").(string)

    if d.HasChange("mode") || d.HasChange("is_hidden") {
        var data = &updateIsiVolumeAttributes{
            Action: "update",
            AttributeMap: []attribute {
                attribute{Name: "mode",      Value: d.Get("mode").(string)},
                attribute{Name: "is_hidden", Value: d.Get("is_hidden").(string)},
            },
        }

        err := client.API.Put(
            context.Background(),
            path.Join("namespace", client.API.VolumesPath()),
            path_name,
            metadataQS,
            nil,
            data,
            nil,
        )
        if err != nil { return diag.FromErr(err) }

    }

    return readVolume(ctx, d, meta)
}

func deleteVolume (ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    client := meta.(*goisilon.Client)

    var diags diag.Diagnostics

    path := d.Get("path").(string)

    err := client.DeleteVolume(context.Background(), path)
    if err != nil { return diag.FromErr(err) }

    d.SetId("")
    return diags
}
