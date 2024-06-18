package constants

import "time"

// TODO move to env variable
const (
	SessionTokenValidTime = 2 * 24 * time.Hour
	RefreshTokenValidTime = 90 * 24 * time.Hour
)
