package repository

import (
	"context"
	"database/sql"

	"github.com/dibikhairurrazi/audio-storage/db"
	"github.com/dibikhairurrazi/audio-storage/module/audio-storage/model"
)

type PostgreRepository struct {
	DB db.DB
}

func NewPostgreSQL(db db.DB) *PostgreRepository {
	return &PostgreRepository{
		DB: db,
	}
}

func (pg *PostgreRepository) FindUser(ctx context.Context, userID int) (model.User, error) {
	var user model.User
	err := pg.DB.ReplicaConn.QueryRow("SELECT id, name, email FROM users where id = $1", userID).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (pg *PostgreRepository) FindPhraseMetadata(ctx context.Context, userID, audioID int) (model.Phrase, error) {
	var phrase model.Phrase
	err := pg.DB.ReplicaConn.QueryRow("SELECT id, user_id, original_filename, filepath FROM phrases where id = $1 and user_id = $2", audioID, userID).
		Scan(&phrase.ID, &phrase.UserID, &phrase.OriginalFileName, &phrase.FilePath)
	if err != nil {
		return model.Phrase{}, err
	}

	return phrase, nil
}

func (pg *PostgreRepository) SavePhraseMetadata(ctx context.Context, phrase model.Phrase) error {
	_, err := pg.DB.MasterConn.Exec(`INSERT INTO phrases (id, user_id, original_filename, filepath) 
	VALUES ($1, $2, $3, $4) ON CONFLICT(id) DO UPDATE SET original_filename = $3, filepath = $4, updated_at = NOW()`, phrase.ID, phrase.UserID, phrase.OriginalFileName, phrase.FilePath)
	return err
}

func (pg *PostgreRepository) CreateTx(ctx context.Context) (*sql.Tx, error) {
	return pg.DB.MasterConn.BeginTx(ctx, &sql.TxOptions{})
}
