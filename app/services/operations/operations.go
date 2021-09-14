//go:generate mockgen -source=operations.go -destination=./mocks.go -package=operations

package operations

import (
	"github.com/bgoldovsky/casher/app/logger"
	"github.com/bgoldovsky/casher/app/models"
)

const (
	pageSize = 5
)

type repository interface {
	Create(operation *models.Operation) error
	Remove(operationID int64) error
	Get(userID, page, size int64) (*models.OperationPaginator, error)
}

// Service Сервис управления финансовыми операциями
type Service struct {
	repo repository
}

// New Возвращает инициализированный экземпляр сервиса
func New(repo repository) *Service {
	return &Service{repo: repo}
}

// Get Возвращает список операций с пагинацией
// Если параметры пагинации не указаны, то вернет все операции
func (s *Service) Get(userID int64, page int64) (*models.OperationPaginator, error) {
	paginator, err := s.repo.Get(userID, page, pageSize)
	if err != nil {
		logger.Log.WithError(err).WithField("userID", userID).Errorf("get paginator error")
		return nil, err
	}

	return paginator, nil
}

// Create Создает новую операцию
func (s *Service) Create(userID int64, subject string, amount int64, operationType models.OperationType, msg string) error {
	operation := &models.Operation{
		UserID:  userID,
		Subject: subject,
		Amount:  amount,
		Type:    operationType,
		Message: msg,
	}

	err := s.repo.Create(operation)
	if err != nil {
		logger.Log.WithError(err).WithField("operation", operation).Errorf("create operations error")
		return err
	}

	return nil
}

// Remove Удаляет операцию
func (s *Service) Remove(operationID int64) error {
	err := s.repo.Remove(operationID)
	if err != nil {
		logger.Log.WithError(err).WithField("operationID", operationID).Errorf("remove operations error")
		return err
	}

	return nil
}
