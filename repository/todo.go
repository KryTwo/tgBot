package repository

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	e "main/lib"
	"main/structs"
)

// FindUser searches for a user in the database, returns true if found
func FindUser(db *sqlx.DB, uName string) (bool, error) {
	var in interface{}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.
		Select("username").
		From("users").
		Where("username = ?")
	stmt, _, err := builder.ToSql()
	if err != nil {
		return false, e.Wrap("Can't build squirrel, ", err)
	}

	err2 := db.Get(&in, stmt, uName)

	if err2 != nil && err2.Error() != "sql: no rows in result set" {
		return false, e.Wrap("Can't do query get to sqlx", err2)
	}

	if in == nil {
		return false, nil
	}

	return true, nil
}

// InsertUserToDB add user to database
func InsertUserToDB(db *sqlx.DB, uName string) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.
		Insert("users").
		Columns("username").
		Values("?")
	stmt, _, err := builder.ToSql()
	if err != nil {
		return e.Wrap("Can't build squirrel", err)
	}

	if _, err := db.Exec(stmt, uName); err != nil {
		return e.Wrap("Can't do Exec to db", err)
	}

	return nil
}

// GetStats get statistics to all users
func GetStats(db *sqlx.DB) ([]structs.Users, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Select("username, count").
		From("users")
	req, _, err := builder.ToSql()
	if err != nil {
		return nil, e.Wrap("Can't build request", err)
	}

	var in []structs.Users
	if err := db.Select(&in, req); err != nil {
		return nil, e.Wrap("Can't build request", err)
	}

	return in, nil
}

// GetStatsMe get statistics for one user
func GetStatsMe(db *sqlx.DB, username string) (structs.Users, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Select("username, count").
		From("users").
		Where("username = ?")
	req, _, err := builder.ToSql()
	if err != nil {
		return structs.Users{}, e.Wrap("Can't build builder", err)
	}

	var u structs.Users
	err2 := db.Get(&u, req, username)
	if err2 != nil {
		return structs.Users{}, e.Wrap("Can't get query to db", err2)
	}

	return u, nil
}

// IncrCounter was increment counter to random user (id = i)
func IncrCounter(db *sqlx.DB, i int) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Update("users").
		Set("count", sq.Expr("count+1")).
		Where("id = ?")
	req, _, err := builder.ToSql()

	if err != nil {
		return e.Wrap("Can't build builder", err)
	}

	if _, err := db.Exec(req, i); err != nil {
		return e.Wrap("Can't do query to db", err)
	}

	return nil
}
