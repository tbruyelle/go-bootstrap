package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/rnubel/pgmgr/pgmgr"
	"os"
)

func displayErrorOrMessage(err error, args ...interface{}) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err)
		os.Exit(1)
	} else {
		fmt.Println(args...)
	}
}

func displayVersion(config *pgmgr.Config) {
	v, err := pgmgr.Version(config)
	if v < 0 {
		displayErrorOrMessage(err, "Database has no schema_migrations table; run `pgmgr db migrate` to create it.")
	} else {
		displayErrorOrMessage(err, "Latest migration version:", v)
	}
}

func main() {
	config := &pgmgr.Config{}
	app := cli.NewApp()

	app.Name = "pgmgr"
	app.Usage = "manage your app's Postgres database"
	app.Version = "0.0.1"

	s := make([]string, 0)

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config-file, c",
			Value:  ".pgmgr.json",
			Usage:  "set the path to the JSON configuration file specifying your DB parameters",
			EnvVar: "PGMGR_CONFIG_FILE",
		},
		cli.StringFlag{
			Name:   "database, d",
			Value:  "",
			Usage:  "the database name which pgmgr will connect to or try to create",
			EnvVar: "PGMGR_DATABASE",
		},
		cli.StringFlag{
			Name:   "username, u",
			Value:  "",
			Usage:  "the username which pgmgr will connect with",
			EnvVar: "PGMGR_USERNAME",
		},
		cli.StringFlag{
			Name:   "password, P",
			Value:  "",
			Usage:  "the password which pgmgr will connect with",
			EnvVar: "PGMGR_PASSWORD",
		},
		cli.StringFlag{
			Name:   "host, H",
			Value:  "",
			Usage:  "the host which pgmgr will connect to",
			EnvVar: "PGMGR_HOST",
		},
		cli.IntFlag{
			Name:   "port, p",
			Value:  0,
			Usage:  "the port which pgmgr will connect to",
			EnvVar: "PGMGR_PORT",
		},
		cli.StringFlag{
			Name:   "url",
			Value:  "",
			Usage:  "connection URL or DSN containing connection info; will override the other params if given",
			EnvVar: "PGMGR_URL",
		},
		cli.StringFlag{
			Name:   "dump-file",
			Value:  "",
			Usage:  "where to dump or load the database structure and contents to or from",
			EnvVar: "PGMGR_DUMP_FILE",
		},
		cli.StringFlag{
			Name:   "migration-folder",
			Value:  "",
			Usage:  "folder containing the migrations to apply",
			EnvVar: "PGMGR_MIGRATION_FOLDER",
		},
		cli.StringSliceFlag{
			Name:   "seed-tables",
			Value:  (*cli.StringSlice)(&s),
			Usage:  "list of tables (or globs matching table names) to dump the data of",
			EnvVar: "PGMGR_SEED_TABLES",
		},
	}

	app.Before = func(c *cli.Context) error {
		// TODO: LoadConfig should validate some basic properties of a valid config,
		// like that the database name is set, and return an error if not.
		pgmgr.LoadConfig(config, c)

		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:  "migration",
			Usage: "generates a new migration with the given name",
			Action: func(c *cli.Context) {
				if len(c.Args()) == 0 {
					println("migration name not given! try `pgmgr migration NameGoesHere`")
				} else {
					pgmgr.CreateMigration(config, c.Args()[0])
				}
			},
		},
		{
			Name:  "config",
			Usage: "displays the current configuration as seen by pgmgr",
			Action: func(c *cli.Context) {
				fmt.Printf("%+v\n", config)
			},
		},
		{
			Name:  "db",
			Usage: "manage your database. use 'pgmgr db help' for more info",
			Subcommands: []cli.Command{
				{
					Name:  "create",
					Usage: "creates the database if it doesn't exist",
					Action: func(c *cli.Context) {
						displayErrorOrMessage(pgmgr.Create(config), "Database", config.Database, "created successfully.")
					},
				},
				{
					Name:  "drop",
					Usage: "drops the database (all sessions must be disconnected first. this command does not force it)",
					Action: func(c *cli.Context) {
						displayErrorOrMessage(pgmgr.Drop(config), "Database", config.Database, "dropped successfully.")
					},
				},
				{
					Name:  "dump",
					Usage: "dumps the database schema and contents to the dump file (see --dump-file)",
					Action: func(c *cli.Context) {
						err := pgmgr.Dump(config)
						displayErrorOrMessage(err, "Database dumped to", config.DumpFile, "successfully")
					},
				},
				{
					Name:  "load",
					Usage: "loads the database schema and contents from the dump file (see --dump-file)",
					Action: func(c *cli.Context) {
						err := pgmgr.Load(config)
						displayErrorOrMessage(err, "Database loaded successfully.")
						displayVersion(config)
					},
				},
				{
					Name:  "version",
					Usage: "returns the current schema version",
					Action: func(c *cli.Context) {
						displayVersion(config)
					},
				},
				{
					Name:  "migrate",
					Usage: "applies any un-applied migrations in the migration folder (see --migration-folder)",
					Action: func(c *cli.Context) {
						err := pgmgr.Migrate(config)
						if err != nil {
							fmt.Fprintln(os.Stderr, "Error during migration:", err)
							os.Exit(1)
						}
					},
				},
				{
					Name:  "rollback",
					Usage: "rolls back the latest migration",
					Action: func(c *cli.Context) {
						pgmgr.Rollback(config)
					},
				},
			},
		},
	}

	app.Action = func(c *cli.Context) {
		app.Command("help").Run(c)
	}

	app.Run(os.Args)
}
