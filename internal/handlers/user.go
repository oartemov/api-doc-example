package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"

	"api-doc-example/internal/models"
)

type UserHandler struct {
	mu     sync.RWMutex
	users  map[int]models.User
	nextID int
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		users:  make(map[int]models.User),
		nextID: 1,
	}
}

// ListUsers возвращает список всех пользователей
// @Summary Получить список всех пользователей
// @Description Возвращает массив объектов User в структуре SuccessResponse
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} models.SuccessResponse{data=[]models.User}
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/v1/users [get]
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var users []models.User
	for _, user := range h.users {
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(models.SuccessResponse{
		Success: true,
		Data:    users,
	})
}

// GetUser возвращает пользователя по ID
// @Summary Получить пользователя по ID
// @Description Возвращает данные пользователя по указанному идентификатору
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} models.SuccessResponse{data=models.User}
// @Failure 400 {object} models.ErrorResponse "Invalid ID"
// @Failure 404 {object} models.ErrorResponse "User not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	h.mu.RLock()
	user, exists := h.users[id]
	h.mu.RUnlock()

	if !exists {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(models.SuccessResponse{
		Success: true,
		Data:    user,
	})
}

// CreateUser создает нового пользователя
// @Summary Создать нового пользователя
// @Description Создает нового пользователя с указанными email и именем. Возвращает созданного пользователя с присвоенным ID
// @Tags users
// @Accept json
// @Produce json
// @Param request body models.CreateUserRequest true "Данные для создания пользователя"
// @Success 201 {object} models.SuccessResponse{data=models.User}
// @Failure 400 {object} models.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/v1/users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Невалидный запрос", http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	user := models.User{
		ID:        h.nextID,
		Email:     req.Email,
		Name:      req.Name,
		CreatedAt: time.Now(),
	}
	h.users[h.nextID] = user
	h.nextID++
	h.mu.Unlock()

	log.Printf("Создан новый пользователь: %+v", user)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(models.SuccessResponse{
		Success: true,
		Data:    user,
	})
}

// UpdateUser обновляет данные пользователя
// @Summary Обновить данные пользователя
// @Description Обновляет email и/или имя пользователя по указанному ID. Возвращает обновленные данные пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Param request body models.UpdateUserRequest true "Данные для обновления (все поля опциональны)"
// @Success 200 {object} models.SuccessResponse{data=models.User}
// @Failure 400 {object} models.ErrorResponse "Invalid ID or request body"
// @Failure 404 {object} models.ErrorResponse "User not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/v1/users/{id} [put]
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Невалидный запрос", http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	user, exists := h.users[id]
	if !exists {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Name != nil {
		user.Name = *req.Name
	}

	h.users[id] = user

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(models.SuccessResponse{
		Success: true,
		Data:    user,
	})
}

// DeleteUser удаляет пользователя по ID
// @Summary Удалить пользователя
// @Description Удаляет пользователя по указанному идентификатору. Возвращает 204 No Content при успехе
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 204 "No Content"
// @Failure 400 {object} models.ErrorResponse "Invalid ID"
// @Failure 404 {object} models.ErrorResponse "User not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.users[id]; !exists {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	delete(h.users, id)

	w.WriteHeader(http.StatusNoContent)
}
