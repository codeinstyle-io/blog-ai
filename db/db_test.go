package db

import (
	"os"
	"testing"

	"codeinstyle.io/captain/models"
	"github.com/stretchr/testify/assert"
)

func TestInsertTestData(t *testing.T) {
	db := SetupTestDB()

	// Create data directory if it doesn't exist
	if err := os.MkdirAll("data", 0755); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		setup   func() error
		wantErr bool
	}{
		{
			name: "Insert test data successfully",
			setup: func() error {
				// Create test data file
				data := `{
					"tags": ["test1", "test2"],
					"posts": [
						{
							"title": "Test Post",
							"slug": "test-post",
							"content": "Content",
							"publishedAt": "-1d",
							"visible": true,
							"excerpt": "Excerpt"
						}
					]
				}`
				return os.WriteFile("data/test_posts.json", []byte(data), 0644)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.setup(); err != nil {
				t.Fatal(err)
			}
			err := InsertTestData(db)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify data was inserted
				var tagCount int64
				db.Model(&models.Tag{}).Count(&tagCount)
				assert.Equal(t, int64(2), tagCount)

				var postCount int64
				db.Model(&models.Post{}).Count(&postCount)
				assert.Equal(t, int64(1), postCount)
			}
		})
	}

	// Cleanup
	defer os.RemoveAll("data")
}
