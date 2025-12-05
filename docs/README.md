# Swagger Documentation

## Quick Access

Once the server is running, access Swagger UI at:
```
http://localhost:8080/swagger/index.html
```

## Available Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/url/shorten` | Create a shortened URL |
| GET | `/url/{alias}` | Redirect to original URL |
| GET | `/url/links/{alias}` | Get URL information |
| GET | `/url/links` | List all URLs |

## Regenerating Documentation

After making changes to API endpoints or annotations:
```bash
$HOME/go/bin/swag init -g cmd/api/main.go
```

## Example: Create Short URL

**Request:**
```bash
curl -X POST http://localhost:8080/url/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.google.com"}'
```

**Response:**
```json
{
  "short_url": "http://localhost:8080/abc123",
  "original_url": "https://www.google.com"
}
```
