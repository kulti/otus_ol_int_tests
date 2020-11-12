package gameserver

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
)

type store struct {
	conn *pgxpool.Pool
}

func newStore(connStr string) (*store, error) {
	conn, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return &store{conn: conn}, nil
}

func (s *store) JoinOrCreateGame(ctx context.Context, userID string) (_ string, _ bool, resultErr error) {
	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return "", false, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		rollbackErr := tx.Rollback(ctx)

		// if resultErr == nil && rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed) {
		if resultErr == nil && rollbackErr != nil && rollbackErr.Error() != pgx.ErrTxClosed.Error() {
			resultErr = fmt.Errorf("failed to rollback: %w", rollbackErr)
		}
	}()

	row := tx.QueryRow(ctx, "SELECT id FROM games WHERE black_user_id IS NULL")

	var gameID int
	err = row.Scan(&gameID)
	if err == nil {
		_, err = tx.Exec(ctx, "UPDATE games SET black_user_id = $2 WHERE id = $1", gameID, userID)
		if err != nil {
			return "", false, fmt.Errorf("failed to join game %q: %w", gameID, err)
		}
		return strconv.Itoa(gameID), true, tx.Commit(ctx)
	}

	// if !errors.Is(err, pgx.ErrNoRows) {
	if err.Error() != pgx.ErrNoRows.Error() {
		return "", false, fmt.Errorf("failed to find game: %w", err)
	}

	row = tx.QueryRow(ctx, "INSERT INTO games (white_user_id) VALUES ($1) RETURNING id", userID)
	err = row.Scan(&gameID)
	if err != nil {
		return "", false, fmt.Errorf("failed to create game: %w", err)
	}

	return strconv.Itoa(gameID), false, tx.Commit(ctx)
}

func (s *store) StoreGameMove(ctx context.Context, gameID, userID, move string) error {
	_, err := s.conn.Exec(ctx, "INSERT INTO moves (game_id, user_id, move) VALUES ($1, $2, $3)", gameID, userID, move)
	return err
}

func (s *store) GetGameMoves(ctx context.Context, gameID string) ([]string, error) {
	rows, _ := s.conn.Query(ctx, "SELECT move FROM moves WHERE game_id = $1 ORDER BY id", gameID)

	var moves []string
	err := rows.Err()
	for err == nil && rows.Next() {
		var move string
		err = rows.Scan(&move)
		if err == nil {
			moves = append(moves, move)
		}
	}
	rows.Close()

	if err != nil {
		return nil, err
	}

	return moves, nil
}
