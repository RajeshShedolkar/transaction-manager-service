package events

import (
	"encoding/json"
	"log"
	"transaction-manager/internal/repository"
)

func HandleCardEvent(
	message []byte,
	eventRepo repository.EventRepository,
) {

	var event CardEvent
	err := json.Unmarshal(message, &event)
	if err != nil {
		log.Println("Invalid event json:", err)
		return
	}

	processed, _ := eventRepo.IsProcessed(event.EventID)
	if processed {
		log.Println("Duplicate event ignored:", event.EventID)
		return
	}

	log.Println("Processing event:", event.EventID, event.EventType)

	// STEP 12 will add business logic here
	log.Println("Business handling skipped for STEP 11")

	eventRepo.MarkProcessed(event.EventID, event.EventType)
}
