package request

type RequestReportDTO struct {
	UserID       string `json:"userId" validate:"required,min=1"`
	ReportTypeID int32  `json:"reportTypeId" validate:"required,min=1"`
}
