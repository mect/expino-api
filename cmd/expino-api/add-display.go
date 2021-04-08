package main

import (
	"errors"
	"fmt"

	"github.com/mect/expino-api/pkg/db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(NewAddDisplayCmd())
}

type addDisplayCmdOptions struct {
	Token string
	Name  string

	postgresHost     string
	postgresPort     int
	postgresUsername string
	postgresDatabase string
	postgresPassword string
}

// NewAddDisplayCmd generates the `add-display` command
func NewAddDisplayCmd() *cobra.Command {
	a := addDisplayCmdOptions{}
	c := &cobra.Command{
		Use:     "add-display",
		Short:   "adds a display to the database",
		PreRunE: a.Validate,
		RunE:    a.RunE,
	}
	c.Flags().StringVarP(&a.Token, "password", "t", "", "Token for the display")
	c.Flags().StringVarP(&a.Name, "name", "n", "", "Visible name fore the display")

	c.Flags().StringVar(&a.postgresHost, "postgres-host", "", "PostgreSQL hostname")
	c.Flags().IntVar(&a.postgresPort, "postgres-port", 5432, "PostgreSQL hostname")
	c.Flags().StringVar(&a.postgresUsername, "postgres-username", "", "PostgreSQL hostname")
	c.Flags().StringVar(&a.postgresPassword, "postgres-password", "", "PostgreSQL hostname")
	c.Flags().StringVar(&a.postgresDatabase, "postgres-database", "", "PostgreSQL hostname")

	viper.BindPFlags(c.Flags())

	return c
}

func (a *addDisplayCmdOptions) Validate(cmd *cobra.Command, args []string) error {
	if a.Token == "" {
		return errors.New("need to set --token")
	}

	if a.Name == "" {
		return errors.New("need to set --name")
	}

	return nil
}

func (a *addDisplayCmdOptions) RunE(cmd *cobra.Command, args []string) error {
	dbConn, err := db.NewConnection(db.ConnectionDetails{
		Host:     a.postgresHost,
		Port:     a.postgresPort,
		User:     a.postgresUsername,
		Database: a.postgresDatabase,
		Password: a.postgresPassword,
	})

	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}

	res := dbConn.Create(&db.Display{
		Name:  a.Name,
		Token: a.Token,
	})
	return res.Error
}
