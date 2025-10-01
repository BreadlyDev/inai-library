package user

import (
	"time"

	"github.com/google/uuid"
)

type UserModel struct {
	Id          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	PassHash    string    `json:"pass_hash"`
	JoinedIn    time.Time `json:"joined_in"`
	AccessLevel int       `json:"access_level"`
}

type User struct {
	Id    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Pass  string    `json:"pass"`
}

type UserInfo struct {
	Id          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	JoinedIn    time.Time `json:"joined_in"`
	AccessLevel int       `json:"access_level"`
}

type UserCreate struct {
	Email string `json:"email"`
	Pass  string `json:"pass"`
}

type UserLogin struct {
	Email       string `json:"email"`
	Pass        string `json:"pass"`
	AccessLevel int    `json:"access_level"`
}
