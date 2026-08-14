package config

import (
	"os"
	"path/filepath"
)

var (
	// Sim, eu sei que é uma boa prática colocar variáveis de ambiente, mas isso nem vai pra deploy provavelmente, e se for, é só mudar aqui
	JWTSecret = []byte("IHave11DogsAnd1Cat")
)

func GetDatabasePath() string {
	wd, err := os.Getwd()
	if err == nil {
		path := filepath.Join(
			wd,
			"internal",
			"database",
			"cartalogo.sqlite",
		)

		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	exe, err := os.Executable()
	if err != nil {
		return ""
	}

	return filepath.Join(
		filepath.Dir(exe),
		"database",
		"cartalogo.sqlite",
	)
}
