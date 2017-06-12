package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ijsnow/gorest"
)

type key string

const (
	// SessionKey is the key where we will store sessions
	SessionKey key = "session-key"
	// UserKey is the key where we will store users
	UserKey key = "user-key"
)

type nameRes struct {
	Name string `json:"name"`
}

type user struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Hobbies []string `json:"hobbies"`
}

var users = []user{
	user{
		ID:      1,
		Name:    "Isaac",
		Hobbies: []string{"Rock Climbing", "Skiing", "Learning"},
	},
}

var sessions map[int]int

var sessionCounter = 0

func login(w http.ResponseWriter, r *http.Request) (int, interface{}) {
	user := user{}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return http.StatusInternalServerError, nil
	}

	isUser := false

	for _, u := range users {
		if u.ID == user.ID {
			isUser = true
			break
		}
	}

	if !isUser {
		return http.StatusNotFound, nil
	}

	sessionCounter++

	sessions[sessionCounter] = user.ID

	return http.StatusOK, sessionCounter
}

func getUser(w http.ResponseWriter, r *http.Request) (int, interface{}) {
	sctx := r.Context().Value(UserKey)
	u, ok := sctx.(*user)
	if !ok {
		fmt.Println("Cant conv to user", sctx)
		return http.StatusInternalServerError, nil
	}

	return http.StatusOK, u
}

func panicFunc(w http.ResponseWriter, r *http.Request) (int, interface{}) {
	panic("AAAAAHHHHH!!!!!!!")

	// Shouldn't get to here
	fmt.Println("Oops, it didn't work")
	return http.StatusOK, nil
}

func pointlessMiddleware(next gorest.HandlerFunc) gorest.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) (int, interface{}) {
		fmt.Println("This middleware just prints for fun :)")

		return next(w, r)
	}
}

func requireAuth(next gorest.HandlerFunc) gorest.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) (int, interface{}) {
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			return http.StatusUnauthorized, "Unauthorized"
		}

		token, err := strconv.Atoi(tokenStr)
		if err != nil {
			return http.StatusUnauthorized, "Unauthorized"
		}

		uID := sessions[token]
		if uID == 0 {
			return http.StatusUnauthorized, "Unauthorized"
		}

		var u user
		for _, v := range users {
			if uID == v.ID {
				u = v
				break
			}
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, &u)

		return next(w, r.WithContext(ctx))
	}
}

func main() {
	rt := gorest.NewRouter()

	rt.PostJSON("/login", login)
	rt.GetJSON("/user", getUser, pointlessMiddleware, requireAuth)
	rt.GetJSON("/panic", panicFunc, pointlessMiddleware)

	fmt.Println("Listening at http://localhost:8080")

	sessions = make(map[int]int)

	http.ListenAndServe(":8080", rt.GetHandler())
}
