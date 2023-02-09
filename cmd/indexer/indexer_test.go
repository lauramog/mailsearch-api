package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseEmails(t *testing.T) {
	mailDir := t.TempDir()

	userInboxDirpath := filepath.Join(mailDir, "martin-t", "inbox")
	if err := os.MkdirAll(userInboxDirpath, 0750); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile("mail_sample.txt")
	if err != nil {
		t.Fatal(err)
	}

	mailFilepath := filepath.Join(userInboxDirpath, "mail.txt")
	if err = os.WriteFile(mailFilepath, content, 0750); err != nil {
		t.Fatal(err)
	}

	emails, err := parseEmails(mailDir)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := emails[0][0]["Subject"], "FW: Gas P&L by day"; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}
