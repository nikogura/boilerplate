# {{.ProjectName}}

A Go gRPC service demonstrating JWT-SSH authentication patterns with simple method calls. This example provides a foundation for implementing your own authenticated gRPC services.

## Overview

This service implements:
- gRPC server with JWT-SSH authentication using `github.com/nikogura/jwt-ssh-agent-go`
- Method-level authentication (each RPC method performs its own JWT verification)
- Role-based authorization (users must have appropriate roles to access methods)
- Username-based public key lookup from `users.json`
- Simple foo/bar/baz example methods with different role requirements

## Quick Start

1. **Build the service**:
   ```bash
   go build
   ```

2. **Set up users** (customize users.json):
   ```bash
   # Edit users.json with your usernames, roles, and SSH public keys
   ```

3. **Start the server**:
   ```bash
   ./{{.ProjectName}} server
   ```

4. **Call methods** (in another terminal):
   ```bash
   # Uses current system username
   ./{{.ProjectName}} foo
   
   # Or specify a username
   ./{{.ProjectName}} foo -u myuser
   
   # Other methods
   ./{{.ProjectName}} bar -u myuser
   ./{{.ProjectName}} baz -u myuser
   ```

## Authentication and Authorization Flow

1. Client extracts username (from `-u` flag or current system user)
2. Client creates JWT token signed with SSH private key
3. Client sends username in gRPC request + JWT in Authorization header
4. Server extracts username from request
5. Server looks up that user's public keys from `users.json`
6. Server verifies JWT against only that user's keys
7. Server checks if user has required role for the method:
   - `foo` and `bar` require `user` role (or `admin`)
   - `baz` requires `admin` role
8. On success, server executes the method logic

## Role Requirements

- **foo**: Requires `user` role
- **bar**: Requires `user` role  
- **baz**: Requires `admin` role

Note: Users with `admin` role have access to all methods.

## Replacing Example Methods with Your Own

Follow these steps to replace the foo/bar/baz methods with your own:

### Step 1: Update the Protocol Buffer Definition

Edit `pkg/{{.ProjectPackageName}}/service.proto`:

```protobuf
service {{.ProjectName}} {
    // Replace these with your methods
    rpc MyMethod1(MyMethod1Request) returns (MyMethod1Response);
    rpc MyMethod2(MyMethod2Request) returns (MyMethod2Response);
}

// Replace with your request/response messages
message MyMethod1Request {
    // REQUIRED: Always include username for authentication
    string username = 1;
    // Add your custom fields
    string my_field = 2;
}

message MyMethod1Response {
    // Response message (required for debugging)
    string message = 1;
    // Add your custom fields
    string my_result = 2;
}
```

### Step 2: Regenerate Protocol Buffer Code

```bash
protoc --proto_path=pkg/{{.ProjectPackageName}} \
       --go_out=pkg/{{.ProjectPackageName}} \
       --go-grpc_out=pkg/{{.ProjectPackageName}} \
       service.proto
```

### Step 3: Update Server Methods

Edit `pkg/{{.ProjectPackageName}}/server.go` and replace the Foo/Bar/Baz methods with your own:

