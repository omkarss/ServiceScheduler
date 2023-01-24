package deserializer

type Customer struct {
	FullName    string `json:"FullName" validate:"required"`
	PhoneNumber string `json:"PhoneNumber" validate:"required"`
	Type        string `json:"Type" validate:"required"`
}
