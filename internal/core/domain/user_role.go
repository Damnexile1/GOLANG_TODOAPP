package domain

import "fmt"

type UserRole int

const (
	UserRoleUser    UserRole = 1
	UserRoleManager UserRole = 2
	UserRoleAdmin   UserRole = 3
)

func (ur UserRole) String() string {
	switch ur {
	case UserRoleUser:
		return "user"
	case UserRoleManager:
		return "manager"
	case UserRoleAdmin:
		return "admin"
	default:
		return fmt.Sprintf("unknown(%d)", ur)
	}
}

func (ur UserRole) IsValid() bool {
	return ur == UserRoleUser || ur == UserRoleManager || ur == UserRoleAdmin
}

func UserRoleFromInt(value int) (UserRole, error) {
	role := UserRole(value)
	if !role.IsValid() {
		return 0, fmt.Errorf("invalid user role: %d", value)
	}
	return role, nil
}

func UserRoleFromString(value string) (UserRole, error) {
	switch value {
	case "user":
		return UserRoleUser, nil
	case "manager":
		return UserRoleManager, nil
	case "admin":
		return UserRoleAdmin, nil
	default:
		return 0, fmt.Errorf("invalid user role string: %s", value)
	}
}

// HasPermission проверяет, имеет ли роль необходимые права
func (ur UserRole) HasPermission(requiredRole UserRole) bool {
	// Admin имеет все права
	if ur == UserRoleAdmin {
		return true
	}
	// Manager имеет права Manager и User
	if ur == UserRoleManager && (requiredRole == UserRoleManager || requiredRole == UserRoleUser) {
		return true
	}
	// User имеет только свои права
	if ur == UserRoleUser && requiredRole == UserRoleUser {
		return true
	}
	return false
}
