//nolint:staticcheck // Package name matches service name requirement from prompt.xml
package {{.ProjectPackageName}}

import (
	"context"
	"crypto/tls"
	"fmt"
	"os/user"
	"time"

	"github.com/mitchellh/go-homedir"
	"{{.ProjectPackage}}/pkg/jwt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// Client wraps the gRPC client for easier usage.
type Client struct {
	conn   *grpc.ClientConn
	client {{.ProjectPackageName}}Client
	config Config
	logger *zap.Logger
}

// jwtCreds attaches a JWT token to each RPC.
type jwtCreds struct {
	token string
	tls   bool
}

func (j *jwtCreds) GetRequestMetadata(ctx context.Context, uri ...string) (metadata map[string]string, err error) {
	metadata = map[string]string{
		"authorization": "Bearer " + j.token,
	}
	return metadata, err
}

func (j *jwtCreds) RequireTransportSecurity() bool {
	return j.tls
}

// NewClient creates a new gRPC client.
func NewClient(config Config, logger *zap.Logger, username string, pubKeyFile string) (client *Client, err error) {
	// Use provided username or get current user as fallback
	if username == "" {
		currentUser, userErr := user.Current()
		if userErr != nil {
			client = nil
			err = fmt.Errorf("failed to get current user: %w", userErr)
			return client, err
		}
		username = currentUser.Username
	}

	// Use provided pubKeyFile or construct default path
	if pubKeyFile == "" || pubKeyFile == "~/.ssh/id_ed25519.pub" {
		homeDir, homeDirErr := homedir.Dir()
		if homeDirErr != nil {
			client = nil
			err = fmt.Errorf("failed to get home directory: %w", homeDirErr)
			return client, err
		}
		pubKeyFile = fmt.Sprintf("%s/.ssh/id_ed25519.pub", homeDir)
	}

	// Load public key
	pubkey, err := jwt.LoadPubKey(pubKeyFile)
	if err != nil {
		client = nil
		err = fmt.Errorf("failed to load public key from %s: %w", pubKeyFile, err)
		return client, err
	}

	// Create JWT token with server URL as audience (following example pattern)
	var serverURL string
	if config.PlainText {
		serverURL = fmt.Sprintf("http://%s", config.GRPCAddress)
	} else {
		serverURL = fmt.Sprintf("https://%s", config.GRPCAddress)
	}

	token, err := jwt.MakeToken(serverURL, username, pubkey)
	if err != nil {
		client = nil
		err = fmt.Errorf("failed to create JWT token: %w", err)
		return client, err
	}

	logger.Debug("Created JWT token", zap.String("audience", serverURL), zap.String("username", username))

	// Set up transport credentials
	var creds credentials.TransportCredentials
	if config.PlainText {
		creds = insecure.NewCredentials()
	} else {
		tlsConfig := &tls.Config{}
		creds = credentials.NewTLS(tlsConfig)
	}

	// Create gRPC connection
	conn, err := grpc.NewClient(
		config.GRPCAddress,
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(&jwtCreds{
			token: token,
			tls:   !config.PlainText,
		}),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                30 * time.Second,
			Timeout:             20 * time.Second,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		client = nil
		err = fmt.Errorf("failed to connect to server: %w", err)
		return client, err
	}

	client = &Client{
		conn:   conn,
		client: New{{.ProjectPackageName}}Client(conn),
		config: config,
		logger: logger,
	}

	return client, err
}

// Close closes the client connection.
func (c *Client) Close() (err error) {
	err = c.conn.Close()
	return err
}

// CallFoo calls the Foo RPC method.
func (c *Client) CallFoo(ctx context.Context, username string) (response string, err error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.ClientTimeout)
	defer cancel()

	req := &FooRequest{
		Username: username,
	}

	c.logger.Debug("Calling Foo method", zap.String("username", username))

	resp, err := c.client.Foo(ctx, req)
	if err != nil {
		c.logger.Error("Foo call failed", zap.Error(err))
		response = ""
		err = fmt.Errorf("foo call failed: %w", err)
		return response, err
	}

	response = resp.GetMessage()
	return response, err
}

// CallBar calls the Bar RPC method.
func (c *Client) CallBar(ctx context.Context, username string) (response string, err error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.ClientTimeout)
	defer cancel()

	req := &BarRequest{
		Username: username,
	}

	c.logger.Debug("Calling Bar method", zap.String("username", username))

	resp, err := c.client.Bar(ctx, req)
	if err != nil {
		c.logger.Error("Bar call failed", zap.Error(err))
		response = ""
		err = fmt.Errorf("bar call failed: %w", err)
		return response, err
	}

	response = resp.GetMessage()
	return response, err
}

// CallBaz calls the Baz RPC method.
func (c *Client) CallBaz(ctx context.Context, username string) (response string, err error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.ClientTimeout)
	defer cancel()

	req := &BazRequest{
		Username: username,
	}

	c.logger.Debug("Calling Baz method", zap.String("username", username))

	resp, err := c.client.Baz(ctx, req)
	if err != nil {
		c.logger.Error("Baz call failed", zap.Error(err))
		response = ""
		err = fmt.Errorf("baz call failed: %w", err)
		return response, err
	}

	response = resp.GetMessage()
	return response, err
}
