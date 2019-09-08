package vault

import (
	"context"
	"errors"
	"net/http/pprof"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *SystemBackend) pprofPaths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "pprof/?$",

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback:    b.handlePprofIndex,
					Summary:     "Returns the running program's command line.",
					Description: "Returns the running program's command line, with arguments separated by NUL bytes.",
				},
			},
		},
		{
			Pattern: "pprof/cmdline",

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback:    b.handlePprofCmdline,
					Summary:     "Returns the running program's command line.",
					Description: "Returns the running program's command line, with arguments separated by NUL bytes.",
				},
			},
		},
		{
			Pattern: "pprof/goroutine",

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback:    b.handlePprofGoroutine,
					Summary:     "Returns stack traces of all current goroutines.",
					Description: "Returns stack traces of all current goroutines.",
				},
			},
		},
		{
			Pattern: "pprof/heap",

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback:    b.handlePprofHeap,
					Summary:     "Returns a sampling of memory allocations of live object.",
					Description: "Returns a sampling of memory allocations of live object.",
				},
			},
		},
		{
			Pattern: "pprof/profile",

			Fields: map[string]*framework.FieldSchema{
				"seconds": {
					Type:        framework.TypeInt,
					Description: "If provided, specifies the duration to run the profiling command.",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback:    b.handlePprofProfile,
					Summary:     "Returns a pprof-formatted cpu profile payload.",
					Description: "Returns a pprof-formatted cpu profile payload. Profiling lasts for duration specified in seconds GET parameter, or for 30 seconds if not specified.",
				},
			},
		},
		{
			Pattern: "pprof/symbol",

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback:    b.handlePprofSymbol,
					Summary:     "Returns the program counters listed in the request.",
					Description: "Returns the program counters listed in the request.",
				},
			},
		},

		{
			Pattern: "pprof/trace",

			Fields: map[string]*framework.FieldSchema{
				"seconds": {
					Type:        framework.TypeInt,
					Description: "If provided, specifies the duration to run the tracing command.",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback:    b.handlePprofTrace,
					Summary:     "Returns the execution trace in binary form.",
					Description: "Returns  the execution trace in binary form. Tracing lasts for duration specified in seconds GET parameter, or for 1 second if not specified.",
				},
			},
		},
	}
}

func (b *SystemBackend) handlePprofIndex(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	if err := checkRequestHandlerParams(req); err != nil {
		return nil, err
	}

	pprof.Index(req.ResponseWriter, req.HTTPRequest)
	return nil, nil
}

func (b *SystemBackend) handlePprofCmdline(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	if err := checkRequestHandlerParams(req); err != nil {
		return nil, err
	}

	pprof.Cmdline(req.ResponseWriter, req.HTTPRequest)
	return nil, nil
}

func (b *SystemBackend) handlePprofGoroutine(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	if err := checkRequestHandlerParams(req); err != nil {
		return nil, err
	}

	pprof.Handler("goroutine").ServeHTTP(req.ResponseWriter, req.HTTPRequest)
	return nil, nil
}

func (b *SystemBackend) handlePprofHeap(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	if err := checkRequestHandlerParams(req); err != nil {
		return nil, err
	}

	pprof.Handler("heap").ServeHTTP(req.ResponseWriter, req.HTTPRequest)
	return nil, nil
}

func (b *SystemBackend) handlePprofProfile(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	if err := checkRequestHandlerParams(req); err != nil {
		return nil, err
	}

	pprof.Profile(req.ResponseWriter, req.HTTPRequest)
	return nil, nil
}

func (b *SystemBackend) handlePprofSymbol(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	if err := checkRequestHandlerParams(req); err != nil {
		return nil, err
	}

	pprof.Symbol(req.ResponseWriter, req.HTTPRequest)
	return nil, nil
}

func (b *SystemBackend) handlePprofTrace(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	if err := checkRequestHandlerParams(req); err != nil {
		return nil, err
	}

	pprof.Trace(req.ResponseWriter, req.HTTPRequest)
	return nil, nil
}

// checkRequestHandlerParams is a helper that checks for the existence of the
// HTTP request and response writer in a logical.Request.
func checkRequestHandlerParams(req *logical.Request) error {
	if req.ResponseWriter == nil {
		return errors.New("no writer for request")
	}

	if req.HTTPRequest == nil || req.HTTPRequest.Body == nil {
		return errors.New("no reader for request")
	}

	return nil
}
