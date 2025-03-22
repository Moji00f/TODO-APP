package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
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
	Id        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Completed bool               `json:"complete"`
	Body      string             `json:"body"`
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

	//app.Use(cors.New(cors.Config{
	//	AllowOrigins:     []string{"http://localhost:5173"},
	//	AllowHeaders:     []string{"Origin, Content-Type, Accept"},
	//	AllowMethods:     []string{"GET, POST, PUT, DELETE, PATCH, OPTIONS"}, // 🔥 اینجا PATCH رو اضافه کردیم
	//	AllowCredentials: true,
	//}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowHeaders: []string{"Origin,Content-Type,Accept"},
	}))

	app.Get("/api/todos", GetTodos)
	app.Post("/api/todos", CreateTodo)
	app.Patch("/api/todos/:id", UpdateTodo)
	app.Delete("/api/todos/:id", DeleteTodo)

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

//func UpdateTodo(c fiber.Ctx) error {
//	id := c.Params("id")
//	objectID, err := primitive.ObjectIDFromHex(id)
//
//	if err != nil {
//		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid todo ID"})
//	}
//
//	filter := bson.M{"_id": objectID}
//	update := bson.M{"$set": bson.M{"complete": true}}
//
//	result, err := collection.UpdateOne(context.Background(), filter, update)
//	if err != nil {
//		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
//	}
//
//	if result.ModifiedCount == 0 {
//		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Todo not found"})
//	}
//
//	updatedTodo := bson.M{}
//	err = collection.FindOne(context.Background(), filter).Decode(&updatedTodo)
//	if err != nil {
//		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
//	}
//
//	return c.Status(http.StatusOK).JSON(updatedTodo)
//}

//func UpdateTodo(c fiber.Ctx) error {
//	id := c.Params("id")
//	fmt.Println("Received ID:", id) // لاگ برای ID دریافتی
//
//	objectID, err := primitive.ObjectIDFromHex(id)
//	if err != nil {
//		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid todo ID"})
//	}
//
//	filter := bson.M{"_id": objectID}
//	update := bson.M{"$set": bson.M{"completed": true}}
//
//	result, err := collection.UpdateOne(context.Background(), filter, update)
//	if err != nil {
//		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
//	}
//
//	if result.MatchedCount == 0 {
//		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Todo not found"})
//	}
//	// دریافت اطلاعات به‌روز شده
//	updatedTodo := bson.M{}
//	err = collection.FindOne(context.Background(), filter).Decode(&updatedTodo)
//	if err != nil {
//		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
//	}
//
//	return c.Status(http.StatusOK).JSON(updatedTodo)
//}

func UpdateTodo(c fiber.Ctx) error {
	id := c.Params("id")
	fmt.Println("Received ID:", id) // لاگ برای ID دریافتی

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid todo ID"})
	}

	// ابتدا داده‌ی فعلی را از دیتابیس می‌خوانیم تا مقدار complete را بررسی کنیم
	var todo bson.M
	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&todo)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Todo not found"})
	}

	// اگر complete == true، آن را به false تغییر بدهیم و برعکس
	newCompleteStatus := false
	if todo["completed"] != nil && todo["completed"].(bool) == false {
		newCompleteStatus = true
	}

	// حالا آپدیت کردن وضعیت complete
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"completed": newCompleteStatus}}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if result.MatchedCount == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Todo not found"})
	}

	// دریافت اطلاعات به‌روز شده
	updatedTodo := bson.M{}
	err = collection.FindOne(context.Background(), filter).Decode(&updatedTodo)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(updatedTodo)
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
