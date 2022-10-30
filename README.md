# Pure

> It's a blog based on github discussion that lets you focus on your creativity. All you need to do is open your browser, log in to github and start your creative journey, no more worry about file, sync deployment and image uploads. What's more  you can migrate your blog in the way you like!

# Deploy

Before you do a deployment, you need to set a few system variables

```shell
export GITHUB_USER_NAME=your_github_username
export GITHUB_REPO=your_repo_name
export GITHUB_ACCESS_TOKEN=your_access_token
export CATEGORY_ID=certain_category_id # optional
```
Get your own [GITHUB_ACCESS_TOKEN](https://github.com/settings/tokens) here

Then just run the application

```shell
go build main.go
./main
```

# Screens

<image src="./screens/homepage.png" width="200"/><image src="./screens/postpage.png" width="200"/>
