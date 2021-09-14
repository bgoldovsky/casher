package handlers

import (
	"time"

	"github.com/bgoldovsky/casher/app/models"
)

/*
	Для использования моделей в html/template
	Поля должны быть экспортируемыми
	Сама структура может быть экспортируемой или не экспортируемой
*/

type user struct {
	ID      int64
	Login   string
	Name    string
	Age     uint16
	Balance float64
}

// Конвертирует модель пользователя во view model
func userToView(model *models.User) *user {
	if model == nil {
		return nil
	}

	return &user{
		ID:      model.ID,
		Login:   model.Login,
		Name:    model.Name,
		Age:     getAge(model.Birth, time.Now()),
		Balance: float64(model.Balance) / 100,
	}
}

// Рассчитывает возраст по дате рождения и текущей дате
func getAge(birthdate, today time.Time) uint16 {
	today = today.In(birthdate.Location())
	ty, tm, td := today.Date()
	today = time.Date(ty, tm, td, 0, 0, 0, 0, time.UTC)
	by, bm, bd := birthdate.Date()
	birthdate = time.Date(by, bm, bd, 0, 0, 0, 0, time.UTC)
	if today.Before(birthdate) {
		return 0
	}

	age := ty - by
	anniversary := birthdate.AddDate(age, 0, 0)
	if anniversary.After(today) {
		age--
	}

	return uint16(age)
}

type operation struct {
	ID      int64
	UserID  int64
	Subject string
	Amount  float64
	Type    string
	Message string
	Created time.Time
}

type pagingOperations struct {
	Page       int64
	HasPrev    bool
	HasNext    bool
	Operations []operation
}

// Конвертирует модель операции во view model
func operationToView(model *models.Operation) *operation {
	if model == nil {
		return nil
	}

	return &operation{
		ID:      model.ID,
		UserID:  model.UserID,
		Subject: model.Subject,
		Amount:  float64(model.Amount) / 100,
		Type:    getOperationType(model.Type),
		Message: model.Message,
		Created: model.Created,
	}
}

// Конвертирует тип операции в строку
func getOperationType(model models.OperationType) string {
	if model == models.Deposit {
		return "Пополнение"
	} else if model == models.Withdraw {
		return "Списание"
	}

	return ""
}

// Конвертирует массив моделей операций во view model
func operationsToView(models []models.Operation) []operation {
	res := make([]operation, len(models))

	for idx, val := range models {
		view := operationToView(&val)
		res[idx] = *view
	}

	return res
}

// Конвертирует модель-обертку для операций во view model
func toPagingView(page int64, paginator *models.OperationPaginator) pagingOperations {
	paging := pagingOperations{
		Page:       page,
		HasNext:    paginator.HasMore,
		Operations: operationsToView(paginator.Operations),
	}

	if page <= 1 {
		paging.HasPrev = false
	} else {
		paging.HasPrev = true
	}

	return paging
}
