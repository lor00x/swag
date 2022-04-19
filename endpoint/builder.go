// Copyright 2020 zc2638
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package endpoint

import (
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/zc2638/swag"
)

// Builder uses the builder pattern to generate swagger endpoint definitions
type Builder struct {
	Endpoint *swag.Endpoint
}

// Option represents a functional option to customize the swagger endpoint
type Option func(builder *Builder)

// Apply improves the readability of applied options
func (o Option) Apply(builder *Builder) {
	o(builder)
}

// Handler allows an instance of the web handler to be associated with the endpoint.  This can be especially useful when
// using swag to bind the endpoints to the web router.  See the examples package for how the Handler can be used in
// conjunction with Walk to simplify binding endpoints to a router
func Handler(handler interface{}) Option {
	return func(b *Builder) {
		if v, ok := handler.(func(w http.ResponseWriter, r *http.Request)); ok {
			handler = http.HandlerFunc(v)
		}
		b.Endpoint.Handler = handler
	}
}

// Summary sets the endpoint's summary
func Summary(v string) Option {
	return func(b *Builder) {
		b.Endpoint.Summary = v
	}
}

// Description sets the endpoint's description
func Description(v string) Option {
	return func(b *Builder) {
		b.Endpoint.Description = v
	}
}

// OperationID sets the endpoint's operationId
func OperationID(v string) Option {
	return func(b *Builder) {
		b.Endpoint.OperationID = v
	}
}

// Produces sets the endpoint's produces; by default this will be set to application/json
func Produces(v ...string) Option {
	return func(b *Builder) {
		b.Endpoint.Produces = v
	}
}

// Consumes sets the endpoint's produces; by default this will be set to application/json
func Consumes(v ...string) Option {
	return func(b *Builder) {
		b.Endpoint.Consumes = v
	}
}

func parameter(p swag.Parameter) Option {
	return func(b *Builder) {
		if b.Endpoint.Parameters == nil {
			b.Endpoint.Parameters = []swag.Parameter{}
		}

		b.Endpoint.Parameters = append(b.Endpoint.Parameters, p)
	}
}

// Path defines a path parameter for the endpoint; name, typ, description, and required correspond to the matching
// swagger fields
func Path(name, typ, description string, required bool) Option {
	return PathDefault(name, typ, description, "", required)
}

// PathDefault defines a path parameter for the endpoint; name, typ, description, and required correspond
// to the matching swagger fields
func PathDefault(name, typ, description, defVal string, required bool) Option {
	p := swag.Parameter{
		Name:        name,
		In:          "path",
		Type:        typ,
		Description: description,
		Required:    required,
		Default:     defVal,
	}
	return parameter(p)
}

// Query defines a query parameter for the endpoint; name, typ, description, and required correspond
// to the matching swagger fields
func Query(name, typ, description string, required bool) Option {
	return QueryDefault(name, typ, description, "", required)
}

// QueryDefault defines a query parameter for the endpoint; name, typ, description, and required correspond
// to the matching swagger fields
func QueryDefault(name, typ, description, defVal string, required bool) Option {
	p := swag.Parameter{
		Name:        name,
		In:          "query",
		Type:        typ,
		Description: description,
		Required:    required,
		Default:     defVal,
	}
	return parameter(p)
}

// Body defines a body parameter for the swagger endpoint as would commonly be used for the POST, PUT, and PATCH methods
// prototype should be a struct or a pointer to struct that swag can use to reflect upon the return type
func Body(prototype interface{}, description string, required bool) Option {
	return BodyType(reflect.TypeOf(prototype), description, required)
}

// BodyType defines a body parameter for the swagger endpoint as would commonly be used for the POST, PUT, and PATCH methods
// prototype should be a struct or a pointer to struct that swag can use to reflect upon the return type
// t represents the Type of the body
func BodyType(t reflect.Type, description string, required bool) Option {
	p := swag.Parameter{
		In:          "body",
		Name:        "body",
		Description: description,
		Schema:      swag.MakeSchema(t),
		Required:    required,
	}
	return parameter(p)
}

// Tags allows one or more tags to be associated with the endpoint
func Tags(tags ...string) Option {
	return func(b *Builder) {
		if b.Endpoint.Tags == nil {
			b.Endpoint.Tags = []string{}
		}

		b.Endpoint.Tags = append(b.Endpoint.Tags, tags...)
	}
}

// Security allows a security scheme to be associated with the endpoint.
func Security(scheme string, scopes ...string) Option {
	return func(b *Builder) {
		if b.Endpoint.Security == nil {
			b.Endpoint.Security = &swag.SecurityRequirement{}
		}

		if b.Endpoint.Security.Requirements == nil {
			b.Endpoint.Security.Requirements = []map[string][]string{}
		}

		b.Endpoint.Security.Requirements = append(b.Endpoint.Security.Requirements, map[string][]string{scheme: scopes})
	}
}

// NoSecurity explicitly sets the endpoint to have no security requirements.
func NoSecurity() Option {
	return func(b *Builder) {
		b.Endpoint.Security = &swag.SecurityRequirement{DisableSecurity: true}
	}
}

// ResponseOption allows for additional configurations on responses like header information
type ResponseOption func(response *swag.Response)

// Apply improves the readability of applied options
func (o ResponseOption) Apply(response *swag.Response) {
	o(response)
}

// Schema adds schema definitions to swagger responses
func Schema(schema interface{}) ResponseOption {
	return func(response *swag.Response) {
		response.Schema = swag.MakeSchema(schema)
	}
}

// Header adds header definitions to swagger responses
func Header(name, typ, format, description string) ResponseOption {
	return func(response *swag.Response) {
		if response.Headers == nil {
			response.Headers = map[string]swag.Header{}
		}
		response.Headers[name] = swag.Header{
			Type:        typ,
			Format:      format,
			Description: description,
		}
	}
}

// ResponseType sets the endpoint response for the specified code; may be used multiple times with different status codes
// t represents the Type of the response
func ResponseType(code int, description string, opts ...ResponseOption) Option {
	return func(b *Builder) {
		if b.Endpoint.Responses == nil {
			b.Endpoint.Responses = make(map[string]swag.Response)
		}
		r := swag.Response{
			Description: description,
		}
		for _, opt := range opts {
			opt.Apply(&r)
		}
		b.Endpoint.Responses[strconv.Itoa(code)] = r
	}
}

// Response sets the endpoint response for the specified code; may be used multiple times with different status codes
func Response(code int, description string, opts ...ResponseOption) Option {
	return ResponseType(code, description, opts...)
}

func ResponseSuccess(opts ...ResponseOption) Option {
	return Response(http.StatusOK, "success", opts...)
}

func Deprecated() Option {
	return func(b *Builder) {
		b.Endpoint.Deprecated = true
	}
}

// New constructs a new swagger endpoint using the fields and functional options provided
func New(method, path string, options ...Option) *swag.Endpoint {
	method = strings.ToUpper(method)
	e := &Builder{
		Endpoint: &swag.Endpoint{
			Method:      method,
			Path:        path,
			OperationID: strings.ToLower(method) + camel(path),
			Produces:    []string{"application/json"},
			Consumes:    []string{"application/json"},
		},
	}

	for _, opt := range options {
		opt.Apply(e)
	}
	return e.Endpoint
}
