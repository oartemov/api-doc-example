# Todo API

## Описание

REST API для управления задачами (Todo API). Позволяет создавать, обновлять и удалять задачи, а также выводить список. Также есть возможность добавления и удаления пользователей.

## Требования к окружению

| Инструмент | Версия | Назначение |
|------------|--------|------------|
| Go | 1.21+ | Компиляция и запуск приложения |
| Swaggo CLI | latest | Генерация Swagger-документации |

Установка Swag CLI:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

## Установка и запуск

### 1. Генерация Swagger-документации

```bash
swag init -g main.go -o docs
```

### 2. Запуск приложения

```bash
go run main.go
```

Сервер запустится на `http://localhost:8080`.

## Архитектура

Проект построен на принципах **Clean Architecture** с разделением на слои:

```
main.go                # Точка входа приложения
internal/
├── handlers/          # HTTP-обработчики (Presentation Layer)
│   ├── user.go
│   └── todo.go
└── models/            # Структуры данных (Domain Layer)
    ├── user.go
    └── todo.go
docs/                  # Сгенерированная Swagger-документация
```

| Слой | Директория | Ответственность |
|------|------------|-----------------|
| Presentation | `handlers/` | Обработка HTTP-запросов, валидация входных данных |
| Domain | `models/` | Бизнес-сущности и DTO |
| Entry Point | `/` | Инициализация приложения, настройка роутинга |

## Примеры использования

### Создание пользователя

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "ivan@example.com",
    "name": "Иван Петров"
  }'
```

**Ответ:**

```json
{
  "success": true,
  "data": {
    "id": 1,
    "email": "ivan@example.com",
    "name": "Иван Петров",
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

### Получение списка пользователей

```bash
curl http://localhost:8080/api/v1/users
```

### Создание задачи

```bash
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "title": "Купить продукты",
    "description": "Молоко, хлеб, яйца",
    "done": false
  }'
```

**Ответ:**

```json
{
  "success": true,
  "data": {
    "id": 1,
    "user_id": 1,
    "title": "Купить продукты",
    "description": "Молоко, хлеб, яйца",
    "done": false,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

### Обновление задачи

```bash
curl -X PUT http://localhost:8080/api/v1/todos/1 \
  -H "Content-Type: application/json" \
  -d '{
    "done": true
  }'
```

### Удаление задачи

```bash
curl -X DELETE http://localhost:8080/api/v1/todos/1
```

## Документация API

Swagger-документация доступна после генерации:

- **UI**: `http://localhost:8080/swagger/index.html`
- **JSON**: `http://localhost:8080/swagger/doc.json`

### Основные эндпоинты

#### Users

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/api/v1/users` | Получить список всех пользователей |
| GET | `/api/v1/users/{id}` | Получить пользователя по ID |
| POST | `/api/v1/users` | Создать нового пользователя |
| PUT | `/api/v1/users/{id}` | Обновить данные пользователя |
| DELETE | `/api/v1/users/{id}` | Удалить пользователя |

#### Todos

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/api/v1/todos` | Получить список всех задач |
| GET | `/api/v1/todos/{id}` | Получить задачу по ID |
| POST | `/api/v1/todos` | Создать новую задачу |
| PUT | `/api/v1/todos/{id}` | Обновить задачу |
| DELETE | `/api/v1/todos/{id}` | Удалить задачу |

### Коды ответов

| Код | Значение | Описание |
|-----|----------|----------|
| 200 | OK | Запрос выполнен успешно |
| 201 | Created | Ресурс успешно создан |
| 204 | No Content | Ресурс успешно удален |
| 400 | Bad Request | Неверный формат запроса |
| 404 | Not Found | Ресурс не найден |
| 500 | Internal Server Error | Внутренняя ошибка сервера |
