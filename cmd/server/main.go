package main

import (
	"hubku/lapor_warga_be_v2/internal/database"
	"hubku/lapor_warga_be_v2/internal/routes"
	"log"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/earlydata"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
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

	// Recover middleware
	r.Use(recover.New())

	// Early Data supporting TLS 1.3 0-RTT
	// by default it will only allow early data for GET, HEAD, OPTIONS and TRACE requests (safe method)
	r.Use(earlydata.New())

	// Global Limit
	r.Use(limiter.New(limiter.Config{
		Max:        60,               // 60 requests
		Expiration: 60 * time.Second, // 1 minute
		KeyGenerator: func(c *fiber.Ctx) string {
			ip := c.IP()
			if fwd := c.Get("X-Forwarded-For"); fwd != "" {
				ip = fwd
			}
			return ip
		},
		LimiterMiddleware: limiter.SlidingWindow{},
	}))

	// Request Logging
	r.Use(logger.New(logger.Config{
		Format: "${time} ${ip} ${method} ${path} ${status} ${latency}\n",
	}))

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     viper.GetString("CLIENT_DOMAIN"),
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// Encrypt Cookie
	r.Use(encryptcookie.New(encryptcookie.Config{
		Key: viper.GetString("COOKIE_ENC_KEY"),
	}))

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
