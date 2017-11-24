package models

import (
	"errors"
	"fmt"
	"log"

	"github.com/go-openapi/strfmt"
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

//  _   _
// | | | |___  ___ _ __
// | | | / __|/ _ \ '__|
// | |_| \__ \  __/ |
//  \___/|___/\___|_|
//

// ErrUserAlreadyExist represents a "user already exists" error.
type ErrUserAlreadyExist struct {
	UID string
}

// IsErrUserAlreadyExist checks if an error is a ErrUserAlreadyExists.
func IsErrUserAlreadyExist(err error) bool {
	_, ok := err.(ErrUserAlreadyExist)
	return ok
}

func (err ErrUserAlreadyExist) Error() string {
	return fmt.Sprintf("user already exists [error: %s]", err.UID)
}

// ErrUserNotExist represents a "UserNotExist" kind of error.
type ErrUserNotExist struct {
	UID string
}

// IsErrUserNotExist checks if an error is a ErrUserNotExist.
func IsErrUserNotExist(err error) bool {
	_, ok := err.(ErrUserNotExist)
	return ok
}

func (err ErrUserNotExist) Error() string {
	return fmt.Sprintf("user does not exist [uid: %s]", err.UID)
}

//   ____
//  / ___|_ __ ___  _   _ _ __
// | |  _| '__/ _ \| | | | '_ \
// | |_| | | | (_) | |_| | |_) |
//  \____|_|  \___/ \__,_| .__/
//                       |_|

// ErrGroupNotExist represents a "GroupNotExist" kind of error.
type ErrGroupNotExist struct {
	UID strfmt.UUID
}

// IsErrGroupNotExist checks if an error is a ErrGroupNotExist.
func IsErrGroupNotExist(err error) bool {
	_, ok := err.(ErrGroupNotExist)
	return ok
}

func (err ErrGroupNotExist) Error() string {
	return fmt.Sprintf("group does not exist [uid: %s]", err.UID)
}

//  ____  _                       _               _     _     _
// / ___|| |__   ___  _ __  _ __ (_)_ __   __ _  | |   (_)___| |_
// \___ \| '_ \ / _ \| '_ \| '_ \| | '_ \ / _` | | |   | / __| __|
//  ___) | | | | (_) | |_) | |_) | | | | | (_| | | |___| \__ \ |_
// |____/|_| |_|\___/| .__/| .__/|_|_| |_|\__, | |_____|_|___/\__|
//                   |_|   |_|            |___/

// ErrGroupNotExist represents a "ListItemNotExist" kind of error.
type ErrListItemNotExist struct {
	ID       strfmt.UUID
	GroupUID strfmt.UUID
}

// IsErrListItemNotExist checks if an error is a ErrListItemNotExist.
func IsErrListItemNotExist(err error) bool {
	_, ok := err.(ErrListItemNotExist)
	return ok
}

func (err ErrListItemNotExist) Error() string {
	return fmt.Sprintf("list item does not exist [groupUID: %s, uid: %s]",
		err.GroupUID, err.ID)
}

//  ____  _ _ _
// | __ )(_) | |
// |  _ \| | | |
// | |_) | | | |
// |____/|_|_|_|
