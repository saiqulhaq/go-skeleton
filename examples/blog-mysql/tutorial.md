# Blog API with MySQL - Step-by-Step Guide

**Project Goal:** Build a blog API with posts, comments, categories/tags, and markdown support

**What You'll Learn:**
- Go struct composition and embedding
- GORM associations (has_many, belongs_to, many-to-many)
- Nested JSON responses
- Filtering and pagination
- Markdown processing in Go

---

## Phase 1: Project Setup

### ✅ Step 1.1: Generate the Project
```bash
cd /Users/chucky/projects/github.com/saiqulhaq/go-skeleton
go run .
```

**Answer the prompts:**
- Project name: `blog-mysql`
- Module path: `github.com/saiqulhaq/blog-mysql`
- Database: `1` (MySQL)
- Redis: `y`
- RabbitMQ: `n`
- MongoDB logging: `n`

**Expected Result:** New directory `examples/blog-mysql` created with clean MySQL-only structure

---

### ✅ Step 1.2: Initial Setup and Verification
```bash
cd examples/blog-mysql
cp .env.example .env
```

**Edit `.env` file:**
```env
DB_HOST=localhost
DB_PORT=3306
DB_NAME=blog_db
DB_USER=root
DB_PASSWORD=root
```

**Start the database:**
```bash
docker-compose up -d
```

**Verify services are running:**
```bash
docker-compose ps
# Should show mysql container running
```

**Learning Point:** 
- Understanding Docker Compose for local development
- Environment variable configuration
- Database connection setup

---

## Phase 2: Database Schema Design

### ✅ Step 2.1: Create Posts Migration
**Use the Makefile to generate migration files:**

```bash
make migrate create=create_table_posts
```

This creates two files:
- `database/migration/000004_create_table_posts.up.sql`
- `database/migration/000004_create_table_posts.down.sql`

