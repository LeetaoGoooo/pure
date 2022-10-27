package core

import (
	"context"
	"net/http"
	"os"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type api struct {
	client *githubv4.Client
	owner  string
	repo   string
}

func oauth2Client(accessToken string) *http.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	return oauth2.NewClient(ctx, ts)
}

func NewApi(username string, repo string, accessToken string) api {

	return api{
		client: githubv4.NewClient(oauth2Client(accessToken)),
		owner:  username,
		repo:   repo,
	}
}

// FetchPosts 获取 post 列表
func (api *api) FetchPosts(before string, after string) (posts Discussions, err error) {

	var q struct {
		Resposity struct {
			Discussion Discussions `graphql:"discussions(first:$discussion_first, orderBy:{field:CREATED_AT,direction:DESC}, after:$after, before:$before, categoryId:$categoryId)"`
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}

	binds := map[string]interface{}{
		"discussion_first": githubv4.Int(PER_PAGE_POST_COUNT),
		"owner":            githubv4.String(api.owner),
		"repo":             githubv4.String(api.repo),
		"after":            (*githubv4.String)(nil),
		"before":           (*githubv4.String)(nil),
		"categoryId":       os.Getenv("CATEGORY_ID"),
		"label_first":      githubv4.Int(LABEL_MAX_COUNT),
	}

	if len(after) > 0 {
		binds["after"] = (githubv4.String)(after)
	}

	if len(before) > 0 {
		binds["before"] = (githubv4.String)(before)
	}

	err = api.client.Query(context.Background(), &q, binds)
	if err != nil {
		return Discussions{}, err
	}
	posts = q.Resposity.Discussion
	return posts, nil
}

// FetchPostComments 根据 id 获取 posts 的所有评论
func (api *api) FetchPost(number int) (discussion Node, err error) {
	var q struct {
		Reposity struct {
			Discussion Node `graphql:"discussion(number:$number)"`
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}

	binds := map[string]interface{}{
		"number":      githubv4.Int(number),
		"owner":       githubv4.String(api.owner),
		"repo":        githubv4.String(api.repo),
		"label_first": githubv4.Int(LABEL_MAX_COUNT),
	}

	err = api.client.Query(context.Background(), &q, binds)
	if err != nil {
		return Node{}, err
	}
	post := q.Reposity.Discussion

	return post, err
}

// FetchCategories 获取所有的 Category
func (api *api) FetchCategories(before string, after string) (Categories, error) {
	var q struct {
		Reposity struct {
			DiscussionCategories Categories `graphql:"discussionCategories(first:$category_first, after:$after, before:$before)"`
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}

	binds := map[string]interface{}{
		"owner":          githubv4.String(api.owner),
		"repo":           githubv4.String(api.repo),
		"category_first": githubv4.Int(CATEGORY_MAX_COUNT),
		"after":          (*githubv4.String)(nil),
		"before":         (*githubv4.String)(nil),
	}

	if len(after) > 0 {
		binds["after"] = (githubv4.String)(after)
	}

	if len(before) > 0 {
		binds["before"] = (githubv4.String)(before)
	}

	err := api.client.Query(context.Background(), &q, binds)
	if err != nil {
		return Categories{}, err
	}
	categories := q.Reposity.DiscussionCategories

	return categories, err
}
