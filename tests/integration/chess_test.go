// +build integration

package integration_test

import (
	"context"
	"math/rand"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"

	"github.com/kulti/otus_ol_int_tests/internal/api/chessapi"
	"github.com/kulti/otus_ol_int_tests/internal/api/userstatsapi"
)

type ChessSuite struct {
	suite.Suite
	ctx         context.Context
	gameConn    *grpc.ClientConn
	chessClient chessapi.ChessClient
	statsClient userstatsapi.UserStatsClient
	userWhite   string
	userBlack   string
}

func (s *ChessSuite) SetupSuite() {
	gameHost := os.Getenv("GAME_SERVER_HOST")
	if gameHost == "" {
		gameHost = "127.0.0.1:9001"
	}

	gameConn, err := grpc.Dial(gameHost, grpc.WithInsecure())
	s.Require().NoError(err)

	statHost := os.Getenv("STAT_SERVER_HOST")
	if statHost == "" {
		statHost = "127.0.0.1:9002"
	}

	statConn, err := grpc.Dial(statHost, grpc.WithInsecure())
	s.Require().NoError(err)

	s.ctx = context.Background()
	s.chessClient = chessapi.NewChessClient(gameConn)
	s.statsClient = userstatsapi.NewUserStatsClient(statConn)
}

func (s *ChessSuite) SetupTest() {
	var seed int64 = time.Now().UnixNano()
	rand.Seed(seed)
	s.T().Log("seed:", seed)

	s.userWhite = faker.Word()
	s.userBlack = faker.Word()
	for s.userBlack == s.userWhite {
		s.userBlack = faker.Word()
	}
}

func (s *ChessSuite) TearDownTest() {
}

func (s *ChessSuite) TearDownSuite() {

}

func (s *ChessSuite) TestJoinGame() {
	gameID := s.createGame(s.userWhite)
	s.joinGame(s.userBlack, gameID)
}

func (s *ChessSuite) TestSendSeveralMoves() {
	gameID := s.createGame(s.userWhite)
	s.joinGame(s.userBlack, gameID)

	s.sendMoves(gameID, [...]string{s.userWhite, s.userBlack}, []string{"e4", "e5"})
	s.checkMoves(gameID, "e4", "e5")

	s.sendMoves(gameID, [...]string{s.userWhite, s.userBlack}, []string{"Nf3", "Nf6"})
	s.checkMoves(gameID, "e4", "e5", "Nf3", "Nf6")
}

func (s *ChessSuite) TestSendMoveAfterMate() {
	gameID := s.createGame(s.userWhite)
	s.joinGame(s.userBlack, gameID)

	s.sendMoves(gameID, [...]string{s.userWhite, s.userBlack}, []string{"f3", "e5", "g4", "Qh4"})
	s.checkMoves(gameID, "f3", "e5", "g4", "Qh4")

	s.Require().Error(s.sendMove(gameID, s.userWhite, "a3"))
}

func (s *ChessSuite) TestMate() {
	gameID := s.createGame(s.userWhite)
	s.joinGame(s.userBlack, gameID)

	userBlackStatsOld := s.getUserStats(s.userBlack)

	s.sendMoves(gameID, [...]string{s.userWhite, s.userBlack}, []string{"f3", "e5", "g4", "Qh4"})
	s.checkMoves(gameID, "f3", "e5", "g4", "Qh4")

	userBlackStatsOld.Wins = append(userBlackStatsOld.Wins, &userstatsapi.GameInfo{GameId: gameID})
	s.Require().Eventually(func() bool {
		userBlackStatsNew := s.getUserStats(s.userBlack)
		return reflect.DeepEqual(userBlackStatsOld, userBlackStatsNew)
	}, 5*time.Second, time.Millisecond)
}

func (s *ChessSuite) createGame(userID string) (gameID string) {
	req := &chessapi.JoinGameRequest{
		UserId: userID,
	}
	resp, err := s.chessClient.JoinGame(s.ctx, req)
	s.Require().NoError(err)

	if resp.GetColor() != chessapi.Color_White {
		resp, err := s.chessClient.JoinGame(s.ctx, req)
		s.Require().NoError(err)
		s.Require().Equal(chessapi.Color_White, resp.GetColor())
	}

	return resp.GetGameId()
}

func (s *ChessSuite) joinGame(userID, gameID string) {
	req := &chessapi.JoinGameRequest{
		UserId: userID,
	}
	resp, err := s.chessClient.JoinGame(s.ctx, req)
	s.Require().NoError(err)
	s.Require().Equal(chessapi.Color_Black, resp.GetColor())
	s.Require().Equal(gameID, resp.GetGameId())
}

func (s *ChessSuite) sendMoves(gameID string, userIDs [2]string, moves []string) {
	for i, m := range moves {
		s.Require().NoError(s.sendMove(gameID, userIDs[i%2], m))
	}
}

func (s *ChessSuite) sendMove(gameID, userID, move string) error {
	req := &chessapi.SendMoveRequest{
		UserId: userID,
		Move:   move,
		GameId: gameID,
	}
	_, err := s.chessClient.SendMove(s.ctx, req)
	return err
}

func (s *ChessSuite) checkMoves(gameID string, moves ...string) {
	req := &chessapi.GetMoveRequest{
		GameId: gameID,
	}
	resp, err := s.chessClient.GetMoves(s.ctx, req)
	s.Require().NoError(err)

	s.Require().Equal(moves, resp.GetMoves())
}

func (s *ChessSuite) getUserStats(userID string) *userstatsapi.GetStatsResponse {
	req := &userstatsapi.GetStatsRequest{
		UserId: userID,
	}
	resp, err := s.statsClient.GetStats(s.ctx, req)
	s.Require().NoError(err)

	return resp
}

func TestChessSuite(t *testing.T) {
	suite.Run(t, new(ChessSuite))
}
