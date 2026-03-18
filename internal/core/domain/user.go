package domain

import (
	"fmt"
	"regexp"

	core_errors "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/errors"
)

type User struct {
	ID      int
	Version int

	FullName    string
	PhoneNumber *string
}

func NewUser(id int, version int, fullName string, phoneNumber *string) User {
	return User{
		ID:          id,
		Version:     version,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}
}

func NewUserUninitialized(fullName string, phoneNumber *string) User {
	return NewUser(
		UninitializedId,
		UninitializedVersion,
		fullName,
		phoneNumber,
	)
}

func (u *User) Validate() error {
	fullNameLength := len([]rune(u.FullName))
	if fullNameLength < 3 || fullNameLength > 100 {
		return fmt.Errorf(
			"full name length must be between 3 and 100, got %d: %w",
			fullNameLength,
			core_errors.ErrInvalidArgument,
		)
	}

	if u.PhoneNumber != nil {
		phoneNumberLength := len([]rune(*u.PhoneNumber))
		if phoneNumberLength < 10 || phoneNumberLength > 15 {
			return fmt.Errorf(
				"phoneNumber length must be between 10 and 15, got %d: %w",
				phoneNumberLength,
				core_errors.ErrInvalidArgument,
			)
		}
		re := regexp.MustCompile(`^\+[0-9]+$`)
		if !re.Match([]byte(*u.PhoneNumber)) {
			return fmt.Errorf("invalid phone number: %s: %w", *u.PhoneNumber, core_errors.ErrInvalidArgument)
		}
	}
	return nil
}
