package main

import (
	"bitbucket.org/antinvestor/service-boilerplate/service"
	"log"
	"os"
	"time"

	"bitbucket.org/antinvestor/service-boilerplate/utils"
)

func main() {

	serviceName := "BoilerPlate"

	logger, err := utils.ConfigureLogging(serviceName)
	if err != nil {
		log.Fatal("Failed to configure logging: " + err.Error())
	}

	closer, err := utils.ConfigureJuegler(serviceName)
	if err != nil {
		logger.Fatal("Failed to configure Juegler: " + err.Error())
	}

	defer closer.Close()

	database, err := utils.ConfigureDatabase(logger)
	if err != nil {
		logger.Fatalf("Failed to configure Database: %v", err)
	}

	stdArgs := os.Args[1:]
	if len(stdArgs) > 0 && stdArgs[0] == "migrate" {
		logger.Info("Initiating migrations")

		service.PerformMigration(logger, database)

	} else {
		logger.Infof("Initiating the service at %v", time.Now())

		env := service.Env{
			Logger:          logger,
			ServerPort: utils.GetEnv("SERVER_PORT", "7000"),
		}
		env.SetDb(database)

		service.RunServer(&env)
	}


}
