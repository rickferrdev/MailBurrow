package ports

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Message string
type Code string

type Error struct {
	Trace   error   `json:"-"`
	Message Message `json:"message"`
	Code    Code    `json:"code"`
	Status  int     `json:"status"`
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}

	if e.Trace != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Trace)
	}

	return string(e.Message)
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Trace
}

func New(trace error, message Message, code Code, status int) *Error {
	return &Error{
		Trace:   trace,
		Message: message,
		Code:    code,
		Status:  status,
	}
}

func AsError(err error) *Error {
	if err == nil {
		return nil
	}

	var portErr *Error
	if errors.As(err, &portErr) {
		return portErr
	}

	return NewInternal(err)
}

func FromBody(body []byte, expected *Error) (bool, error) {
	var got *Error
	if err := json.Unmarshal(body, &got); err != nil {
		return false, err
	}

	return expected.Code == got.Code && expected.Message == got.Message, nil
}

const (
	CodeBadRequest   Code = "BAD_REQUEST"
	CodeUnauthorized Code = "UNAUTHORIZED"
	CodeForbidden    Code = "FORBIDDEN"
	CodeNotFound     Code = "NOT_FOUND"
	CodeConflict     Code = "CONFLICT"
	CodeInternal     Code = "INTERNAL_ERROR"

	CodeInvalidID   Code = "INVALID_ID"
	CodeInvalidFrom Code = "INVALID_FROM"
	CodeInvalidTo   Code = "INVALID_TO"

	CodeRabbitMQConnection Code = "RABBITMQ_CONNECTION_ERROR"
	CodeRabbitMQPublish    Code = "RABBITMQ_PUBLISH_ERROR"
	CodeRabbitMQConsume    Code = "RABBITMQ_CONSUME_ERROR"
	CodeRabbitMQAck        Code = "RABBITMQ_ACK_ERROR"
	CodeRabbitMQNack       Code = "RABBITMQ_NACK_ERROR"

	CodeWorkerPayloadInvalid  Code = "CODE_WORKER_PAYLOAD_INVALID"
	CodeLimitAttemptsExceeded Code = "CODE_LIMIT_ATTEMPTS_EXCEEDED"

	CodeDatabaseConnection      Code = "DATABASE_CONNECTION_ERROR"
	CodeDatabasePing            Code = "DATABASE_PING_ERROR"
	CodeDatabaseMigration       Code = "DATABASE_MIGRATION_ERROR"
	CodeDatabaseTransaction     Code = "DATABASE_TRANSACTION_ERROR"
	CodeDatabaseQuery           Code = "DATABASE_QUERY_ERROR"
	CodeDatabaseInsert          Code = "DATABASE_INSERT_ERROR"
	CodeDatabaseUpdate          Code = "DATABASE_UPDATE_ERROR"
	CodeDatabaseDelete          Code = "DATABASE_DELETE_ERROR"
	CodeDatabaseRowsNotFound    Code = "DATABASE_ROWS_NOT_FOUND"
	CodeDatabaseDuplicate       Code = "DATABASE_DUPLICATE_ERROR"
	CodeDatabaseConstraint      Code = "DATABASE_CONSTRAINT_ERROR"
	CodeDatabaseForeignKey      Code = "DATABASE_FOREIGN_KEY_ERROR"
	CodeDatabaseLocked          Code = "DATABASE_LOCKED_ERROR"
	CodeDatabaseBusy            Code = "DATABASE_BUSY_ERROR"
	CodeDatabaseReadonly        Code = "DATABASE_READONLY_ERROR"
	CodeDatabaseInvalidDriver   Code = "DATABASE_INVALID_DRIVER_ERROR"
	CodeDatabaseSQLiteDriver    Code = "DATABASE_SQLITE_DRIVER_ERROR"
	CodeDatabaseInvalidArgument Code = "DATABASE_INVALID_ARGUMENT_ERROR"

	CodeHTTPRequest    Code = "HTTP__REQUEST_ERROR"
	CodeHTTPResponse   Code = "HTTP__RESPONSE_ERROR"
	CodeHTTPBodyParser Code = "HTTP__BODY_PARSER_ERROR"
	CodeHTTPParams     Code = "HTTP__PARAMS_ERROR"

	CodeEmailNotFound       Code = "CODE_EMAIL_NOT_FOUND"
	CodeEmailSend           Code = "EMAIL_SEND_ERROR"
	CodeEmailTemplate       Code = "EMAIL_TEMPLATE_ERROR"
	CodeEmailInvalidAddress Code = "EMAIL_INVALID_ADDRESS"
)

