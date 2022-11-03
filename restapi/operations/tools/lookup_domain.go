// Code generated by go-swagger; DO NOT EDIT.

package tools

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// LookupDomainHandlerFunc turns a function with the right signature into a lookup domain handler
type LookupDomainHandlerFunc func(LookupDomainParams) middleware.Responder

// Handle executing the request and returning a response
func (fn LookupDomainHandlerFunc) Handle(params LookupDomainParams) middleware.Responder {
	return fn(params)
}

// LookupDomainHandler interface for that can handle valid lookup domain params
type LookupDomainHandler interface {
	Handle(LookupDomainParams) middleware.Responder
}

// NewLookupDomain creates a new http.Handler for the lookup domain operation
func NewLookupDomain(ctx *middleware.Context, handler LookupDomainHandler) *LookupDomain {
	return &LookupDomain{Context: ctx, Handler: handler}
}

/*
	LookupDomain swagger:route GET /tools/lookup tools lookupDomain

# Lookup domain

Lookup domain and return all IPv4 addresses
*/
type LookupDomain struct {
	Context *middleware.Context
	Handler LookupDomainHandler
}

func (o *LookupDomain) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewLookupDomainParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