**Edit the `.up.sql` file:**
```sql
CREATE TABLE posts (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    body TEXT NOT NULL,
    markdown_body TEXT NOT NULL,
    author_id BIGINT NOT NULL,
    published_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_author_id (author_id),
    INDEX idx_slug (slug),
    INDEX idx_published_at (published_at),
    
    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

**Edit the `.down.sql` file:**
```sql
DROP TABLE IF EXISTS posts;
```

**Learning Point:**
- Using `make migrate` to generate migration files
- Primary key with AUTO_INCREMENT
- Foreign key relationships
- Indexes for query optimization
- Soft deletes (deleted_at)
- UTF8MB4 for emoji and special characters
- Slug for SEO-friendly URLs

---

### ✅ Step 2.2: Create Comments Migration
```bash
make migrate create=create_table_comments
```

**Edit the `.up.sql` file:**

```sql
CREATE TABLE comments (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    post_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    body TEXT NOT NULL,
    parent_id BIGINT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_post_id (post_id),
    INDEX idx_user_id (user_id),
    INDEX idx_parent_id (parent_id),
    
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES comments(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

**Edit the `.down.sql` file:**
```sql
DROP TABLE IF EXISTS comments;
```

**Learning Point:**
- Self-referencing foreign key (parent_id) for nested comments
- Cascading deletes
- Multiple foreign keys in one table

---

### ✅ Step 2.3: Create Categories Migration
```bash
make migrate create=create_table_categories
```

**Edit the `.up.sql` file:**

```sql
CREATE TABLE categories (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL UNIQUE,
    slug VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_slug (slug)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

**Edit the `.down.sql` file:**
```sql
DROP TABLE IF EXISTS categories;
```

---

### ✅ Step 2.4: Create Tags Migration
```bash
make migrate create=create_table_tags
```

**Edit the `.up.sql` file:**

```sql
CREATE TABLE tags (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL UNIQUE,
    slug VARCHAR(50) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_slug (slug)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

**Edit the `.down.sql` file:**
```sql
DROP TABLE IF EXISTS tags;
```

---

### ✅ Step 2.5: Create Many-to-Many Pivot Tables

**Post Categories (One-to-Many):**

```bash
make migrate create=create_table_post_categories
```

**Edit the `.up.sql` file:**
```sql
CREATE TABLE post_categories (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    post_id BIGINT NOT NULL,
    category_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE KEY unique_post_category (post_id, category_id),
    INDEX idx_post_id (post_id),
    INDEX idx_category_id (category_id),
    
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

**Edit the `.down.sql` file:**
```sql
DROP TABLE IF EXISTS post_categories;
```

**Post Tags (Many-to-Many):**

```bash
make migrate create=create_table_post_tags
```

**Edit the `.up.sql` file:**
```sql
CREATE TABLE post_tags (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    post_id BIGINT NOT NULL,
    tag_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE KEY unique_post_tag (post_id, tag_id),
    INDEX idx_post_id (post_id),
    INDEX idx_tag_id (tag_id),
    
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

**Edit the `.down.sql` file:**
```sql
DROP TABLE IF EXISTS post_tags;
```

**Learning Point:**
- Using `make migrate` for migration file generation
- Many-to-Many relationships require pivot tables
- Composite unique constraints prevent duplicates
- Proper indexing for join queries

---

### ✅ Step 2.6: Run Migrations
```bash
make migrate_up
```

**Verify tables created:**
```bash
# Access MySQL
docker exec -it blog-mysql-mysql-1 mysql -uroot -proot blog_db

# Inside MySQL
SHOW TABLES;
DESCRIBE posts;
DESCRIBE comments;
DESCRIBE categories;
DESCRIBE tags;
DESCRIBE post_categories;
DESCRIBE post_tags;
exit
```

**Learning Point:**
- Database migration workflow
- Verifying schema structure
- Understanding table relationships

---

## Phase 3: Domain Entities (GORM Models)

### ✅ Step 3.0: Install db2struct Tool
We'll use `db2struct` to automatically generate Go structs from our database tables.

```bash
go install github.com/Shelnutt2/db2struct/cmd/db2struct@latest
```

**Important Notes:**
- `db2struct` generates basic structs with GORM tags from existing database tables
- You'll need to manually add GORM associations (foreignKey, many2many, etc.)
- Ensure migrations are run before generating structs
- This saves time on repetitive struct creation while maintaining full control over relationships

### ✅ Step 3.1: Generate Post Entity
Instead of manually creating structs, use db2struct to generate them from the database:

```bash
# Make sure migrations are run first
make migrate_up

# Generate Post struct
db2struct --host localhost \
  --port 3306 \
  --user root \
  --password root \
  --database blog_db \
  --table posts \
  --package entity \
  --gorm \
  --json > internal/repository/mysql/entity/post.go
```

**Then manually add associations** to `internal/repository/mysql/entity/post.go`:

```go
package entity

import (
	"time"
	"gorm.io/gorm"
)

type Post struct {
	ID           int64          `gorm:"column:id;primaryKey"`
	Title        string         `gorm:"column:title;not null"`
	Slug         string         `gorm:"column:slug;uniqueIndex;not null"`
	Body         string         `gorm:"column:body;type:text;not null"`
	MarkdownBody string         `gorm:"column:markdown_body;type:text;not null"`
	AuthorID     int64          `gorm:"column:author_id;not null"`
	PublishedAt  *time.Time     `gorm:"column:published_at"`
	CreatedAt    time.Time      `gorm:"column:created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index"`

	// Associations
	Author     User       `gorm:"foreignKey:AuthorID"`
	Comments   []Comment  `gorm:"foreignKey:PostID"`
	Categories []Category `gorm:"many2many:post_categories;"`
	Tags       []Tag      `gorm:"many2many:post_tags;"`
}

func (Post) TableName() string {
	return "posts"
}
```

**Learning Point:**
- GORM struct tags for column mapping
- Associations: foreignKey, many2many
- Soft deletes with gorm.DeletedAt
- Pointer types for nullable fields (*time.Time)

---

### ✅ Step 3.2: Generate Comment Entity
```bash
db2struct --host localhost \
  --port 3306 \
  --user root \
  --password root \
  --database blog_db \
  --table comments \
  --package entity \
  --gorm \
  --json > internal/repository/mysql/entity/comment.go
```

**Then manually add associations:**

```go
package entity

import (
	"time"
	"gorm.io/gorm"
)

type Comment struct {
	ID        int64          `gorm:"column:id;primaryKey"`
	PostID    int64          `gorm:"column:post_id;not null"`
	UserID    int64          `gorm:"column:user_id;not null"`
	Body      string         `gorm:"column:body;type:text;not null"`
	ParentID  *int64         `gorm:"column:parent_id"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`

	// Associations
	Post    Post      `gorm:"foreignKey:PostID"`
	User    User      `gorm:"foreignKey:UserID"`
	Parent  *Comment  `gorm:"foreignKey:ParentID"`
	Replies []Comment `gorm:"foreignKey:ParentID"`
}

func (Comment) TableName() string {
	return "comments"
}
```

**Learning Point:**
- Self-referencing association (Parent/Replies)
- Nested comments structure
- Has-many and belongs-to relationships

---

### ✅ Step 3.3: Generate Category and Tag Entities
```bash
# Generate Category
db2struct --host localhost \
  --port 3306 \
  --user root \
  --password root \
  --database blog_db \
  --table categories \
  --package entity \
  --gorm \
  --json > internal/repository/mysql/entity/category.go

# Generate Tag
db2struct --host localhost \
  --port 3306 \
  --user root \
  --password root \
  --database blog_db \
  --table tags \
  --package entity \
  --gorm \
  --json > internal/repository/mysql/entity/tag.go
```

**Verify the generated `category.go`:**

```go
package entity

import "time"

type Category struct {
	ID          int64     `gorm:"column:id;primaryKey"`
	Name        string    `gorm:"column:name;uniqueIndex;not null"`
	Slug        string    `gorm:"column:slug;uniqueIndex;not null"`
	Description string    `gorm:"column:description;type:text"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`

	// Associations
	Posts []Post `gorm:"many2many:post_categories;"`
}

func (Category) TableName() string {
	return "categories"
}
```

Create: `internal/repository/mysql/entity/tag.go`

```go
package entity

import "time"

type Tag struct {
	ID        int64     `gorm:"column:id;primaryKey"`
	Name      string    `gorm:"column:name;uniqueIndex;not null"`
	Slug      string    `gorm:"column:slug;uniqueIndex;not null"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`

	// Associations
	Posts []Post `gorm:"many2many:post_tags;"`
}

func (Tag) TableName() string {
	return "tags"
}
```

**Learning Point:**
- Many-to-Many from both sides
- Unique indexes on name and slug
- Bidirectional relationships
- Using db2struct for automatic struct generation
- Manual refinement of associations

### ✅ Step 3.4: Resolve Dependencies
```bash
go mod tidy
```

This will download and add any required modules referenced in the generated code.

---

## Phase 4: Use Case Layer (Business Logic)

### ✅ Step 4.1: Define Post Request/Response DTOs
Create: `internal/usecase/post/entity/post.go`

```go
package entity

import "time"

// Request DTOs
type CreatePostReq struct {
	Title        string   `json:"title" validate:"required,min=3,max=255"`
	Body         string   `json:"body" validate:"required,min=10"`
	MarkdownBody string   `json:"markdown_body" validate:"required"`
	CategoryIDs  []int64  `json:"category_ids" validate:"required,min=1"`
	TagIDs       []int64  `json:"tag_ids"`
	PublishNow   bool     `json:"publish_now"`
}

type UpdatePostReq struct {
	ID           int64    `json:"id" validate:"required"`
	Title        string   `json:"title" validate:"required,min=3,max=255"`
	Body         string   `json:"body" validate:"required,min=10"`
	MarkdownBody string   `json:"markdown_body" validate:"required"`
	CategoryIDs  []int64  `json:"category_ids" validate:"required,min=1"`
	TagIDs       []int64  `json:"tag_ids"`
}

type ListPostsReq struct {
	Page       int     `json:"page" validate:"min=1"`
	Limit      int     `json:"limit" validate:"min=1,max=100"`
	CategoryID *int64  `json:"category_id"`
	TagID      *int64  `json:"tag_id"`
	AuthorID   *int64  `json:"author_id"`
	Search     string  `json:"search"`
}

// Response DTOs
type PostResponse struct {
	ID           int64              `json:"id"`
	Title        string             `json:"title"`
	Slug         string             `json:"slug"`
	Body         string             `json:"body"`
	MarkdownBody string             `json:"markdown_body"`
	Author       AuthorResponse     `json:"author"`
	Categories   []CategoryResponse `json:"categories"`
	Tags         []TagResponse      `json:"tags"`
	CommentsCount int               `json:"comments_count"`
	PublishedAt  *time.Time         `json:"published_at"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
}

type PostListResponse struct {
	ID            int64              `json:"id"`
	Title         string             `json:"title"`
	Slug          string             `json:"slug"`
	Body          string             `json:"body"` // Excerpt only
	Author        AuthorResponse     `json:"author"`
	Categories    []CategoryResponse `json:"categories"`
	Tags          []TagResponse      `json:"tags"`
	CommentsCount int                `json:"comments_count"`
	PublishedAt   *time.Time         `json:"published_at"`
	CreatedAt     time.Time          `json:"created_at"`
}

type AuthorResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CategoryResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type TagResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type PaginatedPostsResponse struct {
	Data       []PostListResponse `json:"data"`
	Pagination PaginationMeta     `json:"pagination"`
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalRows  int64 `json:"total_rows"`
	TotalPages int   `json:"total_pages"`
}
```

**Learning Point:**
- Separation of concerns: DTOs vs Domain entities
- Validation tags for input validation
- Nested response structures
- Pagination metadata
- Different response types (detail vs list)

---

### ✅ Step 4.2: Create Post Repository Interface
Create: `internal/repository/mysql/post.go`

```go
package mysql

import (
	"context"
	"github.com/saiqulhaq/blog-mysql/internal/repository/mysql/entity"
)

type IPostRepository interface {
	Create(ctx context.Context, post *entity.Post) error
	GetByID(ctx context.Context, id int64) (*entity.Post, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Post, error)
	List(ctx context.Context, filters PostFilters) ([]entity.Post, int64, error)
	Update(ctx context.Context, post *entity.Post) error
	Delete(ctx context.Context, id int64) error
	GetByAuthorID(ctx context.Context, authorID int64) ([]entity.Post, error)
}

type PostFilters struct {
	Page       int
	Limit      int
	CategoryID *int64
	TagID      *int64
	AuthorID   *int64
	Search     string
}
```

---

### ✅ Step 4.3: Implement Post Repository
Create: `internal/repository/mysql/post_impl.go`

```go
package mysql

import (
	"context"
	"github.com/saiqulhaq/blog-mysql/internal/repository/mysql/entity"
	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) IPostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(ctx context.Context, post *entity.Post) error {
	return r.db.WithContext(ctx).Create(post).Error
}

func (r *PostRepository) GetByID(ctx context.Context, id int64) (*entity.Post, error) {
	var post entity.Post
	err := r.db.WithContext(ctx).
		Preload("Author").
		Preload("Categories").
		Preload("Tags").
		Preload("Comments").
		First(&post, id).Error
	
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostRepository) GetBySlug(ctx context.Context, slug string) (*entity.Post, error) {
	var post entity.Post
	err := r.db.WithContext(ctx).
		Preload("Author").
		Preload("Categories").
		Preload("Tags").
		Preload("Comments.User").
		Where("slug = ?", slug).
		First(&post).Error
	
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostRepository) List(ctx context.Context, filters PostFilters) ([]entity.Post, int64, error) {
	var posts []entity.Post
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Post{})

	// Apply filters
	if filters.CategoryID != nil {
		query = query.Joins("JOIN post_categories ON posts.id = post_categories.post_id").
			Where("post_categories.category_id = ?", *filters.CategoryID)
	}

	if filters.TagID != nil {
		query = query.Joins("JOIN post_tags ON posts.id = post_tags.post_id").
			Where("post_tags.tag_id = ?", *filters.TagID)
	}

	if filters.AuthorID != nil {
		query = query.Where("author_id = ?", *filters.AuthorID)
	}

	if filters.Search != "" {
		search := "%" + filters.Search + "%"
		query = query.Where("title LIKE ? OR body LIKE ?", search, search)
	}

	// Get total count
	query.Count(&total)

	// Apply pagination
	offset := (filters.Page - 1) * filters.Limit
	err := query.
		Preload("Author").
		Preload("Categories").
		Preload("Tags").
		Offset(offset).
		Limit(filters.Limit).
		Order("created_at DESC").
		Find(&posts).Error

	return posts, total, err
}

func (r *PostRepository) Update(ctx context.Context, post *entity.Post) error {
	return r.db.WithContext(ctx).Save(post).Error
}

func (r *PostRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&entity.Post{}, id).Error
}

func (r *PostRepository) GetByAuthorID(ctx context.Context, authorID int64) ([]entity.Post, error) {
	var posts []entity.Post
	err := r.db.WithContext(ctx).
		Where("author_id = ?", authorID).
		Preload("Categories").
		Preload("Tags").
		Find(&posts).Error
	return posts, err
}
```

**Learning Point:**
- GORM Preload for eager loading
- Join queries for many-to-many filtering
- LIKE queries for search
- Pagination with Offset/Limit
- Soft delete handling (automatic with gorm.DeletedAt)

---

### ✅ Step 4.4: Create Post Usecase
Create: `internal/usecase/post/post_usecase.go`

```go
package post_usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	generalEntity "github.com/saiqulhaq/blog-mysql/entity"
	"github.com/saiqulhaq/blog-mysql/internal/helper"
	"github.com/saiqulhaq/blog-mysql/internal/repository/mysql"
	mentity "github.com/saiqulhaq/blog-mysql/internal/repository/mysql/entity"
	"github.com/saiqulhaq/blog-mysql/internal/usecase/post/entity"
	"github.com/gosimple/slug"
)

