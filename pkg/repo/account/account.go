package account

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgconn"
	log "github.com/sirupsen/logrus"
	"menti/pkg/db"
	"menti/pkg/domain/account"
	interfaces "menti/pkg/repo/account/interface"
)

type accountDataBase struct {
	db *sql.DB
}

func NewAccountDataBase(db *sql.DB) interfaces.AccountRepository {
	return &accountDataBase{
		db: db,
	}
}

func (r *accountDataBase) Migrate(ctx context.Context) error {
	accQuery := `
    CREATE TABLE IF NOT EXISTS account(
       	id SERIAL PRIMARY KEY,
		name text not NULL,
		email text not NULL,
		password text not NULL,
		noteIds integer[]
    );
    `
	_, err := r.db.ExecContext(ctx, accQuery)
	if err != nil {
		message := db.ErrMigrate.Error() + " account"
		log.Fatalf("%q: %s\n", message, err.Error())
		return db.ErrMigrate
	}
	log.Info("account table created")
	return err
}
func (r *accountDataBase) Registration(ctx context.Context) (*account.Account, error) {
	// TODO: хардкод регистрация
	var newAccount = account.Account{
		ID:      1,
		Name:    "test",
		Email:   "test@test.com",
		Pass:    "test",
		NoteIDS: []int64{},
	}
	var existingCount int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM account WHERE ID = $1", newAccount.ID).Scan(&existingCount)
	if err != nil {
		return nil, err
	}

	if existingCount > 0 {
		return nil, db.ErrDuplicate
	}

	var id int64

	err = r.db.QueryRowContext(ctx,
		"INSERT INTO account(name, email, password, noteIds) VALUES($1, $2, $3, $4) RETURNING id",
		newAccount.Name, newAccount.Email, newAccount.Pass, newAccount.NoteIDS).Scan(&id)
	// Check if a user with the same email already exists
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				return nil, db.ErrDuplicate
			}
		}
		return nil, err
	}

	// Add the new account
	requestAccount := &account.Account{
		ID:      id,
		Name:    newAccount.Name,
		Email:   newAccount.Email,
		Pass:    newAccount.Pass,
		NoteIDS: newAccount.NoteIDS,
	}

	return requestAccount, nil
}

func (r *accountDataBase) Login(ctx context.Context, acc account.Login) (int64, error) {
	var id int64
	row := r.db.QueryRowContext(ctx, "SELECT id FROM account WHERE email = $1 and password = $2", acc.Email, acc.Password)

	if err := row.Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error(err)
			return 0, db.ErrNotExist
		}
		log.Error(err)
		return 0, err
	}
	log.Info("login success")
	return id, nil
}

func (r *accountDataBase) CheckAccount(ctx context.Context, id int64) error {
	var accountID int64
	err := r.db.QueryRowContext(ctx, "SELECT id FROM account WHERE id = $1", id).Scan(&accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error(err)
			return db.ErrNotExist
		}
		log.Error(err)
		return err
	}
	log.Info("check account success")
	return nil
}
