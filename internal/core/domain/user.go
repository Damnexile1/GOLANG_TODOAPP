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

type UserPatch struct {
	FullName    Nullable[string]
	PhoneNumber Nullable[string]
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

func (p *UserPatch) Validate() error {
	if p.FullName.Set && p.FullName.Value == nil {
		return fmt.Errorf("UserPatch: full name must be not null: %w", core_errors.ErrInvalidArgument)
	}

	return nil
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

func (u *User) ApplyPatch(patch UserPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate user patch: %w", err)
	}

	tmp := *u

	if patch.FullName.Set {
		tmp.FullName = *patch.FullName.Value
	}

	if patch.PhoneNumber.Set {
		tmp.PhoneNumber = patch.PhoneNumber.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate user patch: %w", err)
	}

	*u = tmp
	return nil
}

func NewUserPatch(fullName Nullable[string], phoneNumber Nullable[string]) UserPatch {
	return UserPatch{
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}
}
