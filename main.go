package main

import (
	"log"
	"os"
	"time"

	"antinvestor.com/service/profile/service"
	"antinvestor.com/service/profile/utils"
)

func main() {

	serviceName := "Profile"

	logger, err := utils.ConfigureLogging(serviceName)
	if err != nil {
		log.Fatal("Failed to configure logging: " + err.Error())
	}

	closer, err := utils.ConfigureJuegler(serviceName)
	if err != nil {
		logger.Fatal("Failed to configure Juegler: " + err.Error())
	}

	defer closer.Close()

	database, err := utils.ConfigureDatabase(logger, false)
	if err != nil {
		logger.WithError(err).Fatal("Could not Configure write database")
	}
	defer database.Close()

	replicaDatabase, err := utils.ConfigureDatabase(logger, true)
	if err != nil {
		logger.WithError(err).Fatal("Could not Configure read database")
	}
	defer replicaDatabase.Close()

	onlyMigrate := utils.GetEnv(utils.EnvOnlyMigrate, "")
	isMigration := utils.GetEnv(utils.EnvMigrate, onlyMigrate)
	stdArgs := os.Args[1:]
	if (len(stdArgs) > 0 && stdArgs[0] == "migrate") || isMigration == "true" {
		logger.Info("Initiating migrations")

		service.PerformMigration(logger, database)

		if onlyMigrate == "true" {
			return
		}

	}

	logger.Infof("Initiating the service at %v", time.Now())

	healthChecker, err := utils.ConfigureHealthChecker(logger, database, replicaDatabase)
	if err != nil {
		logger.Warnf("Error configuring health checks: %v", err)
	}

	env := utils.Env{
		Logger: logger,
		Health: healthChecker,
	}
	env.SetWriteDb(database)
	env.SetReadDb(replicaDatabase)

	service.RunServer(&env)

}
