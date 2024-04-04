package db

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
	"github.com/stas9132/GophKeeper/internal/config"
	"github.com/stas9132/GophKeeper/internal/logger"
)

// DBT struct
type DB struct {
	logger.Logger
	*sql.DB
}

// NewDB constructor
func NewDB(l logger.Logger) (*DB, error) {
	db, err := sql.Open("postgres", config.DatabaseDSN)
	if err != nil {
		l.Error("Error while open db: " + err.Error())
		return nil, err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		l.Error("Error while get driver: " + err.Error())
		return nil, err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/server/db/migration",
		config.DatabaseDSN, driver)
	if err != nil {
		l.Error("Error while create migrate: " + err.Error())
		return nil, err
	} else {
		if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			l.Error("Error while migrate up: " + err.Error())
			return nil, err
		}
	}

	return &DB{
		Logger: l,
		DB:     db,
	}, nil
}

func (d *DB) Register(ctx context.Context, user, password string) (bool, error) {

	ht := md5.Sum([]byte(password))
	hash := hex.EncodeToString(ht[:])

	_, err := d.ExecContext(ctx, "INSERT INTO users(user_id, hash) VALUES ($1,$2)", user, hash)
	if err != nil {
		d.Error("Unable insert record: " + err.Error())
		return false, err
	}
	return true, nil
}

func (d *DB) Login(ctx context.Context, user, password string) (bool, error) {
	row := d.QueryRowContext(ctx, "SELECT hash FROM users WHERE user_id = $1", user)
	var hashDb string
	if err := row.Scan(&hashDb); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			d.Warn("Userid not found: " + user)
			return false, nil
		}
		d.Error("Unable to get user record: " + err.Error())
		return false, err
	}

	ht := md5.Sum([]byte(password))
	hash := hex.EncodeToString(ht[:])

	if hash != hashDb {
		d.Warn("Unauthenticated request: " + user)
		return false, nil
	}

	return true, nil
}

func (d *DB) Logout(ctx context.Context) (bool, error) {
	return true, nil
}