```go
//nolint:dupl // Authentication pattern intentionally duplicated across methods
func (s *Server) MyMethod1(ctx context.Context, req *MyMethod1Request) (response *MyMethod1Response, err error) {
    // AUTHENTICATION FLOW - DO NOT MODIFY THIS SECTION
    requestUsername := req.GetUsername()
    if requestUsername == "" {
        s.logger.Error("no username in request")
        response = &MyMethod1Response{Message: "Auth failed: no username provided"}
        return response, err
    }

    // Get metadata for JWT
    m, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        s.logger.Error("no metadata info")
        response = &MyMethod1Response{Message: fmt.Sprintf("Auth failed for user (%s)", requestUsername)}
        return response, err
    }

    authHeaders := m.Get("authorization")
    if len(authHeaders) == 0 {
        s.logger.Error("no authorization header found")
        response = &MyMethod1Response{Message: fmt.Sprintf("Auth failed for user (%s)", requestUsername)}
        return response, err
    }

    ah := authHeaders[0]
    parts := strings.Split(ah, " ")
    if len(parts) <= 1 {
        s.logger.Error("invalid authorization header format", zap.String("header", ah))
        response = &MyMethod1Response{Message: fmt.Sprintf("Auth failed for user (%s)", requestUsername)}
        return response, err
    }

    tokenString := parts[1]
    logAdapter := &ZapAdapter{logger: s.logger}

    // Verify JWT using request username for public key lookup
    pubKeyFunc := s.PublicKeyFuncForUser(requestUsername)
    jwtUsername, _, err := agentjwt.VerifyToken(tokenString, s.audiences, pubKeyFunc, logAdapter)
    //nolint:nestif // Complex auth flow needed for username extraction
    if err != nil {
        s.logger.Error("invalid token or username not found", zap.Error(err))
        // Extract JWT username for debugging
        jwtUser := unknownUser
        token, parseErr := jwt.Parse(tokenString, nil)
        if parseErr == nil && token.Claims != nil {
            if tokenClaims, claimsOk := token.Claims.(jwt.MapClaims); claimsOk {
                if sub, exists := tokenClaims["sub"]; exists {
                    if subStr, isString := sub.(string); isString {
                        jwtUser = subStr
                    }
                }
            }
        }
        response = &MyMethod1Response{Message: fmt.Sprintf("Auth failed for user (%s) - JWT user (%s)", requestUsername, jwtUser)}
        err = nil // Don't return error details to client
        return response, err
    }

    // Check role authorization - specify required role for your method
    hasRole, err := s.CheckUserRole(jwtUsername, "user") // Change "user" to required role
    if err != nil {
        s.logger.Error("role lookup failed", zap.Error(err))
        response = &MyMethod1Response{Message: fmt.Sprintf("Auth failed for user (%s) - role lookup error", requestUsername)}
        err = nil // Don't return error details to client
        return response, err
    }
    if !hasRole {
        s.logger.Warn("insufficient role", zap.String("user", jwtUsername), zap.String("required_role", "user"))
        response = &MyMethod1Response{Message: fmt.Sprintf("Access denied for user (%s) - insufficient role", jwtUsername)}
        return response, err
    }
    // END AUTHENTICATION AND AUTHORIZATION FLOW

    // 🚀 PUT YOUR BUSINESS LOGIC HERE 🚀
    // Authentication and authorization succeeded - jwtUsername contains the verified user
    s.logger.Info("MY_METHOD_1 CALLED",
        zap.String("method", "MyMethod1"),
        zap.String("request_user", requestUsername),
        zap.String("jwt_user", jwtUsername),
        zap.String("description", "User called MyMethod1"))

    // Example: Process your custom fields
    myField := req.GetMyField()
    
    // TODO: Add your business logic here
    result := fmt.Sprintf("Hello %s, you sent: %s", jwtUsername, myField)

    response = &MyMethod1Response{
        Message:  fmt.Sprintf("You (%s) called MyMethod1", jwtUsername), // Required for debugging
        MyResult: result, // Your custom response
    }

    return response, err
}
```

### Step 4: Update Client Methods

Edit `pkg/{{.ProjectPackageName}}/client.go` to replace the CallFoo/CallBar/CallBaz methods:

```go
// CallMyMethod1 calls the MyMethod1 RPC method.
func (c *Client) CallMyMethod1(ctx context.Context, username string, myField string) (response *MyMethod1Response, err error) {
    ctx, cancel := context.WithTimeout(ctx, c.config.ClientTimeout)
    defer cancel()

    req := &MyMethod1Request{
        Username: username,
        MyField:  myField,
    }

    c.logger.Debug("Calling MyMethod1", zap.String("username", username))

    response, err = c.client.MyMethod1(ctx, req)
    if err != nil {
        c.logger.Error("MyMethod1 call failed", zap.Error(err))
        err = fmt.Errorf("mymethod1 call failed: %w", err)
        return response, err
    }

    return response, err
}
```

