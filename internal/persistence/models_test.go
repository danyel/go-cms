package persistence

import "testing"

func TestModelsUseMigrationTableNames(t *testing.T) {
	tests := map[string]string{
		"users":    (UserModel{}).TableName(),
		"admins":   (AdminModel{}).TableName(),
		"sessions": (SessionModel{}).TableName(),
	}
	for expected, actual := range tests {
		if actual != expected {
			t.Errorf("table name = %q, want %q", actual, expected)
		}
	}
}