type PostUsecase struct {
	postRepo mysql.IPostRepository
}

func NewPostUsecase(postRepo mysql.IPostRepository) IPostUsecase {
	return &PostUsecase{postRepo: postRepo}
}

type IPostUsecase interface {
	Create(ctx context.Context, req entity.CreatePostReq, authorID int64) (*entity.PostResponse, error)
	GetByID(ctx context.Context, id int64) (*entity.PostResponse, error)
	GetBySlug(ctx context.Context, slug string) (*entity.PostResponse, error)
	List(ctx context.Context, req entity.ListPostsReq) (*entity.PaginatedPostsResponse, error)
	Update(ctx context.Context, req entity.UpdatePostReq) error
	Delete(ctx context.Context, id int64) error
}

func (u *PostUsecase) Create(ctx context.Context, req entity.CreatePostReq, authorID int64) (*entity.PostResponse, error) {
	funcName := "PostUsecase.Create"

	// Generate slug from title
	postSlug := slug.Make(req.Title)

	// Prepare post entity
	post := &mentity.Post{
		Title:        req.Title,
		Slug:         postSlug,
		Body:         req.Body,
		MarkdownBody: req.MarkdownBody,
		AuthorID:     authorID,
	}

	// Set published_at if publish_now is true
	if req.PublishNow {
		now := time.Now()
		post.PublishedAt = &now
	}

	// Create post
	if err := u.postRepo.Create(ctx, post); err != nil {
		helper.LogError("postRepo.Create", funcName, err, nil, "")
		return nil, err
	}

	// TODO: Associate categories and tags
	// This requires updating the associations after creation

	// Fetch the created post with associations
	createdPost, err := u.postRepo.GetByID(ctx, post.ID)
	if err != nil {
		helper.LogError("postRepo.GetByID", funcName, err, nil, "")
		return nil, err
	}

	return u.mapToResponse(createdPost), nil
}

