//go:generate mockgen -source=users.go -destination=./mocks.go -package=users

package users

import (
	"errors"
	"time"

	"github.com/bgoldovsky/casher/app/logger"
	"github.com/bgoldovsky/casher/app/models"
	"github.com/bgoldovsky/casher/app/repositories/users"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidPassword = errors.New("invalid user or password error")
	ErrLoginExists     = errors.New("login already exists")
)

type usersRepository interface {
	Create(user *models.User) (int64, error)
	Get(userID int64) (*models.User, error)
	Auth(login string) (*models.User, error)
}

type operationsRepository interface {
	Get(userID, page, size int64) (*models.OperationPaginator, error)
}

// Service Сервис управления пользователями
type Service struct {
	usersRepo      usersRepository
	operationsRepo operationsRepository
}

// New Возвращает инициализированный экземпляр сервиса
func New(usersRepo usersRepository, operationRepo operationsRepository) *Service {
	return &Service{
		usersRepo:      usersRepo,
		operationsRepo: operationRepo,
	}
}

// GetUser Возвращает пользователя по его идентификатору
func (s *Service) GetUser(userID int64) (*models.User, error) {
	user, err := s.usersRepo.Get(userID)
	if err != nil {
		logger.Log.WithError(err).WithField("userID", userID).Errorf("get user error")
		return nil, err
	}

	balance, err := s.getBalance(userID)
	if err != nil {
		return nil, err
	}

	user.Balance = balance
	return user, nil
}

// Auth Возвращает пользователя по его логину и паролю
func (s *Service) Auth(login, password string) (*models.User, error) {
	user, err := s.usersRepo.Auth(login)
	if err != nil {
		logger.Log.WithError(err).WithField("login", login).Errorf("get user by login error")
		return nil, ErrInvalidPassword
	}

	equals := checkPasswordHash(password, user.Password)
	if !equals {
		logger.Log.WithField("user", user).Errorf("invalid password error")
		return nil, ErrInvalidPassword
	}

	balance, err := s.getBalance(user.ID)
	if err != nil {
		return nil, err
	}

	user.Balance = balance
	return user, nil
}

// Create Создает нового пользователя
func (s *Service) Create(login, password, name string, birth time.Time) (int64, error) {
	// В базу сохраняется хеш пароля
	hashedPassword, err := hashPassword(password)
	if err != nil {
		logger.Log.WithError(err).Errorf("hash password error")
		return 0, err
	}

	user := &models.User{
		Login:    login,
		Password: hashedPassword,
		Name:     name,
		Birth:    birth,
	}

	userID, err := s.usersRepo.Create(user)
	if err == users.ErrDuplicateKey {
		logger.Log.WithError(err).Errorf("create user error: login already exists")
		return 0, ErrLoginExists
	}
	if err != nil {
		logger.Log.WithError(err).Errorf("create user error")
		return 0, err
	}

	return userID, nil
}

// Считает баланс указанного пользователя
func (s *Service) getBalance(userID int64) (int64, error) {
	paginator, err := s.operationsRepo.Get(userID, 0, 0)
	if err != nil {
		logger.Log.WithError(err).WithField("userID", userID).Errorf("get balance error")
		return 0, err
	}

	if len(paginator.Operations) == 0 {
		return 0, nil
	}

	balance := int64(0)
	for _, o := range paginator.Operations {
		if o.Type == models.Deposit {
			balance += o.Amount
		} else {
			balance -= o.Amount
		}
	}

	return balance, nil
}

// Берет хеш от пароля
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Сравнивает пароль с его хешем
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
