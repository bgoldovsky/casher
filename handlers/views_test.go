package handlers

import (
	"testing"
	"time"

	"github.com/bgoldovsky/casher/app/models"
	"github.com/stretchr/testify/assert"
)

func Test_UserToView(t *testing.T) {
	model := models.User{
		ID:      123,
		Login:   "test-login",
		Name:    "test-name",
		Balance: 1000,
	}

	act := userToView(&model)

	assert.Equal(t, model.ID, act.ID)
	assert.Equal(t, model.Login, act.Login)
	assert.Equal(t, model.Name, act.Name)
	assert.Equal(t, float64(model.Balance)/100, act.Balance)
}

func Test_OperationToView(t *testing.T) {
	model := models.Operation{
		ID:      1,
		UserID:  2,
		Subject: "test-subj",
		Amount:  1000,
		Type:    1,
		Message: "test-msg",
		Created: time.Now(),
	}

	act := operationToView(&model)

	assert.Equal(t, model.ID, act.ID)
	assert.Equal(t, model.UserID, act.UserID)
	assert.Equal(t, model.Subject, act.Subject)
	assert.Equal(t, float64(model.Amount)/100, act.Amount)
	assert.Equal(t, "Пополнение", act.Type)
	assert.Equal(t, model.Message, act.Message)
	assert.Equal(t, model.Created, act.Created)
}