func (u *PostUsecase) GetByID(ctx context.Context, id int64) (*entity.PostResponse, error) {
	funcName := "PostUsecase.GetByID"

	post, err := u.postRepo.GetByID(ctx, id)
	if err != nil {
		helper.LogError("postRepo.GetByID", funcName, err, nil, "")
		return nil, err
	}

	return u.mapToResponse(post), nil
}

func (u *PostUsecase) GetBySlug(ctx context.Context, slug string) (*entity.PostResponse, error) {
	funcName := "PostUsecase.GetBySlug"

	post, err := u.postRepo.GetBySlug(ctx, slug)
	if err != nil {
		helper.LogError("postRepo.GetBySlug", funcName, err, nil, "")
		return nil, err
	}

	return u.mapToResponse(post), nil
}

func (u *PostUsecase) List(ctx context.Context, req entity.ListPostsReq) (*entity.PaginatedPostsResponse, error) {
	funcName := "PostUsecase.List"

	// Set defaults
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 10
	}

	filters := mysql.PostFilters{
		Page:       req.Page,
		Limit:      req.Limit,
		CategoryID: req.CategoryID,
		TagID:      req.TagID,
		AuthorID:   req.AuthorID,
		Search:     req.Search,
	}

	posts, total, err := u.postRepo.List(ctx, filters)
	if err != nil {
		helper.LogError("postRepo.List", funcName, err, nil, "")
		return nil, err
	}

	// Map to response
	postResponses := make([]entity.PostListResponse, 0, len(posts))
	for _, post := range posts {
		postResponses = append(postResponses, u.mapToListResponse(&post))
	}

	// Calculate pagination
	totalPages := int(total) / req.Limit
	if int(total)%req.Limit > 0 {
		totalPages++
	}

	return &entity.PaginatedPostsResponse{
		Data: postResponses,
		Pagination: entity.PaginationMeta{
			Page:       req.Page,
			Limit:      req.Limit,
			TotalRows:  total,
			TotalPages: totalPages,
		},
	}, nil
}

