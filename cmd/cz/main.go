package main

import (
	"cz/internal/committer"
	"cz/internal/git"
	"fmt"
	"github.com/boostgo/lite/log"
	"os"
)

func main() {
	log.PrettyLog()

	if len(os.Args) < 2 || os.Args[1] != "commit" {
		fmt.Println("Usage: cz commit")
		os.Exit(1)
	}

	// generate commit message
	commitMsg, err := committer.GenerateMessage()
	if err != nil {
		log.Error().Err(err).Msg("Generate commit message")
		os.Exit(1)
	}

	// execute git commit
	if err = git.Commit(commitMsg); err != nil {
		log.Error().Err(err).Msg("Commit")
		os.Exit(1)
	}
}
