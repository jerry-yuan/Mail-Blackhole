package storage

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/jerry-yuan/mail-blackhole/data"
	"github.com/sirupsen/logrus"
)

// Maildir is a maildir storage backend
type Maildir struct {
	Path string
}

// CreateMaildir creates a new maildir storage backend
func CreateMaildir(path string) *Maildir {
	if len(path) == 0 {
		dir, err := os.MkdirTemp("", "mailhog")
		if err != nil {
			panic(err)
		}
		path = dir
	}
	if _, err := os.Stat(path); err != nil {
		err := os.MkdirAll(path, 0770)
		if err != nil {
			panic(err)
		}
	}
	logrus.Println("Maildir path is", path)
	return &Maildir{
		Path: path,
	}
}

// Store stores a message and returns its storage ID
func (maildir *Maildir) Store(m *data.Message) (string, error) {
	b, err := io.ReadAll(m.Raw.Bytes())
	if err != nil {
		return "", err
	}
	err = os.WriteFile(filepath.Join(maildir.Path, string(m.ID)), b, 0660)
	return string(m.ID), err
}

// Count returns the number of stored messages
func (maildir *Maildir) Count() int {
	// FIXME may be wrong, ../. ?
	// and handle error?
	dir, err := os.Open(maildir.Path)
	if err != nil {
		panic(err)
	}
	defer dir.Close()
	n, _ := dir.Readdirnames(0)
	return len(n)
}

// Search finds messages matching the query
func (maildir *Maildir) Search(kind, query string, start, limit int) (*data.Messages, int, error) {
	query = strings.ToLower(query)
	var filteredMessages = make([]data.Message, 0)

	var matched int

	err := filepath.Walk(maildir.Path, func(path string, info os.FileInfo, _ error) error {
		if limit > 0 && len(filteredMessages) >= limit {
			return errors.New("reached limit")
		}

		if info.IsDir() {
			return nil
		}

		msg, err1 := maildir.Load(info.Name())
		if err1 != nil {
			logrus.Println(err1)
			return nil
		}

		switch kind {
		case "to":
			for _, t := range msg.To {
				if strings.Contains(strings.ToLower(t.Mailbox+"@"+t.Domain), query) {
					if start > matched {
						matched++
						break
					}
					filteredMessages = append(filteredMessages, *msg)
					break
				}
			}
		case "from":
			if strings.Contains(strings.ToLower(msg.From.Mailbox+"@"+msg.From.Domain), query) {
				if start > matched {
					matched++
					break
				}
				filteredMessages = append(filteredMessages, *msg)
			}
		case "containing":
			if strings.Contains(strings.ToLower(msg.Raw.Data), query) {
				if start > matched {
					matched++
					break
				}
				filteredMessages = append(filteredMessages, *msg)
			}
		}

		return nil
	})

	if err != nil {
		logrus.Println(err)
	}

	msgs := data.Messages(filteredMessages)
	return &msgs, len(filteredMessages), nil
}

// List lists stored messages by index
func (maildir *Maildir) List(start, limit int) (*data.Messages, error) {
	logrus.Println("Listing messages in", maildir.Path)
	messages := make([]data.Message, 0)

	dir, err := os.Open(maildir.Path)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	n, err := dir.Readdir(0)
	if err != nil {
		return nil, err
	}

	for _, fileinfo := range n {
		if start > 0 {
			start -= 1
			continue
		}
		b, err := os.ReadFile(filepath.Join(maildir.Path, fileinfo.Name()))
		if err != nil {
			return nil, err
		}
		msg := data.FromBytes(b)
		// FIXME domain
		m := *msg.Parse("mailhog.example")
		m.ID = data.MessageID(fileinfo.Name())
		m.Created = fileinfo.ModTime()
		messages = append(messages, m)
		if len(messages) >= limit {
			break
		}
	}

	logrus.Printf("Found %d messages", len(messages))
	msgs := data.Messages(messages)
	return &msgs, nil
}

// DeleteOne deletes an individual message by storage ID
func (maildir *Maildir) DeleteOne(id string) error {
	return os.Remove(filepath.Join(maildir.Path, id))
}

// DeleteAll deletes all in memory messages
func (maildir *Maildir) DeleteAll() error {
	logrus.Infof("Deleting all messages in %s", maildir.Path)

	dir, err := os.Open(maildir.Path)
	if err != nil {
		return err
	}
	defer dir.Close()

	entries, err := dir.ReadDir(0)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		err = os.Remove(filepath.Join(maildir.Path, entry.Name()))
		if err != nil {
			return err
		}
	}

	return nil
}

// Load returns an individual message by storage ID
func (maildir *Maildir) Load(id string) (*data.Message, error) {
	b, err := os.ReadFile(filepath.Join(maildir.Path, id))
	if err != nil {
		return nil, err
	}
	// FIXME domain
	m := data.FromBytes(b).Parse("mailhog.example")
	m.ID = data.MessageID(id)
	return m, nil
}
