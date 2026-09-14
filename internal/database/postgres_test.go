package database

import (
	"testing"
)

func TestSanitizeSQL(t *testing.T) {
	tests := []struct {
		name      string
		query     string
		shouldErr bool
	}{
		{
			name:      "Valid SELECT query",
			query:     "SELECT * FROM peserta_didik WHERE peserta_didik_id = '123'",
			shouldErr: false,
		},
		{
			name:      "Forbidden INSERT",
			query:     "INSERT INTO peserta_didik (nama) VALUES ('Hacker')",
			shouldErr: true,
		},
		{
			name:      "Forbidden UPDATE",
			query:     "UPDATE peserta_didik SET nama = 'Hacker'",
			shouldErr: true,
		},
		{
			name:      "Forbidden DELETE",
			query:     "DELETE FROM peserta_didik WHERE id = 1",
			shouldErr: true,
		},
		{
			name:      "Forbidden DROP TABLE",
			query:     "DROP TABLE peserta_didik CASCADE",
			shouldErr: true,
		},
		{
			name:      "Forbidden TRUNCATE",
			query:     "TRUNCATE TABLE nilai_rapor",
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SanitizeSQL(tt.query)
			if tt.shouldErr && err == nil {
				t.Errorf("expected error for query %q, but got nil", tt.query)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("unexpected error for query %q: %v", tt.query, err)
			}
		})
	}
}

func TestBuildConnString(t *testing.T) {
	connStr := buildConnString("127.0.0.1", 5432, "postgres", "secret", "dapodik_dasmen", "disable")
	if connStr == "" {
		t.Fatal("expected non-empty connection string")
	}
	expectedPrefix := "postgres://postgres:secret@127.0.0.1:5432/dapodik_dasmen"
	if connStr[:len(expectedPrefix)] != expectedPrefix {
		t.Errorf("expected connection string to start with %s, got %s", expectedPrefix, connStr)
	}
}
