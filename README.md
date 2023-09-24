# Pure

> It's a blog based on github discussion that lets you focus on your creativity. All you need to do is open your browser, log in to github and start your creative journey, no more worry about file, sync deployment and image uploads. What's more  you can migrate your blog in the way you like!

# Deploy

Before you do a deployment, you need to config your own `pure.yaml`

```shell
cp pure.sample.yaml pure.yaml
```
Get your own [GITHUB_ACCESS_TOKEN](https://github.com/settings/tokens) here

# [Deploy To Vercel](https://github.com/LeetaoGoooo/pure/tree/vercel)

Before you do a deployment, you need to config `BlogConfig` in `constants.go` 

Get your own [GITHUB_ACCESS_TOKEN](https://github.com/settings/tokens) here


[![Deploy with Vercel](https://vercel.com/button)](https://vercel.com/new/clone?repository-url=https%3A%2F%2Fgithub.com%2FLeetaoGoooo%2Fpure%2Ftree%2Fvercel&env=ACCESS_TOKEN,USER_NAME,REPO,REPO_ID,WEB_HOST,WEB_NAME,WEB_BIO,WEB_EMAIL,CATEGORY_ID,CATEGORY_NAME)

## comments

pure use [gitcus](https://github.com/giscus/giscus) as comment system, visit [gitcus-website](https://giscus.app/) to get your repo id ,which can be found in `data-repo-id`

![gitcus-config](./screens/gitcus-config.png)

Then just run the application

```shell
go build main.go
./main
```

# Screens

<image src="./screens/homepage.png" width="300"/><image src="./screens/postpage.png" width="300"/>
