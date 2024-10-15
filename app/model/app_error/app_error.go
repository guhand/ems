package apperror

import "fmt"

func UniqueKeyError(field string) error {
	return fmt.Errorf("%s already exists", field)
}

func DataNotFoundError(field string) error {
	return fmt.Errorf("%s not found", field)
}
