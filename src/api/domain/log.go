package domain

import (
	"context"
	"time"
	"uuid"
)

type LogEventType string

const (
	LogEventPaymentReceived       LogEventType = "payment_received"
	LogEventReconciliationStarted LogEventType = "reconciliation_started"
	LogEventInvoiceMatched        LogEventType = "invoice_matched"
	LogEventPaymentUnmatched      LogEventType = "payment_unmatched"
	LogEventPaymentAllocated      LogEventType = "payment_allocated"
	LogEventInvoiceStatusChanged  LogEventType = "invoice_status_changed"
	LogEventPaymentReconciled     LogEventType = "payment_reconciled"
	LogEventReconciliationFailed  LogEventType = "reconciliation_failed"
	LogEventPaymentAlreadyHandled LogEventType = "payment_already_processed"
)

type Log struct {
	ID         uuid.UUID
	EventType  LogEventType
	EntityType string
	EntityID   *uuid.UUID
	Message    string
	Metadata   map[string]any
	CreatedAt  time.Time
}

type LogRepository interface {
	Create(ctx context.Context, log Log)
}
