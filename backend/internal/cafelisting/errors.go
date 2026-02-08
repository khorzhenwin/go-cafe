package cafelisting

import "errors"

var ErrNotOwner = errors.New("cafe listing does not belong to this user")
