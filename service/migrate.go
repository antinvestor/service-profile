package service

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
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
	db.AutoMigrate(&AntMigration{}, &ProfileType{},
	&Profile{}, &ContactType{}, &CommunicationLevel{}, &Contact{}, &Country{}, &Address{}, &ProfileAddress{})

	if err := scanForNewMigrations(logger, db, migrationsDirPath); err != nil {
		logger.Warnf("Error scanning for new migrations : %v ", err)
		return
	}

	if err := applyNewMigrations(logger, db); err != nil {
		logger.Warnf("There was an error applying migrations : %v ", err)
	}
}

func scanForNewMigrations(logger *logrus.Entry, db *gorm.DB, migrationsDirPath string) error {

	// Get a list of migration files
	files, err := filepath.Glob(migrationsDirPath + "/*.sql")
	if err != nil {
		logger.Printf("Error running restore %s", err)
		return err
	}

	for _, file := range files {

		var migration AntMigration

		filename := filepath.Base(file)
		filename = strings.Replace(filename, ".sql", "", 1)

		migration.Name = filename
		migrationPatch, err := ioutil.ReadFile(file)

		if db.Where("name = ?", filename).Find(&migration).RecordNotFound() {

			if err != nil {
				logger.Warnf("Problem reading migration file content : %v", err)
				continue
			}
			migration.Patch = string(migrationPatch)
			migration.Version = 0
			db.Create(&migration)
		} else {

			if migration.AppliedAt == nil {

				if migration.Patch != string(migrationPatch) {
					db.Model(&migration).Update("patch", string(migrationPatch))
				}
			}

		}
	}
	return nil
}

func applyNewMigrations(logger *logrus.Entry, db *gorm.DB) error {

	unAppliedMigrations := []AntMigration{}
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
