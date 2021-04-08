package main

import (
	"errors"
	"fmt"

	"github.com/mect/expino-api/pkg/db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	rootCmd.AddCommand(NewAddUserCmd())
}

type addUserCmdOptions struct {
	Username string
	Password string
	Name     string

	postgresHost     string
	postgresPort     int
	postgresUsername string
	postgresDatabase string
	postgresPassword string
}

// NewServeCmd generates the `serve` command
func NewAddUserCmd() *cobra.Command {
	a := addUserCmdOptions{}
	c := &cobra.Command{
		Use:     "add-user",
		Short:   "adds a user to the database",
		PreRunE: a.Validate,
		RunE:    a.RunE,
	}
	c.Flags().StringVarP(&a.Username, "username", "u", "", "Username for the user")
	c.Flags().StringVarP(&a.Password, "password", "p", "", "Password for the user")
	c.Flags().StringVarP(&a.Name, "name", "n", "", "Visible name fore the user")

	c.Flags().StringVar(&a.postgresHost, "postgres-host", "", "PostgreSQL hostname")
	c.Flags().IntVar(&a.postgresPort, "postgres-port", 5432, "PostgreSQL hostname")
	c.Flags().StringVar(&a.postgresUsername, "postgres-username", "", "PostgreSQL hostname")
	c.Flags().StringVar(&a.postgresPassword, "postgres-password", "", "PostgreSQL hostname")
	c.Flags().StringVar(&a.postgresDatabase, "postgres-database", "", "PostgreSQL hostname")

	viper.BindPFlags(c.Flags())

	return c
}

func (a *addUserCmdOptions) Validate(cmd *cobra.Command, args []string) error {
	if a.Username == "" {
		return errors.New("need to set --username")
	}

	if a.Password == "" {
		return errors.New("need to set --password")
	}

	if a.Name == "" {
		return errors.New("need to set --name")
	}

	return nil
}

func (a *addUserCmdOptions) RunE(cmd *cobra.Command, args []string) error {
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	res := dbConn.Create(&db.User{
		Name:     a.Name,
		Username: a.Username,
		Password: string(hashedPassword),
	})
	return res.Error
}
