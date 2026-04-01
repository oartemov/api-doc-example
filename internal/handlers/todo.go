package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"api-doc-example/internal/models"
)

type TodoHandler struct {
	todos  map[int]models.Todo
	nextID int
}

func NewTodoHandler() *TodoHandler {
	return &TodoHandler{
		todos:  make(map[int]models.Todo),
		nextID: 1,
	}
}

// ListTodos возвращает список всех задач
// @Summary Получить список всех задач
// @Description Возвращает массив объектов Todo в структуре SuccessResponse
// @Tags todos
// @Accept json
// @Produce json
// @Success 200 {object} models.SuccessResponse{data=[]models.Todo}
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/v1/todos [get]
func (h *TodoHandler) ListTodos(w http.ResponseWriter, r *http.Request) {
	var todos []models.Todo
	for _, todo := range h.todos {
		todos = append(todos, todo)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(models.SuccessResponse{
		Success: true,
		Data:    todos,
	})
}

// GetTodo возвращает задачу по ID
// @Summary Получить задачу по ID
// @Description Возвращает данные задачи по указанному идентификатору
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "ID задачи"
// @Success 200 {object} models.SuccessResponse{data=models.Todo}
// @Failure 400 {object} models.ErrorResponse "Invalid ID"
// @Failure 404 {object} models.ErrorResponse "Todo not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/v1/todos/{id} [get]
func (h *TodoHandler) GetTodo(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	todo, exists := h.todos[id]
	if !exists {
		http.Error(w, "Задача не найдена", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(models.SuccessResponse{
		Success: true,
		Data:    todo,
	})
}

// CreateTodo создает новую задачу
// @Summary Создать новую задачу
// @Description Создает новую задачу с указанным заголовком, описанием и статусом выполнения. Возвращает созданную задачу с присвоенным ID
// @Tags todos
// @Accept json
// @Produce json
// @Param request body models.CreateTodoRequest true "Данные для создания задачи"
// @Success 201 {object} models.SuccessResponse{data=models.Todo}
// @Failure 400 {object} models.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/v1/todos [post]
func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Невалидный запрос", http.StatusBadRequest)
		return
	}

	todo := models.Todo{
		ID:          h.nextID,
		UserID:      req.UserID,
		Title:       req.Title,
		Description: req.Description,
		Done:        req.Done,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	h.todos[h.nextID] = todo
	h.nextID++

	log.Printf("Создана новая задача: %+v", todo)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(models.SuccessResponse{
		Success: true,
		Data:    todo,
	})
}

// UpdateTodo обновляет задачу
// @Summary Обновить задачу
// @Description Обновляет заголовок, описание и/или статус выполнения задачи по указанному ID. Возвращает обновленные данные задачи
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "ID задачи"
// @Param request body models.UpdateTodoRequest true "Данные для обновления (все поля опциональны)"
// @Success 200 {object} models.SuccessResponse{data=models.Todo}
// @Failure 400 {object} models.ErrorResponse "Invalid ID or request body"
// @Failure 404 {object} models.ErrorResponse "Todo not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/v1/todos/{id} [put]
func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Невалидный запрос", http.StatusBadRequest)
		return
	}

	todo, exists := h.todos[id]
	if !exists {
		http.Error(w, "Задача не найдена", http.StatusNotFound)
		return
	}

	if req.Title != nil {
		todo.Title = *req.Title
	}
	if req.Description != nil {
		todo.Description = *req.Description
	}
	if req.Done != nil {
		todo.Done = *req.Done
	}
	todo.UpdatedAt = time.Now()

	h.todos[id] = todo

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(models.SuccessResponse{
		Success: true,
		Data:    todo,
	})
}

// DeleteTodo удаляет задачу по ID
// @Summary Удалить задачу
// @Description Удаляет задачу по указанному идентификатору. Возвращает 204 No Content при успехе
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "ID задачи"
// @Success 204 "No Content"
// @Failure 400 {object} models.ErrorResponse "Invalid ID"
// @Failure 404 {object} models.ErrorResponse "Todo not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/v1/todos/{id} [delete]
func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	if _, exists := h.todos[id]; !exists {
		http.Error(w, "Задача не найдена", http.StatusNotFound)
		return
	}

	delete(h.todos, id)

	w.WriteHeader(http.StatusNoContent)
}
