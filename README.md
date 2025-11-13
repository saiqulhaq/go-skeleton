
# Go Skeleton - Project Generator

A modern Go REST API project generator with Clean Architecture, similar to `create-next-app` for Next.js.

## 🚀 Quick Start

```bash
# Create a new project
go run github.com/saiqulhaq/go-skeleton/create-go-skeleton@latest

# Or clone and run locally
git clone https://github.com/saiqulhaq/go-skeleton.git
cd go-skeleton
go run .
```

That's it! Answer a few questions and your project is ready in seconds.

## ✨ Features

- 🎯 **Interactive Generator** - Choose only what you need
- 🗄️ **Multiple Databases** - MySQL, PostgreSQL, or MongoDB
- 🐳 **DevContainer Ready** - VS Code DevContainer with your selected services
- ⚡ **Fast Setup** - Ready-to-code project in ~30 seconds
- 🧹 **Clean Output** - No unused files or configurations
- 🏗️ **Clean Architecture** - Best practices built-in
- 🔐 **JWT Authentication** - RS512 algorithm
- 📚 **API Documentation** - Swagger/OpenAPI
- 🧪 **Testing Ready** - Unit tests with mocks
- 🔄 **Background Jobs** - Optional RabbitMQ support
- 💾 **Caching** - Optional Redis support

## 📋 What You Get

A complete Go project with:
- REST API with Fiber framework
- Clean Architecture structure (Entity → Repository → Usecase → Handler)
- Your chosen database (MySQL, PostgreSQL, or MongoDB)
- Optional Redis for caching
- Optional RabbitMQ for message queuing
- Optional MongoDB for centralized logging
- DevContainer configuration with only your selected services
- Database migrations (for SQL databases)
- JWT authentication
- Swagger API documentation
- Unit testing setup with mocks
- Example CRUD operations (Users, TodoList)

## 🎬 Example Session

```
╔══════════════════════════════════════════════════════════╗
║          🚀 Create Go Skeleton Project 🚀              ║
╚══════════════════════════════════════════════════════════╝

✔ What is your project name? (my-go-api): ecommerce-api
✔ What is your Go module path? github.com/mycompany/ecommerce-api

Which database would you like to use?
  1) MySQL/MariaDB (recommended)
  2) PostgreSQL
  3) MongoDB
✔ Select database (1): 1

✔ Would you like to use Redis for caching? (y/N): y
✔ Would you like to use RabbitMQ for message queuing? (y/N): n
✔ Would you like to use MongoDB for logging? (y/N): n

📋 Project Configuration:
  ✓ Project Name: ecommerce-api
  ✓ Database: mysql
  ✓ Redis: Yes
  ✓ RabbitMQ: No

⚠ Create project? (Y/n): y

🔧 Creating project...
  [1/5] Copying template files...
  [2/5] Updating module paths...
  [3/5] Removing unnecessary files...
  [4/5] Generating devcontainer configuration...
  [5/5] Updating configuration files...

✅ Project created successfully!
```

## 📖 Documentation

- **[Generator README](create-go-skeleton/README.md)** - Complete generator documentation
- **[Template Structure](template/)** - The base project template
- **[Cursor Rules](.cursor/general.mdc)** - Code patterns and conventions

## 🏗️ Architecture

The generated project follows Clean Architecture principles:

```
┌─────────────────────────────────────┐
│          HTTP Request               │
└──────────────┬──────────────────────┘
               │
        ┌──────▼────────┐
        │    Handler    │  ← HTTP layer, routing, validation
        └──────┬────────┘
               │
        ┌──────▼────────┐
        │    Usecase    │  ← Business logic
        └──────┬────────┘
               │
        ┌──────▼────────┐
        │  Repository   │  ← Data access
        └──────┬────────┘
               │
        ┌──────▼────────┐
        │   Database    │  ← MySQL/PostgreSQL/MongoDB
        └───────────────┘
```

## 🎯 Why Use This?

### Before (Manual Setup)
1. Clone full template with ALL databases
2. Manually delete unused database code
3. Manually remove unused service configurations
4. Fix broken imports
5. Update docker-compose.yml
6. Update .env
7. Clean up config files
8. Time: **30-60 minutes** ⏰

### After (Generator)
1. Run the generator
2. Answer 6 questions
3. Time: **2 minutes** ⚡

**Time saved: 95%**

## 🛠️ Generated Project Structure

```
your-project/
├── .devcontainer/          # VS Code DevContainer (only selected services!)
├── cmd/
│   ├── api/                # API server entry point
│   ├── worker/             # Background workers
│   └── scheduler/          # Scheduled jobs
├── config/                 # Only configs for selected services
├── internal/
│   ├── http/handler/       # HTTP handlers
│   ├── repository/         # Only your chosen database
│   └── usecase/            # Business logic
├── entity/                 # Domain entities
├── database/migration/     # Database migrations (SQL only)
├── tests/                  # Test utilities and mocks
├── .env.example
├── docker-compose.yaml     # Only selected services
├── go.mod                  # Your module path
├── Makefile
└── README.md
```

