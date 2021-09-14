package users

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/bgoldovsky/casher/app/models"
)

var (
	ErrDuplicateKey = errors.New("duplicate key value error")
)

type queryer interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type repository struct {
	db queryer
}

// New Инициализирует экземпляр репозитория
func New(db queryer) *repository {
	return &repository{db: db}
}

// Create Создает нового пользователя
func (store *repository) Create(user *models.User) (int64, error) {
	row := store.db.QueryRow(
		"insert into users(login, password, name, birth) values ($1,$2,$3,$4) returning id",
		user.Login,
		user.Password,
		user.Name,
		user.Birth,
	)

	var userID int64
	err := row.Scan(&userID)
	if isDuplicateErr(err) {
		return 0, ErrDuplicateKey
	}

	return userID, err
}

// Get Возвращает пользователя по его ID
func (store *repository) Get(userID int64) (*models.User, error) {
	query := "select id, login, password, name, birth, created_at from users where id=$1"

	row := store.db.QueryRow(query, userID)

	u := models.User{}
	if err := row.Scan(&u.ID, &u.Login, &u.Password, &u.Name, &u.Birth, &u.Created); err != nil {
		return nil, err
	}

	return &u, nil
}

// Auth Возвращает пользователя по его логину
func (store *repository) Auth(login string) (*models.User, error) {
	query := "select id, login, password, name, birth, created_at from users where login=$1"

	row := store.db.QueryRow(query, login)

	u := models.User{}
	if err := row.Scan(&u.ID, &u.Login, &u.Password, &u.Name, &u.Birth, &u.Created); err != nil {
		return nil, err
	}

	return &u, nil
}

// Проверяет, является ли ошибка ошибкой дупликации
func isDuplicateErr(err error) bool {
	if err == nil {
		return false
	}

	return strings.Contains(err.Error(), "duplicate key value violates unique constraint")
}
