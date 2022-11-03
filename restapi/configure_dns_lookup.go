// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"dns-lookup-server/dnslookup"
	"dns-lookup-server/models"
	"dns-lookup-server/restapi/operations/history"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	"dns-lookup-server/restapi/operations"
	"dns-lookup-server/restapi/operations/tools"
)

//go:generate swagger generate server --target ../../dns-lookup-server --name DNSLookup --spec ../swagger.json --principal interface{}

func configureFlags(api *operations.DNSLookupAPI) {
	api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{
		{
			ShortDescription: "InfluxDB Configuration",
			LongDescription:  "Configuration of the database to log all successful domain queries",
			Options:          &dnslookup.InfluxDBConfig,
		},
	}
}

func configureAPI(api *operations.DNSLookupAPI) http.Handler {
	domainTools := dnslookup.NewDomainTools()

	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	// api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()

	api.ToolsLookupDomainHandler = tools.LookupDomainHandlerFunc(
		func(params tools.LookupDomainParams) middleware.Responder {
			addresses, err := domainTools.QueryDomain(params.Domain)
			if err != nil {
				return tools.NewLookupDomainBadRequest().WithPayload(&models.UtilsHTTPError{Message: err.Error()})
			}
			response := tools.NewLookupDomainOK().WithPayload(&models.ModelQuery{
				Addresses: addresses,
				ClientIP:  params.HTTPRequest.RemoteAddr,
				CreatedAt: time.Now().Unix(),
				Domain:    params.Domain,
			})
			err = domainTools.SaveDomainOK(response)
			if err != nil {
				/* I decided to not return any results in case the database logging is failing */
				return tools.NewLookupDomainBadRequest().WithPayload(&models.UtilsHTTPError{Message: err.Error()})
			}
			return response
		})
	api.HistoryQueriesHistoryHandler = history.QueriesHistoryHandlerFunc(
		func(params history.QueriesHistoryParams) middleware.Responder {
			response, err := domainTools.DomainQueryHistory()
			if err != nil {
				/* I added a Internal Server Error instead of Bad Request, the error cannot be related to
				the request as there is no parameters */
				return history.NewQueriesHistoryInternalServerError().
					WithPayload(&models.UtilsHTTPError{Message: err.Error()})
			}
			return history.NewQueriesHistoryOK().WithPayload(response)
		})
	api.ToolsValidateIPHandler = tools.ValidateIPHandlerFunc(func(params tools.ValidateIPParams) middleware.Responder {
		isValid := dnslookup.ValidateIpv4(params.Request.IP)
		return tools.NewValidateIPOK().WithPayload(&models.HandlerValidateIPResponse{Status: isValid})
	})

	api.PreServerShutdown = func() {}
	api.ServerShutdown = func() {
		domainTools.Shutdown()
	}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	// TODO: this would be a good place to also register custom metrics for each endpoints
	// Example would be the time it takes for resolving each domain or the time it takes for saving
	// to InfluxDB and the time it takes to read history from it.
	return handler
}

// The middleware configuration happens before anything,
// this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	statusMiddleware := dnslookup.NewStatusMiddleware(handler)
	healthMiddleware := dnslookup.NewHealthMiddleware(statusMiddleware)
	metricMiddleware := dnslookup.NewPrometheusMetricMiddleware(healthMiddleware)
	loggingMiddleware := dnslookup.NewAccessLogMiddleware(metricMiddleware)
	return loggingMiddleware
}
