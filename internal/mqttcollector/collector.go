package mqttcollector

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"oiltech/internal/dtos"
	"oiltech/internal/usecases"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Collector struct {
	datasetUsecase *usecases.DatasetUsecase
	broker         string
	topicPrefix    string
	logger         *slog.Logger
	client         mqtt.Client
	ready          atomic.Bool
}

func New(
	datasetUsecase *usecases.DatasetUsecase,
	broker string,
	topicPrefix string,
	logger *slog.Logger,
) *Collector {
	return &Collector{
		datasetUsecase: datasetUsecase,
		broker:         broker,
		topicPrefix:    strings.TrimSuffix(topicPrefix, "/"),
		logger:         logger,
	}
}

func (c *Collector) Start() error {
	opts := mqtt.NewClientOptions().
		AddBroker(c.broker).
		SetClientID(fmt.Sprintf("collector-%d", time.Now().UnixNano())).
		SetAutoReconnect(true).
		SetConnectRetry(true).
		SetConnectRetryInterval(2 * time.Second).
		SetConnectionLostHandler(func(_ mqtt.Client, err error) {
			c.ready.Store(false)
			c.logger.Error("mqtt connection lost", "error", err)
		}).
		SetOnConnectHandler(func(client mqtt.Client) {
			c.ready.Store(true)
			c.subscribe(client)
		})

	c.client = mqtt.NewClient(opts)
	token := c.client.Connect()
	token.Wait()
	return token.Error()
}

func (c *Collector) Stop() {
	c.ready.Store(false)
	if c.client != nil && c.client.IsConnected() {
		c.client.Disconnect(250)
	}
}

func (c *Collector) Ready() bool {
	return c.ready.Load()
}

func (c *Collector) subscribe(client mqtt.Client) {
	topic := c.topicPrefix + "/+"
	token := client.Subscribe(topic, 1, c.handleMessage)
	token.Wait()
	if token.Error() != nil {
		c.logger.Error("mqtt subscribe failed", "topic", topic, "error", token.Error())
		return
	}
	c.logger.Info("collector subscribed", "topic", topic)
}

func (c *Collector) handleMessage(_ mqtt.Client, message mqtt.Message) {
	datasetCode, err := c.datasetCodeFromTopic(message.Topic())
	if err != nil {
		c.logger.Error("invalid mqtt topic", "topic", message.Topic(), "error", err)
		return
	}

	var msg dtos.MQTTDatasetMessage
	if err := json.Unmarshal(message.Payload(), &msg); err != nil {
		c.logger.Error("invalid mqtt payload", "topic", message.Topic(), "error", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := c.datasetUsecase.StoreMQTTMessage(ctx, datasetCode, msg); err != nil {
		c.logger.Error("dataset message storing failed", "topic", message.Topic(), "dataset_code", datasetCode, "event_id", msg.EventID, "error", err)
		return
	}
	c.logger.Info("dataset message stored", "topic", message.Topic(), "dataset_code", datasetCode, "event_id", msg.EventID)
}

func (c *Collector) datasetCodeFromTopic(topic string) (string, error) {
	prefix := c.topicPrefix + "/"
	if !strings.HasPrefix(topic, prefix) {
		return "", fmt.Errorf("topic must start with %s", prefix)
	}
	datasetCode := strings.TrimPrefix(topic, prefix)
	if datasetCode == "" || strings.Contains(datasetCode, "/") {
		return "", fmt.Errorf("topic must be %s{dataset_code}", prefix)
	}
	return datasetCode, nil
}

func NewHealthHandler(collector *Collector) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if !collector.Ready() {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte(`{"status":"not_ready"}`))
			return
		}
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}
}
