package domain

import "context"

type Upload interface {
	Call(ctx context.Context, in UploadInput) (*UploadOutput, error)
}

type UploadInput struct {
	File []byte `validate:"required"`
}

type UploadOutput struct {
	Path string
}
