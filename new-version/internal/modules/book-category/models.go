package bookcategory

import "time"

type BasicBookCatResponse struct {
	Title       string
	CreatedTime time.Time
}

type BookCatResponse struct {
	Id int
	BasicBookCatResponse
}

type BookCatRequest struct {
	Title string
}
