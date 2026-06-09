# Mail Burrow

Mail Burrow is a study project built with Go to explore asynchronous email processing using HTTP APIs, RabbitMQ, SQLite, and background workers.

The goal of this project is not to be a production-ready email platform, but to practice backend architecture concepts such as message queues, retries, dead-letter queues, dependency injection, storage abstraction, and worker-based processing.

## Project Goals

This project was created to study:

- Building HTTP APIs with Fiber
- Organizing a Go project using internal packages
- Dependency injection with Uber Fx
- Publishing and consuming messages with RabbitMQ
- Implementing retry flows with DLX and TTL queues
- Persisting email processing state with SQLite
- Separating domain, services, ports, infrastructure, and outbound adapters
- Handling errors in a structured way

## How It Works

The main flow is:

```text
POST /emails/publish
  -> saves the email as pending in SQLite
  -> publishes a message to RabbitMQ
  -> worker consumes the message
  -> worker tries to send the email
      -> on success: marks the email as success and ACKs the message
      -> on failure: increments attempts and NACKs without requeue
      -> RabbitMQ sends the message to a retry queue
      -> after TTL, the message returns to the main queue
      -> after max attempts, the email is marked as failed and sent to DLQ
````

## Main Technologies

* Go
* Fiber
* RabbitMQ
* SQLite
* Bun ORM
* Uber Fx
* gomail

## Endpoints

### Publish an email

```http
POST /api/v1/emails/publish
Content-Type: application/json
```

Example body:

```json
{
  "to": "receiver@example.com",
  "from": "sender@example.com",
  "subject": "Hello",
  "body": "This is a test email"
}
```

Example response:

```json
{
  "id": "email-id"
}
```

### Get email status

```http
GET /api/v1/emails?id=email-id
```

Example response:

```json
{
  "id": "email-id",
  "attempts": 0,
  "status": "pending"
}
```

## Environment Variables

Create a `.env` file with the following variables:

```env
AMQP_URL=amqp://guest:guest@localhost:5672/
RABBIT_MQ_PREFETCH=5

SERVER_HOST=0.0.0.0
SERVER_PORT=8080

MAILER_HOST=smtp.example.com
MAILER_PORT=587
MAILER_USERNAME=username
MAILER_PASSWORD=password

DATABASE_URL=database.sqlite3
```

## Running RabbitMQ

Using Docker Compose:

```bash
docker compose -f docker/docker-compose.yml up -d
```

RabbitMQ Management UI:

```text
http://localhost:15672
```

Default credentials:

```text
guest / guest
```

## Running the Application

```bash
go run ./cmd/api
```

## Project Structure

```text
cmd/api
  Application entry point

internal/api
  HTTP handlers and API module

internal/app
  Domain, ports, services, and application logic

internal/config
  Environment, database, and RabbitMQ configuration

internal/infra
  Server, logger, mailer, and infrastructure providers

internal/outbound
  Database adapters, queue publisher, topology, and workers
```

## Study Notes

This project is useful for learning how asynchronous systems behave when something fails.

Some important concepts practiced here:

* A message can be successfully published but fail during processing.
* Failed messages should not be retried immediately forever.
* Retry queues help delay retries.
* Dead-letter queues help isolate messages that could not be processed.
* The database keeps the current state of each email.
* Workers should be idempotent whenever possible.

## Current Status

This is a study project and still evolving.

Automated tests, improved documentation, production hardening, observability, and better local development tooling can be added in future iterations.
