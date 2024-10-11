package config

import (
	"database/sql"
	"errors"
	"time"

	"github.com/gofrs/uuid/v5"
)

var ErrExpiredSession = errors.New("session is expired")

type Session struct {
	ID        string
	Username  string
	UserId    string
	ExpiresAt time.Time
}

type SessionManager struct {
	db *sql.DB
}

func NewSessionManager() {
	SESSION = &SessionManager{
		db: DB,
	}
}

func (s *SessionManager) CreateSession(username string, userId string) (*Session, error) {
	query := `INSERT INTO sessions (id, username, userId, expires_at) VALUES (?, ?, ?, ?)`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	expTime := time.Now().Add(SESSION_EXP_TIME * time.Second)
	stmt.Exec(id.String(), username, userId, expTime)
	session, err := s.GetSession(id.String())
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (s *SessionManager) GetSession(id string) (*Session, error) {
	query := `SELECT * FROM sessions WHERE id = ?`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return nil, NewInternalError(err)
	}
	defer stmt.Close()
	var session Session

	row := stmt.QueryRow(id)
	err = row.Scan(&session.ID, &session.Username, &session.UserId, &session.ExpiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewError(err)
		}
		return nil, NewInternalError(err)
	}
	if time.Now().After(session.ExpiresAt) {
		s.DeleteSession(session.ID)
		return nil, ErrExpiredSession
	}
	return &session, nil
}

func (s *SessionManager) DeleteSession(id string) error {
	query := `DELETE sessions WHERE id = ?`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return NewInternalError(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	return NewInternalError(err)
}
