# String Analyzer Service - HNG13 Backend Stage 1

A RESTful API service that analyzes strings and stores their computed properties, including length, palindrome detection, unique character count, word count, SHA-256 hash, and character frequency mapping.

## 🚀 Features

- **String Analysis**: Automatically computes properties for analyzed strings
- **Duplicate Detection**: Returns 409 Conflict for existing strings
- **Flexible Filtering**: Query strings by multiple criteria
- **Natural Language Queries**: Filter strings using plain English queries
- **Persistent Storage**: PostgreSQL database with GORM ORM
- **RESTful Design**: Clean API endpoints following REST conventions

## 📋 Prerequisites

- Go 1.25.1 or higher
- PostgreSQL 12+ database
- Git

## 🛠️ Tech Stack

- **Language**: Go (Golang)
- **Framework**: Gin Web Framework
- **Database**: PostgreSQL
- **ORM**: GORM
- **Configuration**: godotenv for environment variables

## 📦 Dependencies

```bash
github.com/gin-gonic/gin          # Web framework
github.com/joho/godotenv          # Environment variable loader
gorm.io/gorm                      # ORM library
gorm.io/driver/postgres           # PostgreSQL driver for GORM
```

## ⚙️ Installation & Setup

### 1. Clone the Repository

```bash
git clone https://github.com/whotterre/HNG13_Backend_Task1.git
cd HNG13_Backend_Task1/task_one
```

### 2. Install Dependencies

```bash
go mod download
go mod tidy
```

### 3. Configure Environment Variables

Create a `.env` file in the project root:

```env
DATABASE_URL=postgres://username:password@localhost:5432/strings_db?sslmode=disable
PORT=4000
```

**Note**: Replace `username`, `password`, and `strings_db` with your PostgreSQL credentials and database name.

### 4. Create the Database

```sql
CREATE DATABASE strings_db;
```

The application will automatically create the required tables on startup using GORM migrations.

### 5. Run the Application

```bash
# From the cmd directory
cd cmd
go run .

# Or from the project root
go run cmd/main.go
```

The server will start on `http://localhost:4000` (or the port specified in your `.env` file).

## 🔌 API Endpoints

### 1. Create/Analyze String

**POST** `/strings`

Analyzes and stores a new string.

**Request Body**:
```json
{
  "value": "string to analyze"
}
```

**Success Response (201 Created)**:
```json
{
  "id": "sha256_hash_value",
  "value": "string to analyze",
  "properties": {
    "length": 16,
    "is_palindrome": false,
    "unique_characters": 12,
    "word_count": 3,
    "sha256_hash": "abc123...",
    "character_frequency_map": {
      "s": 2,
      "t": 3,
      "r": 2
    }
  },
  "created_at": "2025-10-21T10:00:00Z"
}
```

**Error Responses**:
- `400 Bad Request`: Missing or empty "value" field
- `409 Conflict`: String already exists in the system
- `422 Unprocessable Entity`: Invalid data type for "value" (must be string)

### 2. Get Specific String

**GET** `/strings/{string_value}`

Retrieves a previously analyzed string by its value.

**Example**: `GET /strings/hello`

**Success Response (200 OK)**:
```json
{
  "id": "sha256_hash_value",
  "value": "hello",
  "properties": {
    "length": 5,
    "is_palindrome": false,
    "unique_characters": 4,
    "word_count": 1,
    "sha256_hash": "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
    "character_frequency_map": {
      "h": 1,
      "e": 1,
      "l": 2,
      "o": 1
    }
  },
  "created_at": "2025-10-21T10:00:00Z"
}
```

**Error Responses**:
- `404 Not Found`: String does not exist in the system

### 3. Get All Strings with Filtering

**GET** `/strings?is_palindrome=true&min_length=5&max_length=20&word_count=2&contains_character=a`

Retrieves all strings with optional filtering.

**Query Parameters**:
- `is_palindrome`: boolean (true/false)
- `min_length`: integer (minimum string length)
- `max_length`: integer (maximum string length)
- `word_count`: integer (exact word count)
- `contains_character`: string (single character to search for)

**Success Response (200 OK)**:
```json
{
  "data": [
    {
      "id": "hash1",
      "value": "string1",
      "properties": { /* ... */ },
      "created_at": "2025-10-21T10:00:00Z"
    }
  ],
  "count": 15,
  "filters_applied": {
    "is_palindrome": true,
    "min_length": 5,
    "max_length": 20,
    "word_count": 2,
    "contains_character": "a"
  }
}
```

