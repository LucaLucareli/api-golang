package request

type LoginRequestDTO struct {
	Document string `json:"document" validate:"required,document"`
	Password string `json:"password" validate:"required,min=8,max=30"`
}
