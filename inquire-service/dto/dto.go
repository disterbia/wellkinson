package dto

type GetInquireParams struct {
	Page      int    `form:"page"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
}
