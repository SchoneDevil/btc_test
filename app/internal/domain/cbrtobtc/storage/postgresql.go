package storage

import (
	"context"
	"log"

	"app/internal/domain/cbrtobtc/model"

	sq "github.com/Masterminds/squirrel"
)

type CbrToBtcStorage struct {
	qB     sq.StatementBuilderType
	client PostgreSQLClient
}

func NewCbrToBtcStorage(client PostgreSQLClient) CbrToBtcStorage {
	return CbrToBtcStorage{
		qB:     sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		client: client,
	}
}

func (s CbrToBtcStorage) Insert(ctx context.Context, ctb model.CbrToBtc) error {
	query := s.qB.Insert("btctocbr").
		Columns("name, value, created_at").
		Values(ctb.Name, ctb.Value, ctb.CreatedAt)

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

func (s CbrToBtcStorage) InsertOrUpdate(ctx context.Context, ctb model.CbrToBtc) error {
	query := s.qB.Insert("btctocbr").
		Columns("name, value, created_at").
		Values(ctb.Name, ctb.Value, ctb.CreatedAt).
		Suffix("ON CONFLICT (name) DO UPDATE SET value = ?", ctb.Value)

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

func (s CbrToBtcStorage) Latest(ctx context.Context) ([]model.CbrToBtc, error) {
	query := s.qB.Select("id, name, value, created_at").
		From("btctocbr").
		OrderBy("name desc")
	sql, args, err := query.ToSql()

	rows, err := s.client.Query(ctx, sql, args...)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	list := make([]model.CbrToBtc, 0)
	for rows.Next() {
		c := model.CbrToBtc{}
		if err = rows.Scan(
			&c.ID, &c.Name, &c.Value, &c.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, c)
	}

	return list, nil
}
