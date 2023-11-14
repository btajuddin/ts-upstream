package ts_upstream

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"net/http"
)

var server = new(TsStruct)

func init() {
	caddy.RegisterModule(&TsUpstreamModule{})
}

type TsUpstreamModule struct {
}

func (m TsUpstreamModule) Cleanup() error {
	return server.Close()
}

func (m TsUpstreamModule) RoundTrip(request *http.Request) (*http.Response, error) {
	return server.Execute(request)
}

func (m TsUpstreamModule) UnmarshalCaddyfile(_ *caddyfile.Dispenser) error {
	// no-op
	return nil
}

func (m *TsUpstreamModule) CaddyModule() caddy.ModuleInfo {

	return caddy.ModuleInfo{
		ID: "http.reverse_proxy.transport.tailscale_http",
		New: func() caddy.Module {
			return m
		},
	}
}

func (m *TsUpstreamModule) Provision(_ caddy.Context) error {
	return server.SetUp()
}
