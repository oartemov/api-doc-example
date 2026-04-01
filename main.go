package main

import (
	"log"
	"net/http"

	"api-doc-example/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	usersRouter := chi.NewRouter()
	userHandler := handlers.NewUserHandler()
	usersRouter.Get("/", userHandler.ListUsers)
	usersRouter.Get("/{id}", userHandler.GetUser)
	usersRouter.Post("/", userHandler.CreateUser)
	usersRouter.Put("/{id}", userHandler.UpdateUser)
	usersRouter.Delete("/{id}", userHandler.DeleteUser)
	r.Mount("/api/v1/users", usersRouter)

	todosRouter := chi.NewRouter()
	todoHandler := handlers.NewTodoHandler()
	todosRouter.Get("/", todoHandler.ListTodos)
	todosRouter.Get("/{id}", todoHandler.GetTodo)
	todosRouter.Post("/", todoHandler.CreateTodo)
	todosRouter.Put("/{id}", todoHandler.UpdateTodo)
	todosRouter.Delete("/{id}", todoHandler.DeleteTodo)
	r.Mount("/api/v1/todos", todosRouter)

	r.Mount("/swagger", httpSwagger.Handler(httpSwagger.URL("/docs/swagger.json")))
	r.Get("/docs/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/swagger.json")
	})

	log.Println("Сервер запущен на http://localhost:8080")
	log.Println("Swagger UI доступен по адресу: http://localhost:8080/swagger/")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
