package isilon

import (
    "crypto/sha256"
    "fmt"
    "context"
    "errors"
    "strconv"

    "github.com/thecodeteam/goisilon"
    "github.com/thecodeteam/goisilon/api"
    "github.com/thecodeteam/goisilon/api/v2"
)

var (
    exportsPath = "platform/2/protocols/nfs/exports"
)

func hashSum(contents interface{}) string {
    return fmt.Sprintf("%x", sha256.Sum256([]byte(contents.(string))))
}

func zoneParam (zone string) api.OrderedValues {
    if zone == "" {
        return api.NewOrderedValues( [][]string{ []string{} } )
    } else {
        return api.NewOrderedValues( [][]string{ []string{"zone", zone} } )
    }
}

func createExport(
    ctx context.Context,
    client *goisilon.Client,
    zone string,
    export *v2.Export)  (int, error) {

    if export.Paths != nil && len(*export.Paths) == 0 {
        return 0, errors.New("no path set")
    }

    var resp v2.Export
    params := zoneParam(zone)
    err := client.API.Post(
        ctx,
        exportsPath,
        "",
        params,
        nil,
        export,
        &resp)
    if err != nil {
        fmt.Printf("Error in utils's createExport: %+v, %+v, %+v\n", err, params, zone)
        return 0, err
    }

    return resp.ID, nil
}

func getExport(
    ctx context.Context,
    client *goisilon.Client,
    zone string,
    id int,
) (*v2.Export, error) {

    params := zoneParam(zone)
    var resp v2.ExportList
    err := client.API.Get(
        ctx,
        exportsPath,
        strconv.Itoa(id),
        params,
        nil,
        &resp)
    if err != nil { return nil, err }
    if len(resp) == 0 { return nil, nil }

    return resp[0], nil
}

func updateExport(
    ctx context.Context,
    client *goisilon.Client,
    zone string,
    export *v2.Export) error {

    params := zoneParam(zone)
    return client.API.Put(
        ctx,
        exportsPath,
        strconv.Itoa(export.ID),
        params,
        nil,
        export,
        nil)
}

func deleteExport(
    ctx context.Context,
    client *goisilon.Client,
    zone string,
    id int) error {

    params := zoneParam(zone)
    return client.API.Delete(
        ctx,
        exportsPath,
        strconv.Itoa(id),
        params,
        nil,
        nil)
}
