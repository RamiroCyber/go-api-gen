# Go API Generator CLI

`go-api-gen` is a command-line interface (CLI) tool written in Go to generate RESTful API modules with a layered architecture (model, repository, service, controller) using the [Fiber](https://gofiber.io) framework. It automates boilerplate code generation, supporting CRUD operations, custom methods, pagination, field validations, custom errors, unit tests, JWT authentication, Swagger documentation, Docker, and CI/CD integration. Ideal for accelerating development of robust and scalable APIs for various applications, from small projects to enterprise systems.

## Features

- **Layered Architecture**: Generates `model`, `repository`, `service`, and `controller` for each module, organized in `internal/modules/<module-name>`.
- **CRUD Operations**: Full support for Create, Read, Update, Delete, with soft delete (via `DeletedAt`).
- **Pagination**: `List` method supports `page` and `limit` parameters, returning JSON like `{data: [], total: int, page: int, limit: int}`.
- **Custom Fields**: Define model fields with types and validations via `--fields` (e.g., `name:string required,email:string email unique`).
- **Database Support**:
  - SQL (PostgreSQL, MySQL, SQLite) via `database/sql`.
  - GORM ORM (`--db gorm`) for query abstraction.
  - Placeholder for MongoDB (`--db mongo`).
- **Validations**: Integrates `github.com/go-playground/validator/v10` for field validation (`--validator`).
- **Custom Errors**: Generated in `pkg/errors` (e.g., `ErrNotFound`, `ErrInvalidEntity`) for robust error handling.
- **Unit Tests**: Generates test files with mocks using `github.com/stretchr/testify` (`--tests`).
- **JWT Authentication**: Generates middleware in `pkg/middleware/auth.go` (`--auth jwt`).
- **Swagger Documentation**: Adds annotations for `swagger.json/yaml` generation with `swag init` (`--swagger`).
- **Docker and CI/CD**: Generates `Dockerfile` and `.github/workflows/ci.yml` for builds and testing (`--docker`, `--ci`).
- **Custom Methods**: Add extra methods via `--methods` (e.g., `FindByEmail`).
- **Security and Performance**: Uses safe queries with placeholders, suggests database indexes, and includes validations in service/controller layers.
- **Version Command**: Supports `go-api-gen --version` to display the CLI version.

## Requirements

- **Go**: 1.20 or higher.
- **Generated Project Dependencies** (installed via `go get` in the generated project):
  - `github.com/gofiber/fiber/v2` (controllers).
  - `github.com/google/uuid` (IDs).
  - `github.com/go-playground/validator/v10` (if `--validator`).
  - `gorm.io/gorm` and `gorm.io/driver/postgres` (or other driver, if `--db gorm`).
  - `github.com/golang-jwt/jwt/v5` (if `--auth`).
  - `github.com/stretchr/testify` (if `--tests`).
  - For Swagger: `github.com/swaggo/swag` and `github.com/swaggo/fiber-swagger`.

## Installation

### Remote Installation (Recommended)

Install `go-api-gen` directly from the latest GitHub release without cloning the repository:

```bash
curl -sSL https://raw.githubusercontent.com/RamiroCyber/go-api-gen/main/install.sh | bash
```

This downloads the appropriate binary (e.g., `go-api-gen-linux-amd64` or `go-api-gen-macos-arm64`) to `$HOME/.local/bin`, makes it executable, and adds the directory to your PATH if needed. After installation, run:

```bash
source ~/.bashrc  # or source ~/.zshrc
go-api-gen --version
```

### Local Installation (Development)

1. Clone the repository:
   ```bash
   git clone https://github.com/RamiroCyber/go-api-gen.git
   cd go-api-gen
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Build and install locally:
   ```bash
   go build -o go-api-gen main.go
   mv go-api-gen $HOME/.local/bin/
   ```

4. Ensure `$HOME/.local/bin` is in your PATH:
   ```bash
   echo 'export PATH="$PATH:$HOME/.local/bin"' >> ~/.bashrc  # or ~/.zshrc
   source ~/.bashrc
   ```

5. Verify:
   ```bash
   go-api-gen --version
   ```

## Usage

The main command is `go-api-gen generate module <module-name> [flags]`, which creates a module in `internal/modules/<module-name>`.

### Available Flags

| Flag             | Description                                                               | Example/Explanation                                                                 |
|------------------|---------------------------------------------------------------------------|-----------------------------------------------------------------------------------|
| `--methods`      | Comma-separated custom methods                                            | `--methods FindByEmail,FindByName`                                                |
| `--fields`       | Comma-separated fields: `name:type [tags]`                                | `--fields "name:string required,email:string email unique,age:int gte=18"`        |
| `--db`           | Database type (default: postgres)                                         | `--db gorm` (or `postgres`, `mysql`, `mongo`)                                     |
| `--validator`    | Enable validations with `validator/v10` (default: true)                   | `--validator false` to disable                                                    |
| `--tests`        | Generate unit tests with mocks (default: false)                           | `--tests true`                                                                    |
| `--auth`         | Authentication type (option: `jwt`)                                       | `--auth jwt` generates middleware in `pkg/middleware/auth.go`                      |
| `--swagger`      | Generate Swagger annotations (default: false)                             | `--swagger true` (requires `swag init` after generation)                          |
| `--docker`       | Generate `Dockerfile` for build and deploy (default: false)               | `--docker true`                                                                   |
| `--ci`           | Generate `.github/workflows/ci.yml` for CI/CD (default: false)            | `--ci true`                                                                       |
| `--root-package` | Root package for the generated project (default: github.com/user/project) | `--root-package github.com/seuusuario/seuprojeto`                                 |

### Example Usage

1. **Basic Module Generation**:
   ```bash
   go-api-gen generate module user
   ```
   Generates a `user` module with basic CRUD, using PostgreSQL (`database/sql`).

2. **Full-Featured Module**:
   ```bash
   go-api-gen generate module user --methods FindByEmail,FindByName --fields "name:string required,email:string email unique,age:int gte=18" --db gorm --validator true --tests true --auth jwt --swagger true --docker true --ci true --root-package github.com/seuusuario/seuprojeto
   ```
   Generates:
   - `internal/modules/user/` with model (fields `Name`, `Email`, `Age`), repository (GORM), service (with validations), controller (with Swagger annotations), and tests.
   - `pkg/errors/errors.go` for custom errors.
   - `pkg/middleware/auth.go` for JWT authentication.
   - `Dockerfile` and `.github/workflows/ci.yml`.

## Generated Project Structure

After running the full-featured command above, the project structure will be:

```
seuprojeto/
├── cmd/
│   └── api/
│       └── main.go  # Example initialization (create manually or use --project-init)
├── internal/
│   └── modules/
│       └── user/
│           ├── model.go
│           ├── repository.go
│           ├── repository_impl.go
│           ├── service.go
│           ├── service_impl.go
│           ├── controller.go
│           └── service_test.go  # If --tests
├── pkg/
│   ├── errors/
│   │   └── errors.go
│   └── middleware/
│       └── auth.go  # If --auth=jwt
├── swagger/
│   ├── swagger.json  # Generated via swag init
│   └── swagger.yaml
├── Dockerfile  # If --docker
├── .github/
│   └── workflows/
│       └── ci.yml  # If --ci
├── go.mod
└── go.sum
```

## Setting Up the Generated Project

1. **Initialize the Project**:
   ```bash
   mkdir seuprojeto && cd seuprojeto
   go mod init github.com/seuusuario/seuprojeto
   ```

2. **Generate a Module**:
   Use the `go-api-gen generate module ...` command as shown above.

3. **Install Dependencies**:
   ```bash
   go get github.com/gofiber/fiber/v2
   go get github.com/google/uuid
   go get github.com/go-playground/validator/v10  # If --validator
   go get gorm.io/gorm gorm.io/driver/postgres  # If --db gorm
   go get github.com/golang-jwt/jwt/v5  # If --auth
   go get github.com/stretchr/testify  # If --tests
   ```

4. **Set Up the Server (cmd/api/main.go)**:
   Example initialization:
   ```go
   package main

   import (
       "database/sql"
       "github.com/gofiber/fiber/v2"
       "github.com/seuusuario/seuprojeto/internal/modules/user"
       "github.com/seuusuario/seuprojeto/pkg/middleware"
       _ "github.com/lib/pq" // or other driver
   )

   func main() {
       app := fiber.New()
       db, _ := sql.Open("postgres", "your-connection-string")
       repo := user.NewUserRepository(db)
       svc := user.NewUserService(repo)
       ctrl := user.NewUserController(svc)
       app.Post("/user", middleware.JWTAuth("secret"), ctrl.Create) // Example with JWT
       app.Listen(":8080")
   }
   ```

5. **Generate Swagger Documentation** (if `--swagger`):
   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   swag init
   ```
   Access at `/swagger/index.html` (requires `github.com/swaggo/fiber-swagger`).

6. **Run Tests** (if `--tests`):
   ```bash
   go test ./...
   ```

7. **Build and Deploy with Docker** (if `--docker`):
   ```bash
   docker build -t seuprojeto .
   docker run -p 8080:8080 seuprojeto
   ```

## Contributing to go-api-gen

1. Clone:
   ```bash
   git clone https://github.com/RamiroCyber/go-api-gen.git
   cd go-api-gen
   ```

2. Add templates in `templates/` or modify `main.go`.

3. Test locally:
   ```bash
   go run main.go generate module user --fields "name:string required" --db gorm
   ```

4. Open issues or PRs for bugs or features (e.g., support for more ORMs, advanced filters).

## Roadmap

- Support for additional ORMs (Ent, Bun).
- Advanced filtering and sorting in `List` (e.g., `--filters true`).
- Redis caching integration.
- GraphQL endpoint generation.
- OpenAPI-based client generation.

## License

MIT