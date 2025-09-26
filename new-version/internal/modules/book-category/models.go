package bookcategory

import "time"

type BasicBookCatResponse struct {
	Title       string    `json:"title"`
	CreatedTime time.Time `json:"created_time"`
}

type BookCatResponse struct {
	Id int `json:"id"`
	BasicBookCatResponse
}

type BookCatRequest struct {
	Title string `json:"title"`
}

func (b BookCat) ToBasicResponse() BasicBookCatResponse {
	return BasicBookCatResponse{
		Title:       b.Title,
		CreatedTime: b.CreatedTime,
	}
}

func (b BookCat) ToResponse() BookCatResponse {
	return BookCatResponse{
		Id:                   b.Id,
		BasicBookCatResponse: b.ToBasicResponse(),
	}
}
