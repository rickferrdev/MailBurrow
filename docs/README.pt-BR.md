# Mail Burrow

Mail Burrow é um projeto de estudo desenvolvido em Go para explorar processamento assíncrono de e-mails usando APIs HTTP, RabbitMQ, SQLite e workers em background.

O objetivo deste projeto não é ser uma plataforma de envio de e-mails pronta para produção, mas sim praticar conceitos de arquitetura backend, como filas, retentativas, dead-letter queues, injeção de dependência, abstração de storage e processamento assíncrono.

## Objetivos do Projeto

Este projeto foi criado para estudar:

- Criação de APIs HTTP com Fiber
- Organização de projetos Go usando pacotes internos
- Injeção de dependência com Uber Fx
- Publicação e consumo de mensagens com RabbitMQ
- Fluxos de retry usando DLX e filas com TTL
- Persistência de estado com SQLite
- Separação entre domínio, serviços, portas, infraestrutura e adapters externos
- Tratamento estruturado de erros

## Como Funciona

O fluxo principal é:

```text
POST /emails/publish
  -> salva o e-mail como pending no SQLite
  -> publica uma mensagem no RabbitMQ
  -> worker consome a mensagem
  -> worker tenta enviar o e-mail
      -> em caso de sucesso: marca como success e confirma a mensagem com ACK
      -> em caso de erro: incrementa attempts e faz NACK sem requeue
      -> RabbitMQ envia a mensagem para uma fila de retry
      -> após o TTL, a mensagem volta para a fila principal
      -> ao exceder o limite de tentativas, o e-mail é marcado como failed e enviado para a DLQ
````

## Principais Tecnologias

* Go
* Fiber
* RabbitMQ
* SQLite
* Bun ORM
* Uber Fx
* gomail

## Endpoints

### Publicar um e-mail

```http
POST /api/v1/emails/publish
Content-Type: application/json
```

Exemplo de body:

```json
{
  "to": "receiver@example.com",
  "from": "sender@example.com",
  "subject": "Olá",
  "body": "Este é um e-mail de teste"
}
```

Exemplo de resposta:

```json
{
  "id": "email-id"
}
```

### Consultar status do e-mail

```http
GET /api/v1/emails?id=email-id
```

Exemplo de resposta:

```json
{
  "id": "email-id",
  "attempts": 0,
  "status": "pending"
}
```

## Variáveis de Ambiente

Crie um arquivo `.env` com as seguintes variáveis:

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

## Subindo o RabbitMQ

Usando Docker Compose:

```bash
docker compose -f docker/docker-compose.yml up -d
```

Interface de gerenciamento do RabbitMQ:

```text
http://localhost:15672
```

Credenciais padrão:

```text
guest / guest
```

## Rodando a Aplicação

```bash
go run ./cmd/api
```

## Estrutura do Projeto

```text
cmd/api
  Ponto de entrada da aplicação

internal/api
  Handlers HTTP e módulo da API

internal/app
  Domínio, portas, serviços e regras de aplicação

internal/config
  Configuração de ambiente, banco de dados e RabbitMQ

internal/infra
  Servidor, logger, mailer e providers de infraestrutura

internal/outbound
  Adapters de banco de dados, publisher de fila, topologia e workers
```

## Notas de Estudo

Este projeto é útil para entender como sistemas assíncronos se comportam quando algo falha.

Alguns conceitos importantes praticados aqui:

* Uma mensagem pode ser publicada com sucesso, mas falhar durante o processamento.
* Mensagens com erro não devem ser reprocessadas imediatamente para sempre.
* Filas de retry ajudam a atrasar novas tentativas.
* Dead-letter queues ajudam a isolar mensagens que não puderam ser processadas.
* O banco de dados mantém o estado atual de cada e-mail.
* Workers devem ser idempotentes sempre que possível.

## Status Atual

Este é um projeto de estudo e ainda está em evolução.

Testes automatizados, documentação mais completa, melhorias para produção, observabilidade e ferramentas de desenvolvimento local podem ser adicionados nas próximas versões.
