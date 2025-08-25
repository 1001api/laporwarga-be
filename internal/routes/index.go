package routes

import (
	"hubku/lapor_warga_be_v2/internal/controllers"
	"hubku/lapor_warga_be_v2/internal/modules/auditlogs"
	"hubku/lapor_warga_be_v2/internal/modules/auth"
	userroles "hubku/lapor_warga_be_v2/internal/modules/user_roles"
	"hubku/lapor_warga_be_v2/internal/modules/users"
	"hubku/lapor_warga_be_v2/pkg"
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

	userRepo := users.NewUserRepository(db)
	roleRepo := userroles.NewUserRolesRepository(db)
	logRepo := auditlogs.NewLogsRepository(db)

	logService := auditlogs.NewLogsService(logRepo)
	userRolesService := userroles.NewUserRolesService(roleRepo, logService)
	userService := users.NewUserService(userRepo, userRolesService, logService, encKey)
	authService := auth.NewAuthService(userService, logService, encKey)

	logsController := controllers.NewLogsController(logService)
	userController := controllers.NewUserController(userService, validator)
	authController := controllers.NewAuthController(authService, validator)
	userRolesController := controllers.NewUserRolesController(userRolesService, validator)

	// Initialize root user
	if err := userService.InitializeRootUser(); err != nil {
		log.Fatal("Failed to initialize root user:", err)
	}

	// API versioning
	versioning := r.Group("/api/v1")

	public := versioning.Group("/auth")
	{
		public.Post("/login", authController.Login)
		public.Post("/refresh", authController.Refresh)
	}

	userRoutes := versioning.Group("/users", JWTMiddleware(authService))
	{
		userRoutes.Get("/me", userController.GetCurrentUser)
		userRoutes.Patch("/me", userController.UpdateCurrentUser)
		userRoutes.Get("/list", RoleMiddleware(string(pkg.RoleAdmin)), userController.GetMasterUser)
		userRoutes.Post("/create", RoleMiddleware(string(pkg.RoleAdmin)), authController.Register)
		userRoutes.Get("/search", RoleMiddleware(string(pkg.RoleAdmin)), userController.SearchUser)
		userRoutes.Get("/:id", RoleMiddleware(string(pkg.RoleAdmin)), userController.GetUserByID)
		userRoutes.Post("/restore/:id", RoleMiddleware(string(pkg.RoleAdmin)), userController.RestoreUser)
		userRoutes.Patch("/:id", RoleMiddleware(string(pkg.RoleAdmin)), userController.UpdateUser)
		userRoutes.Delete("/:id", RoleMiddleware(string(pkg.RoleAdmin)), userController.DeleteUser)
	}

	rolesRoutes := versioning.Group("/roles", JWTMiddleware(authService), RoleMiddleware(string(pkg.RoleAdmin)))
	{
		rolesRoutes.Get("/list", userRolesController.ListAllRoles)
		rolesRoutes.Post("/create", userRolesController.CreateRole)
		rolesRoutes.Post("/assign/:id", userRolesController.AssignRole)
		rolesRoutes.Get("/id/:id", userRolesController.GetRoleByID)
		rolesRoutes.Get("/name/:name", userRolesController.GetRoleByName)
		rolesRoutes.Put("/:id", userRolesController.UpdateRole)
		rolesRoutes.Delete("/:id", userRolesController.RemoveRole)
	}

	logsRoutes := versioning.Group("/logs", JWTMiddleware(authService), RoleMiddleware(string(pkg.RoleAdmin)))
	{
		logsRoutes.Get("/list", logsController.ListLogs)
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
