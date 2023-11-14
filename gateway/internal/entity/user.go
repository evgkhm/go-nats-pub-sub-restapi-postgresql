package user

type User struct {
	ID     uint64 `json:"id" validate:"required"`
	Method string `json:"method" validate:"required"`
}

type UserWithBalance struct {
	ID      uint64  `json:"id" validate:"required"`
	Balance float32 `json:"balance" validate:"required"`
	Method  string  `json:"method" validate:"required"`
}
