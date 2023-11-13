package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID             int        `json:"id"`
	UserName       string     `json:"username"`
	Email          string     `json:"email"`
	Password       string     `json:"password"`
	Patches        []Patch    `json:"patches"`
	CreatedAt      *time.Time `json:"createdAt"`
	LastLoggedInAt *time.Time `json:"lastLoggedInAt"`
}

type UserModel struct {
	DB *sql.DB
}

const getUserByIdQuery string = `
	SELECT *
	FROM users
	JOIN patches ON patches.user_id = users.id
	WHERE id = $1
`

const checkIfUserExistsQuery string = `
	SELECT *
	FROM users
	JOIN patches ON patches.user_id = users.id
	WHERE email = $1 OR username = $2
`

const isEmailValidQuery string = `
	SELECT *
	FROM users
	WHERE email = $1
`

const isUsernameValidQuery string = `
	SELECT *
	FROM users
	WHERE email = $1
`

const createUserQuery string = `
	INSERT INTO users (username, email, password, created_at) 
	VALUES ($1, $2, $3, NOW())
`

const setLastLoginQuery string = `
	UPDATE users 
	SET last_login = NOW()
	WHERE id = $1
`

const setSessionQuery string = `
	UPDATE users
	SET session_token = $2
	WHERE id = $1
`

func (m *UserModel) GetUserById(id int) (User, error) {
	var user User
	row := m.DB.QueryRow(getUserByIdQuery, id)
	err := row.Scan(&user.ID, &user.UserName, &user.Email)
	return user, err
}

func (m *UserModel) GetUserByEmail(email string) (User, error) {
	var user User
	row := m.DB.QueryRow(getUserByIdQuery, email)
	err := row.Scan(&user.ID, &user.UserName, &user.Email)
	return user, err
}

func (m *UserModel) CheckIfUserExists(email string, username string) (bool, error) {
	var user User
	row := m.DB.QueryRow(checkIfUserExistsQuery, email, username)
	err := row.Scan(&user.ID, &user.UserName, &user.Email)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (m *UserModel) IsEmailValid(email string) (bool, error) {
	var user User
	row := m.DB.QueryRow(isEmailValidQuery, email)
	err := row.Scan(&user.ID, &user.UserName, &user.Email)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (m *UserModel) IsUsernameValid(username string) (bool, error) {
	var user User
	row := m.DB.QueryRow(isUsernameValidQuery, username)
	err := row.Scan(&user.ID, &user.UserName, &user.Email)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (m *UserModel) CreateUser(username string, email string, password string) (User, error) {
	var user User
	row := m.DB.QueryRow(createUserQuery, username, email, password)
	err := row.Scan(&user.ID, &user.UserName, &user.Email)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (m *UserModel) SetLastLogin(id int) error {
	_, err := m.DB.Exec(setLastLoginQuery, id)
	return err
}

func (m *UserModel) SetUserSession(id int, sessionToken string) error {
	_, err := m.DB.Exec(setSessionQuery, id, sessionToken)
	return err
}
