README 
URL Shortener
URL Shortener is a microservice for shortening URLs, built with Go. The application allows users to create short versions of long URLs and provides usage statistics for these links.

Table of Contents
Features
Technologies
Installation
Configuration
Running
API Documentation
Usage Examples
Testing
Features
Create shortened links for long URLs
Redirect to the original URL using a short link
Retrieve usage statistics (click counts, creation date, etc.)
Rate limiting for API requests
Request logging and error handling
Technologies
Go
Chi — HTTP framework for routing
PostgreSQL — database for URL storage
Docker — containerization
Swagger — API documentation
Zap — structured logging
Installation
Requirements
Go >= 1.17
Docker (for containerization)
PostgreSQL (if not using Docker)
Clone the Repository
bash
Копіювати код
git clone https://github.com/username/url-shortener.git
cd url-shortener
Build and Run
If you want to build and run the project without Docker:

bash
Копіювати код
go mod download
go build -o url-shortener ./cmd/server
./url-shortener
Run with Docker
Create a .env file based on the example .env.example and set the environment variables.

Start the Docker Compose setup:

bash
Копіювати код
docker-compose up --build
Configuration
Configuration is done through the .env file. Key parameters include:

DB_DRIVER — database driver (default: postgres)
DATABASE_URL — database connection URL
SERVER_PORT — port to run the server on
RATE_LIMIT — request rate limit
API Documentation
The API supports the following actions:

POST /shorten — create a shortened link
GET /{shortURL} — redirect to the original URL
GET /{shortURL}/stats — get statistics for the shortened link
Example Request
Creating a Shortened Link
http
Копіювати код
POST /shorten
Content-Type: application/json

{
  "original_url": "https://example.com/very-long-url"
}
Example Response
json
Копіювати код
{
  "short_url": "http://localhost:8080/abc123",
  "original_url": "https://example.com/very-long-url"
}
Usage Examples
Using curl
Create a Shortened Link:

bash
Копіювати код
curl -X POST http://localhost:8080/shorten -d '{"original_url": "https://example.com/very-long-url"}' -H "Content-Type: application/json"
Redirect to the Original URL:

bash
Копіювати код
curl -L http://localhost:8080/abc123
Get Statistics:

bash
Копіювати код
curl http://localhost:8080/abc123/stats
Testing
The application uses testify for testing.

To run the tests, use:

bash
Копіювати код
go test ./...
