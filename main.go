package main

import (
	"context"
	"fmt"
	"github.com/projectriff/streaming-http-adapter/pkg/proxy"
	"log"
	"os"
	"os/exec"
	"os/signal"
)

func main() {
	grpcAddress := os.Getenv("GRPC_ADDRESS")
	if grpcAddress == "" {
		grpcAddress = ":8081"
	}
	httpAddress := os.Getenv("HTTP_ADDRESS")
	if httpAddress == "" {
		httpAddress = ":8080"
	}
	if len(os.Args) < 2 {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %s invoker-command [invoker-args]...\n", os.Args[0])
		os.Exit(1)
	}

	proxy, err := proxy.NewProxy(grpcAddress, httpAddress)
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
	command.Env = os.Environ()

	if err := command.Start() ; err != nil {
		panic(err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)
	<-stop

	command.Wait()

	proxyShutdown := make(chan struct{})
	go func() {
		if err := proxy.Shutdown(context.Background()) ; err != nil {
			log.Fatalf("error shuting down proxy server %v", err)
		}
		close(proxyShutdown)
	}()

	<-proxyShutdown
}