const (
	MessageBadRequest   Message = "bad request"
	MessageUnauthorized Message = "unauthorized"
	MessageForbidden    Message = "forbidden"
	MessageNotFound     Message = "not found"
	MessageConflict     Message = "conflict"
	MessageInternal     Message = "internal server error"

	MessageInvalidID   Message = "invalid id"
	MessageInvalidFrom Message = "invalid from"
	MessageInvalidTo   Message = "invalid to"

	MessageRabbitMQConnection Message = "rabbitmq connection error"
	MessageRabbitMQPublish    Message = "rabbitmq publish error"
	MessageRabbitMQConsume    Message = "rabbitmq consume error"
	MessageRabbitMQAck        Message = "rabbitmq ack error"
	MessageRabbitMQNack       Message = "rabbitmq nack error"

	MessageWorkerPayloadInvalid  Message = "payload processed by the worker is invalid"
	MessageLimitAttemptsExceeded Message = "limit of attempts exceeded"

	MessageDatabaseConnection      Message = "database connection error"
	MessageDatabasePing            Message = "database ping error"
	MessageDatabaseMigration       Message = "database migration error"
	MessageDatabaseTransaction     Message = "database transaction error"
	MessageDatabaseQuery           Message = "database query error"
	MessageDatabaseInsert          Message = "database insert error"
	MessageDatabaseUpdate          Message = "database update error"
	MessageDatabaseDelete          Message = "database delete error"
	MessageDatabaseRowsNotFound    Message = "database rows not found"
	MessageDatabaseDuplicate       Message = "database duplicate record"
	MessageDatabaseConstraint      Message = "database constraint violation"
	MessageDatabaseForeignKey      Message = "database foreign key constraint violation"
	MessageDatabaseLocked          Message = "database is locked"
	MessageDatabaseBusy            Message = "database is busy"
	MessageDatabaseReadonly        Message = "database is readonly"
	MessageDatabaseInvalidDriver   Message = "database invalid driver"
	MessageDatabaseSQLiteDriver    Message = "sqlite driver error"
	MessageDatabaseInvalidArgument Message = "database invalid argument"

	MessageHTTPRequest    Message = "http request error"
	MessageHTTPResponse   Message = "http response error"
	MessageHTTPBodyParser Message = "http body parser error"
	MessageHTTPParams     Message = "http params error"

	MessageEmailNotFound       Message = "email does not exist"
	MessageEmailSend           Message = "email send error"
	MessageEmailTemplate       Message = "email template error"
	MessageEmailInvalidAddress Message = "email invalid address"
)

func IsCode(err error, code Code) bool {
	var e *Error

	return errors.As(err, &e) && e.Code == code
}

func NewBadRequest(trace error) *Error {
	return New(trace, MessageBadRequest, CodeBadRequest, http.StatusBadRequest)
}

func NewUnauthorized(trace error) *Error {
	return New(trace, MessageUnauthorized, CodeUnauthorized, http.StatusUnauthorized)
}

func NewForbidden(trace error) *Error {
	return New(trace, MessageForbidden, CodeForbidden, http.StatusForbidden)
}

func NewNotFound(trace error) *Error {
	return New(trace, MessageNotFound, CodeNotFound, http.StatusNotFound)
}

func NewConflict(trace error) *Error {
	return New(trace, MessageConflict, CodeConflict, http.StatusConflict)
}

func NewInternal(trace error) *Error {
	return New(trace, MessageInternal, CodeInternal, http.StatusInternalServerError)
}

func NewInvalidID(trace error) *Error {
	if trace == nil {
		trace = errors.New("id is invalid")
	}

	return New(trace, MessageInvalidID, CodeInvalidID, http.StatusBadRequest)
}

func NewInvalidFrom(trace error) *Error {
	if trace == nil {
		trace = errors.New("from is invalid")
	}

	return New(trace, MessageInvalidFrom, CodeInvalidFrom, http.StatusBadRequest)
}

func NewInvalidTo(trace error) *Error {
	if trace == nil {
		trace = errors.New("to is invalid")
	}

	return New(trace, MessageInvalidTo, CodeInvalidTo, http.StatusBadRequest)
}

func NewRabbitMQConnection(trace error) *Error {
	return New(trace, MessageRabbitMQConnection, CodeRabbitMQConnection, http.StatusServiceUnavailable)
}

func NewRabbitMQPublish(trace error) *Error {
	return New(trace, MessageRabbitMQPublish, CodeRabbitMQPublish, http.StatusInternalServerError)
}

func NewRabbitMQConsume(trace error) *Error {
	return New(trace, MessageRabbitMQConsume, CodeRabbitMQConsume, http.StatusInternalServerError)
}

