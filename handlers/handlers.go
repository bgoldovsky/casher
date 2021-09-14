package handlers

import (
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/bgoldovsky/casher/app/logger"
	"github.com/bgoldovsky/casher/app/models"
	"github.com/bgoldovsky/casher/app/services/operations"
	"github.com/bgoldovsky/casher/app/services/users"
	"github.com/bgoldovsky/casher/middleware"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

const (
	userIDKey = "user-id"
)

type PageHandler struct {
	usersSrv      *users.Service
	operationsSrv *operations.Service
	router        *mux.Router
	store         *sessions.CookieStore
}

func New(usersSrv *users.Service, operationsSrv *operations.Service) *PageHandler {
	// Создаем фейковый ключ для хранилища куки
	key := []byte("33446a9dcf9ea060a0a6532b166da32f304af0de")

	handler := &PageHandler{
		usersSrv:      usersSrv,
		operationsSrv: operationsSrv,
		store:         sessions.NewCookieStore(key),
	}

	// Инициализируем и настраиваем роутер
	r := mux.NewRouter()

	// Добавляем доступ к статическим файлам
	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// Настраиваем роутер
	// Роуты авторизации
	r.HandleFunc("/", middleware.Logging(handler.Index)).Methods("GET", "POST")
	r.HandleFunc("/auth/", middleware.Logging(handler.Auth)).Methods("GET", "POST")
	r.HandleFunc("/logout/", middleware.Logging(handler.Logout)).Methods("GET", "POST")
	r.HandleFunc("/registration/", middleware.Logging(handler.Registration)).Methods("GET", "POST")
	// Роуты для работы с операциями
	r.HandleFunc("/operations/", middleware.Logging(handler.Operations)).Methods("GET", "POST")
	r.HandleFunc("/operations/create/", middleware.Logging(handler.Create)).Methods("GET", "POST")
	r.HandleFunc("/operations/delete/{id:[0-9]+}", middleware.Logging(handler.Delete)).Methods("POST")
	// Роуты для обработки ошибок
	r.HandleFunc("/error/", middleware.Logging(handler.Error)).Methods("GET", "POST")
	r.HandleFunc("/error/unauthorized", middleware.Logging(handler.ErrorUnauthorized)).Methods("GET", "POST")

	handler.router = r
	return handler
}

// Router Возвращает сконфигурированный роутер обработчика
func (h *PageHandler) Router() *mux.Router {
	return h.router
}

// Operation handlers

// Index Обработчик главной страницы
func (h *PageHandler) Index(w http.ResponseWriter, r *http.Request) {
	// Проверяем аутентификацию
	userID, isAuth := h.getAuthorizedUserID(r)
	if !isAuth {
		logger.Log.Error("index page handler error: access forbidden")
		http.Redirect(w, r, "/auth/", http.StatusTemporaryRedirect)
		return
	}

	// Получаем текущего пользователя
	u, err := h.usersSrv.GetUser(userID)
	if err != nil {
		logger.Log.WithError(err).Error("index handler error")
		http.Redirect(w, r, "/error/", http.StatusTemporaryRedirect)
		return
	}

	// Парсим шаблон
	// Если хотим паниковать при ошибках парсинга шаблона используем wrapper-функцию template.Must
	tmpl := template.Must(template.ParseFiles(
		"templates/index.html",
		"templates/header.html",
		"templates/footer.html",
	))

	// Рендерим шаблон
	err = tmpl.ExecuteTemplate(w, "index", userToView(u))
	if err != nil {
		logger.Log.WithError(err).Error("index handler error")
		http.Redirect(w, r, "/error/", http.StatusTemporaryRedirect)
		return
	}
}

// Operations Обработчик страницы отображения операций
func (h *PageHandler) Operations(w http.ResponseWriter, r *http.Request) {
	// Проверяем аутентификацию
	userID, isAuth := h.getAuthorizedUserID(r)
	if !isAuth {
		logger.Log.Error("access forbidden error")
		http.Redirect(w, r, "/auth/", http.StatusTemporaryRedirect)
		return
	}

	// Определяем функции для пагинации
	// После парсинга ими можно будет пользоваться внутри шаблона
	funcMap := template.FuncMap{
		"inc": func(i int64) int64 {
			return i + 1
		},
		"dec": func(i int64) int64 {
			return i - 1
		},
	}

	// Парсим шаблон вместе с функциями
	tmpl := template.Must(template.New("wrapper").Funcs(funcMap).ParseFiles(
		"templates/operations.html",
		"templates/header.html",
		"templates/footer.html",
	))

	// Получаем страницу пагинации из запроса
	var nextPage int64 = 1
	var err error
	pageStr := r.URL.Query().Get("page")
	if pageStr != "" {
		nextPage, err = strconv.ParseInt(pageStr, 10, 0)
		if err != nil {
			logger.Log.WithError(err).Error("operations handler error")
			http.Redirect(w, r, "/error/", http.StatusTemporaryRedirect)
			return
		}
	}

	// Получаем список операций
	paginator, err := h.operationsSrv.Get(userID, nextPage)
	if err != nil {
		logger.Log.WithError(err).Error("operations handler error")
		http.Redirect(w, r, "/error/", http.StatusTemporaryRedirect)
		return
	}

	// Рендерим ответ
	err = tmpl.ExecuteTemplate(w, "operations", toPagingView(nextPage, paginator))
	if err != nil {
		logger.Log.WithError(err).Error("operations handler error")
		http.Redirect(w, r, "/error/", http.StatusTemporaryRedirect)
		return
	}
}

// Create Обработчик страницы создания операции
func (h *PageHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Проверяем аутентификацию
	userID, isAuth := h.getAuthorizedUserID(r)
	if !isAuth {
		logger.Log.Error("create handler error: access forbidden")
		http.Redirect(w, r, "/auth/", http.StatusTemporaryRedirect)
		return
	}

	// Парсим шаблон
	tmpl := template.Must(template.ParseFiles(
		"templates/create.html",
		"templates/header.html",
		"templates/footer.html",
	))

	// Если пришел GET запрос, только рендерим шаблон
	if r.Method != http.MethodPost {
		err := tmpl.ExecuteTemplate(w, "create", operationForm{})
		if err != nil {
			logger.Log.WithError(err).Error("create handler error")
			http.Redirect(w, r, "/error/", http.StatusTemporaryRedirect)
			return
		}
		return
	}

	// Если пришел POST запрос, то обрабатываем пришедшую форму
	// Извлекаем данные формы
	amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil {
		logger.Log.WithError(err).Error("create handler error")
		http.Redirect(w, r, "/error/", http.StatusTemporaryRedirect)
		return
	}

	operationType, err := strconv.ParseInt(r.FormValue("type"), 10, 0)
	if err != nil {
		logger.Log.WithError(err).Error("create handler error")
		http.Redirect(w, r, "/error/", http.StatusTemporaryRedirect)
		return
	}

	form := operationForm{
		Subject: r.FormValue("subject"),
		Amount:  amount,
		Type:    operationType,
		Message: r.FormValue("message"),
	}

	// Валидируем данные формы
	if !form.Validate() {
		err := tmpl.ExecuteTemplate(w, "create", form)
		if err != nil {
			logger.Log.WithError(err).WithField("form", form).Error("create handler error")
			http.Redirect(w, r, "/error/", http.StatusTemporaryRedirect)
			return
		}
		return
	}

	// Сохраняем операцию в БД
	err = h.operationsSrv.Create(
		userID,
		form.Subject,
		int64(form.Amount*100),
		models.OperationType(form.Type),
		form.Message,
	)
	if err != nil {
		logger.Log.WithError(err).Error("create handler error")
		http.Redirect(w, r, "/error/", http.StatusTemporaryRedirect)
		return
	}

	// Редиректим на список операций
	http.Redirect(w, r, "/operations/", http.StatusTemporaryRedirect)
}

// Delete Обработчик страницы удаления операции
func (h *PageHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	operationIDStr, ok := vars["id"]
	if !ok {
		logger.Log.Error("delete handler error: operation ID not specified")
		http.Redirect(w, r, "/error/", http.StatusTemporaryRedirect)
		return
	}

	operationID, err := strconv.ParseInt(operationIDStr, 10, 0)
	if err != nil {
		logger.Log.WithError(err).Error("delete handler error")
		http.Redirect(w, r, "/error/", http.StatusTemporaryRedirect)
		return
	}

	err = h.operationsSrv.Remove(operationID)
	if err != nil {
		logger.Log.WithError(err).Error("create handler error")
		http.Redirect(w, r, "/error/", http.StatusTemporaryRedirect)
		return
	}

	http.Redirect(w, r, "/operations/", http.StatusTemporaryRedirect)
}

// Auth handlers

// Auth Обработчик страницы авторизации пользователя
func (h *PageHandler) Auth(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"templates/auth.html",
		"templates/header_unauthorized.html",
		"templates/footer.html",
	))

	form := authForm{}

	// Если пришел GET запрос, только рендерим шаблон и выходим
	if r.Method != http.MethodPost {
		err := tmpl.ExecuteTemplate(w, "auth", form)
		if err != nil {
			logger.Log.WithError(err).Error("auth handler error")
			http.Redirect(w, r, "/error/", http.StatusTemporaryRedirect)
			return
		}
		return
	}

	form.Login = r.FormValue("login")
	form.Password = r.FormValue("password")

	// Валидируем данные формы
	if !form.Validate() {
		err := tmpl.ExecuteTemplate(w, "auth", form)
		if err != nil {
			logger.Log.WithError(err).WithField("form", form).Error("auth handler error")
			http.Redirect(w, r, "/error/", http.StatusTemporaryRedirect)
			return
		}
		return
	}

	// Получаем пользователя по логину и паролю
	u, err := h.usersSrv.Auth(form.Login, form.Password)
	// Если пользователь не найден или пароль не валидирован, то отправляем сообщение пользователю
	if err == users.ErrInvalidPassword {
		form.Errors["Password"] = "Неверное имя пользователя или пароль"
		err := tmpl.ExecuteTemplate(w, "auth", form)
		if err != nil {
			logger.Log.WithError(err).WithField("form", form).Error("registration handler error")
			http.Redirect(w, r, "/error/", http.StatusTemporaryRedirect)
			return
		}
		return
	}
	// Иначе рендерим страницу с ошибкой
	if err != nil {
		logger.Log.WithError(err).WithField("form", form).Error("registration handler error")
		http.Redirect(w, r, "/error/", http.StatusTemporaryRedirect)
		return
	}

	// Авторизуем пользователя
	err = h.authorizeUser(u.ID, w, r)
	if err != nil {
		http.Redirect(w, r, "/error/", http.StatusTemporaryRedirect)
		logger.Log.WithError(err).WithField("form", form).Error("auth handler error")
		return
	}

	// Редиректим пользователя на главную страницу
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

