package service

import (
	"context"

	"github.com/Bhinneka/user-service/src/auth/v1/model"
)

// LDAPService interface for abstracting LDAP service
type LDAPService interface {
	Auth(context.Context, string, string) (*model.LDAPProfile, error)
}
