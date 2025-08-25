//nolint:staticcheck // Package name matches service name requirement from prompt.xml
package {{.ProjectPackageName}}

import (
	"context"
	"crypto"
	"fmt"
	"net"
	"strings"
	"time"

	jwtgo "github.com/golang-jwt/jwt/v4"
	"github.com/nikogura/jwt-ssh-agent-go/pkg/agentjwt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

const unknownUser = "unknown"

// Server implements the {{.ProjectPackageName}} gRPC service.
type Server struct {
	Unimplemented{{.ProjectPackageName}}Server

	config       Config
	logger       *zap.Logger
	grpcServer   *grpc.Server
	healthServer *health.Server
	trustedUsers []TrustedUser
	audiences    []string
}

// NewServer creates a new gRPC server instance.
func NewServer(config Config, logger *zap.Logger, trustedUsers []TrustedUser) (server *Server, err error) {
	// Set up JWT verification (following example pattern)
	signingMethodED25519Agent := &agentjwt.SigningMethodED25519Agent{
		Name: "EdDSA",
		Hash: crypto.SHA256,
	}

	// Register ssh-agent auth
	jwtgo.RegisterSigningMethod(signingMethodED25519Agent.Alg(), func() jwtgo.SigningMethod {
		return signingMethodED25519Agent
	})

	// Build expected audiences (following example pattern)
	var validAudiences []string
	// JWT token creation extracts hostname from URLs, so accept just the hostname
	host := strings.Split(config.GRPCAddress, ":")[0]
	validAudiences = append(validAudiences, host)
	// Also accept the original config audience for backward compatibility
	validAudiences = append(validAudiences, config.Audience)

	// Create gRPC server (no interceptors needed since methods handle authentication directly)
	grpcServer := grpc.NewServer()

	// Create health server
	healthServer := health.NewServer()

	server = &Server{
		config:       config,
		logger:       logger,
		grpcServer:   grpcServer,
		healthServer: healthServer,
		trustedUsers: trustedUsers,
		audiences:    validAudiences,
	}

	// Register services
	Register{{.ProjectPackageName}}Server(grpcServer, server)
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)

	// Enable reflection if configured
	if config.EnableReflection {
		reflection.Register(grpcServer)
		logger.Info("gRPC reflection enabled")
	}

	return server, err
}

