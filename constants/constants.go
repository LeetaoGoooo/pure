package constants

import (
	"os"
	"pure/enties"
	"strconv"
)

var aboutPage int
var BlogConfig enties.PureConfig

func init() {
	aboutPage, err := strconv.ParseUint(os.Getenv("ABOUT"), 10, 64)

	if err != nil {
		panic("about page set failed!")
	}

	BlogConfig = enties.PureConfig{
		About:       aboutPage,
		AccessToken: os.Getenv("ACCESS_TOKEN"),
		UserName:    os.Getenv("USER_NAME"),
		Repo:        os.Getenv("REPO"),
		RepoId:      os.Getenv("REPO_ID"),
		Website: enties.Website{
			Host:  os.Getenv("WEB_HOST"),
			Name:  os.Getenv("WEB_NAME"),
			Bio:   os.Getenv("WEB_BIO"),
			Email: os.Getenv("WEB_EMAIL"),
		},
		//  []enties.Category{{
		// 	Id:   "DIC_kwDOIOm9Ys4CSAcy",
		// 	Name: "随笔",
		// }}
		Categories: []enties.Category{
			{
				Id:   os.Getenv("CATEGORY_ID"),
				Name: os.Getenv("CATEGORY_NAME"),
			},
		},
	}
}

const ErrorPage = `
<!DOCTYPE html>
<html>

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <script src="https://cdn.tailwindcss.com"></script>
</head>

<body>


    <!-- Page Container -->
    <div class="flex items-center justify-center min-h-screen bg-white py-48">
        <div class="flex flex-col">
            <!-- Error Container -->
            {{ if .Message }}
            <div class="flex flex-col items-center">
                <div class="font-bold text-3xl xl:text-7xl lg:text-6xl md:text-5xl mt-10">
                    {{ .Message }}
                </div>
            </div>
            {{ end }}

            <div class="flex flex-col mt-18 justify-center">
                <div class="text-gray-400 font-bold uppercase">
                    Continue With
                </div>

                <div class="flex flex-col items-stretch mt-5">
                    <div class="flex flex-row group px-4 py-8
                    border-t hover:cursor-pointer
                    transition-all duration-200 delay-100">

                        <div class="grow flex flex-col pl-5 pt-2">
                            <div class="font-bold text-sm md:text-lg lg:text-xl group-hover:underline">
                                <a href="/">Home</a>
                            </div>
                            <div class="font-semibold text-sm md:text-md lg:text-lg
                            text-gray-400 group-hover:text-gray-500
                            transition-all duration-200 delay-100">
                                Everything starts here
                            </div>
                        </div>
                    </div>

                </div>
            </div>
        </div>
    </div>
</body>

</html>`

