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
	UserId    int64
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

func (s *SessionManager) CreateSession(username string, userId int64) (*Session, error) {
	err := s.DeleteUserSessions(userId)
	if err != nil {
		return nil, err
	}
	query := `INSERT INTO sessions (id, username, userId, expiresAt) VALUES (?, ?, ?, ?)`
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	expTime := time.Now().Add(SESSION_EXP_TIME * time.Second)
	s.db.Exec(query, id.String(), username, userId, expTime)
	session, err := s.GetSession(id.String())
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (s *SessionManager) GetSession(id string) (*Session, error) {
	query := `SELECT * FROM sessions WHERE id = ?`
	var session Session

	row := s.db.QueryRow(query, id)
	err := row.Scan(&session.ID, &session.Username, &session.UserId, &session.ExpiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if time.Now().After(session.ExpiresAt) {
		s.DeleteSession(session.ID)
		return nil, ErrExpiredSession
	}
	return &session, nil
}

func (s *SessionManager) DeleteSession(id string) error {
	query := `DELETE FROM sessions WHERE userId = (SELECT userId FROM sessions WHERE id = ?)`
	_, err := s.db.Exec(query, id)
	return err
}

func (s *SessionManager) DeleteUserSessions(userId int64) error {
	query := "DELETE FROM sessions WHERE userId = ?"
	_, err := s.db.Exec(query, userId)
	if err != nil {
		return err
	}
	return nil
}