// Start starts the gRPC server.
func (s *Server) Start(ctx context.Context) (err error) {
	lc := &net.ListenConfig{}
	listener, err := lc.Listen(ctx, "tcp", s.config.GRPCAddress)
	if err != nil {
		err = fmt.Errorf("failed to listen on %s: %w", s.config.GRPCAddress, err)
		return err
	}

	s.logger.Info("Starting gRPC server",
		zap.String("address", s.config.GRPCAddress),
		zap.Bool("reflection", s.config.EnableReflection))

	// Set health status to serving
	s.healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	// Start server in goroutine
	go func() {
		serveErr := s.grpcServer.Serve(listener)
		if serveErr != nil {
			s.logger.Error("gRPC server failed", zap.Error(serveErr))
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	// Graceful shutdown
	s.logger.Info("Shutting down gRPC server")
	s.healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_NOT_SERVING)

	shutdownDone := make(chan struct{})
	go func() {
		s.grpcServer.GracefulStop()
		close(shutdownDone)
	}()

	// Force stop if graceful shutdown takes too long
	select {
	case <-shutdownDone:
		s.logger.Info("gRPC server shutdown complete")
	case <-time.After(s.config.ShutdownTimeout):
		s.logger.Warn("Graceful shutdown timeout, forcing stop")
		s.grpcServer.Stop()
	}

	return err
}

// PublicKeyFunc returns public keys for a given username.
func (s *Server) PublicKeyFunc(subject string) (pubkeys []string, err error) {
	s.logger.Debug("pubkey lookup attempt", zap.String("username", subject))
	for _, user := range s.trustedUsers {
		if user.Username == subject {
			s.logger.Debug("pubkey lookup successful", zap.String("username", subject))
			pubkeys = user.PublicKeys
			return pubkeys, err
		}
	}
	s.logger.Debug("user not found", zap.String("username", subject))
	pubkeys = nil
	err = fmt.Errorf("user not found: %s", subject)
	return pubkeys, err
}

// PublicKeyFuncForUser returns public keys for a specific username only.
func (s *Server) PublicKeyFuncForUser(requestedUsername string) func(subject string) ([]string, error) {
	return func(subject string) (pubkeys []string, err error) {
		s.logger.Debug("pubkey lookup attempt", zap.String("requested_user", requestedUsername), zap.String("jwt_subject", subject))
		// Only allow the requested username
		if subject != requestedUsername {
			s.logger.Debug("username mismatch", zap.String("requested", requestedUsername), zap.String("jwt_subject", subject))
			pubkeys = nil
			err = fmt.Errorf("username mismatch: requested %s, jwt subject %s", requestedUsername, subject)
			return pubkeys, err
		}
		// Look up keys for the requested user
		return s.PublicKeyFunc(requestedUsername)
	}
}

// GetUserRole returns the role for a given username.
func (s *Server) GetUserRole(username string) (role string, err error) {
	for _, user := range s.trustedUsers {
		if user.Username == username {
			role = user.Role
			return role, err
		}
	}
	role = ""
	err = fmt.Errorf("user not found: %s", username)
	return role, err
}

// CheckUserRole verifies that a user has the required role.
func (s *Server) CheckUserRole(username string, requiredRole string) (hasRole bool, err error) {
	userRole, err := s.GetUserRole(username)
	if err != nil {
		hasRole = false
		return hasRole, err
	}

	// Admin role has access to everything
	if userRole == "admin" {
		hasRole = true
		return hasRole, err
	}

	// Check exact role match
	hasRole = userRole == requiredRole
	return hasRole, err
}

//nolint:dupl,gocognit,funlen // Authentication pattern intentionally duplicated across methods following example
func (s *Server) Foo(ctx context.Context, req *FooRequest) (response *FooResponse, err error) {
	// Extract username from request
	requestUsername := req.GetUsername()
	if requestUsername == "" {
		s.logger.Error("no username in request")
		response = &FooResponse{Message: "Auth failed: no username provided"}
		return response, err
	}

	// Get metadata or return unsuccessful
	m, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		s.logger.Error("no metadata info")
		response = &FooResponse{Message: fmt.Sprintf("Auth failed for user (%s)", requestUsername)}
		return response, err
	}

	authHeaders := m.Get("authorization")

	// Parse out the bearer token
	if len(authHeaders) == 0 {
		s.logger.Error("no authorization header found")
		response = &FooResponse{Message: fmt.Sprintf("Auth failed for user (%s)", requestUsername)}
		return response, err
	}

	ah := authHeaders[0]
	parts := strings.Split(ah, " ")
	// Get the JWT string
	if len(parts) <= 1 {
		s.logger.Error("invalid authorization header format", zap.String("header", ah))
		response = &FooResponse{Message: fmt.Sprintf("Auth failed for user (%s)", requestUsername)}
		return response, err
	}

	tokenString := parts[1]

	// Have to wrap the zap logger cos the Debug signature is different than what is expected by agentjwt
	logAdapter := &ZapAdapter{logger: s.logger}

	// Verify JWT using request username for public key lookup
	pubKeyFunc := s.PublicKeyFuncForUser(requestUsername)
	jwtUsername, _, err := agentjwt.VerifyToken(tokenString, s.audiences, pubKeyFunc, logAdapter)
	//nolint:nestif // Complex auth flow needed for username extraction
	if err != nil {
		s.logger.Error("invalid token or username not found", zap.Error(err))
		// Try to extract JWT username for error message
		jwtUser := unknownUser
		token, parseErr := jwtgo.Parse(tokenString, nil)
		if parseErr == nil && token.Claims != nil {
			if tokenClaims, claimsOk := token.Claims.(jwtgo.MapClaims); claimsOk {
				if sub, exists := tokenClaims["sub"]; exists {
					if subStr, isString := sub.(string); isString {
						jwtUser = subStr
					}
				}
			}
		}
		response = &FooResponse{Message: fmt.Sprintf("Auth failed for user (%s) - JWT user (%s)", requestUsername, jwtUser)}
		err = nil // Don't return error details to client
		return response, err
	}

	// Check role authorization - Foo requires 'user' role
	hasRole, err := s.CheckUserRole(jwtUsername, "user")
	if err != nil {
		s.logger.Error("role lookup failed", zap.Error(err))
		response = &FooResponse{Message: fmt.Sprintf("Auth failed for user (%s) - role lookup error", requestUsername)}
		err = nil // Don't return error details to client
		return response, err
	}
	if !hasRole {
		s.logger.Warn("insufficient role", zap.String("user", jwtUsername), zap.String("required_role", "user"))
		response = &FooResponse{Message: fmt.Sprintf("Access denied for user (%s) - insufficient role", jwtUsername)}
		return response, err
	}

	// 🚀 PUT YOUR BUSINESS LOGIC HERE 🚀
	// Authentication and authorization succeeded - jwtUsername contains the verified user

	// Log who called which method - clearly annotated for Foo
	s.logger.Info("FOO METHOD CALLED",
		zap.String("method", "Foo"),
		zap.String("request_user", requestUsername),
		zap.String("jwt_user", jwtUsername),
		zap.String("description", "User called the Foo method"))

	// TODO: Replace this example logic with your own business logic
	// Return response saying who called what method (pass back username for error checking)
	message := fmt.Sprintf("You (%s) called Foo", jwtUsername)

	response = &FooResponse{
		Message: message,
	}

	return response, err
}

