package helper

import (
	"context"
	"database/sql"
	"ecommerce-api/pkg/constant"
	"github.com/Masterminds/squirrel"
	"reflect"
)

type Callback func(ctx context.Context, tx *sql.Tx) error

func Transaction(ctx context.Context, db *sql.DB, cb Callback) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if err = cb(ctx, tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	_ = tx.Commit()
	return nil
}

func PageQuery(ctx context.Context, builder squirrel.SelectBuilder, pagination constant.Pagination, data interface{}) (err error) {
	builder = builder.Limit(uint64(pagination.PageSize)).Offset(getOffset(pagination))

	rows, err := builder.QueryContext(ctx)
	if err != nil {
		return
	}

	// data: pointer to nil slice
	v := reflect.ValueOf(data)
	v = reflect.Indirect(v)

	for rows.Next() {
		newVal := reflect.New(v.Type().Elem()).Elem()
		fields := make([]interface{}, newVal.NumField())
		for i := 0; i < newVal.NumField(); i++ {
			fields[i] = newVal.Field(i).Addr().Interface()
		}
		if err = rows.Scan(fields...); err != nil {
			return
		}
		v.Set(reflect.Append(v, newVal))
	}
	if err = rows.Close(); err != nil {
		return
	}
	if err = rows.Err(); err != nil {
		return
	}

	if v.IsNil() {
		s := reflect.MakeSlice(v.Type(), 0, 0)
		v.Set(s)
	}

	return
}

func getOffset(params constant.Pagination) uint64 {
	if params.Page > 0 {
		return uint64((params.Page - 1) * params.PageSize)
	}
	return 0
}

func TotalQuery(ctx context.Context, builder squirrel.SelectBuilder, total *int64) error {
	row := builder.QueryRowContext(ctx)
	if err := row.Scan(total); err != nil {
		return err
	}
	return nil
}

func StructQuery(ctx context.Context, builder squirrel.SelectBuilder, data interface{}) error {
	rows, err := builder.QueryContext(ctx)
	if err != nil {
		return err
	}

	// data: pointer to nil slice
	v := reflect.ValueOf(data)
	v = reflect.Indirect(v)

	for rows.Next() {
		newVal := reflect.New(v.Type().Elem()).Elem()
		fields := make([]interface{}, newVal.NumField())
		for i := 0; i < newVal.NumField(); i++ {
			fields[i] = newVal.Field(i).Addr().Interface()
		}
		if err = rows.Scan(fields...); err != nil {
			return err
		}
		v.Set(reflect.Append(v, newVal))
	}
	if err = rows.Close(); err != nil {
		return err
	}
	if err = rows.Err(); err != nil {
		return err
	}

	if v.IsNil() {
		s := reflect.MakeSlice(v.Type(), 0, 0)
		v.Set(s)
	}

	return nil
}
