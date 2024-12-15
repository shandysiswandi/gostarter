package inbound

import (
	"context"

	pb "github.com/shandysiswandi/gostarter/api/gen-proto/auth"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
)

type grpcEndpoint struct {
	pb.UnimplementedAuthServiceServer

	telemetry *telemetry.Telemetry

	loginUC          domain.Login
	registerUC       domain.Register
	refreshTokenUC   domain.RefreshToken
	forgotPasswordUC domain.ForgotPassword
	resetPasswordUC  domain.ResetPassword
}

func (g *grpcEndpoint) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	ctx, span := g.telemetry.Tracer().Start(ctx, "auth.inbound.grpcEndpoint.Login")
	defer span.End()

	resp, err := g.loginUC.Call(ctx, domain.LoginInput{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	})
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

func (g *grpcEndpoint) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	ctx, span := g.telemetry.Tracer().Start(ctx, "auth.inbound.grpcEndpoint.Register")
	defer span.End()

	resp, err := g.registerUC.Call(ctx, domain.RegisterInput{
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{Email: resp.Email}, nil
}

func (g *grpcEndpoint) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (
	*pb.RefreshTokenResponse, error,
) {
	ctx, span := g.telemetry.Tracer().Start(ctx, "auth.inbound.grpcEndpoint.RefreshToken")
	defer span.End()

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

func (g *grpcEndpoint) ForgotPassword(ctx context.Context, req *pb.ForgotPasswordRequest) (
	*pb.ForgotPasswordResponse, error,
) {
	ctx, span := g.telemetry.Tracer().Start(ctx, "auth.inbound.grpcEndpoint.ForgotPassword")
	defer span.End()

	resp, err := g.forgotPasswordUC.Call(ctx, domain.ForgotPasswordInput{Email: req.GetEmail()})
	if err != nil {
		return nil, err
	}

	return &pb.ForgotPasswordResponse{
		Email:   resp.Email,
		Message: resp.Message,
	}, nil
}

func (g *grpcEndpoint) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (
	*pb.ResetPasswordResponse, error,
) {
	ctx, span := g.telemetry.Tracer().Start(ctx, "auth.inbound.grpcEndpoint.ResetPassword")
	defer span.End()

	resp, err := g.resetPasswordUC.Call(ctx, domain.ResetPasswordInput{
		Token:    req.GetToken(),
		Password: req.GetPassword(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.ResetPasswordResponse{Message: resp.Message}, nil
}
