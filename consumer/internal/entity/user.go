package user

type User struct {
	ID      uint64  `json:"id"`
	Balance float32 `json:"balance"`
}

type UserWithBalance struct {
	ID      uint64  `json:"id"`
	Balance float32 `json:"balance"`
	Method  string  `json:"method"`
}
