package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/mockz"
	mockClock "github.com/shandysiswandi/gostarter/pkg/clock/mocker"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	mockHash "github.com/shandysiswandi/gostarter/pkg/hash/mocker"
	"github.com/shandysiswandi/gostarter/pkg/jwt"
	mockJwt "github.com/shandysiswandi/gostarter/pkg/jwt/mocker"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	mockValidation "github.com/shandysiswandi/gostarter/pkg/validation/mocker"
	"github.com/stretchr/testify/assert"
)

func TestNewRefreshToken(t *testing.T) {
	type args struct {
		dep Dependency
		s   RefreshTokenStore
	}
	tests := []struct {
		name string
		args args
		want *RefreshToken
	}{
		{
			name: "Success",
			args: args{},
			want: &RefreshToken{tgs: &tokenGenSaver{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewRefreshToken(tt.args.dep, tt.args.s)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRefreshToken_Call(t *testing.T) {
	validToken := "none.eyJzdWIiOiJ0ZXN0IiwiYXV0aF9pZCI6IjEwMSJ9.none"

	type args struct {
		ctx context.Context
		in  domain.RefreshTokenInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.RefreshTokenOutput
		wantErr error
		mockFn  func(a args) *RefreshToken
	}{
		{
			name: "ErrorValidationInput",
			args: args{
				ctx: context.Background(),
				in:  domain.RefreshTokenInput{RefreshToken: "token"},
			},
			want:    nil,
			wantErr: goerror.NewInvalidInput("validation input fail", assert.AnError),
			mockFn: func(a args) *RefreshToken {
				tel := telemetry.NewTelemetry()
				validatorMock := mockValidation.NewMockValidator(t)

				_, span := tel.Tracer().Start(a.ctx, "auth.usecase.RefreshToken")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(assert.AnError)

				return &RefreshToken{
					telemetry: tel,
					validator: validatorMock,
					secHash:   nil,
					jwt:       nil,
					store:     nil,
					clock:     nil,
					tgs:       nil,
				}
			},
		},
		{
			name: "ErrorSecHash",
			args: args{
				ctx: context.Background(),
				in:  domain.RefreshTokenInput{RefreshToken: "token"},
			},
			want:    nil,
			wantErr: goerror.NewServerInternal(assert.AnError),
			mockFn: func(a args) *RefreshToken {
				tel := telemetry.NewTelemetry()
				validatorMock := mockValidation.NewMockValidator(t)
				secHashMock := mockHash.NewMockHash(t)

				_, span := tel.Tracer().Start(a.ctx, "auth.usecase.RefreshToken")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				secHashMock.EXPECT().
					Hash(a.in.RefreshToken).
					Return(nil, assert.AnError)

				return &RefreshToken{
					telemetry: tel,
					validator: validatorMock,
					secHash:   secHashMock,
					jwt:       nil,
					store:     nil,
					clock:     nil,
					tgs:       nil,
				}
			},
		},
		{
			name: "ErrorStoreFindTokenByRefresh",
			args: args{
				ctx: context.Background(),
				in:  domain.RefreshTokenInput{RefreshToken: "token"},
			},
			want:    nil,
			wantErr: goerror.NewServerInternal(assert.AnError),
			mockFn: func(a args) *RefreshToken {
				tel := telemetry.NewTelemetry()
				validatorMock := mockValidation.NewMockValidator(t)
				secHashMock := mockHash.NewMockHash(t)
				storeMock := mockz.NewMockRefreshTokenStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.RefreshToken")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				secHashMock.EXPECT().
					Hash(a.in.RefreshToken).
					Return([]byte("hash_refresh_token"), nil)

				storeMock.EXPECT().
					FindTokenByRefresh(ctx, "hash_refresh_token").
					Return(nil, assert.AnError)

				return &RefreshToken{
					telemetry: tel,
					validator: validatorMock,
					secHash:   secHashMock,
					jwt:       nil,
					store:     storeMock,
					clock:     nil,
					tgs:       nil,
				}
			},
		},
		{
			name: "ErrorStoreFindTokenByRefreshNotFound",
			args: args{
				ctx: context.Background(),
				in:  domain.RefreshTokenInput{RefreshToken: "token"},
			},
			want:    nil,
			wantErr: goerror.NewBusiness("invalid credentials", goerror.CodeUnauthorized),
			mockFn: func(a args) *RefreshToken {
				tel := telemetry.NewTelemetry()
				validatorMock := mockValidation.NewMockValidator(t)
				secHashMock := mockHash.NewMockHash(t)
				storeMock := mockz.NewMockRefreshTokenStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.RefreshToken")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				secHashMock.EXPECT().
					Hash(a.in.RefreshToken).
					Return([]byte("hash_refresh_token"), nil)

				storeMock.EXPECT().
					FindTokenByRefresh(ctx, "hash_refresh_token").
					Return(nil, nil)

				return &RefreshToken{
					telemetry: tel,
					validator: validatorMock,
					secHash:   secHashMock,
					jwt:       nil,
					store:     storeMock,
					clock:     nil,
					tgs:       nil,
				}
			},
		},
		{
			name: "ErrorTokenExpired",
			args: args{
				ctx: context.Background(),
				in:  domain.RefreshTokenInput{RefreshToken: "token"},
			},
			want:    nil,
			wantErr: goerror.NewBusiness("token has expired", goerror.CodeUnauthorized),
			mockFn: func(a args) *RefreshToken {
				tel := telemetry.NewTelemetry()
				validatorMock := mockValidation.NewMockValidator(t)
				secHashMock := mockHash.NewMockHash(t)
				storeMock := mockz.NewMockRefreshTokenStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.RefreshToken")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				secHashMock.EXPECT().
					Hash(a.in.RefreshToken).
					Return([]byte("hash_refresh_token"), nil)

				token := &domain.Token{
					ID:               10,
					UserID:           20,
					AccessToken:      "access",
					RefreshToken:     "refresh",
					AccessExpiredAt:  time.Time{},
					RefreshExpiredAt: time.Time{},
				}
				storeMock.EXPECT().
					FindTokenByRefresh(ctx, "hash_refresh_token").
					Return(token, nil)

				return &RefreshToken{
					telemetry: tel,
					validator: validatorMock,
					secHash:   secHashMock,
					jwt:       nil,
					store:     storeMock,
					clock:     nil,
					tgs:       nil,
				}
			},
		},
		{
			name: "ErrorExtractClaimFromToken",
			args: args{
				ctx: context.Background(),
				in:  domain.RefreshTokenInput{RefreshToken: "token"},
			},
			want:    nil,
			wantErr: goerror.NewBusiness("invalid credentials", goerror.CodeUnauthorized),
			mockFn: func(a args) *RefreshToken {
				tel := telemetry.NewTelemetry()
				validatorMock := mockValidation.NewMockValidator(t)
				secHashMock := mockHash.NewMockHash(t)
				storeMock := mockz.NewMockRefreshTokenStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.RefreshToken")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				secHashMock.EXPECT().
					Hash(a.in.RefreshToken).
					Return([]byte("hash_refresh_token"), nil)

				token := &domain.Token{
					ID:               10,
					UserID:           20,
					AccessToken:      "access",
					RefreshToken:     "refresh",
					AccessExpiredAt:  time.Time{},
					RefreshExpiredAt: time.Now().Add(time.Minute),
				}
				storeMock.EXPECT().
					FindTokenByRefresh(ctx, "hash_refresh_token").
					Return(token, nil)

				return &RefreshToken{
					telemetry: tel,
					validator: validatorMock,
					secHash:   secHashMock,
					jwt:       nil,
					store:     storeMock,
					clock:     nil,
					tgs:       nil,
				}
			},
		},
		{
			name: "ErrorGenerateAndSaveToken",
			args: args{
				ctx: context.Background(),
				in:  domain.RefreshTokenInput{RefreshToken: validToken},
			},
			want:    nil,
			wantErr: goerror.NewServerInternal(assert.AnError),
			mockFn: func(a args) *RefreshToken {
				tel := telemetry.NewTelemetry()
				validatorMock := mockValidation.NewMockValidator(t)
				secHashMock := mockHash.NewMockHash(t)
				storeMock := mockz.NewMockRefreshTokenStore(t)
				clockMock := mockClock.NewMockClocker(t)
				jwtMock := mockJwt.NewMockJWT(t)

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.RefreshToken")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				secHashMock.EXPECT().
					Hash(a.in.RefreshToken).
					Return([]byte("hash_refresh_token"), nil)

				token := &domain.Token{
					ID:               10,
					UserID:           20,
					AccessToken:      "access",
					RefreshToken:     "refresh",
					AccessExpiredAt:  time.Now().Add(time.Minute),
					RefreshExpiredAt: time.Now().Add(time.Minute),
				}
				storeMock.EXPECT().
					FindTokenByRefresh(ctx, "hash_refresh_token").
					Return(token, nil)

				now := time.Time{}
				clockMock.EXPECT().
					Now().
					Return(now)

				email := "test"
				acClaim := jwt.NewClaim(
					token.UserID,
					email,
					now.Add(time.Hour),
					[]string{"gostarter.access.token"},
				)
				jwtMock.EXPECT().
					Generate(acClaim).
					Return("access_token", nil).
					Once()

				refClaim := jwt.NewClaim(
					token.UserID,
					email,
					now.Add(time.Hour*24),
					[]string{"gostarter.refresh.token"},
				)
				jwtMock.EXPECT().
					Generate(refClaim).
					Return("refresh_token", nil).
					Once()

				secHashMock.EXPECT().
					Hash("access_token").
					Return([]byte("hash_access_token"), nil).
					Once()

				secHashMock.EXPECT().
					Hash("refresh_token").
					Return([]byte("hash_refresh_token"), nil).
					Once()

				tokenIn := domain.Token{
					ID:               token.ID,
					UserID:           token.UserID,
					AccessToken:      "hash_access_token",
					RefreshToken:     "hash_refresh_token",
					AccessExpiredAt:  time.Time{}.Add(time.Hour),
					RefreshExpiredAt: time.Time{}.Add(time.Hour * 24),
				}
				storeMock.EXPECT().
					UpdateToken(ctx, tokenIn).
					Return(assert.AnError)

				return &RefreshToken{
					telemetry: tel,
					validator: validatorMock,
					secHash:   secHashMock,
					jwt:       jwtMock,
					store:     storeMock,
					clock:     clockMock,
					tgs: &tokenGenSaver{
						jwt:     jwtMock,
						tel:     tel,
						secHash: secHashMock,
						clock:   clockMock,
						ts:      storeMock,
					},
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				in:  domain.RefreshTokenInput{RefreshToken: validToken},
			},
			want: &domain.RefreshTokenOutput{
				AccessToken:      "access_token",
				RefreshToken:     "refresh_token",
				AccessExpiresIn:  3600,
				RefreshExpiresIn: 86400,
			},
			wantErr: nil,
			mockFn: func(a args) *RefreshToken {
				tel := telemetry.NewTelemetry()
				validatorMock := mockValidation.NewMockValidator(t)
				secHashMock := mockHash.NewMockHash(t)
				storeMock := mockz.NewMockRefreshTokenStore(t)
				clockMock := mockClock.NewMockClocker(t)
				jwtMock := mockJwt.NewMockJWT(t)

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.RefreshToken")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				secHashMock.EXPECT().
					Hash(a.in.RefreshToken).
					Return([]byte("hash_refresh_token"), nil)

				token := &domain.Token{
					ID:               10,
					UserID:           20,
					AccessToken:      "access",
					RefreshToken:     "refresh",
					AccessExpiredAt:  time.Now().Add(time.Minute),
					RefreshExpiredAt: time.Now().Add(time.Minute),
				}
				storeMock.EXPECT().
					FindTokenByRefresh(ctx, "hash_refresh_token").
					Return(token, nil)

				now := time.Time{}
				clockMock.EXPECT().
					Now().
					Return(now)

				email := "test"
				acClaim := jwt.NewClaim(
					token.UserID,
					email,
					now.Add(time.Hour),
					[]string{"gostarter.access.token"},
				)
				jwtMock.EXPECT().
					Generate(acClaim).
					Return("access_token", nil).
					Once()

				refClaim := jwt.NewClaim(
					token.UserID,
					email,
					now.Add(time.Hour*24),
					[]string{"gostarter.refresh.token"},
				)
				jwtMock.EXPECT().
					Generate(refClaim).
					Return("refresh_token", nil).
					Once()

				secHashMock.EXPECT().
					Hash("access_token").
					Return([]byte("hash_access_token"), nil).
					Once()

				secHashMock.EXPECT().
					Hash("refresh_token").
					Return([]byte("hash_refresh_token"), nil).
					Once()

				tokenIn := domain.Token{
					ID:               token.ID,
					UserID:           token.UserID,
					AccessToken:      "hash_access_token",
					RefreshToken:     "hash_refresh_token",
					AccessExpiredAt:  time.Time{}.Add(time.Hour),
					RefreshExpiredAt: time.Time{}.Add(time.Hour * 24),
				}
				storeMock.EXPECT().
					UpdateToken(ctx, tokenIn).
					Return(nil)

				return &RefreshToken{
					telemetry: tel,
					validator: validatorMock,
					secHash:   secHashMock,
					jwt:       jwtMock,
					store:     storeMock,
					clock:     clockMock,
					tgs: &tokenGenSaver{
						jwt:     jwtMock,
						tel:     tel,
						secHash: secHashMock,
						clock:   clockMock,
						ts:      storeMock,
					},
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := tt.mockFn(tt.args)
			got, err := s.Call(tt.args.ctx, tt.args.in)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
