package common

import "errors"

var ErrUserNotFound error = errors.New("user not found")
var ErrTeamNotFound error = errors.New("team not found")
var ErrTeamDuplicate error = errors.New("team already exists")
var ErrPRExists error = errors.New("pull request already exists")
var ErrPRNotFound error = errors.New("pull request not found")
var ErrAuthorNotFound error = errors.New("author not found")
var ErrNotAssignedReviewer error = errors.New("user is not an assigned reviewer for this pull request")
var ErrNoCandidates error = errors.New("no active replacement candidate in team")
