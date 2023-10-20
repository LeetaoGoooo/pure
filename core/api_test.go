package core

import "testing"

func TestApi(t *testing.T) {
	api := NewApi("username", "repor", "token")
	t.Run("TestApiCategories", func(t *testing.T) {
		categories, err := api.FetchCategories("", "")
		if err != nil {
			t.Errorf("FetchCategories error: %v", err)
		}
		t.Logf("categories: %v", categories)
	})
}
