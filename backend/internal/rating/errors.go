package rating

import "errors"

var ErrNotOwner = errors.New("rating does not belong to this user")
