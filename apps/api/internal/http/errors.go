package http

const (
	ErrorCodeInvalidRequest = "invalid_request"
	ErrorCodeNotFound       = "not_found"
	ErrorCodeInternal       = "internal_error"
)

const (
	MessageFailedLoadVenue       = "Failed to load venue."
	MessageFailedLoadVenues      = "Failed to load venues."
	MessageVenueNotFound         = "Venue was not found."
	MessageVenueSlugInvalid      = "Venue slug is required and must be 160 characters or fewer."
	MessageQueryTooLong          = "Query must be 120 characters or fewer."
	MessageCityTooLong           = "City must be 80 characters or fewer."
	MessageDistrictTooLong       = "District must be 80 characters or fewer."
	MessagePlatformInvalid       = "Platform must be one of tiktok, instagram, youtube, facebook, or other."
	MessageLatitudeRangeInvalid  = "Latitude must be between -90 and 90."
	MessageLongitudeRangeInvalid = "Longitude must be between -180 and 180."
	MessageRadiusRangeInvalid    = "Radius must be between 0 and 50000 meters."
	MessageMaxPriceInvalid       = "Max price must be greater than or equal to 0."
	MessageMinPriceInvalid       = "Min price must be greater than or equal to 0."
	MessagePriceRangeInvalid     = "Min price must be less than or equal to max price."
	MessageLimitRangeInvalid     = "Limit must be between 1 and 100."
	MessageOpenNowInvalid        = "Open now must be true or false."
	MessageSortInvalid           = "Sort must be one of trending, videos, distance, or price."
	MessageDistanceSortLocation  = "Distance sort requires lat and lng."
)

func invalidQuery(field string, message string) *errorResponse {
	return &errorResponse{
		Code:    ErrorCodeInvalidRequest,
		Message: message,
		Details: map[string]any{
			"field": field,
		},
	}
}