func (u *PostUsecase) Update(ctx context.Context, req entity.UpdatePostReq) error {
	funcName := "PostUsecase.Update"

	// Fetch existing post
	post, err := u.postRepo.GetByID(ctx, req.ID)
	if err != nil {
		helper.LogError("postRepo.GetByID", funcName, err, nil, "")
		return err
	}

	// Update fields
	post.Title = req.Title
	post.Slug = slug.Make(req.Title)
	post.Body = req.Body
	post.MarkdownBody = req.MarkdownBody

	// TODO: Update categories and tags associations

	if err := u.postRepo.Update(ctx, post); err != nil {
		helper.LogError("postRepo.Update", funcName, err, nil, "")
		return err
	}

	return nil
}

func (u *PostUsecase) Delete(ctx context.Context, id int64) error {
	funcName := "PostUsecase.Delete"

	if err := u.postRepo.Delete(ctx, id); err != nil {
		helper.LogError("postRepo.Delete", funcName, err, nil, "")
		return err
	}

	return nil
}

// Helper mapping functions
func (u *PostUsecase) mapToResponse(post *mentity.Post) *entity.PostResponse {
	response := &entity.PostResponse{
		ID:           post.ID,
		Title:        post.Title,
		Slug:         post.Slug,
		Body:         post.Body,
		MarkdownBody: post.MarkdownBody,
		Author: entity.AuthorResponse{
			ID:    post.Author.ID,
			Name:  post.Author.Name,
			Email: post.Author.Email,
		},
		PublishedAt:   post.PublishedAt,
		CreatedAt:     post.CreatedAt,
		UpdatedAt:     post.UpdatedAt,
		CommentsCount: len(post.Comments),
	}

	// Map categories
	for _, cat := range post.Categories {
		response.Categories = append(response.Categories, entity.CategoryResponse{
			ID:   cat.ID,
			Name: cat.Name,
			Slug: cat.Slug,
		})
	}

	// Map tags
	for _, tag := range post.Tags {
		response.Tags = append(response.Tags, entity.TagResponse{
			ID:   tag.ID,
			Name: tag.Name,
			Slug: tag.Slug,
		})
	}

	return response
}

func (u *PostUsecase) mapToListResponse(post *mentity.Post) entity.PostListResponse {
	// Create excerpt from body (first 200 chars)
	excerpt := post.Body
	if len(excerpt) > 200 {
		excerpt = excerpt[:200] + "..."
	}

	response := entity.PostListResponse{
		ID:    post.ID,
		Title: post.Title,
		Slug:  post.Slug,
		Body:  excerpt,
		Author: entity.AuthorResponse{
			ID:    post.Author.ID,
			Name:  post.Author.Name,
			Email: post.Author.Email,
		},
		PublishedAt:   post.PublishedAt,
		CreatedAt:     post.CreatedAt,
		CommentsCount: len(post.Comments),
	}

	// Map categories
	for _, cat := range post.Categories {
		response.Categories = append(response.Categories, entity.CategoryResponse{
			ID:   cat.ID,
			Name: cat.Name,
			Slug: cat.Slug,
		})
	}

	// Map tags
	for _, tag := range post.Tags {
		response.Tags = append(response.Tags, entity.TagResponse{
			ID:   tag.ID,
			Name: tag.Name,
			Slug: tag.Slug,
		})
	}

	return response
}
```

**Learning Point:**
- Business logic separation from data access
- Slug generation from titles
- Entity to DTO mapping
- Pagination calculation
- Error handling and logging
- Excerpt generation for list views

---

## Phase 5: HTTP Handlers (API Endpoints)

### ✅ Step 5.1: Create Post Handler
Create: `internal/http/handler/post_handler.go`

```go
package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	generalEntity "github.com/saiqulhaq/blog-mysql/entity"
	"github.com/saiqulhaq/blog-mysql/internal/http/middleware"
	postUsecase "github.com/saiqulhaq/blog-mysql/internal/usecase/post"
	"github.com/saiqulhaq/blog-mysql/internal/usecase/post/entity"
)

type PostHandler struct {
	postUsecase postUsecase.IPostUsecase
}

func NewPostHandler(postUsecase postUsecase.IPostUsecase) *PostHandler {
	return &PostHandler{postUsecase: postUsecase}
}

