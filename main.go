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
	UseMongoLogger bool
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
	config.ProjectPath = "./" + config.ProjectName

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
	
	if config.Database != "mongodb" {
		config.UseMongoLogger = promptBool(reader, "Would you like to use MongoDB for logging?")
	}

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
	fmt.Println(ColorGreen + "  ✓ MongoDB Logger: " + ColorReset + boolToYesNo(config.UseMongoLogger))
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
	fmt.Println("  [1/5] Copying template files...")
	if err := copyTemplate(config); err != nil {
		return fmt.Errorf("failed to copy template: %w", err)
	}
	
	// Update module paths
	fmt.Println("  [2/5] Updating module paths...")
	if err := updateModulePaths(config); err != nil {
		return fmt.Errorf("failed to update module paths: %w", err)
	}
	
	// Clean up unnecessary files
	fmt.Println("  [3/5] Removing unnecessary files...")
	if err := cleanupFiles(config); err != nil {
		return fmt.Errorf("failed to cleanup: %w", err)
	}
	
	// Generate devcontainer
	fmt.Println("  [4/5] Generating devcontainer configuration...")
	if err := generateDevcontainer(config); err != nil {
		return fmt.Errorf("failed to generate devcontainer: %w", err)
	}
	
	// Update config files
	fmt.Println("  [5/5] Updating configuration files...")
	if err := updateConfigFiles(config); err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}
	
	return nil
}

func copyTemplate(config *ProjectConfig) error {
	return filepath.Walk("template", func(path string, info os.FileInfo, err error) error {
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
	
	// Remove optional service configs
	if !config.UseRedis {
		os.Remove(filepath.Join(config.ProjectPath, "config/redis.go"))
	}
	
	if !config.UseRabbitMQ {
		os.Remove(filepath.Join(config.ProjectPath, "config/rabbitmq.go"))
	}
	
	if !config.UseMongoLogger && config.Database != "mongodb" {
		os.Remove(filepath.Join(config.ProjectPath, "config/mongodb.go"))
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
	services := `version: '3.8'

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
    volumes:
      - ..:/workspace:cached
    command: sleep infinity
    network_mode: service:db
    depends_on:
      - db
`

	// Add database service
	switch config.Database {
	case "mysql":
		services += `      - redis
      - rabbitmq
` + boolToService(config.UseMongoLogger, "      - mongodb\n") + `

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
		services += `      - redis
      - rabbitmq
` + boolToService(config.UseMongoLogger, "      - mongodb\n") + `

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
		services += `      - redis
      - rabbitmq

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

	// Add MongoDB for logging if needed
	if config.UseMongoLogger && config.Database != "mongodb" {
		services += `
  mongodb:
    image: mongo:6
    restart: unless-stopped
    environment:
      MONGO_INITDB_DATABASE: ` + sanitizeName(config.ProjectName) + `_logs
    volumes:
      - mongodb-data:/data/db
    ports:
      - "27017:27017"
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

	if config.UseMongoLogger && config.Database != "mongodb" {
		services += `  mongodb-data:
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
	if config.Database != "mongodb" && !config.UseMongoLogger {
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
