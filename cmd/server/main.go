package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/hugoaguirre/product-service/internal/adapters/db"
	grpc_adapter "github.com/hugoaguirre/product-service/internal/adapters/grpc"
	"github.com/hugoaguirre/product-service/internal/adapters/rest"
	"github.com/hugoaguirre/product-service/internal/service"
	"github.com/hugoaguirre/product-service/pkg/productapi"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	_ "modernc.org/sqlite"
)

func main() {
	fmt.Printf("Hello Hex!")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	sqliteDB, err := sql.Open("sqlite", "products.db")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer sqliteDB.Close()

	if err := db.RunMigrations(sqliteDB); err != nil {
		log.Fatalf("unable to run migrations: %v", err)
	}

	repo := db.NewSQLiteRepository(sqliteDB)
	svc := service.NewCatalogService(repo)

	g, ctx := errgroup.WithContext(ctx)

	// start gRPC server
	g.Go(func() error {
		listener, err := net.Listen("tcp", ":50051")
		if err != nil {
			return err
		}
		grpcServer := grpc.NewServer()
		productapi.RegisterProductServiceServer(grpcServer, grpc_adapter.New(svc))
		fmt.Println("gRPC server listening on :50051")

		go func() {
			<-ctx.Done()
			fmt.Printf("shutting down gRPC server...")
			grpcServer.GracefulStop()
		}()

		return grpcServer.Serve(listener)
	})

	// start REST server
	g.Go(func() error {
		mux := http.NewServeMux()
		restAdapter := rest.NewHandler(svc)
		mux.HandleFunc("/products/{id}", restAdapter.GetProduct())

		httpServer := &http.Server{
			Addr:    ":8080",
			Handler: mux,
		}
		fmt.Println("HTTP server listening on :8080")

		go func() {
			<-ctx.Done()
			fmt.Println("shutting down HTTP server...")
			if err := httpServer.Shutdown(context.Background()); err != nil {
				fmt.Printf("error while shutting down HTTP server: %v\n", err)
			}
		}()

		return httpServer.ListenAndServe()
	})

	// wait for servers to finish or context to cancel
	if err := g.Wait(); err != nil {
		fmt.Printf("Exit reason: %v\n", err)
	}
}
