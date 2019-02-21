package userkyced

import (
	"fmt"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

const (
	userKycedURL = "user-kyc-url"
)

// NewCliFlags returns cli flags to configure a user kyc status client.
func NewCliFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   userKycedURL,
			Usage:  "user kyced API URL. If this doesn't support, fallback to default Postgres DB for kyced checking",
			EnvVar: "USER_KYCED_URL",
			Value:  "",
		},
	}
}

// NewClientFromContext returns new user kyc client from cli flags.
func NewClientFromContext(sugar *zap.SugaredLogger, c *cli.Context) (*Client, error) {
	userURL := c.String(userKycedURL)
	if userURL == "" {
		return nil, httputil.ErrNoClientURL
	}
	err := validation.Validate(userURL,
		is.URL,
	)
	if err != nil {
		return nil, fmt.Errorf("user kyced url: %s", err.Error())
	}

	return NewClient(sugar, userURL)
}
