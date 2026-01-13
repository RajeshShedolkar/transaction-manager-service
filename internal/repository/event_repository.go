package repository

type EventRepository interface {
	IsProcessed(eventID string) (bool, error)
	MarkProcessed(eventID string, eventType string) error
}
