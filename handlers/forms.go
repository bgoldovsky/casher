package handlers

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/bgoldovsky/casher/app/models"
)

type operationForm struct {
	Subject string
	Amount  float64
	Type    int64
	Message string
	Errors  map[string]string
}

// Validate Валидирует поля формы
// Ошибки валидации добавляются в специальное поле формы, что бы отрендерить их на страничке
func (f *operationForm) Validate() bool {
	f.Errors = map[string]string{}

	if strings.TrimSpace(f.Subject) == "" {
		f.Errors["Subject"] = "введите тему"
	}

	if f.Amount <= 0 {
		f.Errors["Amount"] = "введите сумму"
	}

	if f.Type != int64(models.Deposit) && f.Type != int64(models.Withdraw) {
		f.Errors["Type"] = "выберите тип операции"
	}

	return len(f.Errors) == 0
}

type authForm struct {
	Login    string
	Password string
	Errors   map[string]string
}

// Validate Валидирует поля формы
func (f *authForm) Validate() bool {
	f.Errors = map[string]string{}

	if strings.TrimSpace(f.Login) == "" {
		f.Errors["Login"] = "введите имя пользователя"
	}

	if strings.TrimSpace(f.Password) == "" {
		f.Errors["Password"] = "введите пароль"
	}

	return len(f.Errors) == 0
}

type registrationForm struct {
	Login           string
	Password        string
	ConfirmPassword string
	Name            string
	Birth           time.Time
	Errors          map[string]string
}

// Validate Валидирует поля формы
func (f *registrationForm) Validate() bool {
	f.Errors = map[string]string{}

	if strings.TrimSpace(f.Login) == "" {
		f.Errors["Login"] = "введите имя пользователя"
	}

	if strings.TrimSpace(f.Password) == "" {
		f.Errors["Password"] = "введите пароль"
	}

	if strings.TrimSpace(f.ConfirmPassword) == "" {
		f.Errors["ConfirmPassword"] = "введите пароль еще раз"
	}

	hasSevenOrMore, hasNumber, hasUpper, hasSpecial := verifyPassword(f.Password)

	fmt.Println("VERIFY:", hasSevenOrMore, hasNumber, hasUpper, hasSpecial)

	if !hasSevenOrMore || !hasUpper || !hasNumber || !hasSpecial {
		f.Errors["Password"] = "введенный пароль не надежен"
	}

	if strings.TrimSpace(f.Name) == "" {
		f.Errors["Name"] = "введите настоящее имя"
	}

	return len(f.Errors) == 0
}

// Проверяет требования сложности пароля
// 1. Не менее 7 символов
// 2. Не менее 1 числа
// 2. Не менее 1 заглавной буквы
// 3. Не менее 1 специального символа
func verifyPassword(s string) (sevenOrMore, number, upper, special bool) {
	letters := 0
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
			letters++
		case unicode.IsUpper(c):
			upper = true
			letters++
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
			letters++
		case unicode.IsLetter(c) || c == ' ':
			letters++
		}
	}
	sevenOrMore = letters >= 7
	return
}
