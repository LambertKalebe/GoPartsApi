package auth

type loginRequest struct {
	Username string `json:"username" example:"user"`
	Password string `json:"password" example:"123"`
}

type registerRequest struct {
	Username string `json:"username" example:"user"`
	Password string `json:"password" example:"123"`
}

type loginResponse struct {
	Message string `json:"message" example:"Login realizado com sucesso"`
	Token   string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type registerResponse struct {
	Message string `json:"message" example:"User created"`
}

type logoutResponse struct {
	Message string `json:"message" example:"Logout realizado com sucesso"`
}

type aboutMeResponse struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}
