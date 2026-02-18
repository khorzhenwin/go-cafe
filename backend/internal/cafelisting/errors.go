package cafelisting

import "errors"

var ErrNotOwner = errors.New("cafe listing does not belong to this user")
var ErrInvalidVisitStatus = errors.New("invalid visit_status: must be to_visit or visited")
