package wgplaner

import (
	"errors"
	"log"
)

type ErrorList struct {
	errors []error
}

func (errList *ErrorList) Add(str string) {
	errList.errors = append(errList.errors, errors.New(str))
}

func (errList *ErrorList) AddError(err error) {
	errList.errors = append(errList.errors, err)
}

func (errList *ErrorList) AddList(src *ErrorList) {
	for _, err := range src.errors {
		errList.AddError(err)
	}
}

func (errList *ErrorList) String() string {
	var msg string
	for _, err := range errList.errors {
		msg += err.Error() + "\n"
	}
	return msg
}

func (errList *ErrorList) HasErrors() bool {
	return len(errList.errors) > 0
}

func (errList *ErrorList) Print() {
	for _, err := range errList.errors {
		log.Println(err.Error())
	}
}
