package db_utils

import "math"

type GetAllSerializer struct {
	Data        any
	TotalCount  int64
	Took        int
	CurrentPage int
}

type GetAllResponse struct {
	Data        any `json:"data"`
	Took        int `json:"took"`
	TotalCount  int `json:"totalCount"`
	CurrentPage int `json:"currentPage,omitempty"`
	TotalPages  int `json:"totalPages,omitempty"`
}

func (s *GetAllSerializer) Response() GetAllResponse {

	return GetAllResponse{
		Data:        s.Data,
		Took:        s.Took,
		TotalCount:  int(s.TotalCount),
		CurrentPage: s.CurrentPage,
		TotalPages:  int(math.Ceil(float64(s.TotalCount) / float64(s.Took))),
	}
}