func NewRabbitMQAck(trace error) *Error {
	return New(trace, MessageRabbitMQAck, CodeRabbitMQAck, http.StatusInternalServerError)
}

func NewWorkerPayloadInvalid(trace error) *Error {
	return New(trace, MessageWorkerPayloadInvalid, CodeWorkerPayloadInvalid, http.StatusInternalServerError)
}

func NewLimitAttemptsExceeded(trace error) *Error {
	return New(trace, MessageLimitAttemptsExceeded, CodeLimitAttemptsExceeded, http.StatusUnprocessableEntity)
}

func NewRabbitMQNack(trace error) *Error {
	return New(trace, MessageRabbitMQNack, CodeRabbitMQNack, http.StatusInternalServerError)
}

func NewDatabaseConnection(trace error) *Error {
	return New(trace, MessageDatabaseConnection, CodeDatabaseConnection, http.StatusServiceUnavailable)
}

func NewDatabasePing(trace error) *Error {
	return New(trace, MessageDatabasePing, CodeDatabasePing, http.StatusServiceUnavailable)
}

func NewDatabaseMigration(trace error) *Error {
	return New(trace, MessageDatabaseMigration, CodeDatabaseMigration, http.StatusInternalServerError)
}

func NewDatabaseTransaction(trace error) *Error {
	return New(trace, MessageDatabaseTransaction, CodeDatabaseTransaction, http.StatusInternalServerError)
}

func NewDatabaseRowsNotFound(trace error) *Error {
	if trace == nil {
		trace = sql.ErrNoRows
	}

	return New(trace, MessageDatabaseRowsNotFound, CodeDatabaseRowsNotFound, http.StatusNotFound)
}

func NewDatabaseDuplicate(trace error) *Error {
	return New(trace, MessageDatabaseDuplicate, CodeDatabaseDuplicate, http.StatusConflict)
}

func NewDatabaseConstraint(trace error) *Error {
	return New(trace, MessageDatabaseConstraint, CodeDatabaseConstraint, http.StatusConflict)
}

func NewDatabaseForeignKey(trace error) *Error {
	return New(trace, MessageDatabaseForeignKey, CodeDatabaseForeignKey, http.StatusConflict)
}

func NewDatabaseLocked(trace error) *Error {
	return New(trace, MessageDatabaseLocked, CodeDatabaseLocked, http.StatusServiceUnavailable)
}

func NewDatabaseBusy(trace error) *Error {
	return New(trace, MessageDatabaseBusy, CodeDatabaseBusy, http.StatusServiceUnavailable)
}

func NewDatabaseReadonly(trace error) *Error {
	return New(trace, MessageDatabaseReadonly, CodeDatabaseReadonly, http.StatusInternalServerError)
}

func NewDatabaseInvalidDriver(trace error) *Error {
	return New(trace, MessageDatabaseInvalidDriver, CodeDatabaseInvalidDriver, http.StatusInternalServerError)
}

func NewDatabaseSQLiteDriver(trace error) *Error {
	return New(trace, MessageDatabaseSQLiteDriver, CodeDatabaseSQLiteDriver, http.StatusInternalServerError)
}

func NewDatabaseInvalidArgument(trace error) *Error {
	return New(trace, MessageDatabaseInvalidArgument, CodeDatabaseInvalidArgument, http.StatusBadRequest)
}

func NewHTTPRequest(trace error) *Error {
	return New(trace, MessageHTTPRequest, CodeHTTPRequest, http.StatusBadRequest)
}

func NewHTTPResponse(trace error) *Error {
	return New(trace, MessageHTTPResponse, CodeHTTPResponse, http.StatusInternalServerError)
}

func NewHTTPBodyParser(trace error) *Error {
	return New(trace, MessageHTTPBodyParser, CodeHTTPBodyParser, http.StatusBadRequest)
}

func NewHTTPParams(trace error) *Error {
	return New(trace, MessageHTTPParams, CodeHTTPParams, http.StatusBadRequest)
}

func NewEmailNotFound(trace error) *Error {
	return New(trace, MessageEmailNotFound, CodeEmailNotFound, http.StatusNotFound)
}

func NewEmailSend(trace error) *Error {
	return New(trace, MessageEmailSend, CodeEmailSend, http.StatusInternalServerError)
}

func NewEmailTemplate(trace error) *Error {
	return New(trace, MessageEmailTemplate, CodeEmailTemplate, http.StatusInternalServerError)
}

func NewEmailInvalidAddress(trace error) *Error {
	return New(trace, MessageEmailInvalidAddress, CodeEmailInvalidAddress, http.StatusBadRequest)
}
