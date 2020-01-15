package service

import (
	grpc_health_v1 "antinvestor.com/service/profile/grpc/health"
	"antinvestor.com/service/profile/service/handlers"
	"antinvestor.com/service/profile/utils"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"antinvestor.com/service/profile/grpc/profile"
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

//RunServer Starts a server and waits on it
func RunServer(env *utils.Env) {

	implementation := &handlers.ProfileServer{Env: env}

	serverPort := utils.GetEnv(utils.EnvServerPort, "7005")

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(AuthInterceptor),
	)

	profile.RegisterProfileServiceServer(srv, implementation)
	grpc_health_v1.RegisterHealthServer(srv, implementation)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", serverPort))
	if err != nil {
		env.Logger.Fatalf("Could not start on supplied port %v %v ", serverPort, err)
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {

		env.Logger.Infof("Service running on port : %v", serverPort)

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
