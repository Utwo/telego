package utils

import "errors"

type JsonError struct {
	Message string `json:"message"`
}

var LimitProjectReached = errors.New("the maximum number of projects has been reached for this account")
var LimitCollabReached = errors.New("the maximum number of collaborators and invitations has been reached for this project")
var LimitAssetReached = errors.New("the maximum number of assets has been reached for this project")
var SizeAssetReached = errors.New("the maximum assets size has been reached for this project")