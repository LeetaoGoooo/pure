package enties

type PageQuery struct {
	Next         string `form:"next,omitempty"`
	Pre          string `form:"pre,omitempty"`
	CategoryId   string `uri:"category_id"`
	CategoryName string `uri:"category_name"`
}

type PostQuery struct {
	Id    uint64 `uri:"id" binding:"required"`
	Title string `uri:"title" binding:"required"`
}

type Response[T any] struct {
	Code    int    `json:"code,omitempty"`
	Data    *T     `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

type Category struct {
	Id   string `yaml:"id"`
	Name string `yaml:"name"`
}

type Website struct {
	Host  string `yaml:"host"`
	Bio   string `yaml:"bio"`
	Name  string `yaml:"name"`
	Email string `yaml:"email"`
}

type PureConfig struct {
	UserName    string     `yaml:"username"`
	Repo        string     `yaml:"repo"`
	RepoId      string     `yaml:"repoId"`
	Website     Website    `yaml:"website"`
	AccessToken string     `yaml:"accessToken"`
	Categories  []Category `yaml:"categories"`
	About       uint64     `yaml:"about"`
}