const PostPage = `
<!DOCTYPE html>
<html class="h-full antialiased" lang="en" data-color-mode="auto" data-light-theme="light" data-dark-theme="dark">

<head>
    <meta charSet="utf-8" />
    <meta name="viewport" content="width=device-width" />
    <title>{{ .Post.Title }}</title>
    <meta property="og:description" content="{{  .Post.BodyText | previewContent}}" />
    <meta property="og:title" content="{{  .Post.Title }}" />

    <script src="https://cdn.tailwindcss.com"></script>
</head>

<body class="flex flex-col bg-zinc-50 dark:bg-black">
    <div id="__next">
        <div class="fixed inset-0 flex justify-center sm:px-8">
            <div class="flex w-full max-w-7xl lg:px-8">
                <div class="w-full bg-white ring-1 ring-zinc-100 dark:bg-zinc-900 dark:ring-zinc-300/20"></div>
            </div>
        </div>
        <div class="relative">
            <div class="flex justify-center sm:px-8">
                <div class="relative">
				<nav id="header">
				<div class="w-full flex items-center justify-between mt-0 px-6 py-2">
				   <div class="md:flex md:items-center md:w-auto w-full order-3 md:order-1" id="menu">
					  <nav>
						 <ul
							class="flex rounded-full bg-white/90 px-3 text-sm font-medium text-zinc-800 shadow-lg shadow-zinc-800/5 ring-1 ring-zinc-900/5 backdrop-blur dark:bg-zinc-800/90 dark:text-zinc-200 dark:ring-white/10">
							<li><a class="relative block px-3 py-2 transition hover:text-teal-500 dark:hover:text-teal-400 item-menu"
								  href="/">Home</a></li>
							{{ range $index, $Navbar := .Navbars }}
							<li><a class="relative block px-3 py-2 transition hover:text-teal-500 dark:hover:text-teal-400 item-menu"
								  href="/category/{{.Id}}/{{.Name}}">{{.Name}}</a></li>
							{{ end }}
			 
							{{ if .About }}
							<li><a class="relative block px-3 py-2 transition hover:text-teal-500 dark:hover:text-teal-400 item-menu"
								  href="/about">About</a></li>
							{{ end }}
						 </ul>
			 
					  </nav>
				   </div>
			 
				</div>
			 </nav>
                </div>
            </div>
            <main>
                <div class="sm:px-8 mt-16 lg:mt-32">
                    <div class="mx-auto max-w-7xl lg:px-8">
                        <div class="relative px-4 sm:px-8 lg:px-12">
                            <div class="mx-auto max-w-2xl lg:max-w-5xl">
                                <div class="xl:relative  markdown-body pb-20 pt-20">
                                    <div class="mx-auto max-w-2xl">
                                        <article class="dark:text-white">
                                            <header class="flex flex-col">
                                                <h2
                                                    class="mt-6 text-4xl font-bold tracking-tight text-zinc-800 dark:text-zinc-100 sm:text-5xl">
                                                    {{ .Post.Title }}</h1>
                                                <div class="flex w-full items-center justify-between pb-3 mt-6">
                                                    <div class="flex items-center space-x-3">
                                                        <div class="text-xs text-neutral-500">
                                                            {{ .Post.CreatedAt | formatDate }}</div>
                                                    </div>
                                                    <div class="flex items-center space-x-8">
                                                        {{range .Post.Lables.Nodes}}
                                                        <span
                                                            class="ml-4 text-xs inline-flex items-center leading-sm uppercase px-3 py-1 bg-green-200 text-teal-500 rounded-full dark:bg-slate-800 dark:text-white">
                                                            #{{.Name}}
                                                        </span>

                                                        {{end}}
                                                    </div>
                                                </div>
                                            </header>
                                            <div class="mt-8">
                                                <div
                                                    class="d-block color-fg-default comment-body markdown-body js-comment-body">
                                                    {{ .Post.BodyHTML | unescapeHtml }}
                                                </div>
                                            </div>
                                        </article>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="giscus mx-auto max-w-6xl lg:px-8 mt-10">
                </div>
                    <script src="https://giscus.app/client.js"
                    data-repo="{{$.Repo}}"
                    data-repo-id="{{$.RepoId}}"
                    data-mapping="number"
                    data-term="{{.Post.Number}}"
                    data-reactions-enabled="1"
                    data-emit-metadata="0"
                    data-input-position="top"
                    data-theme="preferred_color_scheme"
                    data-lang="zh-CN"
                    data-loading="lazy"
                    crossorigin="anonymous"
                    async>
            </script>
            </main>
			<footer class="mt-32">
			<div class="sm:px-8">
				<div class="mx-auto max-w-7xl lg:px-8">
					<div class="border-t border-zinc-100 pt-10 pb-16 dark:border-zinc-700/40">
						<div class="relative px-4 sm:px-8 lg:px-12">
							<div class="mx-auto max-w-2xl lg:max-w-5xl">
								<div class="flex flex-col items-center justify-between gap-6 sm:flex-row">
									<div class="flex gap-6 text-sm font-medium text-zinc-800 dark:text-zinc-200">
										<a class="transition hover:text-teal-500 dark:hover:text-teal-400" href="/">Home</a>
										<p class="text-sm text-zinc-400 dark:text-zinc-500">©
											<!-- -->2022
											<!-- --> Leetao. All rights reserved.
										</p>
										<p class="text-sm text-zinc-400 dark:text-zinc-500">
											Powed by <a href="https://github.com/LeetaoGoooo/pure">Pure</a>
										</p>
									</div>
								</div>
							</div>
						</div>
					</div>
				</div>
		</footer>
        </div>
    </div>
    <link rel="stylesheet" href="https://unpkg.com/normalize.css@8.0.1/normalize.css">
	<script src="https://cdn.jsdelivr.net/npm/github-syntax-light@0.5.0/index.min.js"></script>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/github-syntax-light@0.5.0/lib/github-light.min.css">
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/github-markdown-css/5.2.0/github-markdown.min.css">
</body>

</html>
`

const IndexPage = `
<!doctype html>
<html class="h-full antialiased js-focus-visible"
    style="--header-position:sticky; --content-offset:116px; --header-height:180px; --header-mb:-116px; --header-top:0px; --avatar-top:0px; --avatar-image-transform:translate3d(0rem, 0, 0) scale(1); --avatar-border-transform:translate3d(-0.222222rem, 0, 0) scale(1.77778); --avatar-border-opacity:0;"
    data-color-mode="auto" data-light-theme="light" data-dark-theme="dark">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <script src="https://cdn.tailwindcss.com"></script>
</head>

<body class="flex h-full flex-col bg-zinc-50 dark:bg-black">
    <div class="inset-0 flex justify-center sm:px-8">
	<nav id="header">
	<div class="w-full flex items-center justify-between mt-0 px-6 py-2">
	   <div class="md:flex md:items-center md:w-auto w-full order-3 md:order-1" id="menu">
		  <nav>
			 <ul
				class="flex rounded-full bg-white/90 px-3 text-sm font-medium text-zinc-800 shadow-lg shadow-zinc-800/5 ring-1 ring-zinc-900/5 backdrop-blur dark:bg-zinc-800/90 dark:text-zinc-200 dark:ring-white/10">
				<li><a class="relative block px-3 py-2 transition hover:text-teal-500 dark:hover:text-teal-400 item-menu"
					  href="/">Home</a></li>
				{{ range $index, $Navbar := .Navbars }}
				<li><a class="relative block px-3 py-2 transition hover:text-teal-500 dark:hover:text-teal-400 item-menu"
					  href="/category/{{.Id}}/{{.Name}}">{{.Name}}</a></li>
				{{ end }}
 
				{{ if .About }}
				<li><a class="relative block px-3 py-2 transition hover:text-teal-500 dark:hover:text-teal-400 item-menu"
					  href="/about">About</a></li>
				{{ end }}
			 </ul>
 
		  </nav>
	   </div>
 
	</div>
 </nav>
    </div>


    <main>
        <div class="sm:px-8 mt-20">
            <div class="mx-auto max-w-7xl lg:px-8">
                <div class="relative px-4 sm:px-8 lg:px-12">
                    <div class="mx-auto max-w-2xl lg:max-w-5xl">
                        <div class="max-w-2xl">
                            <h1 class="text-4xl font-bold tracking-tight text-zinc-800 dark:text-zinc-100 sm:text-5xl">
                                Leetao</h1>
                            <p class="mt-6 text-base text-zinc-600 dark:text-zinc-400">后端工程师，写有趣的代码，做有趣的事
                                <br>使用的语言: Python、Golang、Dart、Java、TypeScript
                            </p>
                            <div class="mt-6 flex gap-6"><a class="group -m-1 p-1" aria-label="Follow on Twitter"
                                    href="https://twitter.com/LeetaoGoooo"><svg viewBox="0 0 24 24" aria-hidden="true"
                                        class="h-6 w-6 fill-zinc-500 transition group-hover:fill-zinc-600 dark:fill-zinc-400 dark:group-hover:fill-zinc-300">
                                        <path
                                            d="M20.055 7.983c.011.174.011.347.011.523 0 5.338-3.92 11.494-11.09 11.494v-.003A10.755 10.755 0 0 1 3 18.186c.308.038.618.057.928.058a7.655 7.655 0 0 0 4.841-1.733c-1.668-.032-3.13-1.16-3.642-2.805a3.753 3.753 0 0 0 1.76-.07C5.07 13.256 3.76 11.6 3.76 9.676v-.05a3.77 3.77 0 0 0 1.77.505C3.816 8.945 3.288 6.583 4.322 4.737c1.98 2.524 4.9 4.058 8.034 4.22a4.137 4.137 0 0 1 1.128-3.86A3.807 3.807 0 0 1 19 5.274a7.657 7.657 0 0 0 2.475-.98c-.29.934-.9 1.729-1.713 2.233A7.54 7.54 0 0 0 22 5.89a8.084 8.084 0 0 1-1.945 2.093Z">
                                        </path>
                                    </svg></a>
                                <a class="group -m-1 p-1" aria-label="Follow on GitHub"
                                    href="https://github.com/LeetaoGoooo"><svg viewBox="0 0 24 24" aria-hidden="true"
                                        class="h-6 w-6 fill-zinc-500 transition group-hover:fill-zinc-600 dark:fill-zinc-400 dark:group-hover:fill-zinc-300">
                                        <path fill-rule="evenodd" clip-rule="evenodd"
                                            d="M12 2C6.475 2 2 6.588 2 12.253c0 4.537 2.862 8.369 6.838 9.727.5.09.687-.218.687-.487 0-.243-.013-1.05-.013-1.91C7 20.059 6.35 18.957 6.15 18.38c-.113-.295-.6-1.205-1.025-1.448-.35-.192-.85-.667-.013-.68.788-.012 1.35.744 1.538 1.051.9 1.551 2.338 1.116 2.912.846.088-.666.35-1.115.638-1.371-2.225-.256-4.55-1.14-4.55-5.062 0-1.115.387-2.038 1.025-2.756-.1-.256-.45-1.307.1-2.717 0 0 .837-.269 2.75 1.051.8-.23 1.65-.346 2.5-.346.85 0 1.7.115 2.5.346 1.912-1.333 2.75-1.05 2.75-1.05.55 1.409.2 2.46.1 2.716.637.718 1.025 1.628 1.025 2.756 0 3.934-2.337 4.806-4.562 5.062.362.32.675.936.675 1.897 0 1.371-.013 2.473-.013 2.82 0 .268.188.589.688.486a10.039 10.039 0 0 0 4.932-3.74A10.447 10.447 0 0 0 22 12.253C22 6.588 17.525 2 12 2Z">
                                        </path>
                                    </svg></a>
                                <a class="group -m-1 p-1" aria-label="Follow on Atom" href="/atom.xml">
                                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" 
                                    class="h-6 w-6 fill-zinc-500 transition group-hover:fill-zinc-600 dark:fill-zinc-400 dark:group-hover:fill-zinc-300">
                                        <path fill-rule="evenodd" d="M3.75 4.5a.75.75 0 01.75-.75h.75c8.284 0 15 6.716 15 15v.75a.75.75 0 01-.75.75h-.75a.75.75 0 01-.75-.75v-.75C18 11.708 12.292 6 5.25 6H4.5a.75.75 0 01-.75-.75V4.5zm0 6.75a.75.75 0 01.75-.75h.75a8.25 8.25 0 018.25 8.25v.75a.75.75 0 01-.75.75H12a.75.75 0 01-.75-.75v-.75a6 6 0 00-6-6H4.5a.75.75 0 01-.75-.75v-.75zm0 7.5a1.5 1.5 0 113 0 1.5 1.5 0 01-3 0z" clip-rule="evenodd" />
                                      </svg>                                      
                                </a>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <div class="sm:px-8 mt-16">
            <div class="mx-auto max-w-7xl lg:px-8">
                <div class="relative px-4 sm:px-8 lg:px-12">
                    <div class="mx-auto max-w-2xl lg:max-w-5xl">
                        <div class="mx-auto grid max-w-xl grid-cols-1 gap-y-20 lg:max-w-none lg:grid-cols-1">
                            <div class="flex flex-col gap-8">
                                {{range .Posts.Nodes}}
                                    {{ if .Category.Id | isExisted }}
                                <article
                                    class="group relative flex flex-col items-start rounded-xl p-5 bg-white dark:bg-slate-800">
                                    <div class="flex w-full items-center justify-between pb-3">
                                        <div class="flex items-center space-x-3">
                                            <div class="text-xs text-neutral-500">
                                                {{ .CreatedAt | formatDate }}</div>
                                        </div>
                                        <div class="flex items-center space-x-8">
                                            {{range .Lables.Nodes}}
                                            <span
                                                class="ml-4 text-xs inline-flex items-center leading-sm uppercase px-3 py-1 bg-green-200 text-teal-500 rounded-full dark:bg-slate-800 dark:text-white">
                                                #{{.Name}}
                                            </span>

                                            {{end}}
                                        </div>
                                    </div>
                                    <h2 class="text-base font-semibold tracking-tight text-zinc-800 dark:text-zinc-100">
                                        <div
                                            class="absolute -inset-y-6 -inset-x-4 z-0 scale-95 bg-zinc-50 opacity-0 transition  dark:bg-zinc-800/50 sm:-inset-x-6 sm:rounded-2xl">
                                        </div><a href="/post/{{.Number}}/{{.Title}}"><span
                                                class="absolute -inset-y-6 -inset-x-4 z-20 sm:-inset-x-6 sm:rounded-2xl"></span><span
                                                class="relative z-10">{{ .Title}}</span></a>
                                    </h2>
                                    <p class="relative z-10 mt-2 text-sm text-zinc-600 dark:text-zinc-400">
                                        {{ .BodyText | previewContent}}
                                    </p>
                                    <div aria-hidden="true"
                                        class="relative z-10 mt-4 flex items-center text-sm font-medium text-teal-500">
                                        阅读文章<svg viewBox="0 0 16 16" fill="none" aria-hidden="true"
                                            class="ml-1 h-4 w-4 stroke-current">
                                            <path d="M6.75 5.75 9.25 8l-2.5 2.25" stroke-width="1.5"
                                                stroke-linecap="round" stroke-linejoin="round"></path>
                                        </svg></div>

                                </article>
                                    {{ end }}
                                {{end}}

                            </div>

                        </div>
                    </div>
                </div>
            </div>
        </div>

        <div class='flex items-center justify-center mt-10'>
            <div class="flex justify-center items-center space-x-4">
                {{ if .Posts.PageInfo.HasPreviousPage }}
                <a class="px-2 py-1 text-3xl leading-6 text-slate-400 transition cursor-pointer "
                    href="/?pre={{.Posts.PageInfo.StartCursor}}">
                    < </a>
                {{ end }}
                {{ if .Posts.PageInfo.HasNextPage}}
                <a class="px-2 py-1 text-3xl leading-6 text-slate-400 transition  cursor-pointer"
                    href="/?next={{.Posts.PageInfo.EndCursor}}"> > </a>
                {{ end }}  
            </div>
        </div>
    </main>

	<footer class="mt-32">
    <div class="sm:px-8">
        <div class="mx-auto max-w-7xl lg:px-8">
            <div class="border-t border-zinc-100 pt-10 pb-16 dark:border-zinc-700/40">
                <div class="relative px-4 sm:px-8 lg:px-12">
                    <div class="mx-auto max-w-2xl lg:max-w-5xl">
                        <div class="flex flex-col items-center justify-between gap-6 sm:flex-row">
                            <div class="flex gap-6 text-sm font-medium text-zinc-800 dark:text-zinc-200">
                                <a class="transition hover:text-teal-500 dark:hover:text-teal-400" href="/">Home</a>
                                <p class="text-sm text-zinc-400 dark:text-zinc-500">©
                                    <!-- -->2022
                                    <!-- --> Leetao. All rights reserved.
                                </p>
                                <p class="text-sm text-zinc-400 dark:text-zinc-500">
                                    Powed by <a href="https://github.com/LeetaoGoooo/pure">Pure</a>
                                </p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
</footer>
</body>

</html>

`
