package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	h "github.com/AntonyIS/go-loco/api"
	"github.com/AntonyIS/go-loco/app"
	ddr "github.com/AntonyIS/go-loco/repository/dynamodb"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	repo := chooseRepo()
	service := app.NewLocomotiveService(repo)
	handler := h.NewHandler(service)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/{loco_id}", handler.Get)
	// r.Post("/", handler.Post)

	errs := make(chan error, 2)
	go func() {
		fmt.Println("Listening on port :5000")
		errs <- http.ListenAndServe(httpPort(), r)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		errs <- fmt.Errorf("%s", <-c)
	}()

	fmt.Printf("Terminated %s", <-errs)
}

func httpPort() string {
	port := "5000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	return fmt.Sprintf(":%s", port)
}

func chooseRepo() app.LocomotiveRepository {
	switch os.Getenv("DYNAMODB_TABLE") {
	case "redis":
		table := os.Getenv("DYNAMODB_TABLE")
		repo, err := ddr.NewDynamoDBReposistory(table)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	}
	return nil
}
