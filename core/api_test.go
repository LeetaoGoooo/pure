package core

import (
	"testing"
)

func TestApi(t *testing.T) {
	api := NewApi("username", "repo", "token")
	t.Run("TestApiCategories", func(t *testing.T) {
		categories, err := api.FetchCategories("", "")
		if err != nil {
			t.Errorf("FetchCategories error: %v", err)
		}
		t.Logf("categories: %v", categories)
	})

	t.Run("TestApiQuery", func(t *testing.T) {
		var label = "python"
		var keyword = "django"
		queryTestCases := []struct {
			keyword     *string
			label       *string
			categories  *[]string
			resultCount int
		}{
			{&keyword, nil, &[]string{"随笔"}, 3},
			{nil, &label, &[]string{"随笔", "历史存档"}, 10},
			{&keyword, &label, &[]string{"随笔", "历史存档"}, 4},
		}

		for _, queryTestCase := range queryTestCases {
			posts, err := api.QueryPosts(queryTestCase.keyword, queryTestCase.label, queryTestCase.categories)
			if err != nil {
				t.Errorf("QueryPosts error: %v", err)
			}
			if len(posts) != queryTestCase.resultCount {
				t.Errorf("QueryPosts failed, %v the number of results are supposed to be %d, but got %d\n", queryTestCase, queryTestCase.resultCount, len(posts))
			}
		}
	})
}
