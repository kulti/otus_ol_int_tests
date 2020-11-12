package userstats

import (
	"net"
	"sync"

	"google.golang.org/grpc"

	"github.com/kulti/otus_ol_int_tests/internal/api/userstatsapi"
)

type App struct {
	grpcListener net.Listener
	server       *grpc.Server
	stopWg       sync.WaitGroup
	amqpConn     *amqpConn
}

func New() (*App, error) {
	l, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil {
		return nil, err
	}

	grpcServer := grpc.NewServer()
	store, err := newStore("postgres://postgres:postgres@postgres:5432/user_stats_db?sslmode=disable")
	if err != nil {
		return nil, err
	}

	userstatsapi.RegisterUserStatsServer(grpcServer, &server{store})

	amqpConn, err := connectAmqp(store)
	if err != nil {
		return nil, err
	}

	return &App{
		grpcListener: l,
		server:       grpcServer,
		amqpConn:     amqpConn,
	}, nil
}

func (a *App) Start() {
	a.stopWg.Add(2)
	go func() {
		defer a.stopWg.Done()
		a.server.Serve(a.grpcListener)
	}()
	go func() {
		defer a.stopWg.Done()
		a.amqpConn.Run()
	}()
}

func (a *App) Stop() {
	a.server.Stop()
	a.amqpConn.Close()
	a.stopWg.Wait()
}
