# Database Token Storage Example

This example demonstrates database-only token storage with automatic refresh for Salla OAuth tokens.

## Features

✅ **Database-Only Storage** - No file system dependencies  
✅ **Automatic Token Refresh** - Tokens are refreshed automatically when expired  
✅ **Multi-User Support** - Store tokens for multiple users  
✅ **Production-Ready** - Proper error handling and logging  

## Database Setup

The example uses SQLite by default, but you can easily switch to MySQL or PostgreSQL.

### SQLite (Default)
```bash
# No setup needed - database file is created automatically
go run main.go
```

### MySQL
```go
import _ "github.com/go-sql-driver/mysql"

db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/salla")
```

### PostgreSQL
```go
import _ "github.com/lib/pq"

db, err := sql.Open("postgres", "host=localhost user=postgres dbname=salla sslmode=disable")
```

## Database Schema

```sql
CREATE TABLE tokens (
    user_id TEXT PRIMARY KEY,
    access_token TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    token_type TEXT NOT NULL,
    expires_at DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

## Usage

### 1. Install Dependencies

```bash
# For SQLite
go get github.com/mattn/go-sqlite3

# For MySQL
go get github.com/go-sql-driver/mysql

# For PostgreSQL
go get github.com/lib/pq
```

### 2. Set Environment Variables

```bash
export SALLA_CLIENT_ID="your_client_id"
export SALLA_CLIENT_SECRET="your_client_secret"
export SALLA_REDIRECT_URI="your_redirect_uri"
```

### 3. Run the Example

```bash
go run main.go
```

## How It Works

### TokenManager Methods

```go
// Initialize database tables
tokenManager.InitDatabase()

// Save a new token
tokenManager.SaveToken(userID, token)

// Get token (auto-refreshes if expired)
token, err := tokenManager.GetToken(userID)

// Manually refresh a token
newToken, err := tokenManager.RefreshToken(userID)

// Get a ready-to-use client
client, err := tokenManager.GetOrCreateClient(userID)

// Delete a token
tokenManager.DeleteToken(userID)
```

### Automatic Refresh

The `GetToken()` method automatically checks if the token is expired or expiring within 5 minutes. If so, it refreshes the token and updates the database:

```go
token, err := tokenManager.GetToken(userID)
// Token is automatically refreshed if needed
// No manual intervention required!
```

### Multi-User Support

Each user has their own token stored by `user_id`:

```go
// User 1
client1, _ := tokenManager.GetOrCreateClient("user_123")

// User 2
client2, _ := tokenManager.GetOrCreateClient("user_456")
```

## Integration Example

```go
// In your web application handler
func ProductsHandler(w http.ResponseWriter, r *http.Request) {
    // Get user ID from session/JWT
    userID := getUserIDFromSession(r)
    
    // Get client with auto-refreshed token
    client, err := tokenManager.GetOrCreateClient(userID)
    if err != nil {
        http.Error(w, "Unauthorized", 401)
        return
    }
    
    // Use the client
    products, _, err := client.Products.List(nil)
    if err != nil {
        http.Error(w, "Failed to fetch products", 500)
        return
    }
    
    json.NewEncoder(w).Encode(products)
}
```

## Production Considerations

1. **Connection Pooling**: Use `db.SetMaxOpenConns()` and `db.SetMaxIdleConns()`
2. **Encryption**: Encrypt tokens at rest using AES-256
3. **Logging**: Add proper logging for audit trails
4. **Error Handling**: Implement retry logic for transient database failures
5. **Token Rotation**: Implement token rotation policy
6. **Monitoring**: Track token refresh failures

## Notes

- Tokens are automatically refreshed 5 minutes before expiry
- All database operations are logged
- The example includes proper error handling
- No file system dependencies - pure database storage
