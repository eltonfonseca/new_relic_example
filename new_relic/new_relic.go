package newrelic

import (
	"fmt"
	"time"

	"github.com/newrelic/go-agent/v3/integrations/nrlogrus"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

type NewRelic struct {
	app *newrelic.Application
}

func New(appName, license string) *NewRelic {
	newRelic, err := newrelic.NewApplication(
		newrelic.ConfigAppName(appName),
		newrelic.ConfigLicense(license),
		newrelic.ConfigDistributedTracerEnabled(true),
		func(config *newrelic.Config) {
			config.Enabled = true
			logrus.SetLevel(logrus.DebugLevel)
			config.Logger = nrlogrus.StandardLogger()
		},
	)

	if err != nil {
		fmt.Println("failed to create New Relic application")
		return nil
	}

	if err = newRelic.WaitForConnection(3 * time.Second); err != nil {
		fmt.Println("failed to create New Relic application")
		return nil
	}

	return &NewRelic{
		app: newRelic,
	}
}

func (n *NewRelic) StartTransaction(name string) *newrelic.Transaction {
	return n.app.StartTransaction(name)
}

func (n *NewRelic) Application() *newrelic.Application {
	return n.app
}
