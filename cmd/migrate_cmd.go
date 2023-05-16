package cmd

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/spf13/cobra"
	"micro/config"
)

var migrateCMD = cobra.Command{
	Use:     "migrate",
	Long:    "migrate database strucutures. This will migrate tables",
	Aliases: []string{"m"},
	Run:     Runner.Migrate,
}

// migrate database with fake data
func (c *command) Migrate(cmd *cobra.Command, args []string) {
	dir, _ := os.Getwd()
	dsn := "postgresql://" + config.C().Postgres.Username + ":" + config.C().Postgres.Password + "@" + config.C().Postgres.Host + "/" + config.C().Postgres.Schema
	_ = "host=" + config.C().Postgres.Host + " user=" + config.C().Postgres.Username + " password=" + config.C().Postgres.Password + " dbname=" + config.C().Postgres.Schema + " port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	sqlConn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalln("error in connection to database")
	}

	migrations := &migrate.FileMigrationSource{
		Dir: dir + "/migrations",
	}

	cmdArg := args[0]
	switch cmdArg {
	case "up":
		n, err := migrate.Exec(sqlConn, "postgres", migrations, migrate.Up)
		if err != nil {
			log.Fatalf("failed to apply migrations: %v", err)
		}
		log.Printf("Applied %d migrations\n", n)
	case "down":
		n, err := migrate.Exec(sqlConn, "postgres", migrations, migrate.Down)
		if err != nil {
			log.Fatalf("failed to apply migrations: %v", err)
		}
		log.Printf("Applied %d migrations\n", n)
	}

}
