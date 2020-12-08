package view

import (
	"errors"
	"fmt"

	"github.com/crypto-com/chain-indexing/appinterface/rdb"
)

// A "Total" compatible table should have the follow table schema
// | Column   | Type    | Constraint  |
// | -------- | ------- | ----------- |
// | identity | VARCHAR | PRIMARY KEY |
// | total    | BIGINT  | NOT NULL    |

// A generic record total tracker view
type Total struct {
	rdbHandle *rdb.Handle

	tableName string
}

func NewTotal(rdbHandle *rdb.Handle, tableName string) *Total {
	return &Total{
		rdbHandle,

		tableName,
	}
}

func (view *Total) Set(identity string, total int64) error {
	// Postgres UPSERT statement
	sql, sqlArgs, err := view.rdbHandle.StmtBuilder.
		Insert(view.tableName).
		Columns("identity", "total").
		Values(identity, total).
		Suffix("ON CONFLICT (identity) DO UPDATE SET total = EXCLUDED.total").
		ToSql()
	if err != nil {
		return fmt.Errorf("error building total insertion sql: %v: %w", err, rdb.ErrBuildSQLStmt)
	}

	_, err = view.rdbHandle.Exec(sql, sqlArgs...)

	if err != nil {
		return fmt.Errorf("error inserting total: %v: %w", err, rdb.ErrWrite)
	}

	return nil
}

func (view *Total) Increment(identity string, total int64) error {
	// Postgres UPSERT statement
	sql, sqlArgs, err := view.rdbHandle.StmtBuilder.
		Insert(view.tableName+" AS totals").
		Columns("identity", "total").
		Values(identity, total).
		Suffix("ON CONFLICT (identity) DO UPDATE SET total = totals.total + EXCLUDED.total").
		ToSql()
	if err != nil {
		return fmt.Errorf("error building total insertion sql: %v: %w", err, rdb.ErrBuildSQLStmt)
	}

	_, err = view.rdbHandle.Exec(sql, sqlArgs...)

	if err != nil {
		return fmt.Errorf("error inserting total: %v: %w", err, rdb.ErrWrite)
	}

	return nil
}

func (view *Total) FindBy(identity string) (int64, error) {
	sql, sqlArgs, err := view.rdbHandle.StmtBuilder.Select(
		"total",
	).From(
		view.tableName,
	).Where(
		"identity = ?", identity,
	).ToSql()
	if err != nil {
		return int64(0), fmt.Errorf("error preparing total selection SQL: %v", err)
	}

	var total int64
	if err := view.rdbHandle.QueryRow(sql, sqlArgs...).Scan(&total); err != nil {
		if errors.Is(err, rdb.ErrNoRows) {
			return int64(0), nil
		}
		return int64(0), fmt.Errorf("error getting total: %v", err)
	}

	return total, nil
}
