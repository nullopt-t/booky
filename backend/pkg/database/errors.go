package database

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	// Class 08 — Connection Exception
	ErrConnectionException                        = errors.New("connection exception")
	ErrConnectionDoesNotExist                     = errors.New("connection does not exist")
	ErrConnectionFailure                          = errors.New("connection failure")
	ErrSQLClientUnableToEstablishConnection       = errors.New("client unable to establish connection")
	ErrSQLServerRejectedEstablishmentOfConnection = errors.New("server rejected establishment of connection")
	ErrTransactionResolutionUnknown               = errors.New("transaction resolution unknown")
	ErrProtocolViolation                          = errors.New("protocol violation")

	// Class 22 — Data Exception
	ErrInvalidTextRepresentation = errors.New("invalid text representation")
	ErrNumericValueOutOfRange    = errors.New("numeric value out of range")
	ErrInvalidDatetimeFormat     = errors.New("invalid datetime format")
	ErrDivisionByZero            = errors.New("division by zero")
	ErrStringDataRightTruncation = errors.New("string data right truncation")
	ErrNullValueNotAllowed       = errors.New("null value not allowed")
	ErrInvalidParameterValue     = errors.New("invalid parameter value")

	// Class 23 — Integrity Constraint Violation
	ErrIntegrityConstraintViolation = errors.New("integrity constraint violation")
	ErrNotNullViolation             = errors.New("not null violation")
	ErrForeignKeyViolation          = errors.New("foreign key violation")
	ErrConflict                     = errors.New("conflict") // unique_violation
	ErrCheckViolation               = errors.New("check violation")
	ErrExclusionViolation           = errors.New("exclusion violation")

	// Class 25 — Invalid Transaction State
	ErrInvalidTransactionState = errors.New("invalid transaction state")
	ErrActiveSQLTransaction    = errors.New("active sql transaction")
	ErrNoActiveSQLTransaction  = errors.New("no active sql transaction")
	ErrInReadOnlyTransaction   = errors.New("in read-only transaction")
	ErrInFailedSQLTransaction  = errors.New("in failed sql transaction")

	// Class 28 — Invalid Authorization Specification
	ErrInvalidAuthorizationSpecification = errors.New("invalid authorization specification")
	ErrInvalidPassword                   = errors.New("invalid password")

	// Class 40 — Transaction Rollback
	ErrTransactionRollback                     = errors.New("transaction rollback")
	ErrSerializationFailure                    = errors.New("serialization failure")
	ErrDeadlockDetected                        = errors.New("deadlock detected")
	ErrTransactionIntegrityConstraintViolation = errors.New("transaction integrity constraint violation")

	// Class 42 — Syntax Error or Access Rule Violation
	ErrSyntaxError           = errors.New("syntax error")
	ErrUndefinedTable        = errors.New("undefined table")
	ErrUndefinedColumn       = errors.New("undefined column")
	ErrUndefinedFunction     = errors.New("undefined function")
	ErrDuplicateTable        = errors.New("duplicate table")
	ErrDuplicateColumn       = errors.New("duplicate column")
	ErrInsufficientPrivilege = errors.New("insufficient privilege")
	ErrAmbiguousColumn       = errors.New("ambiguous column")

	// Class 53 — Insufficient Resources
	ErrInsufficientResources      = errors.New("insufficient resources")
	ErrDiskFull                   = errors.New("disk full")
	ErrOutOfMemory                = errors.New("out of memory")
	ErrTooManyConnections         = errors.New("too many connections")
	ErrConfigurationLimitExceeded = errors.New("configuration limit exceeded")

	// Class 57 — Operator Intervention
	ErrTimeOut          = errors.New("timeout")
	ErrAdminShutdown    = errors.New("admin shutdown")
	ErrCrashShutdown    = errors.New("crash shutdown")
	ErrCannotConnectNow = errors.New("cannot connect now")
	ErrDatabaseDropped  = errors.New("database dropped")

	// Class 58 — System Error
	ErrIOError       = errors.New("io error")
	ErrUndefinedFile = errors.New("undefined file")
	ErrDuplicateFile = errors.New("duplicate file")

	ErrNotFound = errors.New("not found")
)

func MapError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	if pgErr, ok := err.(*pgconn.PgError); ok {
		switch pgErr.Code {

		// Class 08 — Connection Exception
		case "08000":
			return ErrConnectionException
		case "08003":
			return ErrConnectionDoesNotExist
		case "08006":
			return ErrConnectionFailure
		case "08001":
			return ErrSQLClientUnableToEstablishConnection
		case "08004":
			return ErrSQLServerRejectedEstablishmentOfConnection
		case "08007":
			return ErrTransactionResolutionUnknown
		case "08P01":
			return ErrProtocolViolation

		// Class 22 — Data Exception
		case "22P02":
			return ErrInvalidTextRepresentation
		case "22003":
			return ErrNumericValueOutOfRange
		case "22007":
			return ErrInvalidDatetimeFormat
		case "22012":
			return ErrDivisionByZero
		case "22001":
			return ErrStringDataRightTruncation
		case "22004":
			return ErrNullValueNotAllowed
		case "22023":
			return ErrInvalidParameterValue

		// Class 23 — Integrity Constraint Violation
		case "23000":
			return ErrIntegrityConstraintViolation
		case "23502":
			return ErrNotNullViolation
		case "23503":
			return ErrForeignKeyViolation
		case "23505":
			return ErrConflict
		case "23514":
			return ErrCheckViolation
		case "23P01":
			return ErrExclusionViolation

		// Class 25 — Invalid Transaction State
		case "25000":
			return ErrInvalidTransactionState
		case "25001":
			return ErrActiveSQLTransaction
		case "25P01":
			return ErrNoActiveSQLTransaction
		case "25006":
			return ErrInReadOnlyTransaction
		case "25P02":
			return ErrInFailedSQLTransaction

		// Class 28 — Invalid Authorization
		case "28000":
			return ErrInvalidAuthorizationSpecification
		case "28P01":
			return ErrInvalidPassword

		// Class 40 — Transaction Rollback
		case "40000":
			return ErrTransactionRollback
		case "40001":
			return ErrSerializationFailure
		case "40P01":
			return ErrDeadlockDetected
		case "40002":
			return ErrTransactionIntegrityConstraintViolation

		// Class 42 — Syntax Error or Access Rule Violation
		case "42601":
			return ErrSyntaxError
		case "42P01":
			return ErrUndefinedTable
		case "42703":
			return ErrUndefinedColumn
		case "42883":
			return ErrUndefinedFunction
		case "42P07":
			return ErrDuplicateTable
		case "42701":
			return ErrDuplicateColumn
		case "42501":
			return ErrInsufficientPrivilege
		case "42702":
			return ErrAmbiguousColumn

		// Class 53 — Insufficient Resources
		case "53000":
			return ErrInsufficientResources
		case "53100":
			return ErrDiskFull
		case "53200":
			return ErrOutOfMemory
		case "53300":
			return ErrTooManyConnections
		case "53400":
			return ErrConfigurationLimitExceeded

		// Class 57 — Operator Intervention
		case "57014":
			return ErrTimeOut
		case "57P01":
			return ErrAdminShutdown
		case "57P02":
			return ErrCrashShutdown
		case "57P03":
			return ErrCannotConnectNow
		case "57P04":
			return ErrDatabaseDropped

		// Class 58 — System Error
		case "58030":
			return ErrIOError
		case "58P01":
			return ErrUndefinedFile
		case "58P02":
			return ErrDuplicateFile
		}
	}
	return err
}
