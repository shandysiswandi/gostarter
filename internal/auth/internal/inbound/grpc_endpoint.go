package inbound

import (
	"context"

	pb "github.com/shandysiswandi/gostarter/api/gen-proto/auth"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
)

type GrpcEndpoint struct {
	pb.UnimplementedAuthServiceServer

	loginUC          domain.Login
	registerUC       domain.Register
	refreshTokenUC   domain.RefreshToken
	forgotPasswordUC domain.ForgotPassword
	resetPasswordUC  domain.ResetPassword
}

func (g *GrpcEndpoint) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
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

func (g *GrpcEndpoint) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	resp, err := g.registerUC.Call(ctx, domain.RegisterInput{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{Email: resp.Email}, nil
}

func (g *GrpcEndpoint) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (
	*pb.RefreshTokenResponse, error,
) {
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
	resp, err := g.forgotPasswordUC.Call(ctx, domain.ForgotPasswordInput{Email: req.GetEmail()})
	if err != nil {
		return nil, err
	}

	return &pb.ForgotPasswordResponse{Email: resp.Email, Message: resp.Message}, nil
}

func (g *GrpcEndpoint) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (
	*pb.ResetPasswordResponse, error,
) {
	resp, err := g.resetPasswordUC.Call(ctx, domain.ResetPasswordInput{
		Token: req.GetToken(), Password: req.GetPassword(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.ResetPasswordResponse{Message: resp.Message}, nil
}
