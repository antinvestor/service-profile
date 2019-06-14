package service

import (
	"github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	otgorm "github.com/smacker/opentracing-gorm"
	"time"
	"net/http"
	"os"
	"os/signal"
	"context"
	"fmt"
)

// Error represents a handler error. It provides methods for a HTTP status
// code and embeds the built-in error interface.
type Error interface {
	error
	Status() int
}

// StatusError represents an error with an associated HTTP status code.
type StatusError struct {
	Code int
	Err  error
}

// Allows StatusError to satisfy the error interface.
func (se StatusError) Error() string {
	return se.Err.Error()
}

// Returns our HTTP status code.
func (se StatusError) Status() int {
	return se.Code
}


// Env Context object supplied around the applications lifetime
type Env struct {
	db              *gorm.DB
	Logger          *logrus.Entry
	ServerPort	string

}

func (env *Env) SetDb(db *gorm.DB) {
	env.db = db
}

func (env *Env) GetDb(ctx context.Context) *gorm.DB{
	return otgorm.SetSpanToGorm(ctx, env.db)
}

//RunServer Starts a server and waits on it
func RunServer(env *Env) {

	waitDuration := time.Second * 15
	//router := NewRouterV1(env)

	srv := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%s", env.ServerPort ),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		//Handler:      handlers.RecoveryHandler()(router), // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {

		env.Logger.Infof("Service running on port : %v", env.ServerPort)

		if err := srv.ListenAndServe(); err != nil {
			env.Logger.Fatalf("Service stopping due to error : %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c


	// Create a deadline to wait for.
	env2, cancel := context.WithTimeout(context.Background(), waitDuration)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(env2)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-env.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	env.Logger.Infof("Service shutting down at : %v", time.Now())
}