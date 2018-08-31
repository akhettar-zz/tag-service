package main

import (
	"github.com/tag-service/api"
	_ "github.com/tag-service/docs"
	"github.com/tag-service/logger"
	"github.com/tag-service/repository"
	"github.com/tag-service/vault"
)

// @BasePath /
// @title Tags manager API
// @version 1.0

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /
func main() {

	logger.Info.Println("Starting up the server..")
	api.NewTagHandler(repository.NewRepository(vault.LoadConfig())).CreateRouter().Run(":8080")
	logger.Info.Println("Shutting down the server..")
}
