package repository

import (
	"context"
	"database/sql"
	"errors"
	"go-axiata/config"
	"go-axiata/model"
	"go-axiata/pkg/helper"
	"go-axiata/project"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type repository struct {
	db     *sql.DB
	cfg    *config.Config
	logger zerolog.Logger
}

func NewRepository(db *sql.DB, cfg *config.Config, logger zerolog.Logger) project.Repository {
	return &repository{
		db:     db,
		cfg:    cfg,
		logger: logger,
	}
}

func (r *repository) GetPost(ctx context.Context, req model.Request) (res []*model.GetPost, err error) {
	query, args := helper.AddConditions(getPosts, req)

	// get posts start
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.Error().Msg("key:tQ0oiZKPIFN5b " + err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := &model.GetPost{}
		if err := rows.Scan(
			&user.Id,
			&user.Title,
			&user.Content,
			&user.Status,
			&user.PublishDate,
		); err != nil {
			r.logger.Error().Msg("key:e8o8sBz5eXX " + err.Error())
			return nil, err
		}
		// get posts end

		// get tags start
		rowsTags, err := r.db.QueryContext(ctx, getTags, user.Id)
		if err != nil {
			r.logger.Error().Msg("key:eP4hwstpiwydTDiHNS2L " + err.Error())
			return nil, err
		}
		defer rowsTags.Close()

		for rowsTags.Next() {
			tags := &model.GetTags{}
			if err := rowsTags.Scan(
				&tags.Id,
				&tags.Label,
			); err != nil {
				r.logger.Error().Msg("key:OnX34WOeX5KZjc " + err.Error())
				return nil, err
			}
			user.Tags = append(user.Tags, tags)
		}
		// get tags end

		res = append(res, user)
	}

	return res, nil
}

func (r *repository) DetailPost(ctx context.Context, id uuid.UUID) (res model.GetPost, err error) {
	// check data start
	check := 0
	err = r.db.QueryRowContext(ctx, checkDetailPosts, id).Scan(&check)
	if err != nil {
		r.logger.Error().Msg("key:j8c26nJ7qwnLwC " + err.Error())
		return res, err
	}

	if check <= 0 {
		return res, errors.New("not_found")
	}
	// check data end

	// get posts start
	err = r.db.QueryRowContext(ctx, detailPosts, id).Scan(
		&res.Id,
		&res.Title,
		&res.Content,
		&res.Status,
		&res.PublishDate,
	)
	if err != nil {
		r.logger.Error().Msg("key:dSiQvUi " + err.Error())
		return res, err
	}
	// get posts end

	// get tags start
	rowsTags, err := r.db.QueryContext(ctx, getTags, id)
	if err != nil {
		r.logger.Error().Msg("key:eP4hwstpiwydTDiHNS2L " + err.Error())
		return res, err
	}
	defer rowsTags.Close()

	for rowsTags.Next() {
		tags := &model.GetTags{}
		if err := rowsTags.Scan(
			&tags.Id,
			&tags.Label,
		); err != nil {
			r.logger.Error().Msg("key:OnX34WOeX5KZjc " + err.Error())
			return res, err
		}
		res.Tags = append(res.Tags, tags)
	}
	// get tags end

	return res, nil
}

func (r *repository) DeletePost(ctx context.Context, id uuid.UUID) (err error) {
	// check data start
	check := 0
	err = r.db.QueryRowContext(ctx, checkDetailPosts, id).Scan(&check)
	if err != nil {
		r.logger.Error().Msg("key:j8c26nJ7qwnLwC " + err.Error())
		return err
	}

	if check <= 0 {
		return errors.New("not_found")
	}
	// check data end

	// delete data start
	result, err := r.db.ExecContext(ctx, deletePosts, id)
	if err != nil {
		r.logger.Error().Msg("key:8ZKSXbRmNo " + err.Error())
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error().Msg("key:QF6Y7KmG5A3kRhuBGLWa " + err.Error())
		return err
	}

	if rowsAffected <= 0 {
		return errors.New(" Failed to delete the data. No rows were affected")
	}
	// delete data end

	return nil
}

func (r *repository) CreatePost(ctx context.Context, req model.ReqPost) (id *uuid.UUID, err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Error().Msg("Failed to begin transaction: " + err.Error())
		return nil, err
	}
	defer tx.Rollback()

	// cerate data start
	err = tx.QueryRowContext(ctx, createPost, req.Title, req.Content).Scan(&id)
	if err != nil {
		r.logger.Error().Msg("Failed to create post: " + err.Error())
		return nil, err
	}

	stmt, err := tx.PrepareContext(ctx, creatPostsTag)
	if err != nil {
		r.logger.Error().Msg("key:Q121wDimSRfXQ3pjUEm " + err.Error())
		return nil, err
	}
	defer stmt.Close()

	for _, v := range req.Tags {

		idTags := uuid.New()

		err = tx.QueryRowContext(ctx, getIdTag, v).Scan(&idTags)
		if err != nil {
			r.logger.Error().Msg("key:oJmem0FJTwBpgv2UwB " + err.Error() + " : " + v)
			return nil, err
		}

		_, err = stmt.ExecContext(ctx, id, idTags)
		if err != nil {
			r.logger.Error().Msg("key:JKpKzFVNSQdVJMA " + err.Error())
			return nil, err
		}
	}
	// cerate data end

	err = tx.Commit()
	if err != nil {
		r.logger.Error().Msg("Failed to commit transaction: " + err.Error())
		return nil, err
	}

	return id, nil
}

