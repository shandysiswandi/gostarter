package service

import (
	"context"

	"github.com/shandysiswandi/gostarter/internal/gallery/internal/domain"
)

type UploadStore interface{}

type Upload struct {
	store UploadStore
}

func NewUpload(store UploadStore) *Upload {
	return &Upload{store: store}
}

func (u *Upload) Call(context.Context, domain.UploadInput) (*domain.UploadOutput, error) {
	return &domain.UploadOutput{
		Path: "1234",
	}, nil
}
