package rating

import "errors"

var ErrNotOwner = errors.New("rating does not belong to this user")
var ErrCafeNotVisited = errors.New("cafe must be marked visited before rating")
