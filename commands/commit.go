package commands

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/locke23/git-clone/commit"
	"github.com/locke23/git-clone/index"
	"github.com/spf13/cobra"
)

var Commit = &cobra.Command{
	Use:   "commit",
	Short: "create a commit using staging area",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatalf("not found commit mesage: required argument")
		}

		idx, err := index.Build()
		if err != nil {
			log.Fatalf("could not build index: %v", err)
		}

		authorName := getEnvOrDefault("LIT_AUTHOR_NAME", "lukete dev")
		authorEmail := getEnvOrDefault("LIT_AUTHOR_EMAIL", "lukete@dev.com")

		c := commit.Entry{
			AuthorName:  authorName,
			AuthorEmail: authorEmail,
			AuthorDate:  time.Now(),
			Message:     args[0],
			Index:       idx,
		}

		head, err := commit.GetHEAD()
		if err != nil {
			log.Fatalf("failed to get HEAD: %v", err)
		}

		if head.Hash != "" {
			c.ParentHash = head.Hash
		}

		hash, err := commit.Write(c)
		if err != nil {
			log.Fatalf("could not write commit: %v", err)
		}

		fmt.Printf("[%] %s\n", hash[0:6], c.Message)
	},
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}
