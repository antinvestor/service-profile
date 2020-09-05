package service

import (
	"github.com/antinvestor/service-profile/models"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

/*
This is just documentation on creating a postgresql db and user

postgres=# create database service-profile; create user profile with encrypted password 'profile3'; grant all privileges on database "service-profile" to profile;

*/

// PerformMigration finds missing migrations and records them in the database,
// We use the fragmenta_metadata table to do this
func PerformMigration(logger *logrus.Entry, db *gorm.DB) {

	migrationsDirPath := "./migrations/0001"

	// Migrate the schema
	db.AutoMigrate(&models.Migration{}, &models.ProfileType{},
	&models.Profile{}, &models.ContactType{}, &models.CommunicationLevel{},
	&models.Contact{}, &models.Country{}, &models.Address{}, &models.ProfileAddress{},
	&models.Verification{}, &models.VerificationAttempt{})

	if err := scanForNewMigrations(logger, db, migrationsDirPath); err != nil {
		logger.Warnf("Error scanning for new migrations : %v ", err)
		return
	}

	if err := applyNewMigrations(logger, db); err != nil {
		logger.Warnf("There was an error applying migrations : %v ", err)
	}
}

func scanForNewMigrations(logger *logrus.Entry, db *gorm.DB, migrationsDirPath string) error {

	logger.Info("scanning for new migrations")
	// Get a list of migration files
	files, err := filepath.Glob(migrationsDirPath + "/*.sql")
	if err != nil {
		logger.Printf("Error running restore %s", err)
		return err
	}

	logger.Infof("found %d migrations to process", len(files))

	for _, file := range files {

		var migration models.Migration

		filename := filepath.Base(file)
		filename = strings.Replace(filename, ".sql", "", 1)

		migration.Name = filename
		migrationPatch, err := ioutil.ReadFile(file)

		if db.Where("name = ?", filename).Find(&migration).RecordNotFound() {

			if err != nil {
				logger.Warnf("Problem reading migration file content : %v", err)
				continue
			}
			logger.Infof("migration %s is unapplied", file)
			migration.Patch = string(migrationPatch)

			err = db.Create(&migration).Error
			if err != nil {
				logger.WithError(err).Warnf("There is an error adding migration :%s", file)
			}
		} else {

			if migration.AppliedAt == nil {

				if migration.Patch != string(migrationPatch) {
					err = db.Model(&migration).Update("patch", string(migrationPatch)).Error

					if err != nil {
						logger.WithError(err).Warnf("There is an error updating migration :%s", file)
					}
				}
			}

		}
	}
	return nil
}

func applyNewMigrations(logger *logrus.Entry, db *gorm.DB) error {

	unAppliedMigrations := []models.Migration{}
	if err := db.Where("applied_at IS NULL").Find(&unAppliedMigrations).Error; err != nil {
		return err
	}

	for _, migration := range unAppliedMigrations {

		if err := db.Exec(migration.Patch).Error; err != nil {
			return err
		}

		db.Model(&migration).UpdateColumn("applied_at", time.Now())
		logger.Infof("Successfully applied the file : %v", fmt.Sprintf("%s.sql", migration.Name))
	}

	return nil
}
