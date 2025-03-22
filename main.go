package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
)

type Todo struct {
	Id       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Complete bool               `json:"complete"`
	Body     string             `json:"body"`
}

var collection *mongo.Collection

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	MONGODB_URI := os.Getenv("MONGODB_URI")
	clientOptions := options.Client().ApplyURI(MONGODB_URI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(client, context.Background())

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connect to Mongodb")

	collection = client.Database("GoReact").Collection("todo")

	app := fiber.New()

	app.Get("/api/todos", GetTodos)
	app.Post("/api/todo", CreateTodo)
	app.Patch("/api/todo/:id", UpdateTodo)
	app.Delete("/api/todo/:id", DeleteTodo)

	port := os.Getenv("PORT")

	err = app.Listen(":" + port)
	if err != nil {
		log.Fatal(err)
	}
}

func GetTodos(c fiber.Ctx) error {
	var todos []Todo

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var todo Todo
		if err := cursor.Decode(&todo); err != nil {
			return err
		}
		todos = append(todos, todo)
	}
	return c.Status(200).JSON(todos)
}

func CreateTodo(c fiber.Ctx) error {

	todo := new(Todo)
	body := c.Body()
	if err := json.Unmarshal(body, todo); err != nil {
		return err
	}

	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Todo body cannot be empty"})
	}

	todo.Id = primitive.NewObjectID()

	InsertResult, err := collection.InsertOne(context.Background(), todo)
	if err != nil {
		return err
	}

	todo.Id = InsertResult.InsertedID.(primitive.ObjectID)

	return c.Status(http.StatusCreated).JSON(todo)
}

func UpdateTodo(c fiber.Ctx) error {

	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid todo ID"})
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"completed": "true"}}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"success": "true"})
}

func DeleteTodo(c fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid todo ID"})
	}

	filter := bson.M{"_id": objectID}
	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"success": "true"})
}
