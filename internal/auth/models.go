package auth

type user struct {
	ID       int
	Username string
	PassHash string
	Role     string
}
