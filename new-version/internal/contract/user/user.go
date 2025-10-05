package user

import (
	"time"

	"github.com/google/uuid"
)

type Model struct {
	Id          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	Password    string    `json:"pass_hash"`
	JoinedAt    time.Time `json:"joined_at"`
	AccessLevel int       `json:"access_level"`
}

type Request struct {
	Email    string `json:"email"`
	Password string `json:"pass_hash"`
}

type Response struct {
	Id    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

type InfoResponse struct {
	Id          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	JoinedAt    time.Time `json:"joined_at"`
	AccessLevel int       `json:"access_level"`
}

// type User struct {
// 	Id    uuid.UUID `json:"id"`
// 	Email string    `json:"email"`
// 	Pass  string    `json:"pass"`
// }

// type UserInfo struct {
// 	Id          uuid.UUID `json:"id"`
// 	Email       string    `json:"email"`
// 	JoinedAt    time.Time `json:"joined_in"`
// 	AccessLevel int       `json:"access_level"`
// 	Pass        string    `json:"pass"`
// }

// type UserCreate struct {
// 	Email string `json:"email"`
// 	Pass  string `json:"pass"`
// }

// type UserLogin struct {
// 	Email string `json:"email"`
// 	Pass  string `json:"pass"`
// }

// type UserLoginResp struct {
// 	Email       string `json:"email"`
// 	Pass        string `json:"pass"`
// 	AccessLevel int    `json:"access_level"`
// }
