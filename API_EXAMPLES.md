# JWT Authorization API Examples

This document provides examples of how to use the JWT authorization endpoints.

## Base URL

```
http://localhost:8080
```

## Endpoints

### 1. Register

Create a new user account.

**Request:**

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "newuser",
    "email": "newuser@example.com",
    "password": "mypassword123"
  }'
```

**Response:**

```json
{
  "id": 2,
  "username": "newuser",
  "email": "newuser@example.com",
  "message": "User registered successfully"
}
```

**Error Response (username exists):**

```json
{
  "error": "username already exists"
}
```

**Error Response (missing fields):**

```json
{
  "error": "username, email, and password are required"
}
```

---

### 2. Health Check

Check if the server is running.

**Request:**

```bash
curl -X GET http://localhost:8080/health
```

**Response:**

```json
{
  "status": "ok"
}
```

---

### 3. Login

Authenticate user and receive access and refresh tokens.

**Request:**

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

**Response:**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 900
}
```

**Test Credentials:**

- Username: `testuser`
- Password: `password123`

---

### 4. Access Protected Resource

Use the access token to access a protected endpoint.

**Request:**

```bash
curl -X GET http://localhost:8080/protected \
  -H "Authorization: Bearer <access_token>"
```

Replace `<access_token>` with the token received from the login endpoint.

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

**Error Response (without token):**

```json
{
  "error": "missing authorization header"
}
```

**Error Response (invalid token):**

```json
{
  "error": "invalid token: token is expired"
}
```

---

### 5. Refresh Access Token

Generate a new access token using the refresh token.

**Request:**

```bash
curl -X POST http://localhost:8080/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "<refresh_token>"
  }'
```

Replace `<refresh_token>` with the token received from the login endpoint.

**Response:**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 900
}
```

---

## Complete Flow Example

### Step 1: Login

```bash
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }')

echo $LOGIN_RESPONSE
```

### Step 2: Extract Tokens (using jq)

```bash
ACCESS_TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.access_token')
REFRESH_TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.refresh_token')

echo "Access Token: $ACCESS_TOKEN"
echo "Refresh Token: $REFRESH_TOKEN"
```

### Step 3: Access Protected Resource

```bash
curl -X GET http://localhost:8080/protected \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

### Step 4: Refresh Token (when access token expires)

```bash
REFRESH_RESPONSE=$(curl -s -X POST http://localhost:8080/auth/refresh \
  -H "Content-Type: application/json" \
  -d "{
    \"refresh_token\": \"$REFRESH_TOKEN\"
  }")

NEW_ACCESS_TOKEN=$(echo $REFRESH_RESPONSE | jq -r '.access_token')
echo "New Access Token: $NEW_ACCESS_TOKEN"
```

---

## Architecture

### JWT Manager

- Handles token generation and validation
- Uses HMAC SHA-256 for signing
- Separates access and refresh token secrets for security

### Auth Middleware

- Validates tokens from the `Authorization: Bearer <token>` header
- Stores user claims in request headers for handler access
- Returns 401 Unauthorized for missing or invalid tokens

### Token Claims

The access token contains the following claims:

```json
{
  "user_id": 1,
  "username": "testuser",
  "email": "test@example.com",
  "exp": 1234567890,
  "iat": 1234567890,
  "nbf": 1234567890
}
```

### Token Expiration

- Access Token: 15 minutes (configurable)
- Refresh Token: 7 days (configurable)

---

## Security Notes

1. **Secret Keys:** In production, use strong, randomly generated secret keys and store them in environment variables
2. **HTTPS:** Always use HTTPS in production to prevent token interception
3. **Token Storage:** Store refresh tokens securely on the client (not in localStorage)
4. **Token Rotation:** Consider implementing token rotation on refresh
5. **Rate Limiting:** Implement rate limiting on the login endpoint to prevent brute force attacks

---

## Running the Server

```bash
go run cmd/main.go
```

Server will start on `http://localhost:8080`
