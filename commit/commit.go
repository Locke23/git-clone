package commit

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/locke23/git-clone/hash"
	"github.com/locke23/git-clone/index"
)

type Entry struct {
	AuthorName  string      `json:"author_name"`
	AuthorEmail string      `json:"author_email"`
	AuthorDate  time.Time   `json:"author_date"`
	Message     string      `json:"message"`
	Index       index.Index `json:"index"`
	Hash        string      `json:"hash"`
	ParentHash  string      `json:"parent_hash"`
	Parent      *Entry      `json:"-"`
}

func GetByHash(h string) (Entry, error) {
	if h == "" {
		return Entry{}, nil
	}

	folder := string(h[0:2])
	fileName := string(h[2:])

	file, err := os.OpenFile(filepath.Join(".", ".lit", "objects", folder, fileName), os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return Entry{}, fmt.Errorf("not possible to open commit file: %w", err)
	}

	reader := bufio.NewReader(file)
	contents, err := io.ReadAll(reader)
	if err != nil {
		return Entry{}, fmt.Errorf("not possible to read commit file: %w", err)

	}

	if len(contents) == 0 {
		return Entry{}, nil
	}

	var c Entry
	if err = json.Unmarshal(contents, &c); err != nil {
		return Entry{}, fmt.Errorf("not possible to unmarshal commit file: %w", err)
	}

	return c, nil
}

func Write(c Entry) (string, error) {
	idx := c.Index
	hashes := make([]string, len(idx.Objects))
	var i int
	for _, entry := range idx.Objects {
		hashes[i] = entry.Hash
		i++
	}
	sort.Strings(hashes)
	allHashes := strings.Join(hashes, ".")
	c.Hash = hash.New([]byte(allHashes))

	head, err := GetHEAD()
	if err != nil {
		return "", err
	}
	if head.Hash == c.Hash {
		return "", fmt.Errorf("nothing to commit")
	}

	folder := c.Hash[0:2]
	fileName := c.Hash[2:]
	path := filepath.Join(".", ".lit", "objects", folder)
	if err := os.Mkdir(path, 0755); err != nil && !os.IsExist(err) {
		return "", err
	}

	file, err := os.OpenFile(filepath.Join(path, fileName), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}

	contents, err := json.Marshal(c)
	if err != nil {
		return "", err
	}

	_, err = file.Write(contents)
	if err != nil {
		return "", err
	}

	headFile, err := os.OpenFile(filepath.Join(".", ".lit", "HEAD"), os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}

	defer headFile.Close()
	_, err = headFile.Write([]byte(c.Hash))

	return c.Hash, err
}

func GetHEAD() (Entry, error) {
	headFile, err := os.OpenFile(filepath.Join(".", ".lit", "HEAD"), os.O_RDONLY, 0644)

	if err != nil {
		log.Fatalf("not possible to list commits: %s", err)
	}

	defer headFile.Close()

	headerCommit := make([]byte, 40)

	if _, err = headFile.Read(headerCommit); err != nil {
		if err == io.EOF {
			return Entry{}, nil
		}
		return Entry{}, fmt.Errorf("not possible to get HEAD: %w", err)
	}

	return GetByHash(string(headerCommit))

}
