package operations

import (
	"database/sql"
	"fmt"

	"github.com/bgoldovsky/casher/app/models"
)

type queryer interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

type repository struct {
	db queryer
}

// New Инициализирует экземпляр репозитория
func New(db queryer) *repository {
	return &repository{db: db}
}

// Create Создает новую операцию
func (store *repository) Create(o *models.Operation) error {
	_, err := store.db.Query(
		"insert into operations(user_id, subject, amount, type, message) values ($1,$2,$3,$4,$5)",
		o.UserID,
		o.Subject,
		o.Amount,
		o.Type,
		o.Message,
	)

	return err
}

// Remove Удаляет указанную операцию по ее ID
func (store *repository) Remove(operationID int64) error {
	_, err := store.db.Query("delete from operations where id = $1", operationID)

	return err
}

// Get Возвращает список операций
func (store *repository) Get(userID, page, size int64) (*models.OperationPaginator, error) {
	query := "select id, user_id, subject, amount, type, message, created_at from operations where user_id=$1 order by created_at desc"
	query = addPagination(query, page, size)

	rows, err := store.db.Query(query, userID)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var operations []models.Operation
	for rows.Next() {
		o := models.Operation{}
		if err := rows.Scan(&o.ID, &o.UserID, &o.Subject, &o.Amount, &o.Type, &o.Message, &o.Created); err != nil {
			return nil, err
		}

		operations = append(operations, o)
	}

	// Если пагинация не нужна или количество объектов меньше размера страницы возвращаем все
	if !needPagination(page, size) || len(operations) <= int(size) {
		return &models.OperationPaginator{
			Operations: operations,
		}, nil
	}

	return &models.OperationPaginator{
		Operations: operations[:size],
		HasMore:    true,
	}, nil
}

// Добавляет к строке SQL запроса данные пагинации
func addPagination(query string, page, size int64) string {
	if !needPagination(page, size) {
		return query
	}

	// Запрашиваем на 1 объект больше, что бы проверить есть ли еще данные не запрашивая дополнительное count
	limit, offset := size+1, (page-1)*size
	return fmt.Sprintf("%s limit %d offset %d", query, limit, offset)
}

// Определяет нужна ли в запросе пагинация
func needPagination(page, size int64) bool {
	return page != 0 && size != 0
}
