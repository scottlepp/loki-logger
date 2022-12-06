package log

import (
	"context"
	"net/http"
	"net/http/httputil"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func HandleLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logFn(w, r)
		next.ServeHTTP(w, r)
	})
}

func HandleLogFunc(next http.HandlerFunc) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		logFn(rw, r)
		next.ServeHTTP(rw, r)
	}
}

func logFn(w http.ResponseWriter, r *http.Request) {
	// Debug(r.URL.Path, "body", r.GetBody)
	dumped, err := httputil.DumpRequest(r, true)
	if err != nil {
		Debug("Could not dump http request", "err", err)
	}
	Debug("plugin request", "req", string(dumped))
}

type Handler struct {
	Plugin plugin
}

func (h Handler) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	Debug("query data request", "req", req)
	return h.Plugin.QueryData(ctx, req)
}

func (h Handler) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	Debug("check health", "req", req)
	return h.Plugin.CheckHealth(ctx, req)
}

type plugin interface {
	QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error)
	CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error)
}
