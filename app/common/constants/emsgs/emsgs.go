package emsgs

import "fmt"

const (
	IsEmpty           = "is-empty"
	NotEmail          = "not-email"
	NotURL            = "not-url"
	NotDigitsOnly     = "not-digits-only"
	NotPositiveNumber = "not-positive-number"
	NotNumber         = "not-number"
)

const (
	NotString  = "not-string"
	NotInteger = "not-integer"
	NotFloat   = "not-float"
	NotArray   = "not-array"
	NotObject  = "not-object"
	NotBoolean = "not-boolean"
	NotDate    = "not-date"
)

const (
	ObjectAlreadyExists = "object-already-exists"
	ObjectNotFound      = "object-not-found"
)

const (
	Internal                  = "internal"
	NoPermissions             = "no-permissions"
	Banned                    = "banned"
	TooManyConcurrentRequests = "too-many-concurrent-requests"
)

func GetWrongError(fieldName string) string {
	return fmt.Sprintf("wrong-%s", fieldName)
}

func GetCouldNotBeMoreError(fromProperty, toProperty string) string {
	return fmt.Sprintf("%s-could-not-be-more-than-%s", fromProperty, toProperty)
}

func GetMinLengthError(length int) string {
	return fmt.Sprintf("min-length:%d", length)
}

func GetMaxLengthError(length int) string {
	return fmt.Sprintf("max-length:%d", length)
}

func GetLengthError(min, max int) string {
	return fmt.Sprintf("length:%d:%d", min, max)
}
