package postgres

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/VikaPaz/matchmaker/internal/models"
	"github.com/sirupsen/logrus"
)

type PostgresRepository struct {
	conn *sql.DB
	log  *logrus.Logger
}

func NewRepo(conn *sql.DB, logger *logrus.Logger) *PostgresRepository {
	return &PostgresRepository{
		conn: conn,
		log:  logger,
	}
}

func (r *PostgresRepository) Create(player models.Player) error {
	builder := sq.Insert("users").Columns("id", "name", "skill", "latency", "added")
	builder = builder.Values(player.ID, player.Name, player.Skill, player.Latency, player.Added)
	builder = builder.PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		r.log.Error("Postgres: failed to create player")
		return err
	}

	r.log.Debugf("Executing query: %v", query)
	_, err = r.conn.Exec(query, args...)
	if err != nil {
		r.log.Error("Postgres: failed to create player")
		return err
	}
	return nil
}

func (r *PostgresRepository) Delete(id uint) error {
	builder := sq.Delete("users").Where(sq.Eq{"id": id})
	builder = builder.PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		r.log.Error("Postgres: failed to delete user")
		return err
	}

	r.log.Debugf("Executing query: %v", query)
	_, err = r.conn.Exec(query, args...)
	if err != nil {
		r.log.Error("Postgres: failed to delete user")
		return err
	}
	return nil
}
