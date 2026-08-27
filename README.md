# Payment Reconciliation Service

Backend service for automatic reconciliation of incoming payments with issued invoices.

The service receives invoices and payments, matches payments with invoices, calculates invoice status and keeps an audit log of important actions and discrepancies.

## Features

- Invoice management via API
- Payment management via API
- Automatic payment reconciliation
- Multiple payments per invoice
- Payment allocation
- Idempotent payment processing
- Invoice statuses:
  - `Unmatched`
  - `Partially Paid`
  - `Paid`
  - `Overpaid`
- Payment status management
- Manual payment status changes
- Audit logging
- Background reconciliation processing
- PostgreSQL transactions and row-level locking
- Simple reconciliation report
- Docker-based local environment
- Database migrations

## Technology Stack

- Go
- Echo
- PostgreSQL
- Docker / Docker Compose
- SQL
- Clean Architecture

No external message broker is required. Background processing is implemented using the local `chronos` package.

---

# Architecture

The project follows the principles of Clean Architecture.

The main idea is to separate business logic from HTTP delivery, database access and infrastructure.

The dependency direction is:

```text
Delivery
    ↓
Usecase
    ↓
Domain
    ↑
Repository
```

The domain layer does not depend on HTTP, PostgreSQL or infrastructure implementations.

## Project structure

```text
src/
├── api/
│   ├── domain/
│   │
│   └── {entity}/
│       ├── delivery/
│       │   └── http/
│       │
│       ├── usecase/
│       │
│       └── repository/
│           └── psql/
│
├── assets/
│   └── templates/
│
└── pkg/
    └── ...

config/
├── docker/
│   └── dev/
│       ├── docker-compose.yml
│       └── Dockerfile
│
└── ...

Makefile
go.mod
go.sum
.gitignore
```

### `src/api/domain`

Contains domain entities and domain-level types.

These structures are shared between application layers and describe the business model without depending on infrastructure.

### `api/{entity}/delivery/http`

Contains HTTP handlers for the corresponding entity.

The delivery layer is responsible for:

- parsing HTTP requests
- validating request data
- calling usecases
- returning HTTP responses

### `api/{entity}/usecase`

Contains the main application and business logic.

Examples:

- creating invoices
- receiving payments
- reconciliation
- calculating invoice status
- manually changing payment status
- generating reports

### `api/{entity}/repository/psql`

Contains PostgreSQL-specific repository implementations.

This layer is responsible for database queries and persistence.

### `src/assets/templates`

Contains HTML templates used by the simple reconciliation report interface.

### `src/pkg`

Contains local reusable packages used by the application.

One of the custom packages is `chronos`, which is responsible for scheduling background tasks.

### `config`

Contains application configuration.

Docker-specific configuration is located under:

```text
config/docker
```

---

# Reconciliation

The main business process is automatic reconciliation of payments with invoices.

An invoice can have one or multiple payments.

For example:

```text
Invoice
Amount: €1500

Payment #1: €500
Payment #2: €1000

Total allocated: €1500

Invoice status: Paid
```

Payments are connected to invoices through the `payment_allocation` table.

The allocation allows the system to correctly calculate the total amount paid for an invoice.

## Invoice statuses

The invoice status is calculated from the invoice amount and the total amount allocated to it.

```text
Paid amount = 0
        ↓
Unmatched
```

```text
0 < Paid amount < Invoice amount
        ↓
Partially Paid
```

```text
Paid amount = Invoice amount
        ↓
Paid
```

```text
Paid amount > Invoice amount
        ↓
Overpaid
```

The status is therefore based on the sum of all payments allocated to the invoice.

---

# Reconciliation flow

The background processor periodically looks for pending payments.

The simplified flow is:

```text
Background task
      │
      ▼
Find pending payments
      │
      ▼
Process payment one by one
      │
      ▼
BEGIN TRANSACTION
      │
      ▼
SELECT payment FOR UPDATE
      │
      ▼
Check payment status
      │
      ├── already processed → skip
      │
      ▼
Find invoice
      │
      ▼
SELECT invoice FOR UPDATE
      │
      ▼
Create payment allocation
      │
      ▼
Calculate total allocated amount
      │
      ▼
Calculate invoice status
      │
      ▼
Update invoice status
      │
      ▼
Update payment status
      │
      ▼
Write audit logs
      │
      ▼
COMMIT
```

If any important operation fails, the transaction is rolled back.

---

# Transactions and locking

Reconciliation uses PostgreSQL transactions together with row-level locking.

The payment is locked before processing:

```sql
SELECT ...
FROM payments
WHERE id = $1
FOR UPDATE;
```

The corresponding invoice is also locked:

```sql
SELECT ...
FROM invoices
WHERE reference = $1
FOR UPDATE;
```

This prevents two workers from processing the same payment or modifying the same invoice concurrently.

The transaction guarantees that payment allocation, invoice status update and payment status update are performed atomically.

The following operations belong to the same transaction:

