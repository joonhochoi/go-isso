package db

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"gopkg.in/guregu/null.v3"
)

func initPreference(db *database) error {
	if v, err := db.GetPreference("session-key"); err != nil || v.Valid {
		if _, ok := err.(*AcceptableError); !ok {
			return fmt.Errorf("db.initPreference: %v", err)
		}
	}

	b := make([]byte, 24)
	_, err := rand.Read(b)
	if err != nil {
		return err
	}
	sessionKey := hex.EncodeToString(b)

	err = db.SetPreference("session-key", sessionKey)
	if err != nil {
		return err
	}
	return nil
}

// GetPreference get the corresponding value of the key.
func (db *database) GetPreference(key string) (null.String, error) {
	var value null.String
	err := db.QueryRow("SELECT value FROM preferences WHERE key=?", key).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return null.NewString("", false), newError(fmt.Sprintf("db.GetPreference: %v", err))
		}
		return null.NewString("", false), fmt.Errorf("db.GetPreference: %v", err)
	}
	return value, nil
}

// SetPreference set the corresponding value of the key.
func (db *database) SetPreference(key, value string) error {
	stmt, err := db.Prepare("INSERT INTO preferences (key, value) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("db.SetPreference: %v", err)
	}
	_, err = stmt.Exec(key, value)
	if err != nil {
		return fmt.Errorf("db.SetPreference: %v", err)
	}
	return nil
}
