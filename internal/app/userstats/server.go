package userstats

import (
	"context"

	"github.com/kulti/otus_ol_int_tests/internal/api/userstatsapi"
)

type server struct {
	store *store
}

func (s *server) GetStats(ctx context.Context, req *userstatsapi.GetStatsRequest) (*userstatsapi.GetStatsResponse, error) {
	stats, err := s.store.GetUserStats(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}

	resp := &userstatsapi.GetStatsResponse{
		Wins:  make([]*userstatsapi.GameInfo, 0, len(stats.Wins)),
		Loses: make([]*userstatsapi.GameInfo, 0, len(stats.Loses)),
		Draws: make([]*userstatsapi.GameInfo, 0, len(stats.Draws)),
	}
	for _, info := range stats.Wins {
		resp.Wins = append(resp.Wins, &userstatsapi.GameInfo{GameId: info.gameID})
	}
	for _, info := range stats.Loses {
		resp.Loses = append(resp.Loses, &userstatsapi.GameInfo{GameId: info.gameID})
	}
	for _, info := range stats.Draws {
		resp.Draws = append(resp.Draws, &userstatsapi.GameInfo{GameId: info.gameID})
	}

	return resp, nil
}
