package common

// Access Level uses for authorization.
// If user has equal or higher access than set one,
// she/he has permission for action

type AccessLevel int

const (
	USER_ACCESS_LEVEL  AccessLevel = 50
	ADMIN_ACCESS_LEVEL AccessLevel = 100
)
