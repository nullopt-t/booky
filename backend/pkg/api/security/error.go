package security

import (
	"fmt"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
)

type SecureError struct {
	Code     string
	UserMsg  string
	Internal error
	MetaData map[string]string
}

func (se *SecureError) Error() string {
	return se.UserMsg
}

func (se *SecureError) LogMessage() string {
	return fmt.Sprintf(
		"%s[%s]%s %smsg=%s%s internal=%v meta=%v",
		Yellow,
		se.Code,
		Reset,
		Green,
		se.UserMsg,
		Reset,
		se.Internal,
		se.MetaData,
	)
}