### Step 5: Update Command Files

Replace `cmd/foo.go`, `cmd/bar.go`, `cmd/baz.go` with your own commands:

```go
// cmd/mymethod1.go
package cmd

import (
    "context"
    "fmt"

    "github.com/spf13/cobra"
)

// mymethod1Cmd calls the mymethod1 method
//
//nolint:gochecknoglobals // Cobra boilerplate
var mymethod1Cmd = &cobra.Command{
    Use:   "mymethod1",
    Short: "Call the mymethod1 method",
    Long:  `Calls the mymethod1 method on the server.`,
    Run:   runMyMethod1,
}

//nolint:gochecknoinits // Cobra boilerplate
func init() {
    rootCmd.AddCommand(mymethod1Cmd)
    mymethod1Cmd.Flags().StringP("field", "f", "", "Custom field value")
}

func runMyMethod1(cmd *cobra.Command, args []string) {
    client, err := createClient()
    if err != nil {
        logger.Fatal("Failed to create client", zapError(err))
    }
    defer client.Close()

    // Get effective username (from flag or current user)
    effectiveUsername, err := getEffectiveUsername()
    if err != nil {
        logger.Fatal("Failed to get username", zapError(err))
    }

    // Get custom field
    myField, _ := cmd.Flags().GetString("field")

    ctx := context.Background()
    response, err := client.CallMyMethod1(ctx, effectiveUsername, myField)
    if err != nil {
        logger.Fatal("Failed to call mymethod1", zapError(err))
    }

    fmt.Println(response.GetMessage())
    fmt.Println("Result:", response.GetMyResult())
}
```

## Authentication Error Debugging

When authentication fails, error messages show both usernames for debugging:
- `"Auth failed for user (requested_username) - JWT user (jwt_username)"`

This helps identify mismatches between:
- What username the client thinks it's sending
- What username is in the JWT token
- What username the server extracted from the request

## Key Files to Modify

- `pkg/{{.ProjectPackageName}}/service.proto` - Protocol buffer definitions
- `pkg/{{.ProjectPackageName}}/server.go` - Server method implementations  
- `pkg/{{.ProjectPackageName}}/client.go` - Client method calls
- `cmd/*.go` - Command-line interface

## User Configuration

The `users.json` file defines users with their roles and SSH public keys:

```json
[
  {
    "name": "testuser",
    "role": "user",
    "public_keys": [
        "ssh-ed25519 AAAAABBBBCCCCDDDDEEEFFFF testuser@example.com",
        "ssh-ed25519 AAAAABBBBCCCCDDDDEEEFFFF testuserkey2@example.com"
    ]
  },
  {
    "name": "alice", 
    "role": "admin",
    "public_keys": [
        "ssh-ed25519 AAAAABBBBCCCCDDDDEEEFFFF alicekey1@example.com",
        "ssh-ed25519 AAAAABBBBCCCCDDDDEEEFFFF alicekey2@example.com"
    ]
  }
]
```

**Available roles:**
- `user`: Can access methods requiring user privileges
- `admin`: Can access all methods (has elevated privileges)

## Authentication Requirements

**CRITICAL**: Always include these elements in your methods:

1. **Username in request**: Every request message must have `string username = 1;`
2. **Authentication flow**: Copy the authentication section from example methods
3. **Role authorization**: Add role checking with `s.CheckUserRole(jwtUsername, "required_role")`
4. **Placeholder comment**: Mark where your business logic goes with `// 🚀 PUT YOUR BUSINESS LOGIC HERE 🚀`
5. **Debugging message**: Include `"You (username) called MethodName"` in response message
6. **Proper logging**: Log method calls with user information

## Testing

```bash
# Start server
./{{.ProjectName}} server

# Test your methods
./{{.ProjectName}} mymethod1 -u myuser -f "test data"
```

The authentication flow will automatically handle JWT verification against the user's SSH public keys from `users.json`.