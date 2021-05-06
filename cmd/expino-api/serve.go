package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/mect/expino-api/pkg/api/display"

	socketio "github.com/googollee/go-socket.io"

	"github.com/dgrijalva/jwt-go"

	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/websocket"

	"github.com/mect/expino-api/pkg/api/auth"
	v1 "github.com/mect/expino-api/pkg/api/v1"
	"github.com/mect/expino-api/pkg/db"
)

var io *socketio.Server

// this is used to compare to in case of no user found to keep the response time the same
const dummyHash = `$2a$10$8KqKzq6uHCL72Qhshj9L.uGUz/0lmkjupqYQKCy6th9Rv91k53g82`

func init() {
	rootCmd.AddCommand(NewServeCmd())
}

var protectedEntryPoints = []string{"/v1"}

type serveCmdOptions struct {
	BindAddr string
	Port     int

	jwtSecret string

	db *db.Connection

	postgresHost     string
	postgresPort     int
	postgresUsername string
	postgresDatabase string
	postgresPassword string

	wsClients map[string]chan string
}

// NewServeCmd generates the `serve` command
func NewServeCmd() *cobra.Command {
	s := serveCmdOptions{
		wsClients: map[string]chan string{},
	}
	c := &cobra.Command{
		Use:     "serve",
		Short:   "Serves the HTTP REST endpoint",
		Long:    `Serves the HTTP REST endpoint on the given bind address and port`,
		PreRunE: s.Validate,
		RunE:    s.RunE,
	}
	c.Flags().StringVarP(&s.BindAddr, "bind-address", "b", "0.0.0.0", "address to bind port to")
	c.Flags().IntVarP(&s.Port, "port", "p", 8080, "Port to listen on")

	c.Flags().StringVar(&s.jwtSecret, "jwt-secret", "", "JWT siging key, please do not set in flags")

	c.Flags().StringVar(&s.postgresHost, "postgres-host", "", "PostgreSQL hostname")
	c.Flags().IntVar(&s.postgresPort, "postgres-port", 5432, "PostgreSQL hostname")
	c.Flags().StringVar(&s.postgresUsername, "postgres-username", "", "PostgreSQL hostname")
	c.Flags().StringVar(&s.postgresPassword, "postgres-password", "", "PostgreSQL hostname")
	c.Flags().StringVar(&s.postgresDatabase, "postgres-database", "", "PostgreSQL hostname")
	return c
}

func (s *serveCmdOptions) Validate(cmd *cobra.Command, args []string) error {
	if s.jwtSecret == "" {
		return errors.New("jwt-secret not set")
	}

	if s.postgresUsername == "" || s.postgresPassword == "" || s.postgresDatabase == "" || s.postgresHost == "" {
		return errors.New("PostgreSQL credentials not set")
	}

	return nil
}

func (s *serveCmdOptions) RunE(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	printLogo()

	var err error
	s.db, err = db.NewConnection(db.ConnectionDetails{
		Host:     s.postgresHost,
		Port:     s.postgresPort,
		User:     s.postgresUsername,
		Database: s.postgresDatabase,
		Password: s.postgresPassword,
	})
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}

	err = s.db.DoMigrate()
	if err != nil {
		return fmt.Errorf("error migrating database: %w", err)
	}

	e := echo.New()

	e.HideBanner = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper: func(c echo.Context) bool {
			if strings.HasPrefix(c.Path(), "/ws/") {
				return true
			}
			return false
		},
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))
	e.Use(socketioCORS)
	e.Use(cacheStatic)
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(s.jwtSecret),
		Claims:     &auth.Claim{},
		Skipper: func(c echo.Context) bool {
			// always skip JWT unless path is a protectedPrefix
			for _, protectedPrefix := range protectedEntryPoints {
				if strings.HasPrefix(c.Path(), protectedPrefix) {
					return false
				}
			}
			return true
		},
	}))

	// handlers
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Expino API endpoint")
	})

	e.Any("/ws", s.serveWs)

	http.Handle("/", e)

	e.POST("/login", s.login)
	e.Static("/static", "expino-static")
	v1.NewHTTPHandler(s.db, s.sendBroadcast).Register(e)
	display.NewHTTPHandler(s.db).Register(e)

	go func() {
		log.Printf("Listening on %s:%d\n", s.BindAddr, s.Port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", s.BindAddr, s.Port), nil))
		cancel() // server ended, stop the world
	}()

	go func() {
		for {
			s.sendBroadcast("PING")
			time.Sleep(time.Second)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
			return nil
		}
	}
}

type AuthData struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}

// TODO: move this!
func (s *serveCmdOptions) login(c echo.Context) error {
	data := new(AuthData)
	err := c.Bind(data)
	if err != nil {
		log.Println(err)
	}

	if data.Username == "" || data.Password == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "username or password not specified"})
	}

	data.Username = strings.ToLower(data.Username)

	user := db.User{}
	err = s.db.GetWhereIs(&user, "username", data.Username)
	if errors.Is(err, db.ErrorNotFound) {
		_ = bcrypt.CompareHashAndPassword([]byte(dummyHash), []byte(data.Password))
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "username or password incorrect"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "username or password incorrect"})
	}

	// Set custom claims
	claims := &auth.Claim{
		user.Name,
		user.ID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

func socketioCORS(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Println(c.Path())
		if strings.HasPrefix(c.Path(), "/socket.io/") {
			if origin := c.Request().Header.Get("Origin"); origin != "" {
				c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
				c.Response().Header().Set("Access-Control-Allow-Origin", origin)
			}
		}
		return next(c)
	}
}

func cacheStatic(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if strings.HasPrefix(c.Path(), "/static/") {
			c.Response().Header().Set("Cache-Control", "max-age:290304000, public")
			c.Response().Header().Set("Last-Modified", time.Now().Format(http.TimeFormat))
			c.Response().Header().Set("Expires", time.Now().AddDate(60, 0, 0).Format(http.TimeFormat)) // we use SHAs fo file names 60 years will be fine
		}
		return next(c)
	}
}

func (s *serveCmdOptions) sendBroadcast(data string) {
	for _, ch := range s.wsClients {
		go func() {
			ch <- data
		}()
	}
}

func (s *serveCmdOptions) serveWs(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		id := uuid.New().String()

		updates := make(chan string)
		s.wsClients[id] = updates

		for {
			u := <-updates
			// Write
			err := websocket.Message.Send(ws, u)
			if err != nil {
				c.Logger().Error(err)
				delete(s.wsClients, id)
				break
			}
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}