//nolint:dupl,gocognit,funlen // Authentication pattern intentionally duplicated across methods following example
func (s *Server) Bar(ctx context.Context, req *BarRequest) (response *BarResponse, err error) {
	// Extract username from request
	requestUsername := req.GetUsername()
	if requestUsername == "" {
		s.logger.Error("no username in request")
		response = &BarResponse{Message: "Auth failed: no username provided"}
		return response, err
	}

	// Get metadata or return unsuccessful
	m, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		s.logger.Error("no metadata info")
		response = &BarResponse{Message: fmt.Sprintf("Auth failed for user (%s)", requestUsername)}
		return response, err
	}

	authHeaders := m.Get("authorization")

	// Parse out the bearer token
	if len(authHeaders) == 0 {
		s.logger.Error("no authorization header found")
		response = &BarResponse{Message: fmt.Sprintf("Auth failed for user (%s)", requestUsername)}
		return response, err
	}

	ah := authHeaders[0]
	parts := strings.Split(ah, " ")
	// Get the JWT string
	if len(parts) <= 1 {
		s.logger.Error("invalid authorization header format", zap.String("header", ah))
		response = &BarResponse{Message: fmt.Sprintf("Auth failed for user (%s)", requestUsername)}
		return response, err
	}

	tokenString := parts[1]

	// Have to wrap the zap logger cos the Debug signature is different than what is expected by agentjwt
	logAdapter := &ZapAdapter{logger: s.logger}

	// Verify JWT using request username for public key lookup
	pubKeyFunc := s.PublicKeyFuncForUser(requestUsername)
	jwtUsername, _, err := agentjwt.VerifyToken(tokenString, s.audiences, pubKeyFunc, logAdapter)
	//nolint:nestif // Complex auth flow needed for username extraction
	if err != nil {
		s.logger.Error("invalid token or username not found", zap.Error(err))
		// Try to extract JWT username for error message
		jwtUser := unknownUser
		token, parseErr := jwtgo.Parse(tokenString, nil)
		if parseErr == nil && token.Claims != nil {
			if tokenClaims, claimsOk := token.Claims.(jwtgo.MapClaims); claimsOk {
				if sub, exists := tokenClaims["sub"]; exists {
					if subStr, isString := sub.(string); isString {
						jwtUser = subStr
					}
				}
			}
		}
		response = &BarResponse{Message: fmt.Sprintf("Auth failed for user (%s) - JWT user (%s)", requestUsername, jwtUser)}
		err = nil // Don't return error details to client
		return response, err
	}

	// Check role authorization - Bar requires 'user' role
	hasRole, err := s.CheckUserRole(jwtUsername, "user")
	if err != nil {
		s.logger.Error("role lookup failed", zap.Error(err))
		response = &BarResponse{Message: fmt.Sprintf("Auth failed for user (%s) - role lookup error", requestUsername)}
		err = nil // Don't return error details to client
		return response, err
	}
	if !hasRole {
		s.logger.Warn("insufficient role", zap.String("user", jwtUsername), zap.String("required_role", "user"))
		response = &BarResponse{Message: fmt.Sprintf("Access denied for user (%s) - insufficient role", jwtUsername)}
		return response, err
	}

	// 🚀 PUT YOUR BUSINESS LOGIC HERE 🚀
	// Authentication and authorization succeeded - jwtUsername contains the verified user

	// Log who called which method - clearly annotated for Bar
	s.logger.Info("BAR METHOD CALLED",
		zap.String("method", "Bar"),
		zap.String("request_user", requestUsername),
		zap.String("jwt_user", jwtUsername),
		zap.String("description", "User called the Bar method"))

	// TODO: Replace this example logic with your own business logic
	// Return response saying who called what method (pass back username for error checking)
	message := fmt.Sprintf("You (%s) called Bar", jwtUsername)

	response = &BarResponse{
		Message: message,
	}

	return response, err
}

