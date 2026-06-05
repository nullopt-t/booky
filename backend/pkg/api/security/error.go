package security

import (
	"booky-backend/pkg/api"
	"log/slog"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
)

const (
	CodeValidation   = "VALIDATION_ERROR"
	CodeInternal     = "INTERNAL_ERROR"
	CodeConflict     = "CONFLICT_ERROR"
	CodeAuth         = "AUTH_ERROR"
	CodeNotFound     = "NOT_FOUND"
	CodeUnauthorized = "UNAUTHORIZED"
)

type SecureError struct {
	Status   int
	Code     string
	UserMsg  string
	Internal error
	Fields   []api.FieldError
	MetaData map[string]any
}

func (se *SecureError) Error() string {
	return se.UserMsg
}

func (se *SecureError) Unwrap() error {
	return se.Internal
}

func (se *SecureError) LogMessage() string {
	return slog.GroupValue(
		slog.String("code", se.Code),
		slog.String("message", se.UserMsg),
		slog.Any("meta", se.MetaData),
		slog.Any("internal", se.Internal),
	).String()
}

func NewSecureError(
	status int,
	code,
	public string,
	internal error,
) *SecureError {
	return &SecureError{
		Status:   status,
		Code:     code,
		UserMsg:  public,
		Internal: internal,
	}
}

func (se *SecureError) WithMetaData(
	meta map[string]any,
) *SecureError {
	se.MetaData = meta
	return se
}

func (se *SecureError) WithFields(
	fields []api.FieldError,
) *SecureError {
	se.Fields = fields
	return se
}
