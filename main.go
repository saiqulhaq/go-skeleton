package main

import (
	"bufio"
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

//go:embed template/*
var templateFS embed.FS

const (
	ColorReset  = "\033[0m"
	ColorGreen  = "\033[32m"
	ColorBlue   = "\033[34m"
	ColorYellow = "\033[33m"
	ColorCyan   = "\033[36m"
)

type ProjectConfig struct {
	ProjectName    string
	ProjectPath    string
	ModulePath     string
	Database       string
	UseRedis       bool
	UseRabbitMQ    bool
}

func main() {
	printBanner()
	
	config := collectConfiguration()
	
	printSummary(config)
	
	if !confirm("Create project?") {
		fmt.Println(ColorYellow + "Cancelled." + ColorReset)
		return
	}
	
	if err := createProject(config); err != nil {
		fmt.Printf(ColorYellow+"Error: %v\n"+ColorReset, err)
		os.Exit(1)
	}
	
	printSuccess(config)
}

func printBanner() {
	fmt.Println(ColorCyan + "╔══════════════════════════════════════════════════════════╗" + ColorReset)
	fmt.Println(ColorCyan + "║                                                          ║" + ColorReset)
	fmt.Println(ColorCyan + "║" + ColorYellow + "          🚀 Create Go Skeleton Project 🚀              " + ColorCyan + "║" + ColorReset)
	fmt.Println(ColorCyan + "║                                                          ║" + ColorReset)
	fmt.Println(ColorCyan + "╚══════════════════════════════════════════════════════════╝" + ColorReset)
	fmt.Println()
}

func collectConfiguration() *ProjectConfig {
	reader := bufio.NewReader(os.Stdin)
	config := &ProjectConfig{}

	// Project name
	config.ProjectName = promptString(reader, "What is your project name?", "my-go-api")
	
	// Project path
	defaultPath := "./" + config.ProjectName
	config.ProjectPath = promptString(reader, "Where to create the project?", defaultPath)

	// Module path
	defaultModule := fmt.Sprintf("github.com/yourusername/%s", config.ProjectName)
	config.ModulePath = promptString(reader, "What is your Go module path?", defaultModule)

	// Database
	fmt.Println()
	fmt.Println(ColorBlue + "Which database would you like to use?" + ColorReset)
	fmt.Println("  1) MySQL/MariaDB (recommended)")
	fmt.Println("  2) PostgreSQL")
	fmt.Println("  3) MongoDB")
	
	dbChoice := promptChoice(reader, "Select database", []string{"1", "2", "3"}, "1")
	switch dbChoice {
	case "1":
		config.Database = "mysql"
	case "2":
		config.Database = "postgresql"
	case "3":
		config.Database = "mongodb"
	}

	// Optional services
	config.UseRedis = promptBool(reader, "Would you like to use Redis for caching?")
	config.UseRabbitMQ = promptBool(reader, "Would you like to use RabbitMQ for message queuing?")

	return config
}

func promptString(reader *bufio.Reader, prompt, defaultValue string) string {
	fmt.Print(ColorCyan + "✔ " + prompt)
	if defaultValue != "" {
		fmt.Print(" (" + defaultValue + ")")
	}
	fmt.Print(": " + ColorReset)

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" && defaultValue != "" {
		return defaultValue
	}
	if input == "" {
		return promptString(reader, prompt, defaultValue)
	}
	return input
}

func promptBool(reader *bufio.Reader, prompt string) bool {
	fmt.Print(ColorCyan + "✔ " + prompt + " (y/N): " + ColorReset)
	
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	
	return input == "y" || input == "yes"
}

func promptChoice(reader *bufio.Reader, prompt string, validChoices []string, defaultChoice string) string {
	fmt.Print(ColorCyan + "✔ " + prompt + " (" + defaultChoice + "): " + ColorReset)
	
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	
	if input == "" {
		return defaultChoice
	}
	
	for _, choice := range validChoices {
		if input == choice {
			return input
		}
	}
	
	fmt.Println(ColorYellow + "Invalid choice. Please try again." + ColorReset)
	return promptChoice(reader, prompt, validChoices, defaultChoice)
}

func confirm(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(ColorYellow + "⚠ " + prompt + " (Y/n): " + ColorReset)
	
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	
	return input == "" || input == "y" || input == "yes"
}

