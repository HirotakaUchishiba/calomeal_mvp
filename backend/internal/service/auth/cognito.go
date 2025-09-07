package auth

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

// CognitoService defines the interface for AWS Cognito operations
type CognitoService interface {
	SignUp(ctx context.Context, email, password string) (*cognitoidentityprovider.SignUpOutput, error)
	ConfirmSignUp(ctx context.Context, email, confirmationCode string) error
	SignIn(ctx context.Context, email, password string) (*cognitoidentityprovider.InitiateAuthOutput, error)
	SignOut(ctx context.Context, accessToken string) error
	ResetPassword(ctx context.Context, email string) error
	ConfirmResetPassword(ctx context.Context, email, confirmationCode, newPassword string) error
	GetUser(ctx context.Context, accessToken string) (*cognitoidentityprovider.GetUserOutput, error)
}

type cognitoService struct {
	client       *cognitoidentityprovider.Client
	clientID     string
	userPoolID   string
}

// NewCognitoService creates a new Cognito service
func NewCognitoService() (CognitoService, error) {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create Cognito client
	client := cognitoidentityprovider.NewFromConfig(cfg)

	// Get configuration from environment variables
	clientID := os.Getenv("COGNITO_CLIENT_ID")
	if clientID == "" {
		clientID = "your-cognito-client-id" // Default for development
	}

	userPoolID := os.Getenv("COGNITO_USER_POOL_ID")
	if userPoolID == "" {
		userPoolID = "your-cognito-user-pool-id" // Default for development
	}

	return &cognitoService{
		client:     client,
		clientID:   clientID,
		userPoolID: userPoolID,
	}, nil
}

// SignUp registers a new user with Cognito
func (s *cognitoService) SignUp(ctx context.Context, email, password string) (*cognitoidentityprovider.SignUpOutput, error) {
	input := &cognitoidentityprovider.SignUpInput{
		ClientId: aws.String(s.clientID),
		Username: aws.String(email),
		Password: aws.String(password),
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(email),
			},
		},
	}

	result, err := s.client.SignUp(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to sign up user: %w", err)
	}

	fmt.Printf("CognitoService: User signup initiated for %s\n", email)
	return result, nil
}

// ConfirmSignUp confirms user registration with verification code
func (s *cognitoService) ConfirmSignUp(ctx context.Context, email, confirmationCode string) error {
	input := &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         aws.String(s.clientID),
		Username:         aws.String(email),
		ConfirmationCode: aws.String(confirmationCode),
	}

	_, err := s.client.ConfirmSignUp(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to confirm sign up: %w", err)
	}

	fmt.Printf("CognitoService: User signup confirmed for %s\n", email)
	return nil
}

// SignIn authenticates a user with Cognito
func (s *cognitoService) SignIn(ctx context.Context, email, password string) (*cognitoidentityprovider.InitiateAuthOutput, error) {
	input := &cognitoidentityprovider.InitiateAuthInput{
		ClientId: aws.String(s.clientID),
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		AuthParameters: map[string]string{
			"USERNAME": email,
			"PASSWORD": password,
		},
	}

	result, err := s.client.InitiateAuth(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to sign in user: %w", err)
	}

	fmt.Printf("CognitoService: User signed in successfully: %s\n", email)
	return result, nil
}

// SignOut signs out a user from Cognito
func (s *cognitoService) SignOut(ctx context.Context, accessToken string) error {
	input := &cognitoidentityprovider.GlobalSignOutInput{
		AccessToken: aws.String(accessToken),
	}

	_, err := s.client.GlobalSignOut(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to sign out user: %w", err)
	}

	fmt.Printf("CognitoService: User signed out successfully\n")
	return nil
}

// ResetPassword initiates password reset process
func (s *cognitoService) ResetPassword(ctx context.Context, email string) error {
	input := &cognitoidentityprovider.ForgotPasswordInput{
		ClientId: aws.String(s.clientID),
		Username: aws.String(email),
	}

	_, err := s.client.ForgotPassword(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to initiate password reset: %w", err)
	}

	fmt.Printf("CognitoService: Password reset initiated for %s\n", email)
	return nil
}

// ConfirmResetPassword confirms password reset with verification code
func (s *cognitoService) ConfirmResetPassword(ctx context.Context, email, confirmationCode, newPassword string) error {
	input := &cognitoidentityprovider.ConfirmForgotPasswordInput{
		ClientId:         aws.String(s.clientID),
		Username:         aws.String(email),
		ConfirmationCode: aws.String(confirmationCode),
		Password:         aws.String(newPassword),
	}

	_, err := s.client.ConfirmForgotPassword(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to confirm password reset: %w", err)
	}

	fmt.Printf("CognitoService: Password reset confirmed for %s\n", email)
	return nil
}

// GetUser retrieves user information from Cognito
func (s *cognitoService) GetUser(ctx context.Context, accessToken string) (*cognitoidentityprovider.GetUserOutput, error) {
	input := &cognitoidentityprovider.GetUserInput{
		AccessToken: aws.String(accessToken),
	}

	result, err := s.client.GetUser(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return result, nil
}
