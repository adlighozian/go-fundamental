package project

import (
	"context"
	"go-axiata/model"

	"github.com/google/uuid"
)

type Repository interface {
	GetPost(ctx context.Context, req model.Request) (res []*model.GetPost, err error)
	DetailPost(ctx context.Context, id uuid.UUID) (res model.GetPost, err error)
	DeletePost(ctx context.Context, id uuid.UUID) (err error)
	CreatePost(ctx context.Context, req model.ReqPost) (id *uuid.UUID, err error)
	UpdatePost(ctx context.Context, id uuid.UUID, req model.ReqPost) (err error)
	Register(ctx context.Context, req model.Credentials) (err error)
	Login(ctx context.Context, username string) (res model.Credentials, err error)
}