func (r *repository) UpdatePost(ctx context.Context, id uuid.UUID, req model.ReqPost) (err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Error().Msg("Failed to begin transaction: " + err.Error())
		return err
	}
	defer tx.Rollback()

	// check start
	check := 0
	err = r.db.QueryRowContext(ctx, checkDetailPosts, id).Scan(&check)
	if err != nil {
		r.logger.Error().Msg("key:kCjjVQ912p90 " + err.Error())
		return err
	}

	if check <= 0 {
		return errors.New("not_found")
	}
	// check end

	// get id post tags start
	rowsIdPostTags, err := r.db.QueryContext(ctx, getIdPostsTags, id)
	if err != nil {
		r.logger.Error().Msg("key:bblhSl92fc " + err.Error())
		return err
	}
	defer rowsIdPostTags.Close()

	var tags []*uuid.UUID

	for rowsIdPostTags.Next() {
		tag := &uuid.UUID{}
		if err := rowsIdPostTags.Scan(
			&tag,
		); err != nil {
			r.logger.Error().Msg("key:e8o8sBz5eXX " + err.Error())
			return err
		}

		tags = append(tags, tag)
	}
	// get id post tags end

	// delete post tags start
	stmtDelete, err := tx.PrepareContext(ctx, deletePostsTags)
	if err != nil {
		r.logger.Error().Msg("key:YQntizwnnP8LeCWBrVW " + err.Error())
		return err
	}
	defer stmtDelete.Close()

	for _, v := range tags {
		resultTags, err := stmtDelete.ExecContext(ctx, v)
		if err != nil {
			r.logger.Error().Msg("key:X7x2JfCiZvm " + err.Error())
			return err
		}

		rowsAffectedTags, err := resultTags.RowsAffected()
		if err != nil {
			r.logger.Error().Msg("key:JtROYnulj6bUP5EUPs " + err.Error())
			return err
		}

		if rowsAffectedTags <= 0 {
			return errors.New(" Failed to delete the data. No rows were affected")
		}
	}
	// delete post tags end

	// Update posts start
	resUpdatePosts, err := tx.ExecContext(ctx, updatePosts, req.Title, req.Content, id)
	if err != nil {
		r.logger.Error().Msg("key:Etkz1 " + err.Error())
		return err
	}

	rowsUpdatePosts, err := resUpdatePosts.RowsAffected()
	if err != nil {
		r.logger.Error().Msg("key:UWb5d " + err.Error())
		return err
	}

	if rowsUpdatePosts <= 0 {
		return errors.New(" Failed to update the data. No rows were affected")
	}
	// update posts end

	// add tags start
	stmtAddTags, err := tx.PrepareContext(ctx, creatPostsTag)
	if err != nil {
		r.logger.Error().Msg("key:VVKhehzCVeyzwJAw " + err.Error())
		return err
	}
	defer stmtAddTags.Close()

	for _, v := range req.Tags {

		idTags := uuid.New()

		err = tx.QueryRowContext(ctx, getIdTag, v).Scan(&idTags)
		if err != nil {
			r.logger.Error().Msg("key:9PuZX2s88n " + err.Error() + " : " + v)
			return err
		}

		_, err = stmtAddTags.ExecContext(ctx, id, idTags)
		if err != nil {
			r.logger.Error().Msg("key:TliZ4AZtic " + err.Error())
			return err
		}
	}
	// add tags end

	err = tx.Commit()
	if err != nil {
		r.logger.Error().Msg("Failed to commit transaction: " + err.Error())
		return err
	}

	return nil
}

func (r *repository) Register(ctx context.Context, req model.Credentials) (err error) {
	_, err = r.db.ExecContext(ctx, registerSql, req.Username, req.Password, req.Role)
	if err != nil {
		r.logger.Error().Msg("key:TliZ4AZtic " + err.Error())
		return
	}

	return nil
}

func (r *repository) Login(ctx context.Context, username string) (res model.Credentials, err error) {
	err = r.db.QueryRow(loginSql, username).Scan(
		&res.Username,
		&res.Password,
		&res.Role,
	)
	if err != nil {
		r.logger.Error().Msg("key:BYxJMWEjCLK4EP " + err.Error())
		return res, err
	}

	return res, nil
}
