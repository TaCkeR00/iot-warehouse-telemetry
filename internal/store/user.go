package store

import (
	"iot-warehouse-telemetry/pkg/user"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       uint   `json:"id"`
	UserName string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	IsAdmin  bool   `json:"is_admin"`
}

func (s *Store) IsUserExistsAndIsAdmin(user user.User) (bool, bool) {
	stmt := "SELECT id, username, password, email, is_admin    FROM users WHERE username=?"
	var Tuser User

	row := s.DB.QueryRow(stmt, user.UserName)
	err := row.Scan(&Tuser.ID, &Tuser.UserName, &Tuser.Password, &Tuser.Email, &Tuser.IsAdmin)
	if err != nil {
		return false, false
	}

	err = bcrypt.CompareHashAndPassword([]byte(Tuser.Password), []byte(user.Password))
	if err != nil {
		return false, false
	}
	return true, Tuser.IsAdmin
}

func (s *Store) SignUser(user user.User) error {
	stmt := `INSERT INTO users(username, password, email, is_admin) VALUES(?, ?, ?, ?)`

	HashPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = s.DB.Exec(stmt, user.UserName, HashPass, user.Email, user.IsAdmin)
	if err != nil {
		return err
	}

	return nil
}
