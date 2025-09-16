# Chirpy - A Twitter-like Social Media API

Chirpy is a robust REST API backend for a social media platform similar to Twitter, built with Go and PostgreSQL. Users can create accounts, post short messages (chirps), follow others, and manage their profiles with premium features.

> **Note**: This project was built as part of a guided learning experience from [boot.dev](https://boot.dev), following their Go backend development course.

## 🚀 Features

### Core Functionality
- **User Management**: Registration, authentication, profile updates
- **Chirp System**: Create, read, update, delete short messages (≤140 characters)
- **Authentication**: JWT-based auth with refresh tokens
- **Content Filtering**: Automatic profanity filtering
- **Premium Features**: User upgrades via webhook integration

### API Features
- **Sorting**: Chirps can be sorted by creation date (ascending/descending)
- **Filtering**: Get chirps by specific authors
- **Pagination**: Efficient data retrieval
- **Security**: Password hashing, JWT validation, API key protection
- **Metrics**: Built-in request tracking and admin dashboard

## 🛠 Tech Stack

- **Language**: Go 1.24
- **Database**: PostgreSQL
- **Database Migration**: Goose
- **Query Builder**: SQLC
- **Authentication**: JWT (golang-jwt/jwt)
- **Password Hashing**: bcrypt
- **Environment Management**: godotenv

## 📁 Project Structure

```
Server-GO/
├── main.go              # Application entry point and routing
├── api.go               # API handlers and business logic
├── handlers.go          # Basic HTTP handlers
├── helpers.go           # Utility functions (text cleaning)
├── json.go              # JSON response utilities
├── .env                 # Environment variables
├── go.mod               # Go module dependencies
├── sqlc.yaml            # SQLC configuration
├── assets/              # Static assets
├── internal/
│   ├── auth/            # Authentication utilities
│   │   ├── hash.go      # Password hashing
│   │   ├── jwt.go       # JWT token management
│   │   ├── refresh_token.go # Refresh token handling
│   │   └── api_key.go   # API key validation
│   └── database/        # Generated database models and queries
└── sql/
    ├── schema/          # Database migrations
    └── queries/         # SQL queries for SQLC
```

## 🔧 Installation & Setup

### Prerequisites
- Go 1.24+
- PostgreSQL
- Goose (migration tool)
- SQLC (query generator)

### Setup Steps

1. **Clone the repository**
   ```bash
   git clone https://github.com/max-durnea/Chirpy.git
   cd Chirpy
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   Create a `.env` file:
   ```env
   DB_URL=postgres://username:password@localhost:5432/chirpy?sslmode=disable
   PLATFORM=dev
   SECRET=your-jwt-secret-key
   POLKA_KEY=your-webhook-api-key
   ```

4. **Create database**
   ```bash
   createdb chirpy
   ```

5. **Run migrations**
   ```bash
   goose -dir sql/schema postgres "$DB_URL" up
   ```

6. **Generate database code**
   ```bash
   sqlc generate
   ```

7. **Run the server**
   ```bash
   go run .
   ```

The server will start on `http://localhost:8080`

## 📚 API Documentation

### Authentication Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/users` | Register a new user |
| `POST` | `/api/login` | Login user and get tokens |
| `POST` | `/api/refresh` | Refresh access token |
| `POST` | `/api/revoke` | Revoke refresh token |
| `PUT` | `/api/users` | Update user profile |

### Chirp Endpoints

| Method | Endpoint | Description | Query Parameters |
|--------|----------|-------------|------------------|
| `GET` | `/api/chirps` | Get all chirps | `author_id`, `sort` (asc/desc) |
| `GET` | `/api/chirps/{id}` | Get specific chirp | - |
| `POST` | `/api/chirps` | Create new chirp | - |
| `DELETE` | `/api/chirps/{id}` | Delete chirp | - |

### Admin Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/admin/metrics` | View request metrics |
| `POST` | `/admin/reset` | Reset database (dev only) |

### Utility Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/healthz` | Health check |
| `POST` | `/api/validate_chirp` | Validate chirp content |
| `POST` | `/api/polka/webhooks` | Webhook for user upgrades |

## 🔐 Authentication

Chirpy uses JWT-based authentication with refresh tokens:

1. **Register/Login**: Get access token (1 hour) and refresh token (60 days)
2. **Access Protected Routes**: Include `Authorization: Bearer <token>` header
3. **Token Refresh**: Use refresh token to get new access token
4. **Token Revocation**: Revoke refresh tokens for logout

## 📝 API Examples

### Register User
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "securepassword"}'
```

### Create Chirp
```bash
curl -X POST http://localhost:8080/api/chirps \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{"body": "Hello, Chirpy world!"}'
```

### Get Chirps (Sorted)
```bash
# Get all chirps, newest first
curl "http://localhost:8080/api/chirps?sort=desc"

# Get chirps by specific author
curl "http://localhost:8080/api/chirps?author_id=<user-uuid>&sort=asc"
```

## 🗄️ Database Schema

### Users Table
- `id` (UUID, Primary Key)
- `created_at` (Timestamp)
- `updated_at` (Timestamp)
- `email` (Text, Unique)
- `hashed_password` (Text)
- `is_chirpy_red` (Boolean) - Premium status

### Chirps Table
- `id` (UUID, Primary Key)
- `created_at` (Timestamp)
- `updated_at` (Timestamp)
- `body` (Text, ≤140 chars)
- `user_id` (UUID, Foreign Key)

### Refresh Tokens Table
- `token` (Text, Primary Key)
- `created_at` (Timestamp)
- `updated_at` (Timestamp)
- `user_id` (UUID, Foreign Key)
- `expires_at` (Timestamp)
- `revoked_at` (Timestamp, Nullable)

## 🔧 Development

### Running Tests
```bash
go test ./...
```

### Database Migrations
```bash
# Create new migration
goose -dir sql/schema create migration_name sql

# Run migrations
goose -dir sql/schema postgres "$DB_URL" up

# Rollback
goose -dir sql/schema postgres "$DB_URL" down
```

### Regenerate Database Code
```bash
sqlc generate
```

## 🌟 Features in Detail

### Content Filtering
Chirpy automatically censors inappropriate words:
- "kerfuffle" → "****"
- "sharbert" → "****"  
- "fornax" → "****"

### Premium Features
Users can be upgraded to "Chirpy Red" status via webhook integration with external payment systems.

### Request Metrics
Built-in metrics tracking for monitoring API usage and performance.

## 🚀 Deployment

The application is production-ready and can be deployed to any platform supporting Go applications:

- **Docker**: Containerize with multi-stage builds
- **Cloud Platforms**: Deploy to AWS, GCP, Heroku, etc.
- **Environment**: Set `PLATFORM=production` for production mode

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 👨‍💻 Author

**Max Durnea**
- GitHub: [@max-durnea](https://github.com/max-durnea)

---

Built with ❤️ using Go and PostgreSQL
