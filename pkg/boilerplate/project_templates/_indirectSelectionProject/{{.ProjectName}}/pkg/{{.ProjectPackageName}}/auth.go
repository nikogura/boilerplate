//nolint:staticcheck // Package name matches service name requirement from prompt.xml
package {{.ProjectPackageName}}

import (
	"context"
	"crypto"
	"errors"
	"fmt"
	"strings"

	jwtgo "github.com/golang-jwt/jwt/v4"
	"github.com/nikogura/jwt-ssh-agent-go/pkg/agentjwt"
	"{{.ProjectPackage}}/pkg/jwt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Context key for storing authenticated username.
type contextKey string

const usernameKey contextKey = "username"

// TrustedUser represents a trusted user for authentication.
type TrustedUser struct {
	Username   string   `json:"name"`
	Role       string   `json:"role"`
	PublicKeys []string `json:"public_keys"`
}

// AuthInterceptor handles JWT authentication for gRPC requests.
type AuthInterceptor struct {
	cliUsers  *jwt.CliUsers
	audiences []string
	logger    *zap.Logger
}

// NewAuthInterceptor creates a new authentication interceptor.
func NewAuthInterceptor(trustedUsers []TrustedUser, audiences []string, logger *zap.Logger) (interceptor *AuthInterceptor) {
	// Convert TrustedUser to jwt.CliUser format
	var cliUsers []*jwt.CliUser
	for _, user := range trustedUsers {
		cliUser := &jwt.CliUser{
			Name:    user.Username,
			PubKeys: user.PublicKeys,
		}
		cliUsers = append(cliUsers, cliUser)
	}

	cliUsersStruct := &jwt.CliUsers{
		Users:   cliUsers,
		UserMap: make(map[string]*jwt.CliUser),
	}

	// Create user map for easy lookup
	for _, user := range cliUsers {
		cliUsersStruct.UserMap[user.Name] = user
	}

	// Set up JWT verification with SSH agent
	signingMethodED25519Agent := &agentjwt.SigningMethodED25519Agent{
		Name: "EdDSA",
		Hash: crypto.SHA256,
	}
	jwtgo.RegisterSigningMethod(signingMethodED25519Agent.Alg(), func() jwtgo.SigningMethod {
		return signingMethodED25519Agent
	})

	interceptor = &AuthInterceptor{
		cliUsers:  cliUsersStruct,
		audiences: audiences,
		logger:    logger,
	}

	return interceptor
}

// UnaryInterceptor provides JWT authentication for unary gRPC calls.
func (a *AuthInterceptor) UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	ctx, err = a.authenticate(ctx, info.FullMethod)
	if err != nil {
		resp = nil
		return resp, err
	}

	resp, err = handler(ctx, req)
	return resp, err
}

// StreamInterceptor provides JWT authentication for streaming gRPC calls.
func (a *AuthInterceptor) StreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	_, err = a.authenticate(ss.Context(), info.FullMethod)
	if err != nil {
		return err
	}

	return handler(srv, ss)
}

// authenticate validates the JWT token in the request metadata.
func (a *AuthInterceptor) authenticate(ctx context.Context, method string) (newCtx context.Context, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		a.logger.Warn("Missing metadata in request", zap.String("method", method))
		newCtx = nil
		err = status.Error(codes.Unauthenticated, "missing metadata")
		return newCtx, err
	}

	authHeaders := md.Get("authorization")
	if len(authHeaders) == 0 {
		a.logger.Warn("Missing authorization header", zap.String("method", method))
		newCtx = nil
		err = status.Error(codes.Unauthenticated, "missing authorization header")
		return newCtx, err
	}

	authHeader := authHeaders[0]
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		a.logger.Warn("Invalid authorization header format", zap.String("method", method))
		newCtx = nil
		err = status.Error(codes.Unauthenticated, "invalid authorization header format")
		return newCtx, err
	}

	tokenString := parts[1]

	// Create logger adapter for agentjwt
	logAdapter := &ZapAdapter{logger: a.logger}

	// Verify JWT token using agentjwt
	username, _, err := agentjwt.VerifyToken(tokenString, a.audiences, a.publicKeyFunc, logAdapter)
	if err != nil {
		a.logger.Warn("Invalid JWT token",
			zap.String("method", method),
			zap.Error(err))
		newCtx = nil
		err = status.Error(codes.Unauthenticated, "invalid token")
		return newCtx, err
	}

	a.logger.Debug("Authentication successful", zap.String("method", method), zap.String("username", username))

	// Store username in context
	newCtx = context.WithValue(ctx, usernameKey, username)
	return newCtx, err
}

// publicKeyFunc returns public keys for a given username.
func (a *AuthInterceptor) publicKeyFunc(username string) (pubkeys []string, err error) {
	a.logger.Debug("Public key lookup", zap.String("username", username))

	user, ok := a.cliUsers.UserMap[username]
	if ok {
		a.logger.Debug("Public key lookup successful", zap.String("username", username))
		pubkeys = user.PubKeys
		return pubkeys, err
	}

	a.logger.Debug("User not found", zap.String("username", username))
	pubkeys = nil
	err = fmt.Errorf("user not found: %s", username)
	return pubkeys, err
}

// ZapAdapter adapts zap.Logger to the interface expected by agentjwt.
type ZapAdapter struct {
	logger *zap.Logger
}

//nolint:goprintffuncname // Interface required by agentjwt.Logger
func (z *ZapAdapter) Debug(format string, args ...interface{}) {
	z.logger.Debug(fmt.Sprintf(format, args...))
}

// LoadTrustedUsers loads trusted users from a list.
func LoadTrustedUsers(users []TrustedUser) (userMap map[string]TrustedUser, err error) {
	userMap = make(map[string]TrustedUser)

	for _, user := range users {
		if user.Username == "" {
			userMap = nil
			err = errors.New("user missing username")
			return userMap, err
		}
		if len(user.PublicKeys) == 0 {
			userMap = nil
			err = fmt.Errorf("user %s missing public keys", user.Username)
			return userMap, err
		}
		userMap[user.Username] = user
	}

	return userMap, err
}
