package ts_upstream

import (
	"errors"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/thanhpk/randstr"
	"net/http"
	"os"
	"tailscale.com/tsnet"
)

var hostnameBase, hostnameOk = os.LookupEnv("TS_BASE_HOSTNAME")
var authKey, authKeyOk = os.LookupEnv("TS_AUTH_KEY")
var hostname = ""

func init() {
	if hostnameOk {
		hostname = hostnameBase + "-" + randstr.Base62(5)
	}
	caddy.RegisterModule(TsUpstreamModule{})
}

type TsUpstreamModule struct {
	server *tsnet.Server
	client *http.Client
}

func (m TsUpstreamModule) Cleanup() error {
	return m.server.Close()
}

func (m TsUpstreamModule) RoundTrip(request *http.Request) (*http.Response, error) {
	return m.client.Transport.RoundTrip(request)
}

func (m TsUpstreamModule) UnmarshalCaddyfile(_ *caddyfile.Dispenser) error {
	// no-op
	return nil
}

func (m TsUpstreamModule) CaddyModule() caddy.ModuleInfo {

	return caddy.ModuleInfo{
		ID: "http.reverse_proxy.transport.tailscale_http",
		New: func() caddy.Module {
			return m
		},
	}
}

func (m *TsUpstreamModule) Provision(_ caddy.Context) error {
	if !hostnameOk || !authKeyOk {
		return errors.New("TS_BASE_HOSTNAME and TS_AUTH_KEY are required")
	}

	m.server = &tsnet.Server{
		AuthKey:   authKey,
		Dir:       "/var/run/tailscale",
		Ephemeral: true,
		Hostname:  hostname,
	}

	m.client = m.server.HTTPClient()

	return nil
}
