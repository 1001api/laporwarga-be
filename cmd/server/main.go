package main

import (
	"hubku/lapor_warga_be_v2/internal/database"
	"hubku/lapor_warga_be_v2/internal/routes"
	"log"
	"runtime"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

var (
	db *pgxpool.Pool
)

func init() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("No .env file found, continuing with environment variables")
	}
	viper.AutomaticEnv()
}

func main() {
	// connect to database
	db = database.ConnectPG()

	r := fiber.New()

	// CORS
	r.Use(cors.New())

	// ROUTING
	routes.Routing(r, db)

	r.Get("/health", monitor.New(monitor.Config{
		Title: "Lapor Warga BE V2",
	}))

	port := viper.GetString("PORT")
	if port == "" {
		port = "8181"
	}

	log.Println("Server Succesfully to Listed in Port:", port)
	log.Println("Go version:", runtime.Version())

	log.Fatal(r.Listen(":" + port))
}
