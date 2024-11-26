package db

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsertTestData(t *testing.T) {
	db := SetupTestDB()

	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Insert test data successfully",
			setup: func() {
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
				os.WriteFile("data/test_posts.json", []byte(data), 0644)
			},
			wantErr: false,
		},
		{
			name: "Missing test data file",
			setup: func() {
				os.Remove("data/test_posts.json")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := InsertTestData(db)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify data was inserted
				var tagCount int64
				db.Model(&Tag{}).Count(&tagCount)
				assert.Equal(t, int64(2), tagCount)

				var postCount int64
				db.Model(&Post{}).Count(&postCount)
				assert.Equal(t, int64(1), postCount)
			}
		})
	}
}
