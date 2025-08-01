package domain

type DomainError struct {
	Err  error
	Code int
}