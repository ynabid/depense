package db

import (
	"database/sql"
	"errors"
	"github.com/ynabid/math/rand"
	"strconv"
	"time"
)

type User struct {
	Id         int64
	Email      string
	Password   string
	First_name string
	Last_name  string
	Birthdate  string
}

type Session struct {
	Id         string
	UserId     int64
	ExpiryDate int64
}

func Login(email, password string) (*Session, error) {
	var session *Session = &Session{}
	var userId int64
	rows, err := db.Query(
		"SELECT id FROM user WHERE email = ? and password = ?",
		email,
		password,
	)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		rows.Scan(&userId)
	} else {
		return nil, errors.New("User or Password is incorrect")
	}
	defer rows.Close()
	session.Id = rand.RandString(20)
	_, err = SessionRead(session.Id)

	for err == nil {
		session.Id = rand.RandString(20)
		_, err = SessionRead(session.Id)
	}
	session.ExpiryDate = time.Now().Unix() + 7776000 //+3 month
	session.UserId = userId
	_, err = session.Insert()
	if err != nil {
		return nil, err
	}
	return session, nil
}
func IsAuthentified(userId int64, sessionId string) (*Session, error) {
	session, err := SessionRead(sessionId)
	if err != nil {
		return nil, err
	}
	if session.UserId != userId {
		return nil, errors.New("Session User ID error : " + strconv.FormatInt(session.UserId, 10))
	}
	if session.ExpiryDate <= time.Now().Unix() {
		session.Delete()
		return nil, errors.New("Session expired")
	}
	return session, nil
}

func SessionRead(id string) (*Session, error) {
	var session Session
	rows, err := db.Query(
		"SELECT * FROM session WHERE id = ?",
		id,
	)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		err = rows.Scan(&session.Id, &session.UserId, &session.ExpiryDate)
		if err != nil {
			return nil, err
		}

	} else {
		return nil, errors.New("Session not found")
	}
	defer rows.Close()

	return &session, nil
}

func ReadSessionByUserId(userId int64) (*Session, error) {
	var session Session
	rows, err := db.Query(
		"SELECT * FROM session WHERE user_id = ?",
		userId,
	)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		rows.Scan(&session.Id)
		rows.Scan(&session.UserId)
		rows.Scan(&session.ExpiryDate)
	}
	defer rows.Close()
	return &session, nil
}

func (s *Session) Insert() (sql.Result, error) {
	return db.Exec(
		"INSERT INTO session (id,user_id,expiry_date) VALUES (?,?,?)",
		s.Id,
		s.UserId,
		s.ExpiryDate,
	)
}
func (s *Session) Delete() (sql.Result, error) {
	return db.Exec(
		"DELETE FROM session WHERE id = ?",
		s.Id,
	)
}
