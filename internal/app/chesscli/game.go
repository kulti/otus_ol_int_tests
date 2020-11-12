package chesscli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/notnil/chess"

	"github.com/kulti/otus_ol_int_tests/internal/api/chessapi"
)

func runNewGame(ctx context.Context, gameID, userID string, chessClient chessapi.ChessClient) error {
	moveCh := make(chan string)
	go handleUserInput(moveCh)

	getMovesReq := &chessapi.GetMoveRequest{GameId: gameID}
	sendMoveReq := &chessapi.SendMoveRequest{GameId: gameID, UserId: userID}

	game := chess.NewGame()
	fmt.Println(game.Position().Board().Draw())

	localMoves := []string{}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		select {
		case <-ctx.Done():
			return nil
		case move := <-moveCh:
			sendMoveReq.Move = move
			_, err := chessClient.SendMove(ctx, sendMoveReq)
			if err != nil {
				return err
			}
		case <-ticker.C:
			resp, err := chessClient.GetMoves(ctx, getMovesReq)
			if err != nil {
				return err
			}

			moves := resp.GetMoves()
			for i := len(localMoves); i < len(moves); i++ {
				if err := game.MoveStr(moves[i]); err != nil {
					return err
				}

				localMoves = append(localMoves, moves[i])
				fmt.Println(game.Position().Board().Draw())
			}
		}
	}
}

func handleUserInput(moveCh chan string) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		moveCh <- scanner.Text()
	}
}
