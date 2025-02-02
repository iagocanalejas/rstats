package db

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/iagocanalejas/rstats/internal/utils/assert"
)

type EntityRow struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

func (r *Repository) GetClubByID(clubID int64) (*EntityRow, error) {
	query, args, err := sq.
		Select("e.id as id", "e.name as name").
		From("entity e").
		Where(sq.Eq{"e.id": clubID, "e.type": "CLUB"}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	assert.NoError(err, "building query=%s args=%s", query, args)

	var club EntityRow
	if err = r.db.Get(&club, query, args...); err != nil {
		return nil, err
	}

	return &club, nil
}
