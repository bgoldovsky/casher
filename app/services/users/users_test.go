package users

import (
	"errors"
	"testing"
	"time"

	"github.com/bgoldovsky/casher/app/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	user = models.User{
		Login:    "test-user",
		Password: "qwerty",
		Name:     "jon doe",
		Birth:    time.Now(),
	}

	operation = models.Operation{
		UserID:  123,
		Subject: "test-subj",
		Amount:  1000,
		Type:    1,
		Message: "test-msg",
	}
)

func TestService_Get_UsersError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usersRepo := NewMockusersRepository(ctrl)
	operationsRepo := NewMockoperationsRepository(ctrl)

	expErr := errors.New("test error")

	usersRepo.EXPECT().Get(gomock.Any()).Return(nil, expErr)

	service := New(usersRepo, operationsRepo)

	act, err := service.GetUser(user.ID)

	assert.Nil(t, act)
	assert.ErrorIs(t, err, expErr)
}

func TestService_Get_OperationsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usersRepo := NewMockusersRepository(ctrl)
	operationsRepo := NewMockoperationsRepository(ctrl)

	expErr := errors.New("test error")

	usersRepo.EXPECT().Get(gomock.Any()).Return(&user, nil)
	operationsRepo.EXPECT().Get(gomock.Any(), int64(0), int64(0)).Return(nil, expErr)

	service := New(usersRepo, operationsRepo)

	act, err := service.GetUser(user.ID)

	assert.Nil(t, act)
	assert.ErrorIs(t, err, expErr)
}

func TestService_Get_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usersRepo := NewMockusersRepository(ctrl)
	operationsRepo := NewMockoperationsRepository(ctrl)

	operations := models.OperationPaginator{
		Operations: []models.Operation{operation},
		HasMore:    false,
	}

	usersRepo.EXPECT().Get(gomock.Any()).Return(&user, nil)
	operationsRepo.EXPECT().Get(gomock.Any(), int64(0), int64(0)).Return(&operations, nil)

	service := New(usersRepo, operationsRepo)

	act, err := service.GetUser(user.ID)

	assert.NoError(t, err)
	assert.Equal(t, &user, act)
}

func TestService_Login_UsersError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usersRepo := NewMockusersRepository(ctrl)
	operationsRepo := NewMockoperationsRepository(ctrl)

	expErr := errors.New("test error")

	usersRepo.EXPECT().Auth(user.Login).Return(nil, expErr)

	service := New(usersRepo, operationsRepo)

	act, err := service.Auth(user.Login, user.Password)

	assert.Nil(t, act)
	assert.ErrorIs(t, err, ErrInvalidPassword)
}

func TestService_Create_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usersRepo := NewMockusersRepository(ctrl)
	operationsRepo := NewMockoperationsRepository(ctrl)

	expErr := errors.New("test error")

	usersRepo.EXPECT().Create(gomock.Any()).Return(int64(0), expErr)

	service := New(usersRepo, operationsRepo)

	act, err := service.Create(user.Login, user.Password, user.Name, user.Birth)

	assert.Empty(t, act)
	assert.ErrorIs(t, err, expErr)
}

func TestService_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usersRepo := NewMockusersRepository(ctrl)
	operationsRepo := NewMockoperationsRepository(ctrl)

	expID := int64(55)

	usersRepo.EXPECT().Create(gomock.Any()).Return(expID, nil)

	service := New(usersRepo, operationsRepo)

	act, err := service.Create(user.Login, user.Password, user.Name, user.Birth)

	assert.Equal(t, act, expID)
	assert.NoError(t, err)
}
