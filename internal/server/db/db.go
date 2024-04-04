package db

import (
	"context"
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/stas9132/GophKeeper/internal/config"
	"github.com/stas9132/GophKeeper/internal/logger"
)

// DBT struct
type DB struct {
	logger logger.Logger
	db     *sql.DB
}

// NewDB constructor
func NewDB(l logger.Logger) (*DB, error) {
	db, err := sql.Open("pgx", config.DatabaseDSN)
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
		"pgx://"+config.DatabaseDSN, driver)
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
		logger: l,
		db:     db,
	}, nil
}

func (d *DB) Register(ctx context.Context, user, password string) (bool, error) {
	//if err := s.MakeBucket(ctx, user, minio.MakeBucketOptions{Region: config.S3Location, ObjectLocking: true}); err != nil {
	//	return false, err
	//}
	//hash := md5.Sum([]byte(password))
	//info, err := s.PutObject(ctx, authBucketName, user, bytes.NewReader(hash[:]), int64(len(hash)), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	//if err != nil {
	//	s.Error(err.Error())
	//	return false, err
	//}
	//s.Info(info.Key + " stored in auth-storage")
	return true, nil
}

func (d *DB) Login(ctx context.Context, user, password string) (bool, error) {
	//data, err := s.GetObject(ctx, authBucketName, user, minio.GetObjectOptions{})
	//if err != nil {
	//	s.Error(err.Error())
	//	return false, err
	//}
	//hash := md5.Sum([]byte(password))
	//buf, err := io.ReadAll(data)
	//if err != nil {
	//	s.Error(err.Error())
	//	return false, err
	//}
	//if c := bytes.Compare(hash[:], buf); c != 0 {
	//	s.Warn("Invalid credentials for user: " + user)
	//	return false, err
	//}
	//s.Info("Login complete")
	return true, nil
}

func (d *DB) Logout(ctx context.Context) (bool, error) {
	return true, nil
}
