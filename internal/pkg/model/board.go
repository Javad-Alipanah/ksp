package model

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	database "ksp/internal/pkg/db/mysql"
)

type Board struct {
	ID string
	Size int
	Start string
	Target string
	Path string
}

func (b *Board) Save() int64 {
	statement, err := database.Db.Prepare("INSERT INTO Boards(size, start, target, path) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	res, err := statement.Exec(b.Size, b.Start, b.Target, b.Path)
	if err != nil {
		log.Fatal(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("row inserted")

	return id
}

func (b *Board) Update(tx *sql.Tx) int64 {
	statement, err := tx.Prepare("UPDATE Boards SET size = ?, start = ?, target = ?, path = ? WHERE id = ?")
	if err != nil {
		_ = tx.Rollback()
		log.Fatal(err)
	}

	res, err := statement.Exec(b.Size, b.Start, b.Target, b.Path, b.ID)
	if err != nil {
		_ = tx.Rollback()
		log.Fatal(err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		log.Fatal(err)
	}

	return n
}

func Get(id uint64) *Board {
	statement, err := database.Db.Prepare("SELECT id, size, start, target, path FROM Boards WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}

	b := new(Board)
	row := statement.QueryRow(id)
	err = row.Scan(&b.ID, &b.Size, &b.Start, &b.Target, &b.Path)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatal(err)
		}
		return nil
	}
	return b
}

func GetAll() *[]Board {
	statement, err := database.Db.Prepare("SELECT id, size, start, target, path FROM Boards")
	if err != nil {
		log.Fatal(err)
	}

	defer statement.Close()
	rows, err := statement.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var boards []Board
	for rows.Next() {
		var b Board
		err := rows.Scan(&b.ID, &b.Size, &b.Start, &b.Target, &b.Path)
		if err != nil{
			log.Fatal(err)
		}
		boards = append(boards, b)
	}
	return &boards
}

func SelectForUpdate(id uint64, tx *sql.Tx) *Board {
	statement, err := tx.Prepare("SELECT size, start, target, path FROM Boards WHERE id = ? FOR UPDATE ")
	if err != nil {
		_ = tx.Rollback()
		log.Fatal(err)
	}

	b := new(Board)
	row := statement.QueryRow(id)
	err = row.Scan(&b.Size, &b.Start, &b.Target, &b.Path)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatal(err)
		}
		return nil
	}
	return b
}

func Delete(id string) int64 {
	statement, err := database.Db.Prepare("DELETE FROM Boards WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}

	res, err := statement.Exec(id)
	if err != nil {
		log.Fatal(err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	return n
}