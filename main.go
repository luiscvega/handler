package main

import (
	"fmt"
	"net/http"

	"github.com/luiscvega/routes"
)

type specialNetHttpGreetingHandler struct {
	greeting string
}

func (h specialNetHttpGreetingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bs := []byte(fmt.Sprintf("%s, %s!\n", h.greeting, r.URL.Path[3:]))

	w.Write(bs)
}

type specialLuiscvegaRoutesGreetingHandler struct {
	key      string
	greeting string
}

func (h specialLuiscvegaRoutesGreetingHandler) Serve(w http.ResponseWriter, r *http.Request, params map[string]string) {
	value := params[h.key]

	bs := []byte(fmt.Sprintf("%s, %s!\n", h.greeting, value))

	w.Write(bs)
}

func main() {
	// =================================
	// Using net/http (standard library)
	// =================================

	// http.HandlerFunc type casts the anonymous function into valid http.Handler interface
	http.HandleFunc("/a/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bs := []byte(fmt.Sprintf("Hello, %s!\n", r.URL.Path[3:]))

		w.Write(bs)
	}))

	// Use named func that implements the http.Handler interface
	hi := specialNetHttpGreetingHandler{"Hi"}

	// Note that since we use a struct that implements that http.Handler interface (see specialNetHttpGreetingHandler), we
	// will use the http.Handle which accepts an interface (as opposed to http.HandlerFunc which accepts a func that
	// will by type casted into a http.Handler Interface
	http.Handle("/b/", hi)

	// =================================
	// Using github.com/luiscvega/routes
	// =================================
	rs := routes.Routes{}

	// Same as above, my routes packages converts an anonymous function to a type that implements my routes.Handler interface
	rs.Add("GET", "/c/:foo", routes.HandlerFunc(func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		bs := []byte(fmt.Sprintf("Goodbye, %s!\n", params["foo"]))

		w.Write(bs)
	}))

	// Same as above, but now I can pass in some state like the key that parameter will be named after
	ciao := specialLuiscvegaRoutesGreetingHandler{
		key:      "fullname",
		greeting: "Ciao",
	}

	rs.Add("GET", "/d/:fullname", ciao)

	// rs is a routes.Routes struct, which implements the http.Handler interface, there we can pass this in same as in /b"
	http.Handle("/", rs)

	if err := http.ListenAndServe(":8342", nil); err != nil {
		panic(err)
	}
}
