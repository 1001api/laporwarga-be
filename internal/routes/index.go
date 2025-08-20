package routes

import (
	"hubku/lapor_warga_be_v2/internal/controllers"
	"hubku/lapor_warga_be_v2/internal/modules/auth"
	"hubku/lapor_warga_be_v2/internal/modules/users"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

func Routing(r fiber.Router, db *pgxpool.Pool) {
	var encKey = viper.GetString("ENC_KEY")

	if encKey == "" {
		log.Fatal("ENC_KEY is not set")
	}

	validator := validator.New()

	userRepo := users.NewUserRepository(db, encKey)

	userService := users.NewUserService(userRepo)
	authService := auth.NewAuthService(userService)

	userController := controllers.NewUserController(userService, validator)
	authController := controllers.NewAuthController(authService, validator)

	// API versioning
	versioning := r.Group("/api/v1")

	public := versioning.Group("/auth")
	{
		public.Post("/login", authController.Login)
		public.Post("/refresh", authController.Refresh)
	}

	userRoutes := versioning.Group("/users", JWTMiddleware(authService))
	{
		userRoutes.Get("/list", userController.GetMasterUser)
		userRoutes.Post("/create", RoleMiddleware("admin", "superadmin"), authController.Register)
		userRoutes.Get("/me", userController.GetCurrentUser)
		userRoutes.Patch("/me", userController.UpdateCurrentUser)
		userRoutes.Get("/search", RoleMiddleware("admin", "superadmin"), userController.SearchUser)
		userRoutes.Get("/:id", RoleMiddleware("admin", "superadmin"), userController.GetUserByID)
		userRoutes.Post("/restore/:id", RoleMiddleware("admin", "superadmin"), userController.RestoreUser)
		userRoutes.Patch("/:id", RoleMiddleware("admin", "superadmin"), userController.UpdateUser)
		userRoutes.Delete("/:id", RoleMiddleware("admin", "superadmin"), userController.DeleteUser)
	}
}

func JWTMiddleware(authService auth.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header required",
			})
		}

		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format",
			})
		}

		tokenString := authHeader[7:]

		// Validate token
		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

func RoleMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("role")
		if userRole == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User role not found in context",
			})
		}

		role := cast.ToString(userRole)
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Insufficient permissions",
		})
	}
}
