![banner](https://i.imgur.com/QMg08hA.png)

# Lapor Warga - Backend API

> ⚠️ **Work in Progress**: This project is under heavy development. Features, APIs, and documentation are subject to change or addition without prior notice.

A modern, secure, and scalable backend system for a public citizen reporting platform. Built with Go, PostgreSQL, and Fiber framework to enable citizens to report issues and government officials to manage and respond to community concerns.

## 🌟 Features

### Core Functionality
- **🔐 Authentication & Authorization**
  - JWT-based authentication with access and refresh tokens
  - Secure cookie-based session management for web clients
  - Bearer token authentication for mobile clients
  - Role-based access control (Admin, Official, Citizen)
  - Password hashing with bcrypt

- **👥 User Management**
  - User registration and profile management
  - Email and phone verification support
  - Credibility scoring system
  - User status management (probation, regular, suspended)
  - OAuth provider support (extensible)
  - Failed login attempt tracking with account locking
  - Soft delete with restoration capability

- **🗺️ Geographic Area Management**
  - PostGIS-powered geospatial data handling
  - Multi-polygon boundary support for administrative areas
  - Automatic center point calculation
  - Area hierarchy support (provinsi, kabupaten, kecamatan)
  - Spatial indexing with GIST for efficient queries
  - Configurable boundary simplification (off, simple, detail)

- **📋 Audit Logging**
  - Comprehensive activity tracking
  - JSONB metadata storage for flexible log data
  - Entity-based logging (users, roles, areas)
  - Action tracking (create, update, delete, assign, restore, login)

- **🛡️ Security Features**
  - AES-GCM encryption for sensitive data (email, phone, fullname)
  - SHA-256 hashing for searchable encrypted fields
  - Encrypted cookies
  - Rate limiting (60 requests per minute per IP)
  - CORS configuration
  - TLS 1.3 0-RTT early data support

### Technical Features
- **High Performance**: Built with Fiber framework for optimal performance
- **Structured Logging**: Zerolog integration with Axiom for centralized logging
- **Database**: PostgreSQL with PostGIS extension
- **Type-Safe Queries**: SQLC for compile-time SQL validation
- **Hot Reload**: Air for development with live reload
- **Clean Architecture**: Separation of concerns with controllers, services, and repositories
- **API Versioning**: `/api/v1` prefix for future compatibility
- **Mobile API**: Dedicated mobile endpoints with custom authentication

## 🏗️ Architecture

```
lapor_warga_be/
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── controllers/     # HTTP handlers
│   ├── database/
│   │   ├── generated/   # SQLC generated code
│   │   ├── migrations/  # Database migrations
│   │   └── queries/     # SQL queries for SQLC
│   ├── modules/         # Business logic
│   │   ├── areas/       # Area management
│   │   ├── auditlogs/   # Audit logging
│   │   ├── auth/        # Authentication
│   │   ├── user_roles/  # Role management
│   │   └── users/       # User management
│   └── routes/          # Route definitions & middleware
├── pkg/                 # Shared utilities
└── scripts/             # Helper scripts
```

## 🚀 Getting Started

### Prerequisites

- **Go**: 1.24.0 or higher
- **PostgreSQL**: 14+ with PostGIS extension
- **golang-migrate**: For database migrations
- **sqlc**: For generating type-safe Go code from SQL
- **Air** (optional): For hot reload during development

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/lapor_warga_be.git
   cd lapor_warga_be
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   
   Create a `.env` file in the root directory:
   ```env
   # Server Configuration
   PORT=8181
   APP_DOMAIN=localhost
   ENV_PROD=false
   CLIENT_DOMAIN=http://localhost:3000

   # Database
   DATABASE_URL=postgresql://user:password@localhost:5432/lapor_warga?sslmode=disable

   # Security Keys (generate secure random strings)
   ENC_KEY=your-32-byte-encryption-key-here
   COOKIE_ENC_KEY=your-cookie-encryption-key

   # JWT Configuration
   JWT_EXPIRY=15        # minutes
   JWT_REFRESH_EXPIRY=4320  # minutes (3 days)

   # Mobile API
   MOBILE_KEY=your-mobile-api-key

   # Logging (Axiom)
   AXIOM_TOKEN=your-axiom-token
   AXIOM_DATASET=your-dataset-name
   ```

4. **Run database migrations**
   ```bash
   chmod +x scripts/migrate.sh
   ./scripts/migrate.sh up
   ```

5. **Generate SQLC code** (if you modify SQL queries)
   ```bash
   sqlc generate
   ```

6. **Run the application**
   
   **Development (with hot reload):**
   ```bash
   air
   ```
   
   **Production:**
   ```bash
   go run cmd/server/main.go
   ```

The server will start on `http://localhost:8181` (or your configured PORT).

## 📡 API Endpoints

### Authentication
- `POST /api/v1/auth/login` - User login (web)
- `POST /api/v1/auth/refresh` - Refresh access token
- `GET /api/v1/auth/session` - Get current session info

### Users
- `GET /api/v1/users/me` - Get current user profile
- `PATCH /api/v1/users/me` - Update current user profile
- `GET /api/v1/users/list` - List all users (Admin only)
- `POST /api/v1/users/create` - Create new user (Admin only)
- `GET /api/v1/users/search` - Search users (Admin only)
- `GET /api/v1/users/:id` - Get user by ID (Admin only)
- `PATCH /api/v1/users/:id` - Update user (Admin only)
- `DELETE /api/v1/users/:id` - Soft delete user (Admin only)
- `POST /api/v1/users/restore/:id` - Restore deleted user (Admin only)

### Roles
- `GET /api/v1/roles/list` - List all roles (Admin only)
- `POST /api/v1/roles/create` - Create new role (Admin only)
- `POST /api/v1/roles/assign/:id` - Assign role to user (Admin only)
- `GET /api/v1/roles/id/:id` - Get role by ID (Admin only)
- `GET /api/v1/roles/name/:name` - Get role by name (Admin only)
- `PUT /api/v1/roles/:id` - Update role (Admin only)
- `DELETE /api/v1/roles/:id` - Delete role (Admin only)

### Areas
- `POST /api/v1/areas/create` - Create new area (Admin only)
- `GET /api/v1/areas/list` - List all areas with pagination
- `GET /api/v1/areas/boundary/:id` - Get area boundary geometry
- `PATCH /api/v1/areas/toggle-status/:id` - Toggle area active status (Admin only)

### Audit Logs
- `GET /api/v1/logs/list` - List audit logs (Admin only)

### Mobile API
- `POST /api/v1/m/auth/login` - Mobile login
- `POST /api/v1/m/auth/refresh` - Mobile token refresh

### Health Check
- `GET /health` - Server health and monitoring dashboard

## 🔧 Development

### Database Migrations

**Create a new migration:**
```bash
migrate create -ext sql -dir internal/database/migrations -seq migration_name
```

**Apply migrations:**
```bash
./scripts/migrate.sh up
```

**Rollback last migration:**
```bash
./scripts/migrate.sh down
```

**Go to specific version:**
```bash
./scripts/migrate.sh goto <version>
```

**Force version (use with caution):**
```bash
./scripts/migrate.sh force <version>
```

### SQLC Code Generation

After modifying SQL queries in `internal/database/queries/`, regenerate Go code:
```bash
sqlc generate
```

### Project Structure

- **Controllers**: Handle HTTP requests and responses
- **Services**: Contain business logic
- **Repositories**: Database access layer
- **Middleware**: JWT authentication, role-based access, rate limiting
- **Migrations**: Database schema versioning
- **Queries**: Type-safe SQL queries with SQLC

## 🔐 Security

### Data Encryption
- **Sensitive fields** (email, phone, fullname) are encrypted using AES-GCM
- **Hash fields** enable searching without decryption using SHA-256
- **Passwords** are hashed with bcrypt (cost factor: 12)

### Authentication
- **JWT tokens** with configurable expiration
- **Refresh tokens** for extended sessions
- **Cookie-based** authentication for web clients
- **Bearer token** authentication for mobile clients

### Rate Limiting
- **60 requests per minute** per IP address
- **Sliding window** algorithm
- Configurable per-route limits

## 🌍 Timezone

The application uses **Asia/Jakarta (WIB)** timezone for all timestamps and logging.

## 📊 Monitoring

- **Health endpoint**: `/health` provides server metrics
- **Structured logging**: JSON logs with Zerolog
- **Centralized logging**: Axiom integration for log aggregation
- **Request logging**: Automatic logging of all HTTP requests

## 📦 Dependencies

### Core
- **Fiber v2**: Fast HTTP framework
- **pgx/v5**: PostgreSQL driver
- **SQLC**: Type-safe SQL code generation
- **Viper**: Configuration management
- **JWT**: JSON Web Token implementation

### Security
- **bcrypt**: Password hashing
- **AES-GCM**: Data encryption
- **validator/v10**: Input validation

### Geospatial
- **PostGIS**: PostgreSQL spatial extension
- **paulmach/orb**: Geometry handling in Go

### Logging & Monitoring
- **Zerolog**: Structured logging
- **Axiom**: Log aggregation and analytics

## 🤝 Contributing

This project is currently under heavy development. Contribution guidelines will be added soon.

## 📞 Support

For issues and questions, please open an issue on GitHub.

---

**Note**: This is a backend API service. It requires a frontend application and/or mobile app to provide a complete user experience.