package isilon

import (
    "context"
    "strings"
    "strconv"

    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/thecodeteam/goisilon"
    "github.com/thecodeteam/goisilon/api/v2"
)

func export_v2() *schema.Resource {
    return &schema.Resource{
        CreateContext: createExportv2,
        ReadContext: readExportv2,
        UpdateContext: updateExportv2,
        DeleteContext: deleteExportv2,
        Importer: &schema.ResourceImporter{
            State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
                splits := strings.Split(d.Id(), "/")
                zone := splits[0]
                id := splits[1]
                d.Set("zone", zone)
                d.SetId(id)
                return []*schema.ResourceData{d}, nil
            },
        },
        Schema: map[string]*schema.Schema{
            "paths": &schema.Schema{
                Type: schema.TypeList,
                Required: true,
                Elem: &schema.Schema{ Type: schema.TypeString },
            },
            "zone": &schema.Schema{
                Type: schema.TypeString,
                Optional: true,
                ForceNew: true,
                Default: "",
            },
            "clients": &schema.Schema{
                Type: schema.TypeList,
                Optional: true,
                Elem: &schema.Schema{ Type: schema.TypeString },
            },
        },
    }
}

func expandElements(elements []interface{}) []string {
    var output_elements []string
    for _, element := range elements {
        output_elements = append(output_elements, element.(string))
    }
    return output_elements
}

func createExportv2 (ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    client := meta.(*goisilon.Client)

    var diags diag.Diagnostics

    paths := expandElements(d.Get("paths").([]interface{}))

    clients, ok := d.GetOk("clients")
    var real_clients []string
    if !ok {
        real_clients = []string{}
    } else {
        real_clients = expandElements(clients.([]interface{}))
    }

    export := &v2.Export{
        Paths: &paths,
        Clients: &real_clients,
    }

    if export.Paths != nil && len(*export.Paths) == 0 {
        return diag.Errorf("No path set")
    }

    zone := d.Get("zone").(string)
    var resp v2.Export
    params := zoneParam(zone)
    err := client.API.Post(
        context.Background(),
        "platform/2/protocols/nfs/exports",
        "",
        params,
        nil,
        export,
        &resp)
    if err != nil {
        return diag.FromErr(err)
    }

    id_str := strconv.Itoa(resp.ID)
    d.SetId(id_str)

    return diags
}

func readExportv2 (ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    client := meta.(*goisilon.Client)

    var diags diag.Diagnostics

    id, err := strconv.Atoi(d.Id())
    if err != nil { return diag.FromErr(err) }

    zone := d.Get("zone").(string)
    export, err := getExport(ctx, client, zone, id)
    if err != nil { return diag.FromErr(err) }

    d.Set("paths", export.Paths)
    d.Set("clients", export.Clients)

    return diags
}

func updateExportv2 (ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    client := meta.(*goisilon.Client)

    var diags diag.Diagnostics

    id, err := strconv.Atoi(d.Id())
    if err != nil { return diag.FromErr(err) }

    paths := expandElements(d.Get("paths").([]interface{}))

    clients, ok := d.GetOk("clients")
    var real_clients []string
    if !ok {
        real_clients = []string{}
    } else {
        real_clients = expandElements(clients.([]interface{}))
    }

    export := &v2.Export{
        ID: id,
        Paths: &paths,
        Clients: &real_clients,
    }

    zone := d.Get("zone").(string)
    err = updateExport(ctx, client, zone, export)
    if err != nil { return diag.FromErr(err) }

    return diags
}

func deleteExportv2 (ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    client := meta.(*goisilon.Client)

    var diags diag.Diagnostics

    id, err := strconv.Atoi(d.Id())
    if err != nil { return diag.FromErr(err) }

    zone := d.Get("zone").(string)
    err = deleteExport(ctx, client, zone, id)
    if err != nil { return diag.FromErr(err) }

    return diags
}