func printSummary(config *ProjectConfig) {
	fmt.Println()
	fmt.Println(ColorBlue + "📋 Project Configuration:" + ColorReset)
	fmt.Println(ColorGreen + "  ✓ Project Name: " + ColorReset + config.ProjectName)
	fmt.Println(ColorGreen + "  ✓ Module Path: " + ColorReset + config.ModulePath)
	fmt.Println(ColorGreen + "  ✓ Database: " + ColorReset + config.Database)
	fmt.Println(ColorGreen + "  ✓ Redis: " + ColorReset + boolToYesNo(config.UseRedis))
	fmt.Println(ColorGreen + "  ✓ RabbitMQ: " + ColorReset + boolToYesNo(config.UseRabbitMQ))
	fmt.Println()
}

func boolToYesNo(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}

func createProject(config *ProjectConfig) error {
	fmt.Println(ColorBlue + "🔧 Creating project..." + ColorReset)
	
	// Create project directory
	if err := os.MkdirAll(config.ProjectPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	
	// Copy template files
	fmt.Println("  [1/6] Copying template files...")
	if err := copyTemplate(config); err != nil {
		return fmt.Errorf("failed to copy template: %w", err)
	}
	
	// Create go.mod file
	fmt.Println("  [2/6] Creating go.mod file...")
	if err := createGoMod(config); err != nil {
		return fmt.Errorf("failed to create go.mod: %w", err)
	}
	
	// Update module paths
	fmt.Println("  [3/6] Updating module paths...")
	if err := updateModulePaths(config); err != nil {
		return fmt.Errorf("failed to update module paths: %w", err)
	}
	
	// Clean up unnecessary files
	fmt.Println("  [4/6] Removing unnecessary files...")
	if err := cleanupFiles(config); err != nil {
		return fmt.Errorf("failed to cleanup: %w", err)
	}
	
	// Generate devcontainer
	fmt.Println("  [5/6] Generating devcontainer configuration...")
	if err := generateDevcontainer(config); err != nil {
		return fmt.Errorf("failed to generate devcontainer: %w", err)
	}
	
	// Update config files
	fmt.Println("  [6/6] Updating configuration files...")
	if err := updateConfigFiles(config); err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}
	
	// Update environment files
	if err := updateEnvFiles(config); err != nil {
		return fmt.Errorf("failed to update env files: %w", err)
	}
	
	return nil
}

func copyTemplate(config *ProjectConfig) error {
	err := filepath.Walk("template", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Skip the template root directory
		if path == "template" {
			return nil
		}
		
		// Get relative path
		relPath, err := filepath.Rel("template", path)
		if err != nil {
			return err
		}
		
		destPath := filepath.Join(config.ProjectPath, relPath)
		
		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}
		
		// Copy file
		return copyFile(path, destPath)
	})
	
	if err != nil {
		return err
	}
	
	// Copy .gitignore separately as it's not included in filepath.Walk by default
	gitignoreSrc := "template/.gitignore"
	gitignoreDst := filepath.Join(config.ProjectPath, ".gitignore")
	if _, err := os.Stat(gitignoreSrc); err == nil {
		return copyFile(gitignoreSrc, gitignoreDst)
	}
	
	return nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()
	
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()
	
	_, err = io.Copy(destFile, sourceFile)
	return err
}

func createGoMod(config *ProjectConfig) error {
	goModContent := `module ` + config.ModulePath + `

go 1.24.1

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/bxcodec/faker v2.0.1+incompatible
	github.com/go-co-op/gocron/v2 v2.11.0
	github.com/go-playground/locales v0.14.1
	github.com/go-playground/universal-translator v0.18.1
	github.com/go-playground/validator/v10 v10.14.1
	github.com/gofiber/fiber/v2 v2.52.5
	github.com/gofiber/swagger v1.1.0
	github.com/golang-jwt/jwt/v4 v4.5.2
	github.com/google/uuid v1.6.0
	github.com/joeshaw/envdecode v0.0.0-20200121155833-099f1fc765bd
	github.com/pkg/errors v0.9.1
	github.com/rabbitmq/amqp091-go v1.8.1
	github.com/redis/go-redis/v9 v9.3.0
	github.com/stretchr/testify v1.9.0
	github.com/subosito/gotenv v1.4.2
	github.com/swaggo/swag v1.16.3
	github.com/valyala/fasthttp v1.51.0
	go.mongodb.org/mongo-driver v1.11.7
	go.uber.org/zap v1.27.0
	golang.org/x/crypto v0.36.0
	gorm.io/driver/mysql v1.5.1
	gorm.io/driver/postgres v1.5.9
	gorm.io/gorm v1.25.10
)
`

	goModPath := filepath.Join(config.ProjectPath, "go.mod")
	return os.WriteFile(goModPath, []byte(goModContent), 0644)
}

