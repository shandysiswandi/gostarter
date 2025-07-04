package usecase

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/shandysiswandi/goreng/goerror"
	"github.com/shandysiswandi/goreng/mocker"
	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/mockz"
	"github.com/shandysiswandi/gostarter/internal/lib"
	"github.com/stretchr/testify/assert"
)

func TestNewLogin(t *testing.T) {
	tests := []struct {
		name string
		dep  Dependency
		s    LoginStore
		want *Login
	}{
		{
			name: "Success",
			dep:  Dependency{},
			s:    nil,
			want: &Login{tgs: &tokenGenSaver{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewLogin(tt.dep, tt.s)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLogin_Call(t *testing.T) {
	type args struct {
		ctx context.Context
		in  domain.LoginInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.LoginOutput
		wantErr error
		mockFn  func(a args) *Login
	}{
		{
			name: "ErrorValidationInput",
			args: args{
				ctx: context.Background(),
				in: domain.LoginInput{
					Email:    "email",
					Password: "password",
				},
			},
			want:    nil,
			wantErr: goerror.NewInvalidInput("Invalid request payload", assert.AnError),
			mockFn: func(a args) *Login {
				tel := telemetry.NewTelemetry()
				validatorMock := mocker.NewMockValidator(t)

				_, span := tel.Tracer().Start(a.ctx, "auth.usecase.Login")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(assert.AnError)

				return &Login{
					tel:       tel,
					validator: validatorMock,
					hash:      nil,
					secHash:   nil,
					jwt:       nil,
					store:     nil,
				}
			},
		},
		{
			name: "ErrorStoreUserByEmail",
			args: args{
				ctx: context.Background(),
				in: domain.LoginInput{
					Email:    "email",
					Password: "password",
				},
			},
			want:    nil,
			wantErr: goerror.NewServerInternal(assert.AnError),
			mockFn: func(a args) *Login {
				tel := telemetry.NewTelemetry()
				validatorMock := mocker.NewMockValidator(t)
				storeMock := mockz.NewMockLoginStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.Login")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				storeMock.EXPECT().
					UserByEmail(ctx, a.in.Email).
					Return(nil, assert.AnError)

				return &Login{
					tel:       tel,
					validator: validatorMock,
					hash:      nil,
					secHash:   nil,
					jwt:       nil,
					store:     storeMock,
				}
			},
		},
		{
			name: "ErrorStoreUserByEmailNotFound",
			args: args{
				ctx: context.Background(),
				in: domain.LoginInput{
					Email:    "email",
					Password: "password",
				},
			},
			want:    nil,
			wantErr: goerror.NewBusiness("Invalid credentials", goerror.CodeUnauthorized),
			mockFn: func(a args) *Login {
				tel := telemetry.NewTelemetry()
				validatorMock := mocker.NewMockValidator(t)
				storeMock := mockz.NewMockLoginStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.Login")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				storeMock.EXPECT().
					UserByEmail(ctx, a.in.Email).
					Return(nil, nil)

				return &Login{
					tel:       tel,
					validator: validatorMock,
					hash:      nil,
					secHash:   nil,
					jwt:       nil,
					store:     storeMock,
				}
			},
		},
		{
			name: "ErrorStoreFindUserByEmailNotVerified",
			args: args{
				ctx: context.Background(),
				in: domain.LoginInput{
					Email:    "email",
					Password: "password",
				},
			},
			want:    nil,
			wantErr: goerror.NewBusiness("Account isn't verified, please check email", goerror.CodeForbidden),
			mockFn: func(a args) *Login {
				tel := telemetry.NewTelemetry()
				validatorMock := mocker.NewMockValidator(t)
				storeMock := mockz.NewMockLoginStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.Login")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				user := &domain.User{
					ID:         10,
					Name:       "",
					Email:      "email",
					Password:   "password",
					VerifiedAt: sql.Null[time.Time]{Valid: false},
				}
				storeMock.EXPECT().
					UserByEmail(ctx, a.in.Email).
					Return(user, nil)

				return &Login{
					tel:       tel,
					validator: validatorMock,
					hash:      nil,
					secHash:   nil,
					jwt:       nil,
					store:     storeMock,
				}
			},
		},
		{
			name: "ErrorVerifyPassword",
			args: args{
				ctx: context.Background(),
				in: domain.LoginInput{
					Email:    "email",
					Password: "password",
				},
			},
			want:    nil,
			wantErr: goerror.NewBusiness("Invalid credentials", goerror.CodeUnauthorized),
			mockFn: func(a args) *Login {
				tel := telemetry.NewTelemetry()
				validatorMock := mocker.NewMockValidator(t)
				storeMock := mockz.NewMockLoginStore(t)
				hashMock := mocker.NewMockHash(t)

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.Login")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				user := &domain.User{
					ID:         10,
					Name:       "",
					Email:      "email",
					Password:   "password",
					VerifiedAt: sql.Null[time.Time]{Valid: true},
				}
				storeMock.EXPECT().
					UserByEmail(ctx, a.in.Email).
					Return(user, nil)

				hashMock.EXPECT().
					Verify(user.Password, a.in.Password).
					Return(false)

				return &Login{
					tel:       tel,
					validator: validatorMock,
					hash:      hashMock,
					secHash:   nil,
					jwt:       nil,
					store:     storeMock,
				}
			},
		},
		{
			name: "ErrorStoreTokenByUserID",
			args: args{
				ctx: context.Background(),
				in: domain.LoginInput{
					Email:    "email",
					Password: "password",
				},
			},
			want:    nil,
			wantErr: goerror.NewServerInternal(assert.AnError),
			mockFn: func(a args) *Login {
				tel := telemetry.NewTelemetry()
				validatorMock := mocker.NewMockValidator(t)
				storeMock := mockz.NewMockLoginStore(t)
				hashMock := mocker.NewMockHash(t)

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.Login")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				user := &domain.User{
					ID:         10,
					Name:       "",
					Email:      "email",
					Password:   "password",
					VerifiedAt: sql.Null[time.Time]{Valid: true},
				}
				storeMock.EXPECT().
					UserByEmail(ctx, a.in.Email).
					Return(user, nil)

				hashMock.EXPECT().
					Verify(user.Password, a.in.Password).
					Return(true)

				storeMock.EXPECT().
					TokenByUserID(ctx, user.ID).
					Return(nil, assert.AnError)

				return &Login{
					tel:       tel,
					validator: validatorMock,
					hash:      hashMock,
					secHash:   nil,
					jwt:       nil,
					store:     storeMock,
				}
			},
		},
		{
			name: "ErrorJWTGenerateAccessToken",
			args: args{
				ctx: context.Background(),
				in: domain.LoginInput{
					Email:    "email",
					Password: "password",
				},
			},
			want:    nil,
			wantErr: goerror.NewServerInternal(assert.AnError),
			mockFn: func(a args) *Login {
				tel := telemetry.NewTelemetry()
				validatorMock := mocker.NewMockValidator(t)
				storeMock := mockz.NewMockLoginStore(t)
				hashMock := mocker.NewMockHash(t)
				secHashMock := mocker.NewMockHash(t)
				clockMock := mocker.NewMockClocker(t)
				jwtMock := mocker.NewMockJWT(t)

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.Login")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				user := &domain.User{
					ID:         10,
					Name:       "",
					Email:      "email",
					Password:   "password",
					VerifiedAt: sql.Null[time.Time]{Valid: true},
				}
				storeMock.EXPECT().
					UserByEmail(ctx, a.in.Email).
					Return(user, nil)

				hashMock.EXPECT().
					Verify(user.Password, a.in.Password).
					Return(true)

				storeMock.EXPECT().
					TokenByUserID(ctx, user.ID).
					Return(nil, nil)

				now := time.Time{}
				clockMock.EXPECT().
					Now().
					Return(now)

				acClaim := lib.NewJWTClaim(
					user.ID,
					a.in.Email,
					now.Add(time.Hour),
					[]string{"gostarter.access.token"},
				)
				jwtMock.EXPECT().
					Generate(acClaim).
					Return("", assert.AnError).
					Once()

				return &Login{
					tel:       tel,
					validator: validatorMock,
					hash:      hashMock,
					secHash:   nil,
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
			name: "ErrorJWTGenerateRefreshToken",
			args: args{
				ctx: context.Background(),
				in: domain.LoginInput{
					Email:    "email",
					Password: "password",
				},
			},
			want:    nil,
			wantErr: goerror.NewServerInternal(assert.AnError),
			mockFn: func(a args) *Login {
				tel := telemetry.NewTelemetry()
				validatorMock := mocker.NewMockValidator(t)
				storeMock := mockz.NewMockLoginStore(t)
				hashMock := mocker.NewMockHash(t)
				secHashMock := mocker.NewMockHash(t)
				clockMock := mocker.NewMockClocker(t)
				jwtMock := mocker.NewMockJWT(t)

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.Login")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				user := &domain.User{
					ID:         10,
					Name:       "",
					Email:      "email",
					Password:   "password",
					VerifiedAt: sql.Null[time.Time]{Valid: true},
				}
				storeMock.EXPECT().
					UserByEmail(ctx, a.in.Email).
					Return(user, nil)

				hashMock.EXPECT().
					Verify(user.Password, a.in.Password).
					Return(true)

				storeMock.EXPECT().
					TokenByUserID(ctx, user.ID).
					Return(nil, nil)

				now := time.Time{}
				clockMock.EXPECT().
					Now().
					Return(now)

				acClaim := lib.NewJWTClaim(
					user.ID,
					a.in.Email,
					now.Add(time.Hour),
					[]string{"gostarter.access.token"},
				)
				jwtMock.EXPECT().
					Generate(acClaim).
					Return("access_token", nil).
					Once()

				refClaim := lib.NewJWTClaim(
					user.ID,
					a.in.Email,
					now.Add(time.Hour*24),
					[]string{"gostarter.refresh.token"},
				)
				jwtMock.EXPECT().
					Generate(refClaim).
					Return("", assert.AnError).
					Once()

				return &Login{
					tel:       tel,
					validator: validatorMock,
					hash:      hashMock,
					secHash:   nil,
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
			name: "ErrorSecurityHashAccessToken",
			args: args{
				ctx: context.Background(),
				in: domain.LoginInput{
					Email:    "email",
					Password: "password",
				},
			},
			want:    nil,
			wantErr: goerror.NewServerInternal(assert.AnError),
			mockFn: func(a args) *Login {
				tel := telemetry.NewTelemetry()
				validatorMock := mocker.NewMockValidator(t)
				storeMock := mockz.NewMockLoginStore(t)
				hashMock := mocker.NewMockHash(t)
				secHashMock := mocker.NewMockHash(t)
				clockMock := mocker.NewMockClocker(t)
				jwtMock := mocker.NewMockJWT(t)

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.Login")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				user := &domain.User{
					ID:         10,
					Name:       "",
					Email:      "email",
					Password:   "password",
					VerifiedAt: sql.Null[time.Time]{Valid: true},
				}
				storeMock.EXPECT().
					UserByEmail(ctx, a.in.Email).
					Return(user, nil)

				hashMock.EXPECT().
					Verify(user.Password, a.in.Password).
					Return(true)

				storeMock.EXPECT().
					TokenByUserID(ctx, user.ID).
					Return(nil, nil)

				now := time.Time{}
				clockMock.EXPECT().
					Now().
					Return(now)

				acClaim := lib.NewJWTClaim(
					user.ID,
					a.in.Email,
					now.Add(time.Hour),
					[]string{"gostarter.access.token"},
				)
				jwtMock.EXPECT().
					Generate(acClaim).
					Return("access_token", nil).
					Once()

				refClaim := lib.NewJWTClaim(
					user.ID,
					a.in.Email,
					now.Add(time.Hour*24),
					[]string{"gostarter.refresh.token"},
				)
				jwtMock.EXPECT().
					Generate(refClaim).
					Return("refresh_token", nil).
					Once()

				secHashMock.EXPECT().
					Hash("access_token").
					Return(nil, assert.AnError).
					Once()

				return &Login{
					tel:       tel,
					validator: validatorMock,
					hash:      hashMock,
					secHash:   nil,
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
			name: "ErrorSecurityHashRefreshToken",
			args: args{
				ctx: context.Background(),
				in: domain.LoginInput{
					Email:    "email",
					Password: "password",
				},
			},
			want:    nil,
			wantErr: goerror.NewServerInternal(assert.AnError),
			mockFn: func(a args) *Login {
				tel := telemetry.NewTelemetry()
				validatorMock := mocker.NewMockValidator(t)
				storeMock := mockz.NewMockLoginStore(t)
				hashMock := mocker.NewMockHash(t)
				secHashMock := mocker.NewMockHash(t)
				clockMock := mocker.NewMockClocker(t)
				jwtMock := mocker.NewMockJWT(t)

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.Login")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				user := &domain.User{
					ID:         10,
					Name:       "",
					Email:      "email",
					Password:   "password",
					VerifiedAt: sql.Null[time.Time]{Valid: true},
				}
				storeMock.EXPECT().
					UserByEmail(ctx, a.in.Email).
					Return(user, nil)

				hashMock.EXPECT().
					Verify(user.Password, a.in.Password).
					Return(true)

				storeMock.EXPECT().
					TokenByUserID(ctx, user.ID).
					Return(nil, nil)

				now := time.Time{}
				clockMock.EXPECT().
					Now().
					Return(now)

				acClaim := lib.NewJWTClaim(
					user.ID,
					a.in.Email,
					now.Add(time.Hour),
					[]string{"gostarter.access.token"},
				)
				jwtMock.EXPECT().
					Generate(acClaim).
					Return("access_token", nil).
					Once()

				refClaim := lib.NewJWTClaim(
					user.ID,
					a.in.Email,
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
					Return(nil, assert.AnError).
					Once()

				return &Login{
					tel:       tel,
					validator: validatorMock,
					hash:      hashMock,
					secHash:   nil,
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
			name: "ErrorStoreTokenSave",
			args: args{
				ctx: context.Background(),
				in: domain.LoginInput{
					Email:    "email",
					Password: "password",
				},
			},
			want:    nil,
			wantErr: goerror.NewServerInternal(assert.AnError),
			mockFn: func(a args) *Login {
				tel := telemetry.NewTelemetry()
				validatorMock := mocker.NewMockValidator(t)
				storeMock := mockz.NewMockLoginStore(t)
				hashMock := mocker.NewMockHash(t)
				secHashMock := mocker.NewMockHash(t)
				idnumMock := new(mocker.MockNumberID)
				clockMock := mocker.NewMockClocker(t)
				jwtMock := mocker.NewMockJWT(t)

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.Login")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				user := &domain.User{
					ID:         10,
					Name:       "",
					Email:      "email",
					Password:   "password",
					VerifiedAt: sql.Null[time.Time]{Valid: true},
				}
				storeMock.EXPECT().
					UserByEmail(ctx, a.in.Email).
					Return(user, nil)

				hashMock.EXPECT().
					Verify(user.Password, a.in.Password).
					Return(true)

				storeMock.EXPECT().
					TokenByUserID(ctx, user.ID).
					Return(nil, nil)

				now := time.Time{}
				clockMock.EXPECT().
					Now().
					Return(now)

				acClaim := lib.NewJWTClaim(
					user.ID,
					a.in.Email,
					now.Add(time.Hour),
					[]string{"gostarter.access.token"},
				)
				jwtMock.EXPECT().
					Generate(acClaim).
					Return("access_token", nil).
					Once()

				refClaim := lib.NewJWTClaim(
					user.ID,
					a.in.Email,
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

				idnumMock.EXPECT().
					Generate().
					Return(111)

				tokenIn := domain.Token{
					ID:               111,
					UserID:           10,
					AccessToken:      "hash_access_token",
					RefreshToken:     "hash_refresh_token",
					AccessExpiresAt:  time.Time{}.Add(time.Hour),
					RefreshExpiresAt: time.Time{}.Add(time.Hour * 24),
				}
				storeMock.EXPECT().
					TokenSave(ctx, tokenIn).
					Return(assert.AnError)

				return &Login{
					tel:       tel,
					validator: validatorMock,
					hash:      hashMock,
					secHash:   nil,
					jwt:       jwtMock,
					store:     storeMock,
					clock:     clockMock,
					tgs: &tokenGenSaver{
						uidnumber: idnumMock,
						jwt:       jwtMock,
						tel:       tel,
						secHash:   secHashMock,
						clock:     clockMock,
						ts:        storeMock,
					},
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				in: domain.LoginInput{
					Email:    "email",
					Password: "password",
				},
			},
			want: &domain.LoginOutput{
				AccessToken:      "access_token",
				RefreshToken:     "refresh_token",
				AccessExpiresIn:  3600,
				RefreshExpiresIn: 86400,
			},
			wantErr: nil,
			mockFn: func(a args) *Login {
				tel := telemetry.NewTelemetry()
				validatorMock := mocker.NewMockValidator(t)
				storeMock := mockz.NewMockLoginStore(t)
				hashMock := mocker.NewMockHash(t)
				secHashMock := mocker.NewMockHash(t)
				clockMock := mocker.NewMockClocker(t)
				jwtMock := mocker.NewMockJWT(t)

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.Login")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				user := &domain.User{
					ID:         10,
					Name:       "",
					Email:      "email",
					Password:   "password",
					VerifiedAt: sql.Null[time.Time]{Valid: true},
				}
				storeMock.EXPECT().
					UserByEmail(ctx, a.in.Email).
					Return(user, nil)

				hashMock.EXPECT().
					Verify(user.Password, a.in.Password).
					Return(true)

				token := &domain.Token{ID: 90, UserID: user.ID}
				storeMock.EXPECT().
					TokenByUserID(ctx, user.ID).
					Return(token, nil)

				now := time.Time{}
				clockMock.EXPECT().
					Now().
					Return(now)

				acClaim := lib.NewJWTClaim(
					user.ID,
					a.in.Email,
					now.Add(time.Hour),
					[]string{"gostarter.access.token"},
				)
				jwtMock.EXPECT().
					Generate(acClaim).
					Return("access_token", nil).
					Once()

				refClaim := lib.NewJWTClaim(
					user.ID,
					a.in.Email,
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
					ID:               90,
					UserID:           10,
					AccessToken:      "hash_access_token",
					RefreshToken:     "hash_refresh_token",
					AccessExpiresAt:  time.Time{}.Add(time.Hour),
					RefreshExpiresAt: time.Time{}.Add(time.Hour * 24),
				}
				storeMock.EXPECT().
					TokenUpdate(ctx, tokenIn).
					Return(nil)

				return &Login{
					tel:       tel,
					validator: validatorMock,
					hash:      hashMock,
					secHash:   nil,
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
