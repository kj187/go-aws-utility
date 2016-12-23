package main

import (
	"fmt"
	"os"
	"github.com/kj187/aws-utility/commands"
	"github.com/fatih/color"
)

func main() {
	fmt.Println(`
    /\ \        / / ____| | |  | | | (_) (_) |
   /  \ \  /\  / / (___   | |  | | |_ _| |_| |_ _   _
  / /\ \ \/  \/ / \___ \  | |  | | __| | | | __| | | |
 / ____ \  /\  /  ____) | | |__| | |_| | | | |_| |_| |
/_/    \_\/  \/  |_____/   \____/ \__|_|_|_|\__|\__, |
© Julian Kleinhans - @kj187                      __/ |
                                                |___/
	`)

	if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
		color.Red("No AWS_ACCESS_KEY_ID env var available!")
		return
	}
	if os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
		color.Red("No AWS_SECRET_ACCESS_KEY env var available!")
		return
	}

	if err := commands.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
