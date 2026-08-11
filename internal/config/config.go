package config

var (
	// Sim, eu sei que é uma boa prática colocar variáveis de ambiente, mas isso nem vai pra deploy provavelmente, e se for, é só mudar aqui
	DatabaseUrl = "internal/database/cartalogo.sqlite"
	JWTSecret   = []byte("IHave11DogsAnd1Cat")
)
