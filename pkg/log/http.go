package log

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/experimental/e2e/utils"
)

// HTTPLogger is a http.RoundTripper that logs requests and responses.
type HTTPLogger struct {
	pluginID string
	enabled  func() bool
	proxied  http.RoundTripper
}

type Options struct {
	EnabledFn func() bool
}

// NewHTTPLogger creates a new HTTPLogger.
func NewHTTPLogger(pluginID string, proxied http.RoundTripper, opts ...Options) *HTTPLogger {
	if len(opts) > 1 {
		panic("too many Options arguments provided")
	}

	return &HTTPLogger{
		pluginID: pluginID,
		proxied:  proxied,
		enabled:  enabled,
	}
}

// RoundTrip implements the http.RoundTripper interface.
func (hl *HTTPLogger) RoundTrip(req *http.Request) (*http.Response, error) {
	if !hl.enabled() {
		return hl.proxied.RoundTrip(req)
	}

	buf := []byte{}
	if req.Body != nil {
		if b, err := utils.ReadRequestBody(req); err == nil {
			req.Body = ioutil.NopCloser(bytes.NewReader(b))
			buf = b
		}
	}

	res, err := hl.proxied.RoundTrip(req)
	if err != nil {
		return res, err
	}

	// reset the request body before saving
	if req.Body != nil {
		req.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
	}

	dumped, err2 := httputil.DumpRequest(req, true)
	if err2 != nil {
		Debug("Could not dump http request", "err", err2)
	}
	Debug("http request complete", "req", string(dumped), "err", err)

	if res != nil {
		dumped, err2 = httputil.DumpResponse(res, true)
		if err2 != nil {
			Debug("Could not dump http response", "err", err2)
		}
	}

	Debug("http request complete", "resp", string(dumped), "err", err)

	return res, err
}

func enabled() bool {
	if v, ok := os.LookupEnv("GF_PLUGIN_LOGGER_HTTP"); ok && v == "true" {
		return true
	}
	return false
}
