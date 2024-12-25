package model

type CommonResponse struct {
	Data any `json:"data,omitempty"`
	Meta any `json:"meta,omitempty"`
}

type Meta struct {
	Count int `json:"count,omitempty"`
	Limit int `json:"limit,omitempty"`
	Page  int `json:"page,omitempty"`
	Total int `json:"total,omitempty"`
}