## 🔧 Development Workflow

After generating your project:

1. **Open in DevContainer** (recommended):
   ```bash
   code your-project
   # Then: Cmd+Shift+P → "Dev Containers: Reopen in Container"
   ```

2. **Or start services locally**:
   ```bash
   cd your-project
   docker-compose up -d
   ```

3. **Run migrations** (if SQL database):
   ```bash
   make migrate_up
   ```

4. **Start the API**:
   ```bash
   make run
   ```

5. **View API docs**:
   ```
   http://localhost:7011/apidoc
   ```

## 📚 Template Information

### Principles
- Reusable and Maintainable Code
- Decoupled Code
- Scalable Development

### Included Features
- REST API with Fiber framework
- Clean Architecture
- Swagger API documentation
- Worker/Consumer Queue (optional RabbitMQ)
- Unit Testing (Testify + Mockery)
- JWT Authentication (RS512)
- Structured Logging (Zap)
- Database Migrations

## 👥 Contact
| Name                   | Email                        | Role    |
| ---------------------- | ---------------------------- | ------- |
| Rahmat Ramadhan Putra  | rahmatrdn.dev@gmail.com     | Creator |

## 🤝 Contributing

Feel free to contribute to this repository!

---

## Development Guide (For Template Contributors)

### Prerequisite
- Git (See [Git Installation](https://git-scm.com/downloads))
- Go 1.24+ (See [Golang Installation](https://golang.org/doc/install))
- MySQL / MariaDB / PostgreSQL (Download via Docker or Other sources)
- Mockery (Optional) (See [Mockery Installation](https://github.com/vektra/mockery))
- Go Migrate CLI (Optional) (See [Migrate CLI Installation](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate))
- Redis (Optional based on your requirement) (See [Redis Installation](https://redis.io/docs/getting-started/installation/) or use in Docker)
- RabbitMQ (Optional based on your requirement) (See [RabbitMQ Installation](https://www.rabbitmq.com/download.html) or use in Docker)

#### Windows OS (for a better development experience)

*   Install [Make](https://www.gnu.org/software/make/) (See [Make Installation](https://leangaurav.medium.com/how-to-setup-install-gnu-make-on-windows-324480f1da69)).


### Installation
1. Clone this repo
```sh
git clone https://github.com/rahmatrdn/go-skeleton.git
```
2. Copy `example.env` to `.env`
```sh
cp .env.example .env
```
3. Adjust the `.env` file according to the configuration in your local environment, such as the database or other settings 
4. Create a MySQL database with the name `go_skeleton`
5. Run database migration or Manually run in you SQL Client
```sh
make migrate_up
```
6. Generate `private_key.pem` and `public_key.pem`. You can generate them using an [Online RSA Generator](https://travistidwell.com/jsencrypt/demo/) or other tools. Place the files in the project's root folder.
7. Start the API Service
```sh
go run cmd/api/main.go
```
8. Start the Worker Service (if needed)
```sh
go run cmd/worker/main.go
```

### Api Documentation
For API docs, we are using [Swagger](https://swagger.io/) with [Swag](https://github.com/swaggo/swag) Generator
- Install Swag
```sh
go install github.com/swaggo/swag/cmd/swag@latest
```
- Generate apidoc
```sh
make apidoc
```
- Start API documentations
```sh
go run cmd/api/main.go
```
- Access API Documentation with  browser http://localhost:PORT/apidoc



### Unit test
*tips: if you use `VS Code` as your code editor, you can install extension `golang.go` and follow tutorial [showing code coverage after saving your code](https://dev.to/vuong/golang-in-vscode-show-code-coverage-of-after-saving-test-8g0) to help you create unit test*

- Use [Mockery](https://github.com/vektra/mockery) to generate mock class(es)
```sh
make mock d=DependencyClassName
```
- Run unit test with command below or You can run test per function using Vscode!
```sh
make test
```


### Running In Docker
- Docker Build for API
```sh
docker build -t go-skeleton-api:1.0.1 -f ./deploy/docker/api/Dockerfile .
```
- Docker Build for Worker
```sh
docker build -t go-skeleton-worker:1.0.1 -f ./deploy/docker/worker/Dockerfile .
```
- Run docker compose for API and Workers
```sh
docker-compose -f docker-compose.yaml up -d
```


## Contributing
- Create a new branch with a descriptive name that reflects the changes and switch to the new branch. Use the prefix `feature/` for new features or `fix/` for bug fixes.
```sh
git checkout -b <prefix>/branch-name
```
- Make your change(s) and make the test(s)
- Commit and push your change to upstream repository
```sh
git commit -m "[Type] a meaningful commit message"
git push origin branch-name
```
- Open Merge Request in Repository (Reviewer Check Contact Info)
- Merge Request will be merged only if review phase is passed.

## More Details Information
Contact Creator!
