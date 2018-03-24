package azurets

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/KirillSleta/go-eventhub/eventhub"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"
)

// Client Azure EventHub client definition
type Client struct {
	Sender         eventhub.EventHubClient
	HubName        string
	Logger         log.Logger
	ignoredSamples prometheus.Counter
}

// NewClient Create new storage client
func NewClient(
	name string,
	namespace string,
	sasPolicyName string,
	sasPolicyKey string,
	logLevel string,
	logger log.Logger) (cli *Client, err error) {

	sender := eventhub.NewEventHubClient(1, namespace, sasPolicyName, sasPolicyKey)

	return &Client{
		Sender:  sender,
		HubName: name,
		Logger:  logger,
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
		t := model.Time.Time(s.Timestamp)
		message := make(map[string]string)
		message["Timestamp"] = t.Format(time.RFC3339)
		for key, value := range s.Metric {
			message[string(key)] = string(value)
		}
		message["Value"] = strconv.FormatFloat(float64(s.Value), 'f', -1, 64)
		m, err := json.Marshal(message)
		if err != nil {
			level.Error(c.Logger).Log("msg", "Cannot marshal incoming message", "err", err.Error())
			return err
		}

		level.Debug(c.Logger).Log("msg", "Message", "payload", string(m))
		err = c.Sender.Send(c.HubName, &eventhub.Message{Body: m})
		if err != nil {
			level.Error(c.Logger).Log("msg", "Cannot send metrics", "err", err.Error())
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
