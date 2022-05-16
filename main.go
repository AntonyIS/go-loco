package main

import (
	"fmt"
	"log"
	"os"

	h "github.com/AntonyIS/go-loco/api"
	"github.com/AntonyIS/go-loco/app"
	ddr "github.com/AntonyIS/go-loco/repository/dynamodb"
	"github.com/gin-gonic/gin"
)

func main() {
	repo := chooseRepo()
	service := app.NewLocomotiveService(repo)
	handler := h.NewHandler(service)

	r := gin.Default()
	r.POST("/api/v1/add", handler.CreateLoco)
	r.GET("/api/v1/:loco_id", handler.GetLoco)
	r.GET("/api/v1/all", handler.GetAllLoco)
	r.PUT("/api/v1/update/:loco_id", handler.UpdateLoco)
	r.DELETE("/api/v1/delete/:loco_id", handler.DeleteLoco)

	errs := make(chan error, 2)
	go func() {
		fmt.Println("Listening on port :5000")

		errs <- r.Run(":5000")
	}()

	go func() {
		c := make(chan os.Signal, 1)
		errs <- fmt.Errorf("%s", <-c)
	}()

	fmt.Printf("Terminated %s", <-errs)
}

// func httpPort() string {
// 	port := "5000"
// 	if os.Getenv("PORT") != "" {
// 		port = os.Getenv("PORT")
// 	}

// 	return fmt.Sprintf(":%s", port)
// }

func chooseRepo() app.LocomotiveRepository {
	repo, err := ddr.NewDynamoDBReposistory("loco")
	if err != nil {
		log.Fatal(err)
	}
	return repo
}
