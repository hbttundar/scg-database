# SCG-Database: A Contract-Driven Database Toolkit for Go

[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Test Coverage](https://img.shields.io/badge/Coverage-42.9%25-red.svg)](coverage.out)

SCG-Database is a modern, contract-driven database toolkit for Go applications. It provides a clean, extensible architecture that promotes separation of concerns and makes your applications database-agnostic through a powerful adapter pattern.

## 🚀 Features

- **Contract-Based Architecture**: Clean interfaces that decouple your business logic from database implementations
- **Multiple Database Support**: Built-in GORM adapter with support for MySQL, PostgreSQL, SQLite
- **Repository Pattern**: Rich repository interface with query building, CRUD operations, and relationships
- **Migration System**: Integrated database migration management with up/down migrations
- **CLI Tools**: Command-line interface for generating models and managing migrations
- **Seeding Support**: Database seeding functionality for development and testing
- **Soft Deletes**: Built-in soft delete support with timestamp management
- **Transaction Support**: Full transaction support with context-aware operations
- **Extensible**: Easy to add custom adapters for other databases or ORMs

## 📦 Installation

```bash
go get github.com/hbttundar/scg-database
```

## 🏗️ Architecture

The library is built around several core contracts:

- **Connection**: Database connection management
- **Repository**: Data access operations
- **Model**: Entity definitions with relationships
- **Migrator**: Database schema management
- **Seeder**: Database seeding operations

## 🚀 Quick Start

### 1. Basic Usage

```go
package main

import (
    "context"
    "log"
    
    "github.com/hbttundar/scg-database/config"
    "github.com/hbttundar/scg-database/db"
)

func main() {
    // Configure database connection
    cfg := config.Config{
        Driver: "gorm:sqlite",
        DSN:    "app.db",
    }
    
    // Connect to database
    conn, err := db.Connect(cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()
    
    // Use the connection...
}
```

### 2. Define Models

```go
package user

import (
    "github.com/hbttundar/scg-database/contract"
)

type User struct {
    contract.BaseModel
    Name  string
    Email string
}

func (u *User) TableName() string {
    return "users"
}

func (u *User) Relationships() map[string]contract.Relationship {
    return map[string]contract.Relationship{
        // Example relationships:
        // "Profile": contract.NewHasOne(&Profile{}, "user_id", "id"),
        // "Orders": contract.NewHasMany(&Order{}, "user_id", "id"),
        // "Roles": contract.NewBelongsToMany(&Role{}, "user_roles"),
    }
}
```

### 3. Repository Operations

```go
// Create repository
userRepo, err := conn.NewRepository(&user.User{})
if err != nil {
    log.Fatal(err)
}

ctx := context.Background()

// Create user
newUser := &user.User{
    Name:  "John Doe",
    Email: "john@example.com",
}
err = userRepo.Create(ctx, newUser)

// Find user
foundUser, err := userRepo.Find(ctx, newUser.ID)

// Query with conditions
users, err := userRepo.Where("name LIKE ?", "John%").
    OrderBy("created_at", "DESC").
    Limit(10).
    Get(ctx)

// Update user
newUser.Name = "John Smith"
err = userRepo.Update(ctx, newUser)

// Soft delete
err = userRepo.Delete(ctx, newUser)
```

## 🛠️ CLI Tools

The package includes a powerful CLI tool for code generation and migration management.

### Generate Models

```bash
go run ./cmd/scg-db make model User
```

This creates a new model file with the proper structure and contracts.

### Migration Management

```bash
# Create a new migration
go run ./cmd/scg-db migrate make create_users_table

# Run pending migrations
go run ./cmd/scg-db migrate up

# Rollback migrations
go run ./cmd/scg-db migrate down

# Fresh migration (drop all tables and re-run)
go run ./cmd/scg-db migrate fresh
```

### Configuration

Create a `config.yaml` file:

```yaml
database:
  default: gorm:sqlite
  connections:
    gorm:sqlite:
      dsn: app.db
    gorm:mysql:
      dsn: "user:password@tcp(localhost:3306)/database?charset=utf8mb4&parseTime=True&loc=Local"
    gorm:postgres:
      dsn: "host=localhost user=username password=password dbname=database port=5432 sslmode=disable"
  paths:
    models: "domain"
    migrations: "database/migrations"
```

## 📚 Advanced Usage

### Transactions

```go
err := conn.Transaction(ctx, func(txConn contract.Connection) error {
    userRepo, _ := txConn.NewRepository(&user.User{})
    
    // All operations within this function are part of the transaction
    err := userRepo.Create(ctx, &user.User{Name: "Alice"})
    if err != nil {
        return err // This will rollback the transaction
    }
    
    // More operations...
    return nil // This will commit the transaction
})
```

### Relationships and Eager Loading

```go
// Load user with related data
users, err := userRepo.With("Profile", "Orders").Get(ctx)
```

### Batch Operations

```go
users := []contract.Model{
    &user.User{Name: "Alice", Email: "alice@example.com"},
    &user.User{Name: "Bob", Email: "bob@example.com"},
}

// Create in batches
err := userRepo.CreateInBatches(ctx, users, 100)
```

### Custom Queries

```go
// Raw SQL queries
results, err := conn.Select(ctx, "SELECT * FROM users WHERE created_at > ?", time.Now().AddDate(0, -1, 0))

// Execute statements
result, err := conn.Statement(ctx, "UPDATE users SET status = ? WHERE last_login < ?", "inactive", cutoffDate)
```

## 🧪 Testing

The library includes comprehensive test coverage and powerful testing utilities designed specifically for microservices that need to test against real databases running in Docker containers.

### Basic Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### DatabaseTestSuite for Microservices

The `testing` package provides a comprehensive testing framework for microservices that need to test against real databases:

```go
package user_test

import (
    "testing"
    "time"
    
    "github.com/hbttundar/scg-database/testing"
    "github.com/stretchr/testify/suite"
)

type UserServiceTestSuite struct {
    *testing.DatabaseTestSuite
    userService *UserService
}

func (suite *UserServiceTestSuite) SetupSuite() {
    // Configure test database (works with Docker containers)
    cfg := testing.DatabaseTestConfig{
        Driver:          "gorm:postgres", // or mysql, sqlite
        DSN:             "postgres://user:pass@localhost:5432/testdb",
        MigrationsPath:  "/migrations",
        CleanupStrategy: testing.CleanupTruncate,
        Timeout:         30 * time.Second,
    }
    
    suite.DatabaseTestSuite = testing.NewDatabaseTestSuite(&cfg)
    suite.DatabaseTestSuite.SetupSuite()
    
    // Initialize your service with the test database connection
    suite.userService = NewUserService(suite.Connection)
}

func (suite *UserServiceTestSuite) TestCreateUser() {
    // Test your service methods
    user, err := suite.userService.CreateUser("John Doe", "john@example.com")
    suite.NoError(err)
    suite.NotNil(user)
    
    // Use built-in assertion helpers
    suite.AssertRecordExists(&User{}, user.ID)
}

func (suite *UserServiceTestSuite) TestUserRepository() {
    // Create repository for direct testing
    userRepo := suite.CreateRepository(&User{})
    
    // Seed test data
    testUser := &User{Name: "Test User", Email: "test@example.com"}
    suite.SeedData(testUser)
    
    // Test repository operations
    found, err := userRepo.Find(context.Background(), testUser.ID)
    suite.NoError(err)
    suite.Equal("Test User", found.(*User).Name)
}

func (suite *UserServiceTestSuite) TestTransactions() {
    // Test transaction handling
    err := suite.ExecuteInTransaction(func(conn contract.Connection) error {
        userRepo, _ := conn.NewRepository(&User{})
        return userRepo.Create(context.Background(), &User{Name: "TX User"})
    })
    suite.NoError(err)
}

func (suite *UserServiceTestSuite) TearDownTest() {
    // Clean up after each test
    suite.TruncateTable("users")
}

func TestUserServiceSuite(t *testing.T) {
    suite.Run(t, new(UserServiceTestSuite))
}
```

### Testing Features

#### Cleanup Strategies

Choose how to clean up test data between tests:

```go
// Truncate tables after each test (fastest)
CleanupStrategy: testing.CleanupTruncate

// Wrap each test in a transaction and rollback (isolated)
CleanupStrategy: testing.CleanupTransaction

// Drop and recreate database (thorough but slow)
CleanupStrategy: testing.CleanupRecreate

// No cleanup (useful for debugging)
CleanupStrategy: testing.CleanupNone
```

#### Database Readiness Checking

The test suite automatically waits for your Docker database to be ready:

```bash
# Set environment variables for database timeout
export TEST_DB_TIMEOUT=60s
export TEST_MIGRATIONS_PATH=/path/to/migrations
```

#### Built-in Assertions

```go
// Assert record exists
suite.AssertRecordExists(&User{}, userID)

// Assert record doesn't exist
suite.AssertRecordNotExists(&User{}, userID)

// Assert table is empty
suite.AssertTableEmpty("users")

// Seed test data
suite.SeedData(&User{Name: "Test"}, &User{Name: "Test2"})

// Truncate specific tables
suite.TruncateTable("users")

// Get raw SQL connection for advanced operations
rawDB := suite.GetRawConnection()
```

#### Docker Integration Example

Use with Docker Compose for integration testing:

```yaml
# docker-compose.test.yml
version: '3.8'
services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: testdb
      POSTGRES_USER: testuser
      POSTGRES_PASSWORD: testpass
    ports:
      - "5432:5432"
    
  test:
    build: .
    depends_on:
      - postgres
    environment:
      TEST_DB_DSN: "postgres://testuser:testpass@postgres:5432/testdb"
    command: go test ./...
```

### Testing Best Practices

1. **Use Real Databases**: Test against the same database type you use in production
2. **Isolate Tests**: Each test should be independent and clean up after itself
3. **Seed Consistently**: Use the `SeedData` method for consistent test data
4. **Test Transactions**: Verify your transaction handling works correctly
5. **Performance Testing**: Use the testing suite for performance benchmarks

## 🔧 Extending the Library

### Custom Adapters

You can create custom adapters by implementing the `contract.DBAdapter` interface:

```go
type MyCustomAdapter struct{}

func (a *MyCustomAdapter) Name() string {
    return "mycustom"
}

func (a *MyCustomAdapter) Connect(cfg config.Config) (contract.Connection, error) {
    // Implement your connection logic
}

// Register your adapter
db.RegisterAdapter("mycustom", &MyCustomAdapter{})
```

## 📁 Project Structure

```
scg-database/
├── adapter/gorm/          # GORM database adapter
├── cmd/scg-db/           # CLI application
├── config/               # Configuration management
├── contract/             # Interface definitions
├── db/                   # Core database functionality
├── example/              # Usage examples
├── migration/            # Migration system
├── seeder/              # Database seeding
└── testing/             # Testing utilities
```

## 🎯 Example Application

Run the complete example to see the library in action:

```bash
cd example
chmod +x example.sh
./example.sh
```

This script demonstrates:
- Model generation via CLI
- Migration creation and execution
- Repository operations
- Error handling

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with [GORM](https://gorm.io/) for the default database adapter
- Uses [golang-migrate](https://github.com/golang-migrate/migrate) for migration management
- CLI powered by [Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/viper)
