package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
)

// Post represents a post.
type Post struct {
	ID         string  `json:"id"`
	UserID     string  `json:"user_id"`
	Title      string  `json:"title"`
	Body       string  `json:"body"`
	CategoryID int     `json:"category_id"`
	Price      uint    `json:"price"`
	X          float64 `json:"x"`
	Y          float64 `json:"y"`
	CreatedAt  int64   `json:"created_at"`
}

// Match selects best users to resolve p.
func (p *Post) Match(users *sync.Map, url string) error {
	// send post information to DB server
	postInfo := struct {
		UserID     string  `json:"user_id"`
		CategoryID int     `json:"category_id"`
		X          float64 `json:"x"`
		Y          float64 `json:"y"`
	}{UserID: p.UserID, CategoryID: p.CategoryID, X: p.X, Y: p.Y}

	buf, err := json.Marshal(postInfo)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// get matched users from DB server
	body := new(struct {
		Success  bool     `json:"successs"`
		ErrorMsg string   `json:"error_msg"`
		Users    []string `json:"users"`
	})

	if err = json.NewDecoder(resp.Body).Decode(body); err != nil {
		return err
	}

	if !body.Success {
		return errors.New(body.ErrorMsg)
	}

	// send post information to each client
	for _, id := range body.Users {
		go func(id string) {
			if user, exists := users.Load(id); exists {
				user.(*User).read <- p
			}
		}(id)
	}

	return nil
}
