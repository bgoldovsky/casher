package operations

import (
	"errors"
	"testing"

	"github.com/bgoldovsky/casher/app/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	operation = models.Operation{
		UserID:  123,
		Subject: "test-subj",
		Amount:  1000,
		Type:    1,
		Message: "test-msg",
	}
)

func TestService_Create_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := NewMockrepository(ctrl)

	expErr := errors.New("test error")

	repo.EXPECT().Create(&operation).Return(expErr)

	service := New(repo)
	err := service.Create(operation.UserID, operation.Subject, operation.Amount, operation.Type, operation.Message)

	assert.ErrorIs(t, err, expErr)
}

func TestService_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := NewMockrepository(ctrl)

	repo.EXPECT().Create(&operation).Return(nil)

	service := New(repo)
	err := service.Create(operation.UserID, operation.Subject, operation.Amount, operation.Type, operation.Message)

	assert.NoError(t, err)
}

func TestService_Remove_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := NewMockrepository(ctrl)

	expErr := errors.New("test error")

	repo.EXPECT().Remove(gomock.Any()).Return(expErr)

	service := New(repo)
	err := service.Remove(operation.ID)

	assert.ErrorIs(t, err, expErr)
}

func TestService_Remove_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := NewMockrepository(ctrl)

	repo.EXPECT().Remove(gomock.Any()).Return(nil)

	service := New(repo)
	err := service.Remove(operation.ID)

	assert.NoError(t, err)
}

func TestService_Get_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := NewMockrepository(ctrl)

	expErr := errors.New("test error")

	repo.EXPECT().Get(operation.ID, int64(1), int64(5)).Return(nil, expErr)

	service := New(repo)
	paginator, err := service.Get(operation.ID, 1)

	assert.Nil(t, paginator)
	assert.ErrorIs(t, err, expErr)
}

func TestService_Get_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := NewMockrepository(ctrl)

	exp := &models.OperationPaginator{
		Operations: []models.Operation{operation},
		HasMore:    false,
	}

	repo.EXPECT().Get(operation.ID, int64(1), int64(5)).Return(exp, nil)

	service := New(repo)
	act, err := service.Get(operation.ID, 1)

	assert.Equal(t, exp, act)
	assert.NoError(t, err)
}
