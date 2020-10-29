package main

import (
	"net/http"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	_ "github.com/volatiletech/authboss/v3/auth"
	_ "github.com/volatiletech/authboss/v3/logout"
	"github.com/volatiletech/authboss/v3/remember"
)

func main() {
	ab := SetupAuthboss()

	app := fiber.New()

	app.Use(adaptor.HTTPMiddleware(ab.LoadClientStateMiddleware), adaptor.HTTPMiddleware(remember.Middleware(ab)))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	app.Group("/auth", adaptor.HTTPHandler(http.StripPrefix("/auth", ab.Config.Core.Router)))

	_ = app.Listen(":3000")
}
