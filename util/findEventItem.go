package util

import "github.com/disresc/lib/models"

// FindEventItem returns the matching EventItem
func FindEventItem(event *models.Event, transmitter string, metric string) (*models.EventItem, bool) {
	for _, item := range event.GetItems() {
		if item.GetTransmitter() == transmitter && item.GetMetric() == metric {
			return item, true
		}
	}
	return nil, false
}