// Registration Обработчик страницы регистрации нового пользователя
func (h *PageHandler) Registration(w http.ResponseWriter, r *http.Request) {
	// Создаем регистрационную форму
	local, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		logger.Log.WithError(err).Error("registration handler error")
		http.Redirect(w, r, "/error/", http.StatusTemporaryRedirect)
		return
	}
	birth := time.Date(1986, 4, 15, 0, 0, 0, 0, local)
	form := registrationForm{
		Birth: birth,
	}

	// Парсим шаблон страницы
	tmpl := template.Must(template.ParseFiles(
		"templates/registration.html",
		"templates/header_unauthorized.html",
		"templates/footer.html",
	))

	// Если пришел GET запрос, только рендерим форму и выходим
	if r.Method != http.MethodPost {
		err = tmpl.ExecuteTemplate(w, "registration", form)
		if err != nil {
			logger.Log.WithError(err).Error("registration handler error")
			http.Redirect(w, r, "/error/", http.StatusTemporaryRedirect)
			return
		}
		return
	}

	// Если пришел POST запрос, то обрабатываем форму
	// Загружаем данные из формы в объект
	form.Login = r.FormValue("login")
	form.Password = r.FormValue("password")
	form.ConfirmPassword = r.FormValue("confirm-password")
	form.Name = r.FormValue("name")
	// Отдельно парсим и обрабатываем дату рождения
	birth, err = time.Parse("2006-01-02", r.FormValue("birth"))
	if err != nil {
		logger.Log.WithError(err).WithField("form", form).Error("registration handler error")
		http.Redirect(w, r, "/error/", http.StatusTemporaryRedirect)
		return
	}
	form.Birth = birth

	// Валидируем данные формы
	if !form.Validate() {
		err := tmpl.ExecuteTemplate(w, "registration", form)
		if err != nil {
			logger.Log.WithError(err).WithField("form", form).Error("registration handler error")
			http.Redirect(w, r, "/error/", http.StatusTemporaryRedirect)
			return
		}
		return
	}

	// Создаем пользователя
	userID, err := h.usersSrv.Create(form.Login, form.Password, form.Name, form.Birth)
	// Если пользователь с таким логином уже существует сообщаем об этом
	if err == users.ErrLoginExists {
		form.Errors["Auth"] = "Пользователь с таким именем уже существует"
		err = tmpl.ExecuteTemplate(w, "registration", form)
		if err != nil {
			logger.Log.WithError(err).WithField("form", form).Error("registration handler error")
			http.Redirect(w, r, "/error/", http.StatusTemporaryRedirect)
			return
		}
		return
	}
	// Иначе просто переходим на страницу с ошибкой
	if err != nil {
		logger.Log.WithError(err).WithField("form", form).Error("registration handler error")
		http.Redirect(w, r, "/error/", http.StatusTemporaryRedirect)
		return
	}

	// Авторизуем пользователя
	if err = h.authorizeUser(userID, w, r); err != nil {
		http.Redirect(w, r, "/error/", http.StatusTemporaryRedirect)
		logger.Log.WithError(err).WithField("form", form).Error("registration handler error")
		return
	}

	// Перекидываем пользователя на главную страницу
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

