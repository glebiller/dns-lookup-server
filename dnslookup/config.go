package dnslookup

const httpServerError = 500
const httpGetMethod = "GET"

//nolint:gochecknoglobals // linting ignored for simplicity.
var InfluxDBConfig = struct {
	URL   string `long:"influxdb-url" description:"URL of database" env:"INFLUXDB_URL" default:"http://localhost:8086"`
	Org   string `long:"influxdb-org" description:"Organization to store data to" env:"INFLUXDB_ORG"`
	Token string `long:"influxdb-token" description:"Write access token to the bucket" env:"INFLUXDB_TOKEN"`
}{}
