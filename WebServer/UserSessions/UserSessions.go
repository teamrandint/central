package usersessions

import (
	"seng468/WebServer/Commands"
)

type UserSessions interface {
	HasPendingBuys() bool
	HasPendingSells() bool
	UserId() string
}

type UserSession struct {
	userId       string
	PendingBuys  []*commands.Command
	PendingSells []*commands.Command
}

func NewUserSession(id string) *UserSession {
	session := new(UserSession)
	session.userId = id
	return session
}

func (session *UserSession) HasPendingBuys() bool {
	if len(session.PendingBuys) == 0 {
		return false
	} else {
		return true
	}
}

func (session *UserSession) HasPendingSells() bool {
	if len(session.PendingSells) == 0 {
		return false
	} else {
		return true
	}
}

func (session *UserSession) UserId() string {
	return session.userId
}
