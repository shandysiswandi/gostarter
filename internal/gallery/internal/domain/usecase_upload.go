package domain

import "context"

type Upload interface {
	Call(context.Context, UploadInput) (*UploadOutput, error)
}

type UploadInput struct {
	File []byte `validate:"required"`
}

type UploadOutput struct {
	Path string
}
