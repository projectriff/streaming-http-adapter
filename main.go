package main

import (
	"context"
	"fmt"
	"github.com/projectriff/streaming-http-adapter/pkg/proxy"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func main() {
	// Lookup the gRPC address where our child process expects to run.
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "8081"
	}

	// Per the http-invoker contract, listen for http traffic on a port defined by PORT
	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}
	if len(os.Args) < 2 {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %s invoker-command [invoker-args]...\n", os.Args[0])
		os.Exit(1)
	}

	proxy, err := proxy.NewProxy(fmt.Sprintf(":%d", grpcPort), fmt.Sprintf(":%d", httpPort))
	if err != nil {
		panic(err)
	}
	go func() {
		if err := proxy.Run(); err != nil {
			log.Fatalf("error running proxy %v", err)
		}
	}()

	command := exec.Command(os.Args[1], os.Args[2:]...)
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	command.Stderr = os.Stderr
	// The following makes sure that our child process sees the GRPC_PORT variable too.
	// It should not care about the PORT variable
	command.Env = os.Environ()

	if err := command.Start(); err != nil {
		panic(err)
	}

	done := make(chan struct{}, 2)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	sig := <-stop
	_ = sig

	if err := command.Process.Signal(sig); err != nil {
		panic(err)
	}

	go func() {
		_ = command.Wait()
		done <- struct{}{}
	}()

	go func() {
		if err := proxy.Shutdown(context.Background()); err != nil {
			log.Fatalf("error shuting down proxy server %v", err)
		}
		done <- struct{}{}
	}()

	<-done
	<-done
}
