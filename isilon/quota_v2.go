package isilon

import (
    "context"

    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/thecodeteam/goisilon"
)

func quota_v2() *schema.Resource {
    return &schema.Resource{
        CreateContext: createQuota,
        ReadContext: readQuota,
        UpdateContext: updateQuota,
        DeleteContext: deleteQuota,
        Importer: &schema.ResourceImporter{ State: schema.ImportStatePassthrough },
        Schema: map[string]*schema.Schema{
            "path": &schema.Schema{
                Type: schema.TypeString,
                Required: true,
                ForceNew: true,
            },
            "thresholds": &schema.Schema{
                Type: schema.TypeMap,
                Optional: true,
                Elem: &schema.Schema{
                    Type: schema.TypeInt,
                    Optional: true,
                },
            },
        },
    }
}

func createQuota (ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    client := meta.(*goisilon.Client)

    var diags diag.Diagnostics

    path := d.Get("path").(string)
    hard_threshold := int64(d.Get("thresholds").(map[string]interface{})["hard"].(int))
    err := client.SetQuotaSize(ctx, path, hard_threshold)
    if err != nil { return diag.FromErr(err) }

    d.SetId(path)
    return diags
}

func readQuota (ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    client := meta.(*goisilon.Client)

    var diags diag.Diagnostics

    path := d.Get("path").(string)

    quota, err := client.GetQuota(ctx, path)
    if err != nil { return diag.FromErr(err) }

    d.Set("path", path)
    d.Set("thresholds", map[string]interface{}{"hard": int(quota.Thresholds.Hard)})

    return diags
}

func updateQuota (ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    client := meta.(*goisilon.Client)

    var diags diag.Diagnostics

    path := d.Get("path").(string)
    hard_threshold := int64(d.Get("thresholds").(map[string]interface{})["hard"].(int))

    if d.HasChange("thresholds") {
        err := client.UpdateQuotaSize(ctx, path, hard_threshold)
        if err != nil { return diag.FromErr(err) }
    }

    return diags
}

func deleteQuota (ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    client := meta.(*goisilon.Client)

    var diags diag.Diagnostics

    path := d.Get("path").(string)

    err := client.ClearQuota(ctx, path)
    if err != nil { return diag.FromErr(err) }

    return diags
}
