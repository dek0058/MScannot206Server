package repository

import "MScannot206/shared/entity"

type UserRepository interface {
	Start() error
	Stop() error

	FindUserByUids(uids []string) ([]*entity.User, []string, error)
	InsertUserByUids(uids []string) ([]*entity.User, error)
}
