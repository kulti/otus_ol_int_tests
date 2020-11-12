package chesscli

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"

	"github.com/kulti/otus_ol_int_tests/internal/api/chessapi"
	"github.com/kulti/otus_ol_int_tests/internal/api/userstatsapi"
)

func Execute() error {
	app := cli.App{}
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name: "host",
		},
	}
	app.Commands = []*cli.Command{
		{
			Name: "game",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name: "user",
				},
			},
			Subcommands: []*cli.Command{
				{
					Name: "join",
					Action: func(cliCtx *cli.Context) error {
						conn, err := grpc.Dial(cliCtx.String("host"), grpc.WithInsecure())
						if err != nil {
							return err
						}

						chessClient := chessapi.NewChessClient(conn)
						userID := cliCtx.String("user")
						resp, err := chessClient.JoinGame(cliCtx.Context, &chessapi.JoinGameRequest{UserId: userID})
						if err != nil {
							return err
						}

						fmt.Printf("You are %q\n", resp.Color)

						return runNewGame(cliCtx.Context, resp.GetGameId(), userID, chessClient)
					},
				},
			},
		},
		{
			Name: "user",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name: "user",
				},
			},
			Subcommands: []*cli.Command{
				{
					Name: "stats",
					Action: func(cliCtx *cli.Context) error {
						conn, err := grpc.Dial(cliCtx.String("host"), grpc.WithInsecure())
						if err != nil {
							return err
						}

						userStatsClient := userstatsapi.NewUserStatsClient(conn)
						userID := cliCtx.String("user")
						resp, err := userStatsClient.GetStats(cliCtx.Context, &userstatsapi.GetStatsRequest{UserId: userID})
						if err != nil {
							return err
						}

						fmt.Printf("User stats: %s\n", resp.String())
						return nil
					},
				},
			},
		},
	}
	return app.Run(os.Args)
}
