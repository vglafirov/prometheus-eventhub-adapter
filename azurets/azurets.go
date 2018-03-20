package azurets

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/go-kit/kit/log"

	"github.com/openenergi/go-event-hub/eventhub"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"
)

// Client Azure EventHub client definition
type Client struct {
	Sender         eventhub.Sender
	Logger         log.Logger
	ignoredSamples prometheus.Counter
}

// NewClient Create new storage client
func NewClient(
	name string,
	namespace string,
	sasPolicyName string,
	sasPolicyKey string,
	tokenExpiryInterval time.Duration,
	logLevel string,
	logger log.Logger) (cli *Client, err error) {
	var debug bool
	if strings.ToLower(logLevel) == "debug" {
		debug = true
	} else {
		debug = false
	}
	if logger == nil {
		logger = log.NewNopLogger()
	}
	sender, err := eventhub.NewSender(eventhub.SenderOpts{
		EventHubNamespace:   namespace,
		EventHubName:        name,
		SasPolicyName:       sasPolicyName,
		SasPolicyKey:        sasPolicyKey,
		TokenExpiryInterval: tokenExpiryInterval,
		Debug:               debug,
	})

	if err != nil {
		return nil, err
	}

	return &Client{
		Sender: sender,
		Logger: logger,
		ignoredSamples: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "prometheus_azurets_ignored_samples_total",
				Help: "The total number of samples not sent to AzureTS.",
			},
		),
	}, nil
}

// Write sends a batch of samples to Azure EventHub.
func (c *Client) Write(samples model.Samples) error {
	for _, s := range samples {
		m, _ := json.Marshal(s)
		_, err := c.Sender.Send(string(m))
		if err != nil {
			return err
		}
	}
	return nil
}

// Describe implements prometheus.Collector.
func (c *Client) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.ignoredSamples.Desc()
}

// Collect implements prometheus.Collector.
func (c *Client) Collect(ch chan<- prometheus.Metric) {
	ch <- c.ignoredSamples
}

// Name identifies the client as a AzureTS client.
func (c *Client) Name() string {
	return "azurets"
}