```text
Create payment allocation
        +
Update invoice status
        +
Update payment status
        +
Write audit logs
        ↓
      COMMIT
```

If one of these operations fails, the transaction is rolled back.

---

# Idempotency

Payment processing is idempotent.

Before reconciliation the payment is locked and its current status is checked.

Only payments with:

```text
status = pending
```

are processed.

If a payment has already been processed, the reconciliation operation is skipped.

This prevents the same payment from being processed multiple times.

The database transaction and unique constraints on payment allocation provide an additional layer of protection against duplicate processing.

---

# Background Processing

Background reconciliation is implemented using a custom local package called `chronos`.

`chronos` provides a small in-memory scheduler based on a priority queue.

Tasks contain:

- execution time
- function to execute

The reconciliation usecase schedules the next execution after processing the current batch.

Simplified flow:

```text
Application starts
      │
      ▼
StartBackgroundProcessing()
      │
      ▼
Schedule reconciliation task
      │
      ▼
Chronos waits until execution time
      │
      ▼
processPayments()
      │
      ▼
Process pending payments
      │
      ▼
Schedule next reconciliation
```

The scheduler runs the task in a goroutine so that the scheduler itself is not blocked by the reconciliation process.

For the current implementation, reconciliation is scheduled periodically with a short interval suitable for the test environment.

---

# Audit Logging

Important business actions and discrepancies are stored in the log repository.

Examples include:

- payment received
- reconciliation started
- invoice matched
- payment allocated
- invoice status changed
- payment reconciled
- payment already processed
- unmatched payment
- overpayment
- reconciliation failure
- manual payment status change

Audit logging allows the reconciliation process to be inspected after execution and provides visibility into discrepancies.

Failure to write an audit log does not break the main business operation.

---

# Manual Payment Status Change

The application provides an API endpoint for manually changing a payment status.

Example:

```http
PATCH /api/payments/{id}/status
```

Request:

```json
{
  "status": "completed"
}
```

The operation is handled by the payment usecase rather than directly modifying the database from the HTTP layer.

Manual changes are also recorded in the audit log.

---

# Reports

The application provides a simple reconciliation report.

The report contains two main sections.

## Invoices

The invoice report displays:

- invoice ID
- reference
- invoice amount
- total paid amount
- remaining amount
- invoice status
- due date
- number of payments

This makes unpaid and partially paid invoices easy to identify.

## Payments

The payment report displays:

- payment ID
- payment amount
- reference
- payment date
- payment status
- reconciliation discrepancy
- manual status change action

A discrepancy is presented as additional report information rather than as a payment status.

Examples:

```text
Invoice not found
```

```text
Invoice is partially paid
```

```text
Invoice is overpaid
```

The report is available through:

```http
GET /api/reconciliation/report
```

The HTML interface consumes this endpoint and displays the data in a simple browser interface.

---

# Database

PostgreSQL is used as the primary database.

The main entities are:

```text
invoices
payments
payment_allocation
logs
```

`payment_allocation` represents the relationship between payments and invoices and allows multiple payments to be associated with a single invoice.

Amounts are stored as integer values in the smallest currency unit.

For EUR:

```text
€1423.43 → 142343
€520.30  → 52030
€100.00  → 10000
```

This avoids floating-point precision problems when performing financial calculations.

---

# Database Migrations

Database migrations are included in the project and are executed using the project's migration mechanism.

The application can therefore initialize the required database schema in a reproducible way instead of relying on manually created tables.

---

# Running the project

The project uses Docker Compose for local execution.

## Start development environment

```bash
make dev
```

## Build and start development environment

```bash
make build-dev
```

## Connect to PostgreSQL

```bash
make run_db
```

This executes:

```bash
docker exec -it payment-system-db \
    psql -U postgres \
    -d payment-system
```

---

# Makefile

Available commands:

```text
make dev
```

Starts the development environment.

```text
make build-dev
```

Rebuilds Docker images and starts the development environment.

```text
make run_db
```

Opens a PostgreSQL shell inside the database container.

---

# Local URLs

After starting the application, the reconciliation interface is available at:

```text
/reconciliation
```

The main report API:

```text
GET /api/reconciliation/report
```

Payment status update:

```text
PATCH /api/payments/{id}/status
```

---

# Design Decisions

### Why PostgreSQL?

PostgreSQL provides:

- ACID transactions
- row-level locking
- reliable concurrent data processing
- constraints for data integrity
- aggregation capabilities required for reconciliation

These features are particularly useful for payment processing where consistency is more important than eventual correctness.

### Why a custom scheduler?

The project includes a small `chronos` package instead of introducing an external queue or scheduler.
For the current standalone service this keeps the infrastructure simple while still providing background processing.
If the service becomes part of a larger production system, the scheduler can later be replaced with a distributed job queue without changing the reconciliation business logic.

### Why Clean Architecture?

The reconciliation logic is isolated from HTTP and PostgreSQL.
This makes the core business logic easier to test and allows infrastructure components to be replaced without changing the domain logic.
