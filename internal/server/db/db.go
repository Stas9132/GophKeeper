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
	"github.com/google/uuid"
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

type Meta struct {
	ObjId   string
	ObjType int
}

func (d *DB) PutMeta(ctx context.Context, userId, objName string, objType int) (*Meta, error) {
	tx, err := d.BeginTx(ctx, nil)
	if err != nil {
		d.Error("Begin transaction error: " + err.Error())
		return nil, err
	}
	defer tx.Rollback()
	rows, err := tx.QueryContext(ctx, "SELECT obj_id FROM meta WHERE user_id = $1 AND obj_name = $2", userId, objName)
	if err != nil {
		d.Error("Query error: " + err.Error())
		return nil, err
	}
	var u string
	if rows.Next() {
		if err = rows.Scan(&u); err != nil {
			d.Error("Scan error: " + err.Error())
			return nil, err
		}
	}
	if rows.Next() {
		d.Error("meta corrupted: " + userId + ", " + objName)
		return nil, errors.New("meta corrupted: " + userId + ", " + objName)
	}

	if len(u) == 0 {
		u = uuid.NewString()
		if _, err = tx.ExecContext(ctx, "INSERT INTO meta(user_id, obj_name, obj_id, obj_type) VALUES ($1,$2,$3,$4)", userId, objName, u, objType); err != nil {
			d.Error("SQL exec error: " + err.Error())
			return nil, err
		}

		if err = tx.Commit(); err != nil {
			d.Error("Commit error: " + err.Error())
			return nil, err
		}
	}

	return &Meta{ObjId: u}, nil
}

func (d *DB) GetMeta(ctx context.Context, userId, objName string) (*Meta, error) {
	var objId string
	var objType int
	row := d.QueryRowContext(ctx, "SELECT obj_id, obj_type FROM meta WHERE user_id = $1 AND obj_name = $2", userId, objName)
	if err := row.Scan(&objId, &objType); err != nil {
		d.Error("Query error: " + err.Error())
		return nil, err
	}
	return &Meta{
		ObjId:   objId,
		ObjType: objType,
	}, nil
}
