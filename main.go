package main

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Todo struct {
	ID        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

func main() {

	app := fiber.New()

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")

	todos := []Todo{}

	// Create a Todo
	app.Post("/api/todo", func(c fiber.Ctx) error {
		todo := &Todo{}

		body := c.Body()
		if err := json.Unmarshal(body, todo); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
		}

		if todo.Body == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Todo Body is required"})
		}

		todo.ID = len(todos) + 1
		todos = append(todos, *todo)

		return c.Status(201).JSON(todo)
	})

	// Update a Todo
	app.Patch("/api/todo/:id", func(c fiber.Ctx) error {
		id := c.Params("id")

		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				todos[i].Completed = !todos[i].Completed
				return c.Status(200).JSON(todos[i])
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": "todo is not found"})
	})

	//Get Todos
	app.Get("/api/todos", func(c fiber.Ctx) error {
		return c.Status(200).JSON(todos)
	})

	// Delete a Todo
	app.Delete("/api/todo/:id", func(c fiber.Ctx) error {
		id := c.Params("id")
		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				todos = append(todos[:i], todos[i+1:]...)
				return c.Status(200).JSON(fiber.Map{"success": "true"})
			}
		}
		return c.Status(404).JSON(fiber.Map{"error": "todo is not found"})
	})

	err := app.Listen(":" + port)
	if err != nil {
		log.Fatal(err)
	}

}
