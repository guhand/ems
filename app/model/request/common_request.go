package request

type CommonRequest struct {
	Page   uint   `form:"page"`
	Search string `form:"search"`
}

type DateFilters struct {
	Year  int `form:"year"`
	Month int `form:"month"`
}
type CommonRequestWithDateFilter struct {
	CommonRequest
	DateFilters
}
