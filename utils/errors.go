package utils

import (
	"errors"
	"sdu.store/server/validators"
)

var NoSuchDataWithGivenID error
var InvalidTokenError error
var InvalidStaffAccessError error
var InvalidAdminAccessError error
var NotExistItemError error

func init() {
	NoSuchDataWithGivenID = errors.New("no such data with given id")
	InvalidTokenError = errors.New("invalid token")
	InvalidStaffAccessError = errors.New("invalid staff access token")
	InvalidAdminAccessError = errors.New("invalid admin access token")
	NotExistItemError = errors.New(validators.NotExistItem)
}
