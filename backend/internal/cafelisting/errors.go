package cafelisting

import "errors"

var ErrNotOwner = errors.New("cafe listing does not belong to this user")
var ErrInvalidVisitStatus = errors.New("invalid visit_status: must be to_visit or visited")
var ErrInvalidCafeName = errors.New("cafe name is required")
var ErrInvalidCoordinates = errors.New("latitude and longitude must be provided together and within valid ranges")