// Logout Обработчик нажатия кнопки выхода из системы
func (h *PageHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Сбрасываем сессию пользователя
	err := h.logoutUser(w, r)
	if err != nil {
		logger.Log.WithError(err).Error("logout handler error")
		http.Redirect(w, r, "/error/", http.StatusTemporaryRedirect)
		return
	}

	// Переходим на страницу авторизации
	http.Redirect(w, r, "/auth/", http.StatusTemporaryRedirect)
}

// Error handlers

// Error Обработчик страницы ошибки для авторизованного пользователя
func (h *PageHandler) Error(w http.ResponseWriter, r *http.Request) {
	// Проверяем аутентификацию
	// Если пользователь не авторизован то перекидываем его на страницу ошибки для неавторизованных пользователей
	if _, isAuth := h.getAuthorizedUserID(r); !isAuth {
		logger.Log.Error("error page handler error: access forbidden")
		http.Redirect(w, r, "/error/unauthorized", http.StatusTemporaryRedirect)
		return
	}

	// Парсим шаблон
	tmpl := template.Must(template.ParseFiles(
		"templates/error.html",
		"templates/header.html",
		"templates/footer.html",
	))

	// Рендерим шаблон
	err := tmpl.ExecuteTemplate(w, "error", nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("<h1>Internal server error</h1>"))
		return
	}
}