// @Summary Create a new post
// @Description Create a new blog post with categories and tags
// @Tags Posts
// @Accept json
// @Produce json
// @Param post body entity.CreatePostReq true "Post data"
// @Success 201 {object} entity.PostResponse
// @Failure 400 {object} generalEntity.ResponseFailed
// @Failure 401 {object} generalEntity.ResponseFailed
// @Security BearerAuth
// @Router /posts [post]
func (h *PostHandler) Create(c *fiber.Ctx) error {
	var req entity.CreatePostReq

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(generalEntity.ResponseFailed{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	// Get user ID from JWT token
	userData := c.Locals("userData").(middleware.UserData)

	post, err := h.postUsecase.Create(c.Context(), req, userData.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(generalEntity.ResponseFailed{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to create post",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(generalEntity.ResponseSuccess{
		Code:    fiber.StatusCreated,
		Message: "Post created successfully",
		Data:    post,
	})
}

// @Summary Get post by ID
// @Description Get a single post by its ID
// @Tags Posts
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} entity.PostResponse
// @Failure 404 {object} generalEntity.ResponseFailed
// @Router /posts/{id} [get]
func (h *PostHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(generalEntity.ResponseFailed{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid post ID",
		})
	}

	post, err := h.postUsecase.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(generalEntity.ResponseFailed{
			Code:    fiber.StatusNotFound,
			Message: "Post not found",
			Error:   err.Error(),
		})
	}

	return c.JSON(generalEntity.ResponseSuccess{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    post,
	})
}

// @Summary Get post by slug
// @Description Get a single post by its slug
// @Tags Posts
// @Produce json
// @Param slug path string true "Post slug"
// @Success 200 {object} entity.PostResponse
// @Failure 404 {object} generalEntity.ResponseFailed
// @Router /posts/slug/{slug} [get]
func (h *PostHandler) GetBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")

	post, err := h.postUsecase.GetBySlug(c.Context(), slug)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(generalEntity.ResponseFailed{
			Code:    fiber.StatusNotFound,
			Message: "Post not found",
			Error:   err.Error(),
		})
	}

	return c.JSON(generalEntity.ResponseSuccess{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    post,
	})
}

// @Summary List posts
// @Description Get paginated list of posts with filtering
// @Tags Posts
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param category_id query int false "Filter by category ID"
// @Param tag_id query int false "Filter by tag ID"
// @Param author_id query int false "Filter by author ID"
// @Param search query string false "Search in title and body"
// @Success 200 {object} entity.PaginatedPostsResponse
// @Router /posts [get]
func (h *PostHandler) List(c *fiber.Ctx) error {
	var req entity.ListPostsReq

	// Parse query parameters
	req.Page, _ = strconv.Atoi(c.Query("page", "1"))
	req.Limit, _ = strconv.Atoi(c.Query("limit", "10"))
	req.Search = c.Query("search")

	if categoryID := c.Query("category_id"); categoryID != "" {
		if id, err := strconv.ParseInt(categoryID, 10, 64); err == nil {
			req.CategoryID = &id
		}
	}

	if tagID := c.Query("tag_id"); tagID != "" {
		if id, err := strconv.ParseInt(tagID, 10, 64); err == nil {
			req.TagID = &id
		}
	}

	if authorID := c.Query("author_id"); authorID != "" {
		if id, err := strconv.ParseInt(authorID, 10, 64); err == nil {
			req.AuthorID = &id
		}
	}

	posts, err := h.postUsecase.List(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(generalEntity.ResponseFailed{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to fetch posts",
			Error:   err.Error(),
		})
	}

	return c.JSON(generalEntity.ResponseSuccess{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    posts,
	})
}

// @Summary Update post
// @Description Update an existing post
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param post body entity.UpdatePostReq true "Post data"
// @Success 200 {object} generalEntity.ResponseSuccess
// @Failure 400 {object} generalEntity.ResponseFailed
// @Failure 401 {object} generalEntity.ResponseFailed
// @Security BearerAuth
// @Router /posts/{id} [put]
func (h *PostHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(generalEntity.ResponseFailed{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid post ID",
		})
	}

	var req entity.UpdatePostReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(generalEntity.ResponseFailed{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	req.ID = id

	if err := h.postUsecase.Update(c.Context(), req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(generalEntity.ResponseFailed{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to update post",
			Error:   err.Error(),
		})
	}

	return c.JSON(generalEntity.ResponseSuccess{
		Code:    fiber.StatusOK,
		Message: "Post updated successfully",
	})
}

// @Summary Delete post
// @Description Delete a post by ID
// @Tags Posts
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} generalEntity.ResponseSuccess
// @Failure 401 {object} generalEntity.ResponseFailed
// @Security BearerAuth
// @Router /posts/{id} [delete]
func (h *PostHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(generalEntity.ResponseFailed{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid post ID",
		})
	}

	if err := h.postUsecase.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(generalEntity.ResponseFailed{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to delete post",
			Error:   err.Error(),
		})
	}

	return c.JSON(generalEntity.ResponseSuccess{
		Code:    fiber.StatusOK,
		Message: "Post deleted successfully",
	})
}
```

**Learning Point:**
- Fiber framework request handling
- Query parameter parsing
- Path parameter extraction
- JWT authentication with middleware
- Swagger documentation annotations
- HTTP status codes

---

### ✅ Step 5.2: Register Routes
Update: `cmd/api/main.go`

Add in the route registration section:

```go
// Post routes
postRepo := mysql.NewPostRepository(gormDB)
postUsecase := post_usecase.NewPostUsecase(postRepo)
postHandler := handler.NewPostHandler(postUsecase)

