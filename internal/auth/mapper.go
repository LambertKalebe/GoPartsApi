package auth

func toLoginResponse(token string) loginResponse {
	return loginResponse{
		Message: "Login realizado com sucesso",
		Token:   token,
	}
}

func toRegisterResponse() registerResponse {
	return registerResponse{
		Message: "User created",
	}
}

func toLogoutResponse() logoutResponse {
	return logoutResponse{
		Message: "Logout realizado com sucesso",
	}
}

func toAboutMeResponse(user user) aboutMeResponse {
	return aboutMeResponse{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
	}
}
