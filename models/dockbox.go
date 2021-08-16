package models

import (
	"database/sql"
	"time"

	db "github.com/dockboxhq/server/services"
	"github.com/pkg/errors"
)

type Dockbox struct {
	ID           string         `db:"id"`
	CustomId     sql.NullString `db:"custom_id"`
	ContainerId  sql.NullString `db:"container_id"`
	Source       sql.NullString `db:"source"`
	CreatedAt    time.Time      `db:"created_at"`
	LastModified time.Time      `db:"last_modified"`
}

func GetDockboxById(id string) (Dockbox, error) {
	dockbox := Dockbox{}
	err := db.DbConn.Get(&dockbox, "select * from dockbox where id = ?", id)
	if err != nil {
		return dockbox, err
	}
	return dockbox, nil
}

func CreateDockbox(ID string, source string, container_id sql.NullString, custom_id sql.NullString) (Dockbox, error) {
	dockbox := Dockbox{}
	tx := db.DbConn.MustBegin()
	err := tx.Get(&dockbox, "INSERT INTO dockbox (id, container_id, source, custom_id) VALUES (?, ?, ?, ?) returning *;", ID, container_id, source, custom_id)
	if err != nil {
		tx.Rollback()
		return dockbox, errors.Wrap(err, "create dockbox error")
	}
	tx.Commit()
	return dockbox, nil
}

func UpdateDockbox(dockboxID string, container_id sql.NullString) (Dockbox, error) {
	dockbox := Dockbox{}
	tx := db.DbConn.MustBegin()
	err := tx.Get(&dockbox, "UPDATE dockbox SET container_id = ? where id = ? returning *;", container_id, dockboxID)
	if err != nil {
		tx.Rollback()
		return dockbox, errors.Wrap(err, "create dockbox error")
	}
	tx.Commit()
	return dockbox, nil
}
