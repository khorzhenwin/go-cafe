package rating

import "errors"

var ErrNotOwner = errors.New("rating does not belong to this user")
var ErrCafeNotVisited = errors.New("cafe must be marked visited before rating")
var ErrInvalidRatingValue = errors.New("rating must be between 1 and 5")
var ErrDuplicateRating = errors.New("you already reviewed this cafe")
