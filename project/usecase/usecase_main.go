package usecase

import (
	"context"
	"errors"
	"go-axiata/config"
	"go-axiata/model"
	"go-axiata/project"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type usecase struct {
	cfg        *config.Config
	logger     zerolog.Logger
	repository project.Repository
}

func NewUsecase(cfg *config.Config, logger zerolog.Logger, repository project.Repository) project.Usecase {
	return &usecase{
		cfg:        cfg,
		logger:     logger,
		repository: repository,
	}
}

func (r *usecase) GetPost(ctx context.Context, req model.Request) (res []*model.GetPost, err error) {

	res, err = r.repository.GetPost(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *usecase) DetailPost(ctx context.Context, id uuid.UUID) (res model.GetPost, err error) {

	res, err = r.repository.DetailPost(ctx, id)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (r *usecase) DeletePost(ctx context.Context, id uuid.UUID) (err error) {

	err = r.repository.DeletePost(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *usecase) CreatePost(ctx context.Context, req model.ReqPost) (id *uuid.UUID, err error) {

	if req.Content == "" || req.Title == "" || len(req.Tags) <= 0 {
		return nil, errors.New(" All fields must be filled in")
	}

	id, err = r.repository.CreatePost(ctx, req)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (r *usecase) UpdatePost(ctx context.Context, id uuid.UUID, req model.ReqPost) (err error) {

	if req.Content == "" || req.Title == "" || len(req.Tags) <= 0 {
		return errors.New(" All fields must be filled in")
	}

	err = r.repository.UpdatePost(ctx, id, req)
	if err != nil {
		return err
	}

	return nil
}

func (r *usecase) Register(ctx context.Context, req model.Credentials) (err error) {

	if req.Username == "" || req.Password == "" {
		return errors.New(" All fields must be filled in")
	}

	err = r.repository.Register(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (r *usecase) Login(ctx context.Context, username string) (res model.Credentials, err error) {

	if username == "" {
		return res, errors.New(" All fields must be filled in")
	}

	res, err = r.repository.Login(ctx, username)
	if err != nil {
		return res, err
	}

	return res, nil
}
