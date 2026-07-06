package http

import (
	"fmt"

	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/auth"
)

const (
	ErrorCodeInvalidRequest = "invalid_request"
	ErrorCodeUnauthorized   = "unauthorized"
	ErrorCodeForbidden      = "forbidden"
	ErrorCodeNotFound       = "not_found"
	ErrorCodeConflict       = "conflict"
	ErrorCodeDomainRejected = "domain_rejected"
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

	MessageInvalidRequestBody = "Request body is malformed."
	MessageInvalidCredentials = "Invalid email or password."
	MessageInvalidEmail       = "Email address is invalid."
	MessageEmailTaken         = "Email is already registered."
	MessageSessionExpired     = "Session expired. Please sign in again."
	MessageEmailNotVerified   = "Your Google email is not verified."
	MessageAuthFailed         = "Authentication failed."
	// MessageAccountExistsUsePassword steers a Google sign-in to the existing password
	// login when auto-linking would be unsafe (the existing account is not email-verified).
	MessageAccountExistsUsePassword = "An account with this email already exists. Sign in with your password."
	// MessageGoogleAccountLinked covers the rare race where a Google identity is already
	// linked to a different user.
	MessageGoogleAccountLinked = "This Google account is already linked to another user."
	// MessageEmailVerificationInvalid covers a missing/expired/already-used verification link.
	MessageEmailVerificationInvalid = "This verification link is invalid or has expired."
	// MessageEmailAlreadyVerified is returned when a resend is requested but nothing needs verifying.
	MessageEmailAlreadyVerified = "Your email is already verified."
)

// MessageWeakPassword is derived from auth.MinPasswordLength so the user-facing minimum
// and the validator can never drift apart.
var MessageWeakPassword = fmt.Sprintf("Password must be at least %d characters.", auth.MinPasswordLength)

func invalidQuery(field string, message string) *errorResponse {
	return &errorResponse{
		Code:    ErrorCodeInvalidRequest,
		Message: message,
		Details: map[string]any{
			"field": field,
		},
	}
}
