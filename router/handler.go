package router

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/insanesclub/sasohan-match/model"
	"github.com/labstack/echo/v4"
)

// Connect creates a new connection.
func Connect(users *sync.Map, upgrader websocket.Upgrader) echo.HandlerFunc {
	return func(c echo.Context) error {
		conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}

		// get user ID
		body := new(struct {
			ID string `json:"id"`
		})

		// get payload from socket
		// payload is an in-memory buffer
		// and max size of the buffer is set by Conn.SetReadLimit
		// default max size is 4KB
		// see https://pkg.go.dev/github.com/gorilla/websocket
		//
		// if payload size exceeds limit, ErrReadLimit occurs
		if err = conn.ReadJSON(body); err != nil {
			return err
		}

		// create a user
		user := model.NewUser(body.ID, conn)
		users.Store(body.ID, user)

		go user.Run(users)

		return nil
	}
}

// Disconnect closes a connection.
func Disconnect(users *sync.Map) echo.HandlerFunc {
	return func(c echo.Context) error {
		// get user ID
		body := new(struct {
			ID string `json:"id"`
		})

		if err := c.Bind(body); err != nil {
			return err
		}

		users.Delete(body.ID)

		return nil
	}
}

// Recommend recommends best users to resolve a demand.
func Recommend(users *sync.Map) echo.HandlerFunc {
	return func(c echo.Context) error {
		// get post information
		post := new(model.Post)
		if err := c.Bind(post); err != nil {
			return err
		}

		// find matched users and send post information
		if err := post.Match(users, "http://localhost:3000/match"); err != nil {
			return err
		}

		return nil
	}
}

// Match asks the demander to match with the provider or not.
func Match(users *sync.Map) echo.HandlerFunc {
	return func(c echo.Context) error {
		provider := new(struct {
			User model.User `json:"user"`
			Post model.Post `json:"post"`
		})

		if err := c.Bind(provider); err != nil {
			return err
		}

		if u, exists := users.Load(provider.Post.UserID); exists {
			if err := u.(*model.User).Match(provider.User, provider.Post); err != nil {
				return err
			}
		}

		return nil
	}
}
