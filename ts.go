package ts_upstream

import (
	"errors"
	"github.com/thanhpk/randstr"
	"net/http"
	"os"
	"tailscale.com/tsnet"
)

type TsStruct struct {
	inited bool
	server *tsnet.Server
}

func (t *TsStruct) SetUp() error {
	if t.inited {
		return nil
	}

	hostnameBase, hostnameOk := os.LookupEnv("TS_BASE_HOSTNAME")
	authKey, authKeyOk := os.LookupEnv("TS_AUTH_KEY")

	if !hostnameOk || !authKeyOk {
		return errors.New("TS_BASE_HOSTNAME and TS_AUTH_KEY are required")
	}

	hostname := hostnameBase + "-" + randstr.Base62(5)

	t.server = &tsnet.Server{
		AuthKey:   authKey,
		Dir:       "/var/run/tailscale",
		Ephemeral: true,
		Hostname:  hostname,
	}

	return nil
}

func (t *TsStruct) Close() error {
	return t.server.Close()
}

func (t *TsStruct) Execute(request *http.Request) (*http.Response, error) {
	return t.server.HTTPClient().Transport.RoundTrip(request)
}