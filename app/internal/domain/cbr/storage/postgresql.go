package storage

import (
	"context"
	"log"

	"app/internal/domain/cbr/model"

	sq "github.com/Masterminds/squirrel"
)

type CbrStorage struct {
	qB     sq.StatementBuilderType
	client PostgreSQLClient
}

func NewCbrStorage(client PostgreSQLClient) CbrStorage {
	return CbrStorage{
		qB:     sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		client: client,
	}
}

func (s CbrStorage) Insert(ctx context.Context, cbr model.Cbr) error {
	query := s.qB.Insert("cbr").
		Columns("charcode, name, value, created_at").
		Values(cbr.CharCode, cbr.Name, cbr.Value, cbr.CreatedAt)

	sql, args, err := query.ToSql()

	if err != nil {
		return err
	}

	_, err = s.client.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}
	return nil
}

func (s CbrStorage) SelectLast(ctx context.Context) ([]model.Cbr, error) {
	query := s.qB.Select("id, charcode, name, value, created_at").
		From("(SELECT DISTINCT ON (charcode) * FROM cbr ORDER BY charcode) t").
		OrderBy("charcode ASC")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, _ := s.client.Query(ctx, sql, args...)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	list := make([]model.Cbr, 0)

	for rows.Next() {
		c := model.Cbr{}
		if err = rows.Scan(
			&c.ID, &c.CharCode, &c.Name, &c.Value, &c.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, c)
	}

	return list, nil
}

func (s CbrStorage) SelectUsd(ctx context.Context) (model.Cbr, error) {
	query := s.qB.Select("id, charcode, name, value, created_at").
		From("cbr").
		Where("charcode = 'USD'").OrderBy("created_at DESC")

	sql, args, err := query.ToSql()
	if err != nil {
		return model.Cbr{}, err
	}

	row := s.client.QueryRow(ctx, sql, args...)

	cbr := model.Cbr{}

	err = row.Scan(&cbr.ID, &cbr.CharCode, &cbr.Name, &cbr.Value, &cbr.CreatedAt)
	if err != nil {
		return model.Cbr{}, err
	}

	return cbr, err
}
