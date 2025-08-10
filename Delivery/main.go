package main

import (
	"context"
	"fmt"
	"g6/blog-api/Delivery/bootstrap"
	"g6/blog-api/Delivery/routers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// main.go - Entry point for the blog backend server. Handles server startup and graceful shutdown.

func close_server(srv *http.Server) {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutdown Server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

func main() {
	app := bootstrap.App(".env")
	env := app.Env
	db := app.Mongo.Database(env.DB_Name)
	fmt.Println("✅ Acquired database:", env.DB_Name)
	fmt.Println("📄 Using User Collection:", env.UserCollection)
	fmt.Println("📄 Using Refresh Token Collection:", env.RefreshTokenCollection)
	defer app.CloseDBConnection()

	timeout := time.Duration(env.CtxTSeconds) * time.Second

	router := gin.Default()
	routers.Setup(env, timeout, db, router)

	srv := &http.Server{
		Addr:         env.Port,
		Handler:      router.Handler(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start HTTP Server
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	close_server(srv)
}
