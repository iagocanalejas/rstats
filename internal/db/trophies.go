package db

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/iagocanalejas/rstats/internal/utils/assert"
)

type Trophy struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

func (r *Repository) GetTrophyByID(trophyID int64) (*Trophy, error) {
	query, args, err := sq.
		Select("t.id as id", "t.name as name").
		From("trophy t").
		Where(sq.Eq{"t.id": trophyID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	assert.NoError(err, "building query", "query", query, "args", args)

	var trophy Trophy
	if err = r.db.Get(&trophy, query, args...); err != nil {
		return nil, err
	}

	return &trophy, nil
}
