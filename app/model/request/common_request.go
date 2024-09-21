package request

type CommonRequest struct {
	Page   uint   `form:"page"`
	Search string `form:"search"`
}
