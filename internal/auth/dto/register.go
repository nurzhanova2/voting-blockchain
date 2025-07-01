package dto

// структура, которую получаем от клиента (POST /register)
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
// что вернём клиенту в ответ
type RegisterResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}
