package gameserver

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/notnil/chess"

	"github.com/kulti/otus_ol_int_tests/internal/api/chessapi"
)

func Run() error {
	l, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	store, err := newStore("postgres://postgres:postgres@postgres:5432/game_server_db?sslmode=disable")
	if err != nil {
		return err
	}

	amqpConn, err := connectAmqp()
	if err != nil {
		return err
	}

	chessapi.RegisterChessServer(grpcServer, &chessServer{store, amqpConn})

	return grpcServer.Serve(l)
}

type chessServer struct {
	store    *store
	amqpConn *amqpConn
}

func (s *chessServer) JoinGame(ctx context.Context, req *chessapi.JoinGameRequest) (*chessapi.JoinGameResponse, error) {
	gameID, gameExists, err := s.store.JoinOrCreateGame(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	resp := &chessapi.JoinGameResponse{
		GameId: gameID,
	}
	if gameExists {
		resp.Color = chessapi.Color_Black
	}
	return resp, nil
}

func (s *chessServer) SendMove(ctx context.Context, req *chessapi.SendMoveRequest) (*chessapi.SendMoveResponse, error) {
	moves, err := s.store.GetGameMoves(ctx, req.GetGameId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	moves = append(moves, req.GetMove())

	game := chess.NewGame()
	for _, m := range moves {
		if err := game.MoveStr(m); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
	}

	switch game.Outcome() {
	case chess.BlackWon:
		s.amqpConn.Send(msg{GameID: req.GetGameId(), UserID: req.GetUserId(), Outcome: "win"})
	case chess.Draw:
	case chess.WhiteWon:
	}

	err = s.store.StoreGameMove(ctx, req.GetGameId(), req.GetUserId(), req.GetMove())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &chessapi.SendMoveResponse{}, nil
}

func (s *chessServer) GetMoves(ctx context.Context, req *chessapi.GetMoveRequest) (*chessapi.GetMoveResponse, error) {
	moves, err := s.store.GetGameMoves(ctx, req.GetGameId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &chessapi.GetMoveResponse{Moves: moves}, nil
}
