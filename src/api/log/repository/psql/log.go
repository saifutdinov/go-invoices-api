package psql

import (
	"context"
	"encoding/json"
	"log"
	"uuid"

	"github.com/saifutdinov/go-invoices-api/api/domain"
)

func (lr *LogRepository) Create(ctx context.Context, logRecord domain.Log) {
	id := uuid.New()

	metadataJson, err := json.Marshal(logRecord.Metadata)
	if err != nil {
		log.Println(err)
	}

	_, err = lr.ExecContext(ctx, "INSERT INTO logs (id, event_type, entity_type, entity_id, message, metadata) VALUES ($1, $2, $3, $4, $5, $6)", id, logRecord.EventType, logRecord.EntityType, logRecord.EntityID, logRecord.Message, string(metadataJson))
	if err != nil {
		log.Println(err)
	}
}
