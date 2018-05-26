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

// ErrUserInvalidUID represents a "UserNotExist" kind of error.
type ErrUserInvalidUID struct {
	UID string
}

// IsErrUserInvalidUID checks if an error is a ErrUserInvalidUID.
func IsErrUserInvalidUID(err error) bool {
	_, ok := err.(ErrUserInvalidUID)
	return ok
}

func (err ErrUserInvalidUID) Error() string {
	return fmt.Sprintf("invalid user id format [UID: %s]", err.UID)
}

// ErrUserMissingProperty represents a "UserNotExist" kind of error.
type ErrUserMissingProperty struct {
	Field string
}

// IsErrUserMissingProperty checks if an error is a ErrUserMissingProperty.
func IsErrUserMissingProperty(err error) bool {
	_, ok := err.(ErrUserMissingProperty)
	return ok
}

func (err ErrUserMissingProperty) Error() string {
	return fmt.Sprintf("missing required property [%s]", err.Field)
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

// ErrGroupCodeNotExist represents a "CodeNotExist" kind of error.
type ErrGroupCodeNotExist struct {
	Code string
}

// IsErrGroupCodeNotExist checks if an error is a ErrGroupCodeNotExist.
func IsErrGroupCodeNotExist(err error) bool {
	_, ok := err.(ErrGroupCodeNotExist)
	return ok
}

func (err ErrGroupCodeNotExist) Error() string {
	return fmt.Sprintf("group code does not exist [code: %s]", err.Code)
}

// ErrGroupInvalidUUID represents a "Invalid Group UUID" kind of error.
type ErrGroupInvalidUUID struct {
	UID string
}

// IsErrGroupInvalidUUID checks if an error is a ErrGroupInvalidUUID.
func IsErrGroupInvalidUUID(err error) bool {
	_, ok := err.(ErrGroupCodeNotExist)
	return ok
}

func (err ErrGroupInvalidUUID) Error() string {
	return fmt.Sprintf("invalid group UUID [%s]", err.UID)
}

//  ____  _                       _               _     _     _
// / ___|| |__   ___  _ __  _ __ (_)_ __   __ _  | |   (_)___| |_
// \___ \| '_ \ / _ \| '_ \| '_ \| | '_ \ / _` | | |   | / __| __|
//  ___) | | | | (_) | |_) | |_) | | | | | (_| | | |___| \__ \ |_
// |____/|_| |_|\___/| .__/| .__/|_|_| |_|\__, | |_____|_|___/\__|
//                   |_|   |_|            |___/

// ErrListItemNotExist represents a "ListItemNotExist" kind of error.
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

// ErrListItemHasBill represents a "ListItemNotExist" kind of error.
type ErrListItemHasBill struct {
	ID       strfmt.UUID
	GroupUID strfmt.UUID
}

// IsErrListItemHasBill checks if an error is a ErrListItemNotExist.
func IsErrListItemHasBill(err error) bool {
	_, ok := err.(ErrListItemHasBill)
	return ok
}

func (err ErrListItemHasBill) Error() string {
	return fmt.Sprintf("list item cannot be 'un-bought'. Already has a bill [groupUID: %s, uid: %s]",
		err.GroupUID, err.ID)
}

//  ____  _ _ _
// | __ )(_) | |
// |  _ \| | | |
// | |_) | | | |
// |____/|_|_|_|
