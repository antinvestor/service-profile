package service

import (
	"bitbucket.org/antinvestor/service-profile/profile"
	"context"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	otgorm "github.com/smacker/opentracing-gorm"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
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
	wDb        *gorm.DB
	rDb        *gorm.DB
	Logger     *logrus.Entry
	ServerPort string
}

func (env *Env) SetWriteDb(db *gorm.DB) {
	env.wDb = db
}

func (env *Env) SetReadDb(db *gorm.DB) {
	env.rDb = db
}

func (env *Env) GeWtDb(ctx context.Context) *gorm.DB {
	return otgorm.SetSpanToGorm(ctx, env.wDb)
}

func (env *Env) GetRDb(ctx context.Context) *gorm.DB {
	return otgorm.SetSpanToGorm(ctx, env.rDb)
}

//RunServer Starts a server and waits on it
func RunServer(env *Env) {

	implementation := &ProfileServer{Env: env}

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(AuthInterceptor),
	)

	profile.RegisterProfileServiceServer(srv, implementation)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", env.ServerPort))
	if err != nil {
		env.Logger.Fatalf("Could not start on supplied port %v %v ", env.ServerPort, err)
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {

		env.Logger.Infof("Service running on port : %v", env.ServerPort)

		// start the server
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %s", err)
		}

	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	srv.Stop()
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-env.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	env.Logger.Infof("Service shutting down at : %v", time.Now())
}
