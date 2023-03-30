package utils

import "errors"

var NoSuchDataWithGivenID error
var InvalidTokenError error
var InvalidStaffAccessError error
var InvalidAdminAccessError error

func init() {
	NoSuchDataWithGivenID = errors.New("no such data with given id")
	InvalidTokenError = errors.New("invalid token")
	InvalidStaffAccessError = errors.New("invalid staff access token")
	InvalidAdminAccessError = errors.New("invalid admin access token")

}
