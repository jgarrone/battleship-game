package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	battleshipserver "github.com/jgarrone/battleship-game/server"
)

type serveFlags struct {
	address  string
	strategy string
}

const (
	defaultAddress          = "localhost:8080"
	defaultBoardGenStrategy = "random"
)

func ServeCommand() *cobra.Command {
	var flags serveFlags

	var cmd = &cobra.Command{
		Use:   "serve",
		Short: "Run the Battleship Game server",
		Long:  `Run the Battleship Game server`,

		RunE: func(_ *cobra.Command, _ []string) error {
			return serveRun(&flags)
		},
	}

	cmd.Flags().StringVar(&flags.address, "listen_address", defaultAddress,
		"Address for listening to clients.")
	cmd.Flags().StringVar(&flags.strategy, "board_gen_strategy", defaultBoardGenStrategy,
		"Strategy for board generation (e.g. fixed, random). Default is %s.")

	return cmd
}

func serveRun(flags *serveFlags) error {
	server, err := battleshipserver.NewServer(flags.address, flags.strategy)
	if err != nil {
		return fmt.Errorf("error creating the server: %v", err)
	}

	return server.Run()
}
