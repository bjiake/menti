package note

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgconn"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"menti/pkg/db"
	"menti/pkg/domain/account"
	"menti/pkg/domain/note"
	interfaces "menti/pkg/repo/note/interface"
)

type noteDataBase struct {
	db *sql.DB
}

func NewNoteDataBase(db *sql.DB) interfaces.NoteRepository {
	return &noteDataBase{
		db: db,
	}
}

func (r *noteDataBase) Migrate(ctx context.Context) error {
	accQuery := `
    CREATE TABLE IF NOT EXISTS note(
       	id SERIAL PRIMARY KEY,
		name text not NULL,
		content text not NULL
    );
    `
	_, err := r.db.ExecContext(ctx, accQuery)
	if err != nil {
		message := db.ErrMigrate.Error() + " note"
		log.Fatalf("%q: %s\n", message, err.Error())
		return db.ErrMigrate
	}
	log.Info("note table created")
	return err
}

func (r *noteDataBase) GetAll(ctx context.Context, userID int64) ([]note.Note, error) {
	rows := r.db.QueryRowContext(ctx, "SELECT * FROM account where id = $1", userID)

	var acc account.Account
	var noteIDs []int64
	if err := rows.Scan(&acc.ID, &acc.Name, &acc.Email, &acc.Pass, pq.Array(&noteIDs)); err != nil {
		log.Error(err)
		return nil, err
	}
	acc.NoteIDS = noteIDs

	rowsNotes, err := r.db.QueryContext(ctx, "SELECT * FROM note where id = ANY($1)", acc.NoteIDS)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer rowsNotes.Close()
	var notes []note.Note
	for rowsNotes.Next() {
		var note note.Note
		if err = rowsNotes.Scan(&note.ID, &note.Name, &note.Content); err != nil {
			log.Error(err)
			return nil, err
		}
		notes = append(notes, note)
	}
	if err = rowsNotes.Err(); err != nil {
		log.Error(err)
		return nil, err
	}
	log.Info("notes found\tlen:%v\tnotes:%v", len(notes), notes)
	return notes, nil
}

func (r *noteDataBase) Post(ctx context.Context, userId int64, newPeople note.Note) (int64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		log.Error(err)
		return 0, err
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(ctx, "SELECT * FROM account where id = $1", userId)
	var acc account.Account
	var noteIDs []int64
	if err := row.Scan(&acc.ID, &acc.Name, &acc.Email, &acc.Pass, pq.Array(&noteIDs)); err != nil {
		log.Error(err)
		return 0, err
	}
	acc.NoteIDS = noteIDs

	var id int64

	err = tx.QueryRowContext(ctx, "INSERT INTO note(name, content) values($1, $2) RETURNING id",
		newPeople.Name, newPeople.Content).Scan(&id)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				log.Info(err)
				return 0, db.ErrDuplicate
			}
		}
		log.Info(err)
		return 0, err
	}

	noteIDs = append(noteIDs, id)
	_, err = tx.ExecContext(ctx, "UPDATE account SET noteIds = $1 WHERE id = $2", pq.Array(noteIDs), userId)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		log.Error(err)
		return 0, err
	}

	log.Info("note created", id)
	return id, nil
}