api := app.Group("/api/v1")
posts := api.Group("/posts")
{
	posts.Get("/", postHandler.List)           // GET /api/v1/posts
	posts.Get("/:id", postHandler.GetByID)     // GET /api/v1/posts/123
	posts.Get("/slug/:slug", postHandler.GetBySlug) // GET /api/v1/posts/slug/my-post
	
	// Protected routes (require authentication)
	posts.Post("/", middleware.VerifyToken, postHandler.Create)        // POST /api/v1/posts
	posts.Put("/:id", middleware.VerifyToken, postHandler.Update)      // PUT /api/v1/posts/123
	posts.Delete("/:id", middleware.VerifyToken, postHandler.Delete)   // DELETE /api/v1/posts/123
}
```

**Learning Point:**
- Dependency injection pattern
- Route grouping
- Middleware application
- RESTful API design

---

## Phase 6: Testing with Bruno

### ✅ Step 6.1: Run the Application
```bash
make run
```

**Verify:**
- Server starts on port 7011
- No database connection errors
- Swagger docs available at: http://localhost:7011/apidoc

---

### ✅ Step 6.2: Create Bruno Collection for Posts
Create: `api-client/bruno/Post/folder.bru`

```
meta {
  name: Post
  type: folder
}
```

Create: `api-client/bruno/Post/List Posts.bru`

```
meta {
  name: List Posts
  type: http
  seq: 1
}

get {
  url: {{baseUrl}}/posts?page=1&limit=10
  body: none
  auth: none
}

params:query {
  page: 1
  limit: 10
}
```

Create: `api-client/bruno/Post/Create Post.bru`

```
meta {
  name: Create Post
  type: http
  seq: 2
}

post {
  url: {{baseUrl}}/posts
  body: json
  auth: bearer
}

auth:bearer {
  token: {{token}}
}

body:json {
  {
    "title": "My First Blog Post",
    "body": "<h1>Welcome to My Blog</h1><p>This is the rendered HTML content.</p>",
    "markdown_body": "# Welcome to My Blog\n\nThis is the markdown content.",
    "category_ids": [1],
    "tag_ids": [1, 2],
    "publish_now": true
  }
}
```

Create: `api-client/bruno/Post/Get Post By ID.bru`

```
meta {
  name: Get Post By ID
  type: http
  seq: 3
}

get {
  url: {{baseUrl}}/posts/1
  body: none
  auth: none
}
```

**Learning Point:**
- API testing workflow
- Authentication headers
- JSON request bodies
- Query parameters

---

## Phase 7: Comments Feature (Exercise)

**Now it's your turn!** Apply what you learned to implement Comments CRUD:

### ✅ Step 7.1: Create Comment DTOs
Create: `internal/usecase/comment/entity/comment.go`

**Hint:** Include fields for:
- post_id
- body
- parent_id (for nested comments)

### ✅ Step 7.2: Create Comment Repository
Create: `internal/repository/mysql/comment.go` and `comment_impl.go`

**Methods to implement:**
- `Create`
- `GetByPostID` (with nested comments)
- `Update`
- `Delete`

### ✅ Step 7.3: Create Comment Usecase
Create: `internal/usecase/comment/comment_usecase.go`

**Business logic:**
- Validate post exists before creating comment
- Build nested comment tree structure
- Handle parent_id for replies

### ✅ Step 7.4: Create Comment Handler
Create: `internal/http/handler/comment_handler.go`

**Endpoints:**
- `POST /posts/:postId/comments` - Create comment
- `GET /posts/:postId/comments` - List comments for a post
- `POST /comments/:id/reply` - Reply to a comment
- `DELETE /comments/:id` - Delete comment

---

## Phase 8: Categories & Tags (Exercise)

### ✅ Step 8.1: Implement Category CRUD
Follow the same pattern as Posts:
1. DTOs in `internal/usecase/category/entity/`
2. Repository in `internal/repository/mysql/`
3. Usecase in `internal/usecase/category/`
4. Handler in `internal/http/handler/`

**Endpoints:**
- `GET /categories` - List all
- `POST /categories` - Create
- `GET /categories/:slug/posts` - Get posts by category

### ✅ Step 8.2: Implement Tag CRUD
Same pattern as Categories

---

## Phase 9: Advanced Features

### ✅ Step 9.1: Add Markdown Processing
Install package:
```bash
go get github.com/gomarkdown/markdown
```

Update Post Usecase to convert markdown to HTML:
```go
import (
    "github.com/gomarkdown/markdown"
    "github.com/gomarkdown/markdown/html"
    "github.com/gomarkdown/markdown/parser"
)

