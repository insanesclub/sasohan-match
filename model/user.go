package model

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// User represents a user.
type User struct {
	id   string // user ID
	conn *websocket.Conn
	read chan *Post
	quit chan struct{}
}

// NewUser returns a new user.
func NewUser(id string, conn *websocket.Conn) *User {
	return &User{
		id:   id,
		conn: conn,
		read: make(chan *Post),
		quit: make(chan struct{}),
	}
}

// Run starts user activation.
func (u *User) Run(users *sync.Map) {
	for {
		select {
		case post := <-u.read:
			u.send(post)
		case <-u.quit:
			users.Delete(u.id)
			close(u.read)
			close(u.quit)
			u.conn.Close()
			return
		}
	}
}

// send sends post to client.
func (u *User) send(post *Post) {
	if err := u.conn.WriteJSON(*post); err != nil {
		log.Println(err)
		u.Quit()
	}
}

// Quit alerts u to quit.
func (u *User) Quit() { u.quit <- struct{}{} }

// Match asks u to
func (u *User) Match(provider string, post Post) error {
	if err := u.conn.WriteJSON(struct {
		User string `json:"user"`
		Post Post   `json:"post"`
	}{User: provider, Post: post}); err != nil {
		log.Println(err)
		u.Quit()
		return err
	}

	return nil
}
