<h1 align="center">
Pipet
</h1>

<p align="center">
<a href="https://goreportcard.com/report/github.com/bjesus/pipet"><img src="https://goreportcard.com/badge/github.com/bjesus/pipet" /></a>
  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-yellow.svg" /></a>
  <a href="https://pkg.go.dev/github.com/bjesus/pipet"><img src="https://pkg.go.dev/badge/github.com/bjesus/pipet.svg" alt="Go Reference"></a>
  <br/>
a swiss-army tool for scraping and extracting data from online assets, made for hackers
</p>
<p align="center">
<img src="https://github.com/user-attachments/assets/e23a40de-c391-46a5-a30c-b825cc02ee8a" height="200">
</p>

Pipet is a command line based web scraper. It supports 3 modes of operation - HTML parsing, JSON parsing, and client-side JavaScript evaluation. It relies heavily on existing tools like curl, and it uses unix pipes for extending its built-in capabilities.

You can use Pipet to track a shipment, get notified when concert tickets are available, stock price changes, and any other kind of information that appears online.

# Try it out!
1. Create a `hackernews.pipet` file containing this:
```
curl https://news.ycombinator.com/
.title .titleline
  span > a
  .sitebit a
```
2. Run `go run github.com/bjesus/pipet/cmd/pipet@latest hackernews.pipet` or install Pipet and run `pipet hackernews.pipet`
3. See all of the latest hacker news in your terminal!

<details><summary>Get as JSON</summary>
  
Use the `--json` flag to make Pipet collect the results into a nice JSON.  For example, run `pipet --json hackernews.pipet` to a JSON representation of the above results.</details>
<details><summary>Render to a template</summary>

Add a template file called `hackernews.tpl` next to your `hackernews.pipet` file with this content:
```
<ul>
  {{range $index, $item := index (index . 0) 0}}
    <li>{{index $item 0}} ({{index $item 1}})</li>
  {{end}}
</ul>
```

Now run `pipet hackernews.pipet` again and Pipet will automatically detect your template file, and render the results to it.
</details>
<details><summary>Use pipes</summary>

