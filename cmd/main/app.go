package main

import (
	"context"
	"fmt"
	"go_advanced/internal/config"
	"go_advanced/internal/user"
	"go_advanced/internal/user/db"
	"go_advanced/pkg/client/mongodb"
	"go_advanced/pkg/logging"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("create router")
	router := httprouter.New()

	cfg := config.GetConfig()

	cfgM := cfg.MongoDB
	mongoDBClient, err := mongodb.NewClient(context.Background(),
		cfgM.Host, cfgM.Port, cfgM.Username, cfgM.Password, cfgM.Database, cfgM.AuthDB)
	if err != nil {
		panic(err)
	}

	storage := db.NewStorage(mongoDBClient, cfg.MongoDB.Collection, logger)

	user1 := user.User{
		ID:           "",
		Username:     "User 1",
		PasswordHash: "12345",
		Email:        "user1@mail.com",
	}
	user1Id, err := storage.Create(context.Background(), user1)
	if err != nil {
		panic(err)
	}

	logger.Info(user1Id)

	logger.Info("register handler")
	handler := user.NewHandler(logger)
	handler.Register(router)

	start(router, cfg)

}

func start(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("start application")

	var listner net.Listener
	var listnerErr error

	if cfg.Listen.Type == "sock" {

		addDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}
		logger.Info("create socket")
		socketPath := path.Join(addDir, "app.sock")
		logger.Debugf("socket path: %s", socketPath)

		logger.Info("listen  unix socket")
		listner, listnerErr = net.Listen("unix", socketPath)
		logger.Infof("server is listening on unix socket: %s", socketPath)
	} else {
		logger.Info("listen tcp")
		listner, listnerErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
		logger.Infof("server is listening on port %s:%s", cfg.Listen.BindIP, cfg.Listen.Port)
	}

	if listnerErr != nil {
		logger.Fatal(listnerErr)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	server.Serve(listner)
}
