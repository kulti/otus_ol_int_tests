package userstats

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type store struct {
	conn *pgxpool.Pool
}

type userStats struct {
	Wins  []gameInfo
	Loses []gameInfo
	Draws []gameInfo
}

type gameInfo struct {
	gameID string
}

func newStore(connStr string) (*store, error) {
	conn, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return &store{conn: conn}, nil
}

func (s *store) SaveUserStats(ctx context.Context, userID, gameID, outcome string) error {
	_, err := s.conn.Exec(ctx, "INSERT INTO user_game_stats(user_id, game_id, outcome) VALUES($1, $2, $3)", userID, gameID, outcome)
	return err
}

func (s *store) GetUserStats(ctx context.Context, userID string) (userStats, error) {
	rows, _ := s.conn.Query(ctx, "SELECT game_id, outcome FROM user_game_stats WHERE user_id = $1", userID)

	stats := userStats{}
	err := rows.Err()
	for err == nil && rows.Next() {
		var gameID, outcome string
		err = rows.Scan(&gameID, &outcome)
		if err == nil {
			switch outcome {
			case "win":
				stats.Wins = append(stats.Wins, gameInfo{gameID: gameID})
			case "lose":
				stats.Loses = append(stats.Loses, gameInfo{gameID: gameID})
			case "draw":
				stats.Draws = append(stats.Draws, gameInfo{gameID: gameID})
			}
		}
	}
	rows.Close()

	if err != nil {
		return userStats{}, err
	}

	return stats, nil
}
