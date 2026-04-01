package models

import "time"

// Todo Модель задачи
//
//	@Description Информация о задаче пользователя
type Todo struct {
	// Уникальный идентификатор задачи
	// @Example 1
	ID int `json:"id"`

	// Идентификатор пользователя-владельца задачи
	// @Example 1
	UserID int `json:"user_id"`

	// Заголовок задачи
	// @Example Купить продукты
	Title string `json:"title"`

	// Подробное описание задачи (опционально)
	// @Example Молоко, хлеб, яйца
	Description string `json:"description"`

	// Флаг выполнения задачи
	// @Example false
	Done bool `json:"done"`

	// Дата и время создания задачи
	// @Example 2024-01-15T10:30:00Z
	CreatedAt time.Time `json:"created_at"`

	// Дата и время последнего обновления задачи
	// @Example 2024-01-15T12:00:00Z
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateTodoRequest Запрос на создание задачи
//
//	@Description Параметры для создания новой задачи
type CreateTodoRequest struct {
	// Идентификатор пользователя-владельца (обязательное поле)
	// @Example 1
	UserID int `json:"user_id" validate:"required"`

	// Заголовок задачи (обязательное поле, макс. 200 символов)
	// @Example Купить продукты
	Title string `json:"title" validate:"required,max=200"`

	// Подробное описание задачи (опционально, макс. 1000 символов)
	// @Example Молоко, хлеб, яйца
	Description string `json:"description" validate:"max=1000"`

	// Флаг выполнения задачи (по умолчанию false)
	// @Example false
	Done bool `json:"done"`
}

// UpdateTodoRequest Запрос на обновление задачи
//
//	@Description Параметры для обновления задачи (все поля опциональны)
type UpdateTodoRequest struct {
	// Новый заголовок задачи (опционально, макс. 200 символов при указании)
	// @Example Обновить заголовок
	Title *string `json:"title" validate:"omitempty,max=200"`

	// Новое описание задачи (опционально, макс. 1000 символов при указании)
	// @Example Новое описание задачи
	Description *string `json:"description" validate:"omitempty,max=1000"`

	// Новый флаг выполнения (опционально)
	// @Example true
	Done *bool `json:"done"`
}
