package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	desc "github.com/uxsnap/auth/pkg/auth_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = 50051

var pool *pgxpool.Pool

type server struct {
	desc.UnimplementedAuthV1Server
}

func (c *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	Id := req.GetId()

	if Id < 0 {
		return nil, errors.New("unsupported id")
	}

	userQuery := sq.Select("*").From("auth").Where(sq.Eq{"id": Id})

	query, args, err := userQuery.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	var resp desc.GetResponse

	rows, err := pool.Query(ctx, query, args...)

	if err != nil {
		log.Fatalf("failed to query user: %v", err)
	}

	err = rows.Scan(&resp)

	if err != nil {
		log.Fatalf("failed to scan user: %v", err)
	}

	return &resp, nil
}

func initEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func initDb() {
	ctx := context.Background()
	dsn := os.Getenv("PG_DSN")

	pool, err := pgxpool.Connect(ctx, dsn)

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	defer pool.Close()
}

func main() {
	initEnv()
	initDb()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))

	if err != nil {
		log.Fatal("Cannot create tcp connection!")
		return
	}

	grpcS := grpc.NewServer()
	reflection.Register(grpcS)

	if err != nil {
		log.Fatal("Cannot create grpc connection!")
		return
	}

	desc.RegisterAuthV1Server(grpcS, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = grpcS.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
