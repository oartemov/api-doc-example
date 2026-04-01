package models

import "time"

// User Модель пользователя
//
//	@Description Информация о пользователе системы
type User struct {
	// Уникальный идентификатор пользователя
	// @Example 1
	ID int `json:"id"`

	// Email адрес пользователя (уникальный)
	// @Example user@example.com
	Email string `json:"email"`

	// Имя пользователя
	// @Example Иван Петров
	Name string `json:"name"`

	// Дата и время создания аккаунта
	// @Example 2024-01-15T10:30:00Z
	CreatedAt time.Time `json:"created_at"`
}

// CreateUserRequest Запрос на создание пользователя
//
//	@Description Параметры для регистрации нового пользователя
type CreateUserRequest struct {
	// Email адрес (обязательное поле, должен быть валидным email)
	// @Example user@example.com
	Email string `json:"email" validate:"required,email"`

	// Имя пользователя (обязательное поле, макс. 100 символов)
	// @Example Иван Петров
	Name string `json:"name" validate:"required,max=100"`
}

// UpdateUserRequest Запрос на обновление пользователя
//
//	@Description Параметры для обновления данных пользователя (все поля опциональны)
type UpdateUserRequest struct {
	// Новый email адрес (опционально, должен быть валидным email при указании)
	// @Example newuser@example.com
	Email *string `json:"email" validate:"omitempty,email"`

	// Новое имя пользователя (опционально, макс. 100 символов при указании)
	// @Example Иван Петров
	Name *string `json:"name" validate:"omitempty,max=100"`
}

// ErrorResponse Ответ с ошибкой
//
//	@Description Структура для возврата информации об ошибке
type ErrorResponse struct {
	// Код ошибки
	// @Example bad_request
	Error string `json:"error"`

	// Сообщение об ошибке
	// @Example Неверный формат email
	Message string `json:"message"`
}

// SuccessResponse Успешный ответ
//
//	@Description Структура для возврата успешного ответа с данными
type SuccessResponse struct {
	// Флаг успешного выполнения
	// @Example true
	Success bool `json:"success"`

	// Данные ответа (опционально)
	Data interface{} `json:"data,omitempty"`
}
