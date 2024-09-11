package service

import (
	"context"
	"log"

	"github.com/shandysiswandi/gostarter/internal/gallery/internal/domain"
)

type GetImageStore interface{}

type GetImage struct {
	store GetImageStore
}

func NewGetImage(store GetImageStore) *GetImage {
	return &GetImage{store: store}
}

func (u *GetImage) Call(ctx context.Context, in domain.GetImageInput) (*domain.GetImageOutput, error) {
	log.Println("ctx", ctx, "in", in)
	return nil, nil
}
