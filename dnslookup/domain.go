package dnslookup

import (
	"context"
	"dns-lookup-server/models"
	"dns-lookup-server/restapi/operations/tools"
	"encoding/json"
	"fmt"
	"net"
	"regexp"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

const numberOfResults = 20

/* Could be dynamic with bucket & limit from numberOfResults & shared with write. */
const domainHistoryQuery = `from(bucket: "successful-queries")
			|> range(start: -1w)
            |> filter(fn: (r) => r._measurement == "log")
			|> filter(fn: (r) => r["_field"] == "json")
			|> limit(n: 20)
			|> sort()`

type DomainTools struct {
	databaseClient influxdb2.Client
}

func NewDomainTools() DomainTools {
	client := influxdb2.NewClient(InfluxDBConfig.URL, InfluxDBConfig.Token)
	_, err := client.Ping(context.Background())
	if err != nil {
		// I prefer to have a fail-fast program in case there is a configuration error
		panic(fmt.Errorf("fail to connect to the database: %v", err))
	}
	fmt.Printf("saving query logs to %s\n", InfluxDBConfig.URL)
	return DomainTools{
		databaseClient: client,
	}
}

func (dt DomainTools) QueryDomain(domain string) ([]*models.ModelAddress, error) {
	/* We should be using a more complete validation of the input domain */
	simpleDomainRegex := regexp.MustCompile(`^(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]$`)
	if !simpleDomainRegex.MatchString(domain) {
		return nil, fmt.Errorf("the domain %s is invalid", domain)
	}

	addresses, err := net.DefaultResolver.LookupIP(context.Background(), "ip4", domain)
	if err != nil {
		return nil, err
	}

	modelAddresses := make([]*models.ModelAddress, len(addresses))
	for i, address := range addresses {
		modelAddresses[i] = &models.ModelAddress{IP: address.String()}
	}

	return modelAddresses, nil
}

func (dt DomainTools) SaveDomainOK(lookupDomainOK *tools.LookupDomainOK) error {
	writeAPI := dt.databaseClient.WriteAPIBlocking(InfluxDBConfig.Org, "successful-queries")

	/*
	 * I used this json marshalling to speed up the encoding & decoding since all
	 * the model is already correctly annotated thanks to swagger
	 * In a real world, it would be better to store in different format.
	 */
	rawJSON, err := json.Marshal(lookupDomainOK.Payload)
	if err != nil {
		return err
	}

	fields := map[string]interface{}{
		"address":  lookupDomainOK.Payload.Addresses,
		"clientIP": lookupDomainOK.Payload.ClientIP,
		"domain":   lookupDomainOK.Payload.Domain,
		"json":     rawJSON,
	}

	point := write.NewPoint("log", nil, fields, time.Unix(lookupDomainOK.Payload.CreatedAt, 0))
	return writeAPI.WritePoint(context.Background(), point)
}

func (dt DomainTools) DomainQueryHistory() ([]*models.ModelQuery, error) {
	queryAPI := dt.databaseClient.QueryAPI(InfluxDBConfig.Org)
	results, err := queryAPI.Query(context.Background(), domainHistoryQuery)
	if err != nil {
		return nil, err
	}

	i := 0
	responses := make([]*models.ModelQuery, 0, numberOfResults)
	for results.Next() {
		record := results.Record()
		modelQuery := &models.ModelQuery{}
		if rawJSON, ok := record.Value().(string); ok {
			err = json.Unmarshal([]byte(rawJSON), &modelQuery)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("failed to cast record value")
		}
		responses = append(responses, modelQuery)
		i++
	}
	return responses, nil
}

func (dt DomainTools) Shutdown() {
	if dt.databaseClient != nil {
		dt.databaseClient.Close()
	}
}
