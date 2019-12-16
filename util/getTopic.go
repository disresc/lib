package util

import (
	"fmt"
	"strings"

	"github.com/disresc/lib/models"
)

// GetTopicFromRequest returns the msg bus topic from a Request
func GetTopicFromRequest(request *models.Request) string {
	topic := fmt.Sprintf("events.%s.%s.%ds", request.Source, request.Transmitter, request.Interval)
	return topic
}

// GetTopicsFromEvent returns the msg bus topic from an Event
func GetTopicsFromEvent(event *models.Event, interval int) []string {
	topicDirect := fmt.Sprintf("events.%s.%s.%ds", event.GetSource(), event.GetItems()[0].GetTransmitter(), interval)
	topics := []string{topicDirect}
	sourceParts := strings.Split(event.GetSource(), "-")
	if len(sourceParts) > 1 && sourceParts[0] == "ve" {
		topicSourceAny := fmt.Sprintf("events.%s.%s.%ds", "ves", event.GetItems()[0].GetTransmitter(), interval)
		topics = append(topics, topicSourceAny)
	} else if len(sourceParts) > 1 && sourceParts[0] == "host" {
		topicSourceAny := fmt.Sprintf("events.%s.%s.%ds", "hosts", event.GetItems()[0].GetTransmitter(), interval)
		topics = append(topics, topicSourceAny)
	}
	return topics
}
