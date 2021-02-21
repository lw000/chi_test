// chi_test project main.go
package main

import (
	"fmt"
	"github.com/asaskevich/EventBus"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/unrolled/render"
	"log"
	"net/http"
)

func calculator(a int, b int) {
	fmt.Printf("%d\n", a+b)
}

func main() {
	r := render.New()
	engine := chi.NewRouter()

	engine.Use(middleware.RequestID)
	engine.Use(middleware.RealIP)
	engine.Use(middleware.Logger)
	engine.Use(middleware.Recoverer)
	//engine.Use(cors.AllowAll().Handler)
	engine.Use(cors.Handler(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	engine.Get("/", func(w http.ResponseWriter, req *http.Request) {
		err := r.JSON(w, http.StatusOK, map[string]interface{}{
			"code": 0,
			"msg":  "ok",
			"data": nil,
		})
		if err != nil {
			log.Fatal(err)
		}
	})

	bus := EventBus.New()
	_ = bus.Subscribe("main:calculator", calculator)
	bus.Publish("main:calculator", 20, 40)
	_ = bus.Unsubscribe("main:calculator", calculator)

	_ = http.ListenAndServe(":8888", engine)
}
