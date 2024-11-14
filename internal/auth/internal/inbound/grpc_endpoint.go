package inbound

import (
	"context"

	pb "github.com/shandysiswandi/gostarter/api/gen-proto/auth"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type GrpcEndpoint struct {
	pb.UnimplementedAuthServiceServer

	telemetry *telemetry.Telemetry
	validator validation.Validator

	loginUC          domain.Login
	registerUC       domain.Register
	refreshTokenUC   domain.RefreshToken
	forgotPasswordUC domain.ForgotPassword
	resetPasswordUC  domain.ResetPassword
}

func NewGrpcEndpoint(
	telemetry *telemetry.Telemetry,
	validator validation.Validator,
	loginUC domain.Login,
	registerUC domain.Register,
	refreshTokenUC domain.RefreshToken,
	forgotPasswordUC domain.ForgotPassword,
	resetPasswordUC domain.ResetPassword,
) *GrpcEndpoint {
	return &GrpcEndpoint{
		telemetry:        telemetry,
		validator:        validator,
		loginUC:          loginUC,
		registerUC:       registerUC,
		refreshTokenUC:   refreshTokenUC,
		forgotPasswordUC: forgotPasswordUC,
		resetPasswordUC:  resetPasswordUC,
	}
}

func (g *GrpcEndpoint) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if err := g.validator.Validate(req); err != nil {
		return nil, goerror.NewInvalidInput("validation failed", err)
	}

	resp, err := g.loginUC.Call(ctx, domain.LoginInput{Email: req.GetEmail(), Password: req.GetPassword()})
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		AccessToken:      resp.AccessToken,
		RefreshToken:     resp.RefreshToken,
		AccessExpiresIn:  resp.AccessExpiresIn,
		RefreshExpiresIn: resp.RefreshExpiresIn,
	}, nil
}

func (g *GrpcEndpoint) Register(ctx context.Context, req *pb.RegisterRequest) (
	*pb.RegisterResponse, error,
) {
	if err := g.validator.Validate(req); err != nil {
		return nil, goerror.NewInvalidInput("validation failed", err)
	}

	_, err := g.registerUC.Call(ctx, domain.RegisterInput{Email: req.GetEmail(), Password: req.GetPassword()})
	if err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{Email: req.GetEmail()}, nil
}

func (g *GrpcEndpoint) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (
	*pb.RefreshTokenResponse, error,
) {
	if err := g.validator.Validate(req); err != nil {
		return nil, goerror.NewInvalidInput("validation failed", err)
	}

	resp, err := g.refreshTokenUC.Call(ctx, domain.RefreshTokenInput{RefreshToken: req.GetRefreshToken()})
	if err != nil {
		return nil, err
	}

	return &pb.RefreshTokenResponse{
		AccessToken:      resp.AccessToken,
		RefreshToken:     resp.RefreshToken,
		AccessExpiresIn:  resp.AccessExpiresIn,
		RefreshExpiresIn: resp.RefreshExpiresIn,
	}, nil
}

func (g *GrpcEndpoint) ForgotPassword(ctx context.Context, req *pb.ForgotPasswordRequest) (
	*pb.ForgotPasswordResponse, error,
) {
	if err := g.validator.Validate(req); err != nil {
		return nil, goerror.NewInvalidInput("validation failed", err)
	}

	_, err := g.forgotPasswordUC.Call(ctx, domain.ForgotPasswordInput{Email: req.GetEmail()})
	if err != nil {
		return nil, err
	}

	return &pb.ForgotPasswordResponse{
		Email:   req.GetEmail(),
		Message: "If an account with this email exists, you'll receive a password reset email shortly.",
	}, nil
}

func (g *GrpcEndpoint) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (
	*pb.ResetPasswordResponse, error,
) {
	if err := g.validator.Validate(req); err != nil {
		return nil, goerror.NewInvalidInput("validation failed", err)
	}

	_, err := g.resetPasswordUC.Call(ctx, domain.ResetPasswordInput{
		Token: req.GetToken(), Password: req.GetPassword(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.ResetPasswordResponse{
		Message: "Your password has been successfully reset.",
	}, nil
}