// ErrorUnauthorized Обработчик страницы ошибки для неавторизованного пользователя
func (h *PageHandler) ErrorUnauthorized(w http.ResponseWriter, _ *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"templates/error_unauthorized.html",
		"templates/header_unauthorized.html",
		"templates/footer.html",
	))

	err := tmpl.ExecuteTemplate(w, "errorUnauthorized", nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("<h1>Internal server error</h1>"))
		return
	}
}

// Common methods

func (h *PageHandler) getAuthorizedUserID(r *http.Request) (int64, bool) {
	// Получаем сессию
	session, err := h.store.Get(r, "cookie-name")
	if err != nil {
		return 0, false
	}

	// Читаем из сессии ID пользователя
	userID, ok := session.Values[userIDKey].(int64)
	if !ok {
		return 0, false
	}

	return userID, true
}

func (h *PageHandler) authorizeUser(userID int64, w http.ResponseWriter, r *http.Request) error {
	// Получаем сессию
	session, err := h.store.Get(r, "cookie-name")
	if err != nil {
		return err
	}

	// Сохраняем в сессию ID пользователя
	session.Values[userIDKey] = userID
	err = session.Save(r, w)
	if err != nil {
		return err
	}

	return nil
}

func (h *PageHandler) logoutUser(w http.ResponseWriter, r *http.Request) error {
	// Получаем сессию
	session, err := h.store.Get(r, "cookie-name")
	if err != nil {
		return err
	}

	// Сохраняем в сессию ID пользователя
	session.Values[userIDKey] = nil
	err = session.Save(r, w)
	if err != nil {
		return err
	}

	return nil
}
