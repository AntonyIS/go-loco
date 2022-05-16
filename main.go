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

	r.Get("/api/v1/get/{loco_id}", handler.GetLoco)
	r.Get("/api/v1/all", handler.GetAllLoco)
	r.Post("/api/v1/add", handler.PostLoco)
	r.Put("/api/v1/update/{loco_id}", handler.UpdateLoco)
	r.Delete("/api/v1/delete/{loco_id}", handler.DeleteLoco)

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
	repo, err := ddr.NewDynamoDBReposistory("loco")
	if err != nil {
		log.Fatal(err)
	}
	return repo
}
