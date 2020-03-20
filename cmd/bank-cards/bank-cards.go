package main

import (
	"bank-cards/cmd/bank-cards/app"
	"bank-cards/pkg/core/cards"
	"context"
	"flag"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jafarsirojov/mux/pkg/mux"
	"log"
	"net"
	"net/http"
)

var (
	host = flag.String("host", "", "Server host")
	port = flag.String("port", "", "Server port")
	dsn  = flag.String("dsn", "", "Postgres DSN")
)

//-host 0.0.0.0 -port 9009 -dsn postgres://user:pass@localhost:5300/app
const ENV_PORT = "PORT"
const ENV_DSN = "DATABASE_URL"
const ENV_HOST = "HOST"

func main() {
	flag.Parse()
	envPort, ok := FromFlagOrEnv(*port, ENV_PORT)
	if !ok {
		log.Println("can't port")
		return
	}
	envDsn, ok := FromFlagOrEnv(*dsn, ENV_DSN)
	if !ok {
		log.Println("can't dsn")
		return
	}
	envHost, ok := FromFlagOrEnv(*host, ENV_HOST)
	if !ok {
		log.Println("can't host")
		return
	}
	addr := net.JoinHostPort(envHost, envPort)
	log.Println("starting server!")
	log.Printf("host = %s, port = %s\n", envHost, envPort)
	pool, err := pgxpool.Connect(
		context.Background(),
		envDsn,
	)
	if err != nil {
		panic(err)
	}
	usersSvc := cards.NewService(pool)
	usersSvc.Start()
	exactMux := mux.NewExactMux()
	server := app.NewMainServer(exactMux, usersSvc)
	exactMux.GET("/api/cards",
		//todo: list my cards
		server.HandleGetAllCards,
		jwtMiddleware,
		requestIdier,
		logger,
	)
	exactMux.GET("/api/cards/{id}",
		server.HandleGetCardById,
		jwtMiddleware,
		requestIdier,
		logger,
	)
	//exactMux.GET("/api/cards/ownerid/{id}",
	//	server.HandleGetCardsByOwnerId,
	//	jwtMiddleware,								delete
	//	requestIdier,
	//	logger,
	//)
	exactMux.POST("/api/cards",
		server.HandlePostCard,
		jwtMiddleware,
		requestIdier,
		logger,
	)
	exactMux.POST("/api/cards/{id}/blocked",
		server.HandleBlockById,
		jwtMiddleware,
		requestIdier,
		logger,
	)
	exactMux.GET("/api/cards/unblocked/{id}",
		server.HandleUnBlockedById,
		jwtMiddleware,
		requestIdier,
		logger,
	)
	exactMux.POST("/api/cards/transmoney/{id}",
		server.HandleTransferMoneyCardToCard,
		jwtMiddleware,
		requestIdier,
		logger,
	)
	panic(http.ListenAndServe(addr, server))
}
