# JWT Authorization Implementation

This Go project implements JWT (JSON Web Token) based authorization without using cookies. All authentication parameters are sent in request bodies.

## Features

✅ JWT token generation and validation
✅ Access tokens (15 minutes) and refresh tokens (7 days)
✅ Request body authentication (no cookies)
✅ Auth middleware for protected endpoints
✅ Password hashing with bcrypt
✅ Separate secrets for access and refresh tokens
✅ Bearer token authentication via Authorization header

## Project Structure

```
.
├── cmd/
│   └── main.go                 # HTTP server and route handlers
├── internal/
│   └── auth/
│       ├── auth.go             # JWT token management
│       ├── user.go             # User model and password utilities
│       └── middleware.go        # Auth middleware for protected routes
├── go.mod                       # Go module definition
├── go.sum                       # Dependency checksums
├── API_EXAMPLES.md              # API usage examples
└── README_JWT_AUTH.md           # This file
```

## Installation & Setup

### Prerequisites

- Go 1.25.4 or higher

### Dependencies

```bash
go mod tidy
```

This will install:

- `github.com/golang-jwt/jwt/v5` - JWT token handling
- `golang.org/x/crypto` - Password hashing with bcrypt

### Running the Server

```bash
go run cmd/main.go
```

The server will start on `http://localhost:8080`

## API Overview

### Authentication Endpoints

#### 1. Login - `/auth/login` (POST)

Authenticate user and receive tokens.

**Request Body:**

```json
{
  "username": "testuser",
  "password": "password123"
}
```

**Response:**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 900
}
```

#### 2. Refresh - `/auth/refresh` (POST)

Generate a new access token using a refresh token.

**Request Body:**

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response:**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 900
}
```

### Protected Endpoints

#### Access Protected Resource - `/protected` (GET)

Example protected endpoint that requires a valid access token.

**Request Header:**

```
Authorization: Bearer <access_token>
```

**Response:**

```json
{
  "message": "This is a protected resource",
  "user": {
    "id": "1",
    "username": "testuser",
    "email": "test@example.com"
  }
}
```

## How It Works

### Authentication Flow

1. **Login:** User sends username and password in request body
2. **Validation:** Server validates credentials
3. **Token Generation:** Server generates access and refresh tokens
4. **Token Usage:** Client sends access token in `Authorization: Bearer <token>` header
5. **Token Refresh:** When access token expires, client uses refresh token to get a new one

### Token Structure

**Access Token Claims:**

```json
{
  "user_id": 1,
  "username": "testuser",
  "email": "test@example.com",
  "exp": 1234567890, // expiration time
  "iat": 1234567890, // issued at
  "nbf": 1234567890 // not before
}
```

**Refresh Token:** Contains only standard claims (exp, iat, nbf)

### Security Mechanisms

1. **HMAC SHA-256 Signing:** Tokens are signed with secret keys
2. **Separate Secrets:** Different secrets for access and refresh tokens
3. **Password Hashing:** Passwords are hashed with bcrypt (cost factor 10)
4. **Token Expiration:** Access tokens expire after 15 minutes
5. **Middleware Validation:** All protected routes validate tokens via middleware

## Usage Examples

### Using cURL

**Login:**

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser", "password": "password123"}'
```

**Access Protected Resource:**

```bash
curl -X GET http://localhost:8080/protected \
  -H "Authorization: Bearer <access_token>"
```

**Refresh Token:**

```bash
curl -X POST http://localhost:8080/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token": "<refresh_token>"}'
```

### Using JavaScript/Fetch

```javascript
// Login
const loginResponse = await fetch("http://localhost:8080/auth/login", {
  method: "POST",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify({
    username: "testuser",
    password: "password123",
  }),
});

const { access_token, refresh_token } = await loginResponse.json();

// Access protected resource
const protectedResponse = await fetch("http://localhost:8080/protected", {
  headers: {
    Authorization: `Bearer ${access_token}`,
  },
});

const data = await protectedResponse.json();
console.log(data);

// Refresh token
const refreshResponse = await fetch("http://localhost:8080/auth/refresh", {
  method: "POST",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify({ refresh_token }),
});

const { access_token: newToken } = await refreshResponse.json();
```

## Test Credentials

The demo includes a pre-configured test user:

- **Username:** `testuser`
- **Password:** `password123`

## Customization

### Changing Token Expiration

In `cmd/main.go`:

```go
jwtManager := auth.NewJWTManager(
  "your-secret-key",
  "your-refresh-secret-key",
  30*time.Minute,      // Access token TTL (change here)
  14*24*time.Hour,     // Refresh token TTL (change here)
)
```

### Changing Secret Keys

⚠️ **Never hardcode secrets in production!** Use environment variables:

```go
accessSecret := os.Getenv("JWT_SECRET")
refreshSecret := os.Getenv("JWT_REFRESH_SECRET")
```

### Adding More Users

Modify the `users` map in `cmd/main.go`:

```go
var users = map[string]*auth.User{
  "testuser": { /* ... */ },
  "newuser": {
    ID:       2,
    Username: "newuser",
    Email:    "newuser@example.com",
    Password: hashedPassword, // Use auth.HashPassword()
  },
}
```

## Key Components

### `internal/auth/auth.go`

- **JWTManager:** Handles token generation and validation
- **GenerateTokens:** Creates access and refresh tokens
- **ValidateToken:** Validates and extracts claims from tokens
- **RefreshAccessToken:** Issues new access token from refresh token

### `internal/auth/user.go`

- **User:** User data structure
- **LoginRequest:** Login request payload
- **RefreshTokenRequest:** Token refresh request payload
- **HashPassword:** Bcrypt password hashing
- **VerifyPassword:** Password verification

### `internal/auth/middleware.go`

- **AuthMiddleware:** HTTP middleware for token validation
- Validates `Authorization: Bearer <token>` header
- Extracts and validates claims
- Passes user info to handlers via request headers

### `cmd/main.go`

- HTTP route handlers
- Server initialization
- Example protected endpoint

## Error Handling

The API returns appropriate HTTP status codes:

| Status | Scenario                                               |
| ------ | ------------------------------------------------------ |
| 200    | Successful authentication or protected resource access |
| 400    | Invalid request body or missing parameters             |
| 401    | Missing/invalid token or credentials                   |
| 405    | Method not allowed (e.g., GET on POST endpoint)        |
| 500    | Server error                                           |

## Production Considerations

1. **Environment Variables:** Store secrets in environment variables, not code
2. **HTTPS:** Always use HTTPS in production
3. **Rate Limiting:** Implement rate limiting on login/refresh endpoints
4. **Token Rotation:** Consider implementing token rotation
5. **Blacklisting:** Implement token blacklist for logout functionality
6. **Database:** Replace in-memory user map with a database
7. **Logging:** Add proper logging for security events
8. **CORS:** Configure CORS headers if needed

## Troubleshooting

### Token expired error

- The access token expires after 15 minutes by default
- Use the refresh token to get a new access token

### Invalid token error

- Ensure the token is being sent in the correct format: `Authorization: Bearer <token>`
- Check that the token hasn't been modified
- Verify the server's secret keys match

### Missing authorization header

- Ensure you're including the `Authorization` header with the token
- Use the format: `Authorization: Bearer <token>` (note the space)

## Additional Resources

See `API_EXAMPLES.md` for detailed API usage examples and complete flow demonstrations.
