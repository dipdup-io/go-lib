package signalr

import (
	"net/http"
	"net/url"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/pkg/errors"
)

// Version - protocol version
type Version string

// Versions
const (
	Version0 Version = "0"
	Version1 Version = "1"
)

// Transport -
type Transport struct {
	url    string
	client *http.Client

	log zerolog.Logger
}

// NewTransport -
func NewTransport(baseURL string) *Transport {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	return &Transport{
		url: baseURL,
		client: &http.Client{
			Timeout:   time.Minute,
			Transport: t,
		},
		log: log.Logger,
	}
}

// Negotiate - is used to establish a connection between the client and the server.
func (t *Transport) Negotiate(version Version) (response NegotiateResponse, err error) {
	u, err := url.Parse(t.url)
	if err != nil {
		return
	}
	u.Path += "/negotiate"
	q := u.Query()
	q.Set("negotiateVersion", string(version))
	u.RawQuery = q.Encode()

	t.log.Debug().Str("url", u.String()).Msg("send negotiate request...")

	req, err := http.NewRequest(http.MethodPost, u.String(), nil)
	if err != nil {
		return
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		err = json.NewDecoder(resp.Body).Decode(&response)
		return
	case http.StatusInternalServerError:
		var e Error
		if err = json.NewDecoder(resp.Body).Decode(&e); err != nil {
			return
		}
		return response, errors.Wrap(ErrNegotiate, e.Error)
	default:
		return response, errors.Wrap(ErrInvalidStatusCode, resp.Status)
	}
}
