package ts_upstream

import (
	"errors"
	"github.com/caddyserver/caddy/v2"
	"github.com/thanhpk/randstr"
	"go.uber.org/zap"
	"net/http"
	"os"
	"tailscale.com/tsnet"
)

type TsStruct struct {
	inited bool
	server *tsnet.Server
}

func (t *TsStruct) SetUp(ctx caddy.Context) error {
	if t.inited {
		return nil
	}

	hostnameBase, hostnameOk := os.LookupEnv("TS_BASE_HOSTNAME")
	authKey, authKeyOk := os.LookupEnv("TS_AUTHKEY")

	if !hostnameOk || !authKeyOk {
		return errors.New("TS_BASE_HOSTNAME and TS_AUTHKEY are required")
	}

	hostname := hostnameBase + "-" + randstr.Base62(5)

	ctx.Logger().Info("", zap.String("base_hostname", hostnameBase), zap.String("hostname", hostname), zap.String("auth_key", authKey))

	t.server = &tsnet.Server{
		Dir:       "/var/run/tailscale",
		Ephemeral: true,
		Hostname:  hostname,
	}

	err := t.server.Start()

	return err
}

func (t *TsStruct) Close() error {
	return t.server.Close()
}

func (t *TsStruct) Execute(request *http.Request) (*http.Response, error) {
	return t.server.HTTPClient().Transport.RoundTrip(request)
}
