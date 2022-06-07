package storage

import (
	"context"
	"fmt"
	"log"

	"app/internal/domain/courses/model"
	"app/internal/handlers"

	sq "github.com/Masterminds/squirrel"
)

type CourseStorage struct {
	qB     sq.StatementBuilderType
	client PostgreSQLClient
}

func NewCourseStorage(client PostgreSQLClient) CourseStorage {
	return CourseStorage{
		qB:     sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		client: client,
	}
}

func (s CourseStorage) Insert(ctx context.Context, course model.Course) error {
	query := s.qB.Insert("courses").
		Columns("symbol, buy, rub, created_at").
		Values(course.Symbol, course.Buy, course.Rub, course.CreatedAt)

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

func (s CourseStorage) SelectLast(ctx context.Context) (model.Course, error) {
	query := s.qB.
		Select("id, symbol, buy, rub, created_at").
		From("courses").
		OrderBy("created_at desc").
		Limit(1)
	sql, args, err := query.ToSql()
	if err != nil {
		return model.Course{}, err
	}
	fmt.Println(sql)
	row := s.client.QueryRow(ctx, sql, args...)

	course := model.Course{}

	err = row.Scan(&course.ID, &course.Symbol, &course.Buy, &course.Rub, &course.CreatedAt)
	if err != nil {
		return model.Course{}, err
	}

	return course, err
}

func (s CourseStorage) Select(ctx context.Context, filter handlers.PostFilter) ([]model.Course, error) {
	query := s.qB.
		Select("id, symbol, buy, rub, created_at").
		From("courses").
		OrderBy("created_at desc")

	if filter.Limit != 0 {
		query.Limit(uint64(filter.Limit))
	}
	if filter.Offset != 0 {
		query.Offset(uint64(filter.Offset))
	}
	if filter.DateStart != "" {
		query.Where("created_at <= ?", filter.DateStart)
	}
	if filter.DateFinish != "" {
		query.Where("created_at >= ?", filter.DateStart)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := s.client.Query(ctx, sql, args...)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	list := make([]model.Course, 0)

	for rows.Next() {
		c := model.Course{}
		if err = rows.Scan(
			&c.ID, &c.Symbol, &c.Buy, &c.Rub, &c.CreatedAt,
		); err != nil {
			return nil, err
		}
	}

	return list, nil
}