func updateModulePaths(config *ProjectConfig) error {
	oldModule := "github.com/rahmatrdn/go-skeleton"
	
	return filepath.Walk(config.ProjectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		
		if !strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, ".mod") {
			return nil
		}
		
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		
		newContent := strings.ReplaceAll(string(content), oldModule, config.ModulePath)
		
		return os.WriteFile(path, []byte(newContent), info.Mode())
	})
}

func cleanupFiles(config *ProjectConfig) error {
	// Remove database config files not being used
	dbConfigs := map[string]string{
		"mysql":      filepath.Join(config.ProjectPath, "config/mysql.go"),
		"postgresql": filepath.Join(config.ProjectPath, "config/postgre.go"),
		"mongodb":    filepath.Join(config.ProjectPath, "config/mongodb.go"),
	}
	
	for db, configFile := range dbConfigs {
		if db != config.Database {
			os.Remove(configFile)
			if db == "mysql" || db == "postgresql" {
				os.Remove(filepath.Join(config.ProjectPath, "config/gorm.go"))
			}
		}
	}
	
	// Remove database repository directories not being used
	dbRepos := map[string]string{
		"mysql":      filepath.Join(config.ProjectPath, "internal/repository/mysql"),
		"postgresql": filepath.Join(config.ProjectPath, "internal/repository/mysql"), // PostgreSQL uses the same mysql folder with GORM
		"mongodb":    filepath.Join(config.ProjectPath, "internal/repository/mongodb"),
	}
	
	for db, repoDir := range dbRepos {
		if db != config.Database {
			// For PostgreSQL, don't remove mysql repo since they share it
			if config.Database == "postgresql" && db == "mysql" {
				continue
			}
			if config.Database == "mysql" && db == "postgresql" {
				continue
			}
			os.RemoveAll(repoDir)
		}
	}
	
	// Remove optional service configs
	if !config.UseRedis {
		os.Remove(filepath.Join(config.ProjectPath, "config/redis.go"))
	}
	
	if !config.UseRabbitMQ {
		os.Remove(filepath.Join(config.ProjectPath, "config/rabbitmq.go"))
	}
	
	return nil
}

func generateDevcontainer(config *ProjectConfig) error {
	devcontainerPath := filepath.Join(config.ProjectPath, ".devcontainer")
	
	// Generate docker-compose.yml for devcontainer
	dockerCompose := generateDevcontainerDockerCompose(config)
	dockerComposePath := filepath.Join(devcontainerPath, "docker-compose.yml")
	
	return os.WriteFile(dockerComposePath, []byte(dockerCompose), 0644)
}

func generateDevcontainerDockerCompose(config *ProjectConfig) string {
	// Build depends_on list dynamically
	dependsOn := []string{"db"}
	if config.UseRedis {
		dependsOn = append(dependsOn, "redis")
	}
	if config.UseRabbitMQ {
		dependsOn = append(dependsOn, "rabbitmq")
	}
	
	// Build depends_on YAML
	dependsOnYaml := ""
	for _, dep := range dependsOn {
		dependsOnYaml += "      - " + dep + "\n"
	}
	
	services := `services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
    volumes:
      - ..:/workspace:cached
    command: sleep infinity
    network_mode: service:db
    depends_on:
` + dependsOnYaml + `
`

	// Add database service
	switch config.Database {
	case "mysql":
		services += `
  db:
    image: mysql:8.0
    restart: unless-stopped
    environment:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_DATABASE: ` + sanitizeName(config.ProjectName) + `
      MYSQL_USER: ` + sanitizeName(config.ProjectName) + `
      MYSQL_PASSWORD: ` + sanitizeName(config.ProjectName) + `
    volumes:
      - mysql-data:/var/lib/mysql
    ports:
      - "3306:3306"
`
	case "postgresql":
		services += `
  db:
    image: postgres:15-alpine
    restart: unless-stopped
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: ` + sanitizeName(config.ProjectName) + `
    volumes:
      - postgres-data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
`
	case "mongodb":
		services += `
  db:
    image: mongo:6
    restart: unless-stopped
    environment:
      MONGO_INITDB_DATABASE: ` + sanitizeName(config.ProjectName) + `
    volumes:
      - mongodb-data:/data/db
    ports:
      - "27017:27017"
`
	}

	// Add Redis if needed
	if config.UseRedis {
		services += `
  redis:
    image: redis:7-alpine
    restart: unless-stopped
    command: redis-server --requirepass ""
    volumes:
      - redis-data:/data
    ports:
      - "6379:6379"
`
	}

	// Add RabbitMQ if needed
	if config.UseRabbitMQ {
		services += `
  rabbitmq:
    image: rabbitmq:3-management-alpine
    restart: unless-stopped
    environment:
      RABBITMQ_DEFAULT_USER: guest
      RABBITMQ_DEFAULT_PASS: guest
    volumes:
      - rabbitmq-data:/var/lib/rabbitmq
    ports:
      - "5672:5672"
      - "15672:15672"
`
	}

	// Add volumes section
	services += `
volumes:
`
	switch config.Database {
	case "mysql":
		services += `  mysql-data:
`
	case "postgresql":
		services += `  postgres-data:
`
	case "mongodb":
		services += `  mongodb-data:
`
	}

	if config.UseRedis {
		services += `  redis-data:
`
	}

	if config.UseRabbitMQ {
		services += `  rabbitmq-data:
`
	}

	return services
}

