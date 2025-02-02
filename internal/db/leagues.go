package db

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/iagocanalejas/rstats/internal/utils/assert"
)

type LeagueRow struct {
	ID       int64   `db:"id"`
	Name     string  `db:"name"`
	Symbol   string  `db:"symbol"`
	Gender   *string `db:"gender"`
	Category *string `db:"category"`
}

func (r *Repository) GetLeagueByID(leagueID int64) (*LeagueRow, error) {
	query, args, err := sq.
		Select("l.id as id", "l.name as name", "l.gender as gender", "l.category as category", "l.symbol as symbol").
		From("league l").
		Where(sq.Eq{"l.id": leagueID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	assert.NoError(err, "building query=%s args=%s", query, args)

	var league LeagueRow
	if err = r.db.Get(&league, query, args...); err != nil {
		return nil, err
	}

	return &league, nil
}
