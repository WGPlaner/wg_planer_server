package models

import (
	"testing"

	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	assert.NoError(t, PrepareTestDatabase())

	uid := "1234567890fakefirebaseid9999"
	user1 := &User{UID: &uid}
	user2 := &User{UID: &uid, DisplayName: swag.String("TestUser")}

	err1 := CreateUser(user1)
	err2 := CreateUser(user2)

	assert.Error(t, err1)
	assert.NoError(t, err2)
}

func TestUpdateUser(t *testing.T) {
	assert.NoError(t, PrepareTestDatabase())
	uid1 := "1234567890fakefirebaseid0001"
	user1 := AssertExistsAndLoadBean(t, &User{UID: &uid1}).(*User)

	*user1.DisplayName = "New User Name"
	err1 := UpdateUser(user1)
	assert.NoError(t, err1)

	userUpdated := AssertExistsAndLoadBean(t, &User{UID: &uid1}).(*User)
	assert.Equal(t, *userUpdated.DisplayName, "New User Name")
}

func TestUpdateUserCols(t *testing.T) {
	assert.NoError(t, PrepareTestDatabase())
	uid1 := "1234567890fakefirebaseid0001"
	user1 := AssertExistsAndLoadBean(t, &User{UID: &uid1}).(*User)

	*user1.DisplayName = "New User Name"
	user1.Email = "new@email.org"
	err1 := UpdateUserCols(user1, `display_name`)
	assert.NoError(t, err1)

	userUpdated := AssertExistsAndLoadBean(t, &User{UID: &uid1}).(*User)
	assert.Equal(t, *userUpdated.DisplayName, "New User Name")
	assert.NotEqual(t, userUpdated.Email, "new@email.org")
}

func TestUser_LeaveGroup(t *testing.T) {
	assert.NoError(t, PrepareTestDatabase())

	uid1 := "1234567890fakefirebaseid0001"
	userOld := AssertExistsAndLoadBean(t, &User{UID: &uid1}).(*User)

	err := userOld.LeaveGroup()
	assert.NoError(t, err)
	assert.Empty(t, userOld.GroupUID)

	userUpdated := AssertExistsAndLoadBean(t, &User{UID: &uid1}).(*User)
	assert.Empty(t, userUpdated.GroupUID)
}

func TestGetUserByUID(t *testing.T) {
	assert.NoError(t, PrepareTestDatabase())
	validUserIDs := []string{
		"1234567890fakefirebaseid0001",
		"1234567890fakefirebaseid0002",
		"1234567890fakefirebaseid0003",
	}
	for _, id := range validUserIDs {
		_, err := GetUserByUID(id)
		assert.NoError(t, err)
	}

	_, err := GetUserByUID("1234567890fakefirebaseid9999")
	assert.True(t, IsErrUserNotExist(err))
}

func TestAreUsersExist(t *testing.T) {
	assert.NoError(t, PrepareTestDatabase())
	validUserIDs := []string{
		"1234567890fakefirebaseid0001",
		"1234567890fakefirebaseid0002",
		"1234567890fakefirebaseid0003",
	}
	exist1, err1 := AreUsersExist(validUserIDs)
	assert.NoError(t, err1)
	assert.True(t, exist1)

	invalidUserIDs := []string{
		"1234567890fakefirebaseid9998",
		"1234567890fakefirebaseid9999",
	}
	exist2, err2 := AreUsersExist(invalidUserIDs)
	assert.NoError(t, err2)
	assert.False(t, exist2)
}

func TestIsUserExist(t *testing.T) {
	assert.NoError(t, PrepareTestDatabase())
	validUserIDs := []string{
		"1234567890fakefirebaseid0001",
		"1234567890fakefirebaseid0002",
		"1234567890fakefirebaseid0003",
	}
	for _, id := range validUserIDs {
		exist, err := IsUserExist(id)
		assert.NoError(t, err)
		assert.True(t, exist)
	}

	notExistUserIDs := []string{
		"1234567890fakefirebaseid9998",
		"1234567890fakefirebaseid9999",
	}
	for _, id := range notExistUserIDs {
		exist, err := IsUserExist(id)
		assert.NoError(t, err)
		assert.False(t, exist)
	}

	invalidUserIDs := []string{
		"",
		"short",
		"1234567890fakefirebaseid9999TooLong",
	}
	for _, id := range invalidUserIDs {
		exist, err := IsUserExist(id)
		assert.True(t, IsErrUserInvalidUID(err))
		assert.False(t, exist)
	}
}