//nolint:dupl,gocognit,funlen // Authentication pattern intentionally duplicated across methods following example
func (s *Server) Baz(ctx context.Context, req *BazRequest) (response *BazResponse, err error) {
	// Extract username from request
	requestUsername := req.GetUsername()
	if requestUsername == "" {
		s.logger.Error("no username in request")
		response = &BazResponse{Message: "Auth failed: no username provided"}
		return response, err
	}

	// Get metadata or return unsuccessful
	m, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		s.logger.Error("no metadata info")
		response = &BazResponse{Message: fmt.Sprintf("Auth failed for user (%s)", requestUsername)}
		return response, err
	}

	authHeaders := m.Get("authorization")

	// Parse out the bearer token
	if len(authHeaders) == 0 {
		s.logger.Error("no authorization header found")
		response = &BazResponse{Message: fmt.Sprintf("Auth failed for user (%s)", requestUsername)}
		return response, err
	}

	ah := authHeaders[0]
	parts := strings.Split(ah, " ")
	// Get the JWT string
	if len(parts) <= 1 {
		s.logger.Error("invalid authorization header format", zap.String("header", ah))
		response = &BazResponse{Message: fmt.Sprintf("Auth failed for user (%s)", requestUsername)}
		return response, err
	}

	tokenString := parts[1]

	// Have to wrap the zap logger cos the Debug signature is different than what is expected by agentjwt
	logAdapter := &ZapAdapter{logger: s.logger}

	// Verify JWT using request username for public key lookup
	pubKeyFunc := s.PublicKeyFuncForUser(requestUsername)
	jwtUsername, _, err := agentjwt.VerifyToken(tokenString, s.audiences, pubKeyFunc, logAdapter)
	//nolint:nestif // Complex auth flow needed for username extraction
	if err != nil {
		s.logger.Error("invalid token or username not found", zap.Error(err))
		// Try to extract JWT username for error message
		jwtUser := unknownUser
		token, parseErr := jwtgo.Parse(tokenString, nil)
		if parseErr == nil && token.Claims != nil {
			if tokenClaims, claimsOk := token.Claims.(jwtgo.MapClaims); claimsOk {
				if sub, exists := tokenClaims["sub"]; exists {
					if subStr, isString := sub.(string); isString {
						jwtUser = subStr
					}
				}
			}
		}
		response = &BazResponse{Message: fmt.Sprintf("Auth failed for user (%s) - JWT user (%s)", requestUsername, jwtUser)}
		err = nil // Don't return error details to client
		return response, err
	}

	// Check role authorization - Baz requires 'admin' role
	hasRole, err := s.CheckUserRole(jwtUsername, "admin")
	if err != nil {
		s.logger.Error("role lookup failed", zap.Error(err))
		response = &BazResponse{Message: fmt.Sprintf("Auth failed for user (%s) - role lookup error", requestUsername)}
		err = nil // Don't return error details to client
		return response, err
	}
	if !hasRole {
		s.logger.Warn("insufficient role", zap.String("user", jwtUsername), zap.String("required_role", "admin"))
		response = &BazResponse{Message: fmt.Sprintf("Access denied for user (%s) - insufficient role", jwtUsername)}
		return response, err
	}

	// 🚀 PUT YOUR BUSINESS LOGIC HERE 🚀
	// Authentication and authorization succeeded - jwtUsername contains the verified user

	// Log who called which method - clearly annotated for Baz
	s.logger.Info("BAZ METHOD CALLED",
		zap.String("method", "Baz"),
		zap.String("request_user", requestUsername),
		zap.String("jwt_user", jwtUsername),
		zap.String("description", "User called the Baz method"))

	// TODO: Replace this example logic with your own business logic
	// Return response saying who called what method (pass back username for error checking)
	message := fmt.Sprintf("You (%s) called Baz", jwtUsername)

	response = &BazResponse{
		Message: message,
	}

	return response, err
}
