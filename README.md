# gorest
A simple router for REST APIs built on top of https://github.com/gowww/router.

## Installation

```
$ go get github.com/ijsnow/gorest
```

```go
import "github.com/ijsnow/gorest"
```

## Usage

1. Create a router

```go
rt := gorest.NewRouter()
```

2. Declare some routes with `Get`, `Post`, `Put`, or `Delete` methods.

Using basic route generator. You have to give the route your [WriteFunc](https://github.com/ijsnow/gorest/blob/master/write.go#L10).

```go
// Return an http status code and the data you would like to serve
func getUser(w http.ResponseWriter, r *http.Request) (int, interface{}) {
	user := User{
    Name: "Isaac",
  }

	return http.StatusOK, user
}

// ...

rt.Get(gorest.JSON, "/user", getUser)
```

Using the JSON helper
```go
rt.GetJSON("/user", getUser)
```

3. Give your server the handler

```go
http.ListenAndServe(":8080", rt.GetHandler())
```

## Additional Info

1. Using middlewares

```go
// I'll get called first
func middleware1(next gorest.HandlerFunc) gorest.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) (int, interface{}) {
		fmt.Println("This is middleware 1")

		return next(w, r)
	}
}

// I'll get called second
func middleware2(next gorest.HandlerFunc) gorest.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) (int, interface{}) {
		fmt.Println("This is middleware 2, let's gamble")
    
    if tossCoin() == "heads" {
      return http.StatusBadRequest, "You lose. No serving cool data for you."
    }

		return next(w, r)
	}
}

// ...

// Execution of middlewares flows from right to left
rt.GetJSON("/user", getUser, middleware2, middleware1)
```
