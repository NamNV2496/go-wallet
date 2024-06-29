package migration

import (
	"context"
	"database/sql"
	"embed"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/namnv2496/go-wallet/config"

	sqlmigrate "github.com/rubenv/sql-migrate"
)

var (
	//go:embed sql/*
	migrationDirectoryMySQL embed.FS
)

func MigrateUp(ctx context.Context, config config.Config, types bool) error {

	if types {
		return GolangMigrate(ctx, config)
	}
	return MigrateSQLUp(ctx, config, sqlmigrate.Up)
}

func GolangMigrate(ctx context.Context, config config.Config) error {
	db, err := sql.Open("postgres", config.DBSource)
	if err != nil {
		log.Fatalln("Cannot connect to database")
		return err
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalln("Cannot get driver of database")
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		config.MigrationURL,
		"postgres",
		driver)
	if err != nil {
		log.Fatalln("Cannot create new instance of migration")
		return err
	}
	// only run 2 first files
	// err = m.Steps(2)
	// to run all file
	err = m.Up()
	log.Println("Migrated Database done: ", err)
	return nil
}

func MigrateSQLUp(ctx context.Context, config config.Config, direction sqlmigrate.MigrationDirection) error {

	db, err := sql.Open("postgres", config.DBSource)
	if err != nil {
		log.Fatalln("Cannot connect to database")
	}
	_, err = sqlmigrate.ExecContext(ctx, db, "postgres", sqlmigrate.EmbedFileSystemMigrationSource{
		FileSystem: migrationDirectoryMySQL,
		Root:       "sql",
	}, direction)
	if err != nil {
		log.Println("failed to execute migration: ", err)
		return err
	}
	log.Println("Executed migration successed!")
	return nil
}
