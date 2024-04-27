package core

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type BlogApi struct {
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

func NewApi(username string, repo string, accessToken string) BlogApi {

	return BlogApi{
		client: githubv4.NewClient(oauth2Client(accessToken)),
		owner:  username,
		repo:   repo,
	}
}

// FetchPosts 获取 post 列表
func (api *BlogApi) FetchPosts(before, after, categoryId string) (posts Discussions, err error) {

	var q struct {
		Resposity struct {
			Discussion Discussions `graphql:"discussions(first:$discussion_first,last: $discussion_last, orderBy:{field:CREATED_AT,direction:DESC}, after:$after, before:$before, categoryId:$categoryId)"`
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}

	binds := map[string]interface{}{
		"discussion_first": githubv4.Int(PER_PAGE_POST_COUNT),
		"discussion_last":  (*githubv4.Int)(nil),
		"owner":            githubv4.String(api.owner),
		"repo":             githubv4.String(api.repo),
		"after":            (*githubv4.String)(nil),
		"before":           (*githubv4.String)(nil),
		"categoryId":       (*githubv4.ID)(nil),
		"label_first":      githubv4.Int(LABEL_MAX_COUNT),
	}

	if len(after) > 0 {
		binds["after"] = (githubv4.String)(after)
	}

	if len(before) > 0 {
		binds["before"] = (githubv4.String)(before)
		binds["discussion_last"] = githubv4.Int(PER_PAGE_POST_COUNT)
		binds["discussion_first"] = (*githubv4.Int)(nil)
	}

	if len(categoryId) > 0 {
		binds["categoryId"] = (githubv4.ID)(categoryId)
	}

	err = api.client.Query(context.Background(), &q, binds)
	if err != nil {
		return Discussions{}, err
	}
	posts = q.Resposity.Discussion
	return posts, nil
}

// FetchPostComments 根据 id 获取 posts 的所有评论
func (api *BlogApi) FetchPost(number uint64) (discussion Node, err error) {
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
func (api *BlogApi) FetchCategories(before string, after string) (Categories, error) {
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

// FetchPostsByLabel 根据 label 和 category 获取对应的 posts
func (api *BlogApi) QueryPosts(keyword string, label string, categories []string) (SearchResults, error) {
	var q struct {
		Search struct {
			PageInfo PageInfo
			Nodes    []struct {
				Node `graphql:"... on Discussion"`
			}
		} `graphql:"search(query: $query first: $first, type: $type)"`
	}

	var query = fmt.Sprintf("repo:%s/%s ", api.owner, api.repo)

	if len(strings.Trim(keyword, "")) != 0 {
		query = fmt.Sprintf("%s %s ", query, keyword)
	}

	if len(strings.Trim(label, "")) != 0 {
		query = fmt.Sprintf("%s label:\"%s\" ", query, label)
	}

	if len(categories) != 0 {
		for _, category := range categories {
			query = fmt.Sprintf("%s category:\"%s\" ", query, category)
		}
	}

	binds := map[string]interface{}{
		"first":       githubv4.Int(PER_PAGE_POST_COUNT),
		"query":       githubv4.String(query),
		"type":        githubv4.SearchTypeDiscussion,
		"label_first": githubv4.Int(LABEL_MAX_COUNT),
	}

	err := api.client.Query(context.Background(), &q, binds)
	if err != nil {
		return SearchResults{}, err
	}
	var posts []Node
	for _, node := range q.Search.Nodes {
		posts = append(posts, node.Node)
	}
	return SearchResults{
		PageInfo: q.Search.PageInfo,
		Nodes:    posts,
	}, nil
}

// fetch all labels from discussion 获取所有的 labels
func (api *BlogApi) FetchAllLabels() ([]Label, error) {
	var q struct {
		Reposity struct {
			Labels struct {
				Edges []struct {
					Node Label
				}
			} `graphql:"labels(first: $first)"`
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}

	binds := map[string]interface{}{
		"first": githubv4.Int(MAX_LABELS_COUNT),
		"owner": githubv4.String(api.owner),
		"repo":  githubv4.String(api.repo),
	}

	err := api.client.Query(context.Background(), &q, binds)
	if err != nil {
		return []Label{}, err
	}

	var labels []Label
	for _, node := range q.Reposity.Labels.Edges {
		labels = append(labels, node.Node)
	}
	return labels, nil
}
