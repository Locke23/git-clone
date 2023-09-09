package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/locke23/git-clone/index"
	"github.com/spf13/cobra"
)

var Init = &cobra.Command{
	Use:   "init",
	Short: "Start a new lit repository",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := os.Stat(".lit")
		if !os.IsNotExist(err) {
			log.Fatal("Failed to init o lit repository: .lit folder already exists")
		}
		if err := os.Mkdir(".lit", 0755); err != nil {
			log.Fatal("failed to initialize lit repository: %s", err)
		}

		err = func() error {
			if err := os.Mkdir(filepath.Join(".", ".lit", "objects"), 0755); err != nil {
				return err
			}

			file, err := os.OpenFile(filepath.Join(".", ".lit", "index"), os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}
			defer file.Close()

			head, err := os.OpenFile(filepath.Join(".", ".lit", "HEAD"), os.O_CREATE, 0644)
			if err != nil {
				return err
			}
			defer head.Close()

			var idx index.Index
			idx.Objects = make(map[string]index.Entry)
			content, err := json.Marshal(idx)
			if err != nil {
				return err
			}

			_, err = file.Write(content)
			return err
		}()

		if err != nil {
			err2 := os.RemoveAll(".lit")
			if err2 != nil {
				log.Fatalf("failed to initialize lit repository: corrupted .lit folder: %s", err)
			}
			log.Fatalf("failed to initialize lit repository: %s", err)
		}

		pwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("failed to initialize lit repository: %s", err)
		}
		fmt.Printf("Initialized empty lit repository in %s \n", filepath.Join(pwd, ".lit"))
	},
}
