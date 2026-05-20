package simulator

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"oiltech/internal/dtos"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Simulator struct {
	broker      string
	topicPrefix string
	logger      *slog.Logger
	client      mqtt.Client
}

func New(broker, topicPrefix string, logger *slog.Logger) *Simulator {
	return &Simulator{
		broker:      broker,
		topicPrefix: topicPrefix,
		logger:      logger,
	}
}

func (s *Simulator) Start() error {
	opts := mqtt.NewClientOptions().
		AddBroker(s.broker).
		SetClientID(fmt.Sprintf("mqtt-simulator-%d", time.Now().UnixNano())).
		SetAutoReconnect(true)
	s.client = mqtt.NewClient(opts)
	token := s.client.Connect()
	token.Wait()
	return token.Error()
}

func (s *Simulator) Stop() {
	if s.client != nil && s.client.IsConnected() {
		s.client.Disconnect(250)
	}
}

func (s *Simulator) PublishLoop(stop <-chan struct{}) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			s.publish("meter_readings_raw", meterReadingMessage())
			s.publish("process_measurements_raw", processMeasurementMessage())
			s.publish("equipment_states_raw", equipmentStateMessage())
		}
	}
}

func (s *Simulator) publish(datasetCode string, msg dtos.MQTTDatasetMessage) {
	payload, err := json.Marshal(msg)
	if err != nil {
		s.logger.Error("marshal message failed", "error", err)
		return
	}
	topic := fmt.Sprintf("%s/%s", s.topicPrefix, datasetCode)
	token := s.client.Publish(topic, 1, false, payload)
	token.Wait()
	if token.Error() != nil {
		s.logger.Error("publish failed", "topic", topic, "error", token.Error())
		return
	}
	s.logger.Info("published", "topic", topic, "event_id", msg.EventID)
}

var facilities = []string{"FAC-001", "FAC-002", "FAC-003"}
var equipmentIDs = []string{"PUMP-01", "PUMP-02", "ESP-01", "ESP-02", "COMP-01"}
var equipmentTypes = []string{"pump", "esp", "compressor"}
var resourceTypes = []string{"oil", "gas", "water_produced"}
var measurementKinds = []string{"volume", "mass", "energy"}
var qualityValues = []string{"good", "uncertain", "bad"}
var statusValues = []string{"ok", "manual", "substituted"}
var units = []string{"m3", "t", "GJ"}
var sensorIDs = []string{"SNS-T-01", "SNS-P-01", "SNS-T-02", "SNS-P-02", "SNS-F-01"}
var parameterIDs = []string{"temperature", "pressure", "flow_rate", "rpm", "power"}
var stateValues = []string{"running", "stopped", "maintenance", "fault"}
var modeValues = []string{"auto", "manual", "standby"}

func meterReadingMessage() dtos.MQTTDatasetMessage {
	now := time.Now().UTC()
	fac := facilities[rand.Intn(len(facilities))]
	eq := equipmentIDs[rand.Intn(len(equipmentIDs))]
	eqType := equipmentTypes[rand.Intn(len(equipmentTypes))]
	return dtos.MQTTDatasetMessage{
		EventID: fmt.Sprintf("mr-%d", now.UnixNano()),
		Values: map[string]any{
			"source_system":    "SCADA",
			"source_tag":       fmt.Sprintf("TAG.%s.%s.METER", fac, eq),
			"meter_id":         fmt.Sprintf("MTR-%s-%s", fac, eq),
			"facility_id":      fac,
			"equipment_id":     eq,
			"equipment_type":   eqType,
			"resource_type":    resourceTypes[rand.Intn(len(resourceTypes))],
			"measurement_kind": measurementKinds[rand.Intn(len(measurementKinds))],
			"moment":           now.Format(time.RFC3339),
			"period_start":     now.Add(-time.Hour).Format(time.RFC3339),
			"period_end":       now.Format(time.RFC3339),
			"value":            round(rand.Float64() * 1000),
			"unit":             units[rand.Intn(len(units))],
			"quality":          qualityValues[rand.Intn(len(qualityValues))],
			"status":           statusValues[rand.Intn(len(statusValues))],
		},
	}
}

func processMeasurementMessage() dtos.MQTTDatasetMessage {
	now := time.Now().UTC()
	fac := facilities[rand.Intn(len(facilities))]
	eq := equipmentIDs[rand.Intn(len(equipmentIDs))]
	eqType := equipmentTypes[rand.Intn(len(equipmentTypes))]
	sns := sensorIDs[rand.Intn(len(sensorIDs))]
	param := parameterIDs[rand.Intn(len(parameterIDs))]
	return dtos.MQTTDatasetMessage{
		EventID: fmt.Sprintf("pm-%d", now.UnixNano()),
		Values: map[string]any{
			"source_system":  "DCS",
			"source_tag":     fmt.Sprintf("TAG.%s.%s.%s", fac, eq, param),
			"sensor_id":      sns,
			"facility_id":    fac,
			"well_id":        fmt.Sprintf("WELL-%s-01", fac),
			"equipment_id":   eq,
			"equipment_type": eqType,
			"parameter_id":   param,
			"moment":         now.Format(time.RFC3339),
			"value":          round(rand.Float64() * 500),
			"unit":           paramUnit(param),
			"quality":        qualityValues[rand.Intn(len(qualityValues))],
			"status":         statusValues[rand.Intn(len(statusValues))],
		},
	}
}

func equipmentStateMessage() dtos.MQTTDatasetMessage {
	now := time.Now().UTC()
	fac := facilities[rand.Intn(len(facilities))]
	eq := equipmentIDs[rand.Intn(len(equipmentIDs))]
	eqType := equipmentTypes[rand.Intn(len(equipmentTypes))]
	return dtos.MQTTDatasetMessage{
		EventID: fmt.Sprintf("es-%d", now.UnixNano()),
		Values: map[string]any{
			"source_system":  "SCADA",
			"facility_id":    fac,
			"equipment_id":   eq,
			"equipment_type": eqType,
			"state":          stateValues[rand.Intn(len(stateValues))],
			"mode":           modeValues[rand.Intn(len(modeValues))],
			"scheme_id":      fmt.Sprintf("SCH-%d", 1+rand.Intn(5)),
			"moment":         now.Format(time.RFC3339),
			"quality":        qualityValues[rand.Intn(len(qualityValues))],
		},
	}
}

func paramUnit(param string) string {
	switch param {
	case "temperature":
		return "degC"
	case "pressure":
		return "bar"
	case "flow_rate":
		return "m3/h"
	case "rpm":
		return "rpm"
	case "power":
		return "kW"
	default:
		return "unit"
	}
}

func round(v float64) float64 {
	return float64(int(v*100)) / 100
}