func (u *PostUsecase) Create(ctx context.Context, req entity.CreatePostReq, authorID int64) (*entity.PostResponse, error) {
    // Convert markdown to HTML
    md := []byte(req.MarkdownBody)
    extensions := parser.CommonExtensions | parser.AutoHeadingIDs
    parser := parser.NewWithExtensions(extensions)
    htmlRenderer := html.NewRenderer(html.RendererOptions{})
    htmlBody := markdown.ToHTML(md, parser, htmlRenderer)
    
    post := &mentity.Post{
        Title:        req.Title,
        Slug:         slug.Make(req.Title),
        Body:         string(htmlBody), // Rendered HTML
        MarkdownBody: req.MarkdownBody,  // Original markdown
        AuthorID:     authorID,
    }
    
    // ... rest of the code
}
```

**Learning Point:**
- Third-party package integration
- Markdown to HTML conversion
- Storing both raw and rendered content

---

### ✅ Step 9.2: Add Full-Text Search
Update migration to add FULLTEXT index:

```sql
ALTER TABLE posts ADD FULLTEXT INDEX ft_title_body (title, body);
```

Update repository to use MATCH AGAINST:
```go
if filters.Search != "" {
    query = query.Where("MATCH(title, body) AGAINST(? IN BOOLEAN MODE)", filters.Search)
}
```

**Learning Point:**
- MySQL FULLTEXT indexes
- Full-text search queries
- Boolean mode search

---

### ✅ Step 9.3: Add Post Views Counter
Add migration:
```sql
ALTER TABLE posts ADD COLUMN views_count INT DEFAULT 0;
```

Update entity and add method:
```go
func (r *PostRepository) IncrementViews(ctx context.Context, id int64) error {
    return r.db.WithContext(ctx).
        Model(&entity.Post{}).
        Where("id = ?", id).
        UpdateColumn("views_count", gorm.Expr("views_count + 1")).
        Error
}
```

Call in GetByID handler:
```go
go postUsecase.IncrementViews(c.Context(), id) // Async
```

**Learning Point:**
- Atomic updates with SQL expressions
- Async operations with goroutines
- Performance optimization

---

## Phase 10: Testing

### ✅ Step 10.1: Write Unit Tests for Post Usecase
Create: `internal/usecase/post/post_usecase_test.go`

```go
package post_usecase

import (
	"context"
	"testing"
	"time"

	"github.com/saiqulhaq/blog-mysql/internal/repository/mysql/entity"
	"github.com/saiqulhaq/blog-mysql/internal/usecase/post/entity"
	"github.com/saiqulhaq/blog-mysql/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPostUsecase_Create(t *testing.T) {
	mockRepo := new(mocks.PostRepository)
	usecase := NewPostUsecase(mockRepo)

	ctx := context.Background()
	req := postEntity.CreatePostReq{
		Title:        "Test Post",
		Body:         "<p>Test body</p>",
		MarkdownBody: "Test body",
		CategoryIDs:  []int64{1},
		PublishNow:   true,
	}
	authorID := int64(1)

	// Mock expectations
	mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.Post")).
		Return(nil).
		Run(func(args mock.Arguments) {
			post := args.Get(1).(*entity.Post)
			post.ID = 1
		})

	mockRepo.On("GetByID", ctx, int64(1)).
		Return(&entity.Post{
			ID:           1,
			Title:        "Test Post",
			Slug:         "test-post",
			Body:         "<p>Test body</p>",
			MarkdownBody: "Test body",
			AuthorID:     authorID,
			Author: entity.User{
				ID:    1,
				Name:  "Test User",
				Email: "test@example.com",
			},
			PublishedAt: &time.Time{},
		}, nil)

	// Execute
	result, err := usecase.Create(ctx, req, authorID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Post", result.Title)
	assert.Equal(t, "test-post", result.Slug)
	mockRepo.AssertExpectations(t)
}
```

**Learning Point:**
- Unit testing with mocks
- Testify assertions
- Mock expectations and verifications

---

### ✅ Step 10.2: Generate Mocks for Post Repository
```bash
mockery --name=IPostRepository --dir=internal/repository/mysql --output=tests/mocks
```

---

## Phase 11: Documentation

### ✅ Step 11.1: Update Swagger Docs
```bash
# Install swag if not already installed
go install github.com/swaggo/swag/cmd/swag@latest

# Generate swagger docs (use the Makefile command)
make apidoc
```

**Note:** The Makefile uses the correct swag command: `swag init -d cmd/api,internal/http/handler --parseInternal --pd`

**Verify:** Visit http://localhost:7011/apidoc

**Learning Point:**
- API documentation automation
- Swagger/OpenAPI specification
- Developer documentation best practices

---

## Summary & Next Steps

### 🎯 What You've Accomplished:
✅ Built a complete blog API with:
- Posts with CRUD operations
- Comments (nested/threaded)
- Categories and Tags (many-to-many)
- Markdown support
- Pagination and filtering
- Full-text search
- Authentication
- API documentation
- Unit tests

### 📚 Key Concepts Learned:
1. **Go Struct Composition** - Embedding and associations
2. **GORM Associations** - has_many, belongs_to, many2many
3. **Nested JSON Responses** - DTOs and mapping
4. **Pagination** - Offset/Limit pattern
5. **Filtering** - Dynamic query building
6. **Clean Architecture** - Separation of concerns
7. **Repository Pattern** - Data access abstraction
8. **Dependency Injection** - Testable code
9. **Testing** - Unit tests with mocks

### 🚀 Challenge Yourself:
1. Add user profiles with bio and avatar
2. Implement post likes/reactions
3. Add bookmarking/favorites
4. Create a feed algorithm (trending posts)
5. Add image upload for posts
6. Implement post scheduling (publish_at)
7. Add email notifications
8. Create analytics (views, engagement)
9. Add rate limiting
10. Deploy to production!

### 📖 Further Reading:
- [GORM Associations](https://gorm.io/docs/associations.html)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Go Testing](https://golang.org/pkg/testing/)
- [Fiber Framework](https://docs.gofiber.io/)

---

**Congratulations!** 🎉 You've built a production-ready blog API and learned key Go concepts along the way!