Use Unix pipes after your queries, as if they were running in your shell. For example, count the characters in each title (with `wc`) and extract the full article URL (with [htmlq](https://github.com/mgdm/htmlq)):

```
curl https://news.ycombinator.com/
.title .titleline
  span > a
  span > a | wc -c
  .sitebit a
  .sitebit a | htmlq --attribute href a
```
</details>
<details><summary>Monitor for changes</summary>
  
Set an interval and a command to run on change, and have Pipet notify you when something happened. For example, get a notification whenever the Hacker News #1 story is different:

```
curl https://news.ycombinator.com/
.title .titleline a
```

Run it with `pipet --interval 60 --on-change "notify-send {}" hackernews.pipet`

</details>

# Installation

## Pre-built
Download the latest release from the [Releases](https://github.com/bjesus/pipet/releases/) page. `chmod +x pipet` and run `./pipet`.

## Compile
This installation method requires Go to be installed on your system.
You can use Go to install Pipet using `go install https://github.com/bjesus/pipet@latest`.  Otherwise you can run it without installing using `go run`.

## Distros
Packages are currently only available for [Arch Linux](https://aur.archlinux.org/packages/pipet-git) and Homebrew (`brew tap bjesus/pipet && brew install pipet`).

# Usage

The only required argument for Pipet is the path to your `.pipet` file. Other than this, the `pipet` command accepts the following flags:

- `--json`, `-j` - Output as JSON (default: false)
- `--template value`, `-t value` - Specify a path to a template file. You can also simply name the file like your `.pipet` file but with a `.tpl` extension for it to be auto-detected.
- `--separator value`, `-s value` - Set a separator for text output (can be used multiple times for setting different separators for different levels of data nesting)
- `---max-pages value`, `-p value` - Maximum number of pages to scrape (default: 3)
- `--interval value`, `-i value` - Rerun Pipet after X seconds. Use 0 to disable (default: 0)
- `--on-change value`, `-c value` - A command to run when the pipet result is new
- `--verbose`, `-v` - Enable verbose logging (default: false)
- `--help`, `-h` - Show help

# Pipet files
Pipet files describe where and how to get the data you are interested in. They are normal text files containing one or more blocks separated by an empty line. Lines beginning with `//` are ignored and can be used for comments. Every block can have 3 sections:

1. **Resource** - The first line containing the URL and the tool we are using for scraping
2. **Queries** - The following lines describing the selectors reaching the data we would like scrap
3. **Next page** - An _optional_ last line starting with `>` describing the selector pointing to the "next page" of data

Below is an example Pipet file.

```
// Read Wikipedia's "On This Day" and the subject of today's featured article
curl https://en.wikipedia.org/wiki/Main_Page
div#mp-otd li
  body
div#mp-tfa > p > b > a

// Get the weather in Alert, Canada
curl https://wttr.in/Alert%20Canada?format=j1
current_condition.0.FeelsLikeC
current_condition.0.FeelsLikeF

// Check how popular the Pipet repo is
playwright https://github.com/bjesus/pipet
Array.from(document.querySelectorAll('.about-margin .Link')).map(e => e.innerText.trim()).filter(t=> /^\d/.test(t) )
```

##  Resource

Resource lines can start with either `curl` or `playwright`.

### curl

Resource lines starting with `curl` will be executed using curl. This is meant so that you can use your browser to find the request containing the information you are interested in, right click it, choose "Copy as cURL", and paste in your Pipet file. This ensures that your headers and cookies are all the same, making it very easy to get data that is behind a login page or hidden from bots. For example, this is a perfectly valid first line for a block: `curl 'https://news.ycombinator.com/' --compressed -H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:131.0) Gecko/20100101 Firefox/131.0' -H 'Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/png,image/svg+xml,*/*;q=0.8' -H 'Accept-Language: en-US,en;q=0.5' -H 'Accept-Encoding: gzip, deflate, br, zstd' -H 'DNT: 1' -H 'Sec-GPC: 1' -H 'Connection: keep-alive' -H 'Upgrade-Insecure-Requests: 1' -H 'Sec-Fetch-Dest: document' -H 'Sec-Fetch-Mode: navigate' -H 'Sec-Fetch-Site: none' -H 'Sec-Fetch-User: ?1' -H 'Priority: u=0, i' -H 'Pragma: no-cache' -H 'Cache-Control: no-cache' -H 'TE: trailers'`.

### Playwright

Resource lines starting with `playwright` will use a headless browser to navigate to the specified URL. If you don't have a headless browser installed, Pipet will attempt to download one for you.

## Queries

Query lines define 3 things:
1. The way to the exact pieces of data you would like to extract (e.g. using CSS selectors)
2. The data structure your output will use (e.g. every title and URL should be grouped together by item)
3. The way the data will be processed (e.g. using Unix pipes) before it is printed

Pipet uses 3 different query types - for HTML, for JSON, and for when loading pages with Playwright.

### HTML Queries
HTML Queries use CSS Selectors to select specific elements. Whitespace nesting is used for iterations - parent lines will run as iterators, making their children lines run for each occurance of the parent selector. This means that you can use nesting to determine the structure of your final output. See the following 3 examples:

<details><summary>Get only the first title and first URL</summary>
  
```
curl https://news.ycombinator.com/
.title .titleline > a
.sitebit a
```

</details><details><summary>Get all the titles, and then get all URLs</summary>
  
```
curl https://news.ycombinator.com/
.title .titleline
  span > a
.title .titleline
  .sitebit a
```

</details><details><summary>Get all the title and URL for each story</summary>
  
```
curl https://news.ycombinator.com/
.title .titleline
  span > a
  .sitebit a
```
</details>

When writing your child selectors, note that the whole document isn't available anymore. Pipet is passing only your parent HTML to the child iterations.

By default, Pipet will return the `innerText` of your elements. If you need to another piece of data, use Unix pipes. When piping HTML elements, Pipet will pipe the element's complete HTML. For example, you can use `| htmq --attr href a` to extract the `href` attribute from links.

### JSON Queries

JSON Queries use the [GJSON syntax](https://github.com/tidwall/gjson/blob/master/SYNTAX.md) to select specific elements. Here too, whitespace nesting is used for iterations - parent lines will run as iterators, making their children lines run for each occurance of the parent selector. If you don't like GJSON, that's okay. For example, you can use `jq` by passing parts or the complete JSON to it using Unix pipes, like `@this | jq '.[].firstName'`.

When using pipes, Pipet will attempt to parse the returned string. If it's valid JSON, it will be parsed and injected as an object into the Pipet result.

### Playwright Queries

Playwright Queries are different and do not use whitespace nesting. Instead, queries here are simply JavaScript code that will be evaluated after the webpage loaded. If the JavaScript code returns something that can be serialized as JSON, it will be included in Pipet's output. Otherwise, you can write JavaScript that will click, scroll or perform any other action you might want.

## Next page

The Next Page line lets you specify a CSS selector that will be used to determine the link to the next page of data. Pipet will then follow it and execute the same queries over it. For example, see this `hackernews.pipet` file:
```
curl https://news.ycombinator.com/
.title .titleline
  span > a
  .sitebit a
> a.morelink
```

The Next Page line is currently only available when processing HTML files.
