package common

import "errors"

var ErrUserNotFound error = errors.New("user not found")
var ErrTeamNotFound error = errors.New("team not found")
var ErrTeamDuplicate error = errors.New("team already exists")