**Error Responses**:
- `400 Bad Request`: Invalid query parameter values or types

### 4. Natural Language Filtering

**GET** `/strings/filter-by-natural-language?query=all%20single%20word%20palindromic%20strings`

Filter strings using natural language queries.

**Supported Queries**:
- "all single word palindromic strings" → `word_count=1, is_palindrome=true`
- "strings longer than 10 characters" → `min_length=11`
- "palindromic strings that contain the first vowel" → `is_palindrome=true, contains_character=a`
- "strings containing the letter z" → `contains_character=z`

**Success Response (200 OK)**:
```json
{
  "data": [ /* array of matching strings */ ],
  "count": 3,
  "interpreted_query": {
    "original": "all single word palindromic strings",
    "parsed_filters": {
      "word_count": 1,
      "is_palindrome": true
    }
  }
}
```

**Error Responses**:
- `400 Bad Request`: Unable to parse natural language query
- `422 Unprocessable Entity`: Query parsed but resulted in conflicting filters

### 5. Delete String

**DELETE** `/strings/{string_value}`

Deletes a string from the system.

**Example**: `DELETE /strings/hello`

**Success Response (204 No Content)**: Empty response body

**Error Responses**:
- `404 Not Found`: String does not exist in the system

## 📂 Project Structure

```
task_one/
├── cmd/
│   └── main.go              # Application entry point
├── config/
│   └── config.go            # Configuration loader (env vars)
├── dto/
│   └── dto.go               # Data Transfer Objects
├── handlers/
│   └── handlers.go          # HTTP request handlers
├── initializers/
│   └── connectDB.go         # Database connection & migration
├── models/
│   └── string.go            # Database models
├── repository/
│   └── repository.go        # Data access layer
├── routes/
│   └── routes.go            # Route definitions
├── services/
│   ├── services.go          # Business logic
│   └── string_helpers.go    # String analysis helper functions
├── .env                     # Environment variables (not committed)
├── go.mod                   # Go module dependencies
├── go.sum                   # Dependency checksums
└── README.md                # This file
```

## 🧪 Testing the API

### Using cURL

```bash
# Create a string
curl -X POST http://localhost:4000/strings \
  -H "Content-Type: application/json" \
  -d '{"value":"hello world"}'

# Get a string
curl http://localhost:4000/strings/hello%20world

# Get all strings with filters
curl "http://localhost:4000/strings?is_palindrome=false&min_length=5"

# Delete a string
curl -X DELETE http://localhost:4000/strings/hello%20world
```

### Using Postman/Thunder Client

Import the following collection or create requests manually:

1. **POST** `http://localhost:4000/strings` with JSON body `{"value": "test"}`
2. **GET** `http://localhost:4000/strings/test`
3. **GET** `http://localhost:4000/strings?is_palindrome=true`
4. **DELETE** `http://localhost:4000/strings/test`

## 🔧 Configuration

### Environment Variables

| Variable       | Description                          | Default                                                    |
|----------------|--------------------------------------|------------------------------------------------------------|
| `DATABASE_URL` | PostgreSQL connection string         | `postgres://postgres:password@localhost:5432/strings?sslmode=disable` |
| `PORT`         | Server port                          | `4000`                                                     |

## 🚢 Deployment

### Prerequisites
- A PostgreSQL database (e.g., from Railway, Supabase, or AWS RDS)
- A hosting platform (Railway, Heroku, AWS, etc.)

### Steps

1. **Set Environment Variables** on your hosting platform:
   ```
   DATABASE_URL=postgres://user:pass@host:5432/dbname?sslmode=require
   PORT=4000
   ```

2. **Deploy**:
   - Railway: Connect your GitHub repo and deploy
   - Heroku: `git push heroku main`
   - AWS: Use Elastic Beanstalk or ECS

3. **Verify**: Test your deployed API endpoints

## 🐛 Troubleshooting

### Database Connection Issues

```bash
# Ensure PostgreSQL is running
sudo service postgresql status

# Check credentials in .env file
cat .env

# Test connection manually
psql -h localhost -U your_username -d strings_db
```

### Port Already in Use

```bash
# Find and kill process using port 4000
# Linux/Mac:
lsof -ti:4000 | xargs kill -9

# Windows (PowerShell):
Get-Process -Id (Get-NetTCPConnection -LocalPort 4000).OwningProcess | Stop-Process
```


### Code Organization

- **Handlers**: HTTP request/response handling
- **Services**: Business logic and string analysis
- **Repository**: Database queries and data persistence
- **Models**: Database schema definitions
- **DTOs**: Request/response structures