func updateConfigFiles(config *ProjectConfig) error {
	// Update config.go to only include selected options
	configPath := filepath.Join(config.ProjectPath, "config/config.go")
	
	// Read current config
	content, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	
	configStr := string(content)
	
	// Remove unused database options
	if config.Database != "mysql" {
		configStr = removeLines(configStr, "MysqlOption")
	}
	if config.Database != "postgresql" {
		configStr = removeLines(configStr, "PostgreSqlOption")
	}
	if config.Database != "mongodb" {
		configStr = removeLines(configStr, "MongodbOption")
	}
	
	// Remove unused service options
	if !config.UseRedis {
		configStr = removeLines(configStr, "RedisOption")
	}
	if !config.UseRabbitMQ {
		configStr = removeLines(configStr, "RabbitMQOption")
	}
	
	return os.WriteFile(configPath, []byte(configStr), 0644)
}

func updateEnvFiles(config *ProjectConfig) error {
	// Update .env.devcontainer with actual project database name
	envDevcontainerPath := filepath.Join(config.ProjectPath, ".devcontainer/.env.devcontainer")
	
	content, err := os.ReadFile(envDevcontainerPath)
	if err != nil {
		return err
	}
	
	dbName := sanitizeName(config.ProjectName)
	envContent := strings.ReplaceAll(string(content), "PROJECT_DB_NAME", dbName)
	
	return os.WriteFile(envDevcontainerPath, []byte(envContent), 0644)
}

func removeLines(content, pattern string) string {
	lines := strings.Split(content, "\n")
	result := []string{}
	
	for _, line := range lines {
		if !strings.Contains(line, pattern) {
			result = append(result, line)
		}
	}
	
	return strings.Join(result, "\n")
}

func boolToService(condition bool, service string) string {
	if condition {
		return service
	}
	return ""
}

func sanitizeName(name string) string {
	return strings.ReplaceAll(strings.ToLower(name), "-", "_")
}

func printSuccess(config *ProjectConfig) {
	fmt.Println()
	fmt.Println(ColorGreen + "✅ Project created successfully!" + ColorReset)
	fmt.Println()
	fmt.Println(ColorBlue + "📝 Next steps:" + ColorReset)
	fmt.Println()
	fmt.Println("  1. Navigate to your project:")
	fmt.Println(ColorCyan + "     cd " + config.ProjectName + ColorReset)
	fmt.Println()
	fmt.Println("  2. Open in VS Code DevContainer:")
	fmt.Println(ColorCyan + "     code " + config.ProjectName + ColorReset)
	fmt.Println("     Then: Cmd+Shift+P → 'Dev Containers: Reopen in Container'")
	fmt.Println()
	fmt.Println("  3. Or start services locally:")
	fmt.Println(ColorCyan + "     docker-compose up -d" + ColorReset)
	fmt.Println()
	
	if config.Database != "mongodb" {
		fmt.Println("  4. Run database migrations:")
		fmt.Println(ColorCyan + "     make migrate_up" + ColorReset)
		fmt.Println()
	}
	
	fmt.Println("  5. Start the API:")
	fmt.Println(ColorCyan + "     make run" + ColorReset)
	fmt.Println()
	fmt.Println(ColorYellow + "🎉 Happy coding!" + ColorReset)
}
