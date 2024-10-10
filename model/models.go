package model

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type (
	Response struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Info    any    `json:"info"`
		Data    any    `json:"data"`
	}

	Request struct {
		Limit string `json:"limit" url:"limit"`
		Page  string `json:"page" url:"page"`
		Tag   string `json:"tag" url:"tag"`
	}

	GetPost struct {
		Id          *uuid.UUID `json:"id"`
		Title       *string    `json:"title"`
		Content     *string    `json:"content"`
		Status      *string    `json:"status"`
		PublishDate *time.Time `json:"publish_date"`
		Tags        []*GetTags `json:"tags"`
	}

	GetTags struct {
		Id    *uuid.UUID `json:"id"`
		Label *string    `json:"label"`
	}

	ReqPost struct {
		Title   string   `json:"title"`
		Content string   `json:"content"`
		Tags    []string `json:"tags"`
	}

	Credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	Claims struct {
		Username string `json:"username"`
		Role     string `json:"role"`
		jwt.RegisteredClaims
	}
)
