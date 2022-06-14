package request

import validation "github.com/go-ozzo/ozzo-validation/v4"

type CurdDTO struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	City  string `json:"city"`
}

func (i CurdDTO) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.Name, validation.Required, validation.Length(1, 100)))
}
