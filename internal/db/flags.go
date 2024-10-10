package db

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/iagocanalejas/rstats/internal/utils/assert"
)

type FlagRow struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

func (r *Repository) GetFlagByID(flagID int64) (*FlagRow, error) {
	query, args, err := sq.
		Select("f.id as id", "f.name as name").
		From("flag f").
		Where(sq.Eq{"f.id": flagID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	assert.NoError(err, "building query", "query", query, "args", args)

	var flag FlagRow
	if err = r.db.Get(&flag, query, args...); err != nil {
		return nil, err
	}

	return &flag, nil
}
