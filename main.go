package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"go_back/internal/controllers"
	"go_back/internal/initializers"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatalln("Failed to load environment variables! \n", err.Error())
	}

	initializers.ConnectDB(&config)
}

func main() {
	app := fiber.New()
	micro := fiber.New()

	app.Mount("/api", micro)

	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${ip} ${method} ${status}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
		Output:     os.Stdout,
	}))

	// Используйте middleware CORS
	micro.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:8081, http://192.168.1.166:8000",
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Content-Type,Authorization",
	}))

	micro.Route("/user", func(router fiber.Router) {
		router.Post("/login", controllers.LoginUser)
		router.Post("/register", controllers.CreateUser)
		router.Post("/send-code", controllers.SendNewUserCode)
		router.Post("/confirm-code", controllers.ConfirmNewUser)
		router.Post("/change-password", controllers.ChangePassword)
		router.Post("/confirm-change-password", controllers.ConfirmChangePassword)
		router.Delete("/:userId", controllers.DeleteUser)
		router.Get("/", controllers.FindUsers)
		router.Get("/:userId", controllers.FindUserById)
	})

	micro.Route("/data", func(router fiber.Router) {
		router.Get("/schedule/:userId", controllers.GetSchedule)
		router.Post("/new-schedule/:userId", controllers.NewSchedule)
		router.Post("/update-schedule/:userId/:id", controllers.UpdateSchedule)
		router.Delete("/delete-schedule/:userId/:id", controllers.DeleteSchedule)

		router.Get("/tasks/:userId", controllers.GetTasks)
		router.Post("/new-task/:userId", controllers.NewTask)
		router.Post("/update-task/:userId/:id", controllers.UpdateTask)
		router.Delete("/delete-task/:userId/:id", controllers.DeleteTask)

		router.Get("/calendar/:userId", controllers.GetCalendar)
		router.Post("/new-calendar/:userId", controllers.NewCalendar)
		router.Post("/update-calendar/:userId/:id", controllers.UpdateCalendar)
		router.Delete("/delete-calendar/:userId/:id", controllers.DeleteCalendar)
	})

	log.Fatal(app.Listen(":8000"))
}
