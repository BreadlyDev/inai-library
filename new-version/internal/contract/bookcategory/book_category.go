package bookcategory

import "time"

type Request struct {
	Title string `json:"title"`
}

type Response struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	CreatedTime time.Time `json:"created_time"`
}
