package core

import (
	"fmt"
	"os"
	"testing"
)

func TestNewApi(t *testing.T) {
	apiClient := NewApi(os.Getenv("GITHUB_USER_NAME"), os.Getenv("GITHUB_REPO"), os.Getenv("GITHUB_ACCESS_TOKEN"))
	fmt.Printf("apiClient.client: %v\n", apiClient.client)
}

func TestFetchPosts(t *testing.T) {
	apiClient := NewApi(os.Getenv("GITHUB_USER_NAME"), os.Getenv("GITHUB_REPO"), os.Getenv("GITHUB_ACCESS_TOKEN"))
	disuccsions, err := apiClient.FetchPosts("", "")
	if err != nil {
		t.Fatal("fetch posts errors:", err)
	}
	if disuccsions.TotalCount == 0 {
		t.Failed()
	}
	if len(disuccsions.Nodes) == 0 {
		t.Failed()
	}

	if disuccsions.PageInfo.HasNextPage {
		t.Fail()
	}
	fmt.Printf("get disussion:%v\n", disuccsions)
}
