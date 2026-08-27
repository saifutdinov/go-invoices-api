-- recived invoice from API
CREATE TABLE IF NOT EXISTS invoice (
    id UUID PRIMARY KEY NOT NULL,
    external_id varchar NOT NULL UNIQUE,
    amount BIGINT NOT NULL,
    due_date date NOT NULL,
    createdat timestamp without time zone NOT NULL DEFAULT now(),
    updatedat timestamp without time zone,
    status varchar NOT NULL DEFAULT 'Unmatched'
);

-- recived payment from API
CREATE TABLE IF NOT EXISTS payment (
    id UUID PRIMARY KEY,
    external_id varchar NOT NULL UNIQUE,
    amount BIGINT NOT NULL,
    payment_date date NOT NULL,
    reference varchar(255) NOT NULL,
    source varchar(100) NOT NULL,
    status varchar(100) NOT NULL,
    createdat timestamp without time zone NOT NULL DEFAULT now(),
    updatedat timestamp without time zone
);

-- mapping invoice and payment
CREATE TABLE IF NOT EXISTS payment_allocation (
    payment_id varchar NOT NULL,
    invoice_id varchar NOT NULL,
    amount BIGINT NOT NULL,
    createdat timestamp without time zone NOT NULL DEFAULT now(),
    PRIMARY KEY (payment_id, invoice_id)
);

-- log
CREATE TABLE logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type varchar(100) NOT NULL,
    entity_type varchar(50) NOT NULL,
    entity_id UUID,
    message TEXT,
    metadata JSONB,
    created_at timestamp without time zone NOT NULL DEFAULT now()
);

CREATE INDEX idx_logs_entity
    ON logs (entity_type, entity_id);

CREATE INDEX idx_logs_event_type
    ON logs (event_type);

CREATE INDEX idx_logs_created_at
    ON logs (created_at);