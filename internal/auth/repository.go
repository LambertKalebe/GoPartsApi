package auth

import "g0/internal/database"

const checkUserQuery = `
	SELECT id, username, pass_hash, role
	FROM users
	WHERE username = ?;
`

const insertUserQuery = `
	INSERT INTO users (username, pass_hash, role)
	VALUES (?, ?, ?)
`

func checkUser(username string) (user, error) {
	var u user
	err := database.DB.QueryRow(
		checkUserQuery,
		username,
	).Scan(
		&u.ID,
		&u.Username,
		&u.PassHash,
		&u.Role,
	)
	return u, err
}

func insertUser(username, passHash, role string) error {

	_, err := database.DB.Exec(
		insertUserQuery,
		username,
		passHash,
		role,
	)
	return err
}
