package models

import "time"

const (
	Deposit  OperationType = 1
	Withdraw OperationType = 2
)

// OperationType Тип финансовой операции
type OperationType int64

// Operation Модель финансовой операции
type Operation struct {
	ID      int64
	UserID  int64
	Subject string
	Amount  int64
	Type    OperationType
	Message string
	Created time.Time
}

// OperationPaginator Обертка для пагинации данных о финансовых операциях
type OperationPaginator struct {
	Operations []Operation
	HasMore    bool
}
