package main

import (
	"github.com/JustGritt/go-grpc/database"
	_ "github.com/JustGritt/go-grpc/docs"
	"github.com/JustGritt/go-grpc/routes"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/gofiber/swagger"
	"github.com/golang-jwt/jwt/v4" // gin-swagger middleware
)

func welcome(c *fiber.Ctx) error {
	return c.SendString("Welcome to Fiber!")
}

func restricted(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	c.SendString("Welcome " + name + "!")
	return c.Next()
}

func setupRoutes(app *fiber.App) {
	// Welcome
	app.Get("/api", welcome)
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Login route
	app.Post("/login", routes.Login)
	app.Get("/api/stream", routes.GetStream)
	app.Post("/api/users", routes.CreateUser)

	app.Use(jwtware.New(jwtware.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Token is missing, returns with error code 400 "Missing or malformed JWT"
			return c.Status(400).JSON(fiber.Map{
				"message": "Missing or malformed JWT",
			})
		},
		SigningKey: []byte("secret"),
	}))

	// Restricted Routes
	app.Get("/restricted", restricted)

	// GET routes
	// =================
	app.Get("/api/users", restricted, routes.GetUsers)
	app.Get("/api/users/:id", restricted, routes.GetUser)
	app.Get("/api/products", restricted, routes.GetProducts)
	app.Get("/api/products/:id", restricted, routes.GetProduct)
	app.Get("/api/payments", restricted, routes.GetPayments)
	app.Get("/api/payments/:id", restricted, routes.GetPayment)
	// app.Get("/api/payments/user/:id", routes.GetPaymentsByUser)
	// app.Get("/api/payments/product/:id", routes.GetPaymentsByProduct)

	// POST routes
	// =================
	app.Post("/api/products", restricted, routes.CreateProduct)
	app.Post("/api/payments", restricted, routes.CreatePayment)

	// PUT routes
	// =================
	app.Put("/api/users/:id", restricted, routes.UpdateUser)
	app.Put("/api/products/:id", restricted, routes.UpdateProduct)
	app.Put("/api/payments/:id", restricted, routes.UpdatePayment)

	// DELETE routes
	// =================
	app.Delete("/api/users/:id", restricted, routes.DeleteUser)
	app.Delete("/api/products/:id", restricted, routes.DeleteProduct)
	app.Delete("/api/payments/:id", restricted, routes.DeletePayment)

}

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/
// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      localhost:3000
// @BasePath  /
func main() {
	database.Connect()
	app := fiber.New()

	setupRoutes(app)

	app.Listen(":3000")
}
