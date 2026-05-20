package domain

import "fmt"

type ErrorCode string

const (
	ErrNotFoundCode      ErrorCode = "NOT_FOUND"
	ErrAlreadyExistsCode ErrorCode = "ALREADY_EXISTS"
	ErrValidationCode    ErrorCode = "VALIDATION_ERROR"
	ErrRepositoryCode    ErrorCode = "REPOSITORY_ERROR"
	ErrAlgorithmRunCode  ErrorCode = "ALGORITHM_RUN_ERROR"
)

type DomainError struct {
	code    ErrorCode
	message string
}

func NewDomainError(code ErrorCode) *DomainError {
	return &DomainError{
		code:    code,
		message: string(code),
	}
}

func NewDomainErrorWithMessage(code ErrorCode, message string) *DomainError {
	return &DomainError{
		code:    code,
		message: message,
	}
}

func (e *DomainError) Error() string {
	return fmt.Sprintf("%s: %s", e.code, e.message)
}

func (e *DomainError) Code() ErrorCode {
	return e.code
}

func (e *DomainError) Message() string {
	return e.message
}
