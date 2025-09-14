package db

import (
	"context"
	"memecoin_trading_bot/app/app_errors"
	"os"
	"path"
	"runtime"
)

func (db *DB) Migrate(migrations_dir string) error {
	_, fileName, _, ok := runtime.Caller(0)
	if !ok {
		return app_errors.ErrSourceFilePath
	}

	baseDir := path.Join(fileName, "..", "..", "..", migrations_dir)
	files, err := os.ReadDir(baseDir)
	if err != nil {
		return err
	}

	tx, err := db.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	for _, file := range files {
		migration, err := os.ReadFile(path.Join(baseDir, file.Name()))
		if err != nil {
			return err
		}
		_, err = tx.Exec(string(migration))
		if err != nil {
			err = tx.Rollback()
			if err != nil {
				return err
			}
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
