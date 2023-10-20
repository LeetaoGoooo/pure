package core

import (
	"github.com/shurcooL/githubv4"
)

type LabelNodes struct {
	Nodes []Label `json:"nodes,omitempty"`
}

type Label struct {
	Id   githubv4.String `json:"id,omitempty"`
	Name githubv4.String `json:"name,omitempty"`
}

type Category struct {
	Id   githubv4.String `json:"id,omitempty"`
	Name githubv4.String `json:"name,omitempty"`
}

type Categories struct {
	PageInfo PageInfo
	Nodes    []Category
}

type Node struct {
	Number    githubv4.Int      `json:"number,omitempty"`
	Id        githubv4.String   `json:"id,omitempty"`
	Title     githubv4.String   `json:"title,omitempty"`
	Body      githubv4.String   `json:"body,omitempty"`                        // markdown 的 body
	BodyHTML  githubv4.HTML     `json:"bodyHTML,omitempty" graphql:"bodyHTML"` // bodyHTML
	BodyText  githubv4.String   `json:"bodyText,omitempty"`
	CreatedAt githubv4.DateTime `json:"created_at,omitempty"`
	// UpdatedAt githubv4.DateTime `json:"updated_at,omitempty"`
	Category Category   `json:"category,omitempty"`
	Lables   LabelNodes `graphql:"labels(first: $label_first)" json:"lables,omitempty"` // 一个节点最大的 label 不应该超过 LABEL_MAX_COUNT
}

type PageInfo struct {
	EndCursor       githubv4.String  `json:"end_cursor,omitempty"`
	HasNextPage     githubv4.Boolean `json:"has_next_page"`
	HasPreviousPage githubv4.Boolean `json:"has_previous_page"`
	StartCursor     githubv4.String  `json:"start_cursor,omitempty"`
}

type Discussions struct {
	TotalCount int      `json:"total_count,omitempty"`
	PageInfo   PageInfo `json:"page_info,omitempty"`
	Nodes      []Node   `json:"nodes,omitempty"`
}

type Author struct {
	AvatarUrl githubv4.String `json:"avatar_url,omitempty"`
	Login     githubv4.String `json:"login,omitempty"`
}

type Comment struct {
	Id     githubv4.String `json:"id,omitempty"`
	Body   githubv4.String `json:"body,omitempty"`
	Author Author          `json:"author,omitempty"`
	Repies Repies          `graphql:"labels(first: $replies_first)" json:"repies,omitempty"`
}

type Comments struct {
	TotalCount int       `json:"total_count,omitempty"`
	PageInfo   PageInfo  `json:"page_info,omitempty"`
	Nodes      []Comment `json:"nodes,omitempty"`
}

type Repies struct {
	PageInfo PageInfo  `json:"page_info,omitempty"`
	Nodes    []Comment `json:"nodes,omitempty"`
}
