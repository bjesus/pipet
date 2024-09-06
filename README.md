<h1 align="center">
Pipet
</h1>

<p align="center">
Pipet a swiss-army tool for scraping and extracting data from online assets, made for hackers
</p>
<p align="center">
<img src="https://github.com/user-attachments/assets/e23a40de-c391-46a5-a30c-b825cc02ee8a" height="200">
</p>

Pipet is a command line based web scraper. It supports three mode of operation - HTML parsing, JSON parsing, and client-side JavaScript evaluation. It relies heavily on existing tools like curl, and it uses unix pipes for extending its built-in capabilities.

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

Add a tepmlate file called `hackernews.tpl` next to your `hackernews.pipet` file with this content:
```
<ul>
  {{range $index, $item := index (index . 0) 0}}
    <li>
  {{index $item 0}} ({{index $item 1}})</li>{{end}}
</ul>

<p>{{ .timestamp }}</p>
```

Now run `pipet hackernews.pipet` again and Pipet will automatically detect your template file, and render the results to it.
</details>
<details><summary>Use pipes</summary>

Use Unix pipes after your queries, as if they were running in your shell. For example, count the charaters in each title (with `wc`) and extract the full article URL (with [htmlq](https://github.com/mgdm/htmlq)):

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
You will need to have Go installed for this installation method.
You can use Go to install Pipet using `go install https://github.com/bjesus/pipet@latest`.  Otherwise you can run it without installing using `go run`.

## Distros
Packages are currently only available for [Arch Linux](https://aur.archlinux.org/packages/pipet-git).

# Usage

The only required argument for Pipet is the path to your `.pipet` file. Other than this, the `pipet` command accepts the following flags:

- `--json`, `-j` - Output as JSON (default: false)
-  `--template value`, `-t value` - Specify a path to a template file. You can also simply name the file like your `.pipet` file but with a `.tpl` extension for it to be auto-detected.
-  `--separator value`, `-s value` - Set a separator for text output (can be used multiple times for setting different separators for different levels of data nesting)
-  `---max-pages value`, `-p value` - Maximum number of pages to scrape (default: 3)
-  `--interval value`, `-i value` - Rerun pipet after X seconds. Use 0 to disable (default: 0)
-  `--on-change value`, `-c value` - A command to run when the pipet result is new
-  `--verbose`, `-v` - Enable verbose logging (default: false)
-  `--help`, `-h` - Show help

# Pipet files
Pipet files describe where and how to get the data you are interested in. They are normal text files containing one or more blocks, separated with an empty line. Line beginning with `//` are ignored and can be used for comments. Every block has at least 2 sections - the first line containing the URL and the tool we are using for scraping, and the following lines describing the selectors reaching the data we would like scrap. Some blocks can end with a special last line pointing to the "next page" selector - more on that later.

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

Blocks can start with either `curl` or `playwright`. Pipet doesn't just call these things `curl` because it's cool - it actually uses curl to fetch the resource. This might sound weird, but it's meant so that you can use your browser to find the request containing the information you are interested in, right click it, choose "Copy as cURL", and paste in your Pipet file. This ensures that your headers and cookies are all the same, making it very easy to get data which is behind a login page or is hidden from bots.

Starting a block with `playwright` will use a headless browser to navigate to the specified URL.

The lines following the first line are your _queries_. There are 3 different type of queries - for HTML files, for JSON files, and for websites loaded using `playwright`.

## HTML Queries
HTML Queries use CSS Selectors to point select specific elements. Whitespace nesting is used for iterations - parent lines will run as iterators, making their children lines run for each occurance of the parent selector. This means that you can use nesting to determine the structure of your final output. When writing your child selectors, note that the whole document isn't available anymore, and only the parent document is present during the iteration.

By defult, Pipet will return the `innerText` of your elements. If you need to another piece of data, use Pipes. When piping HTML elements, Pipet will pipe the element's complete HTML to the receiving program.

## JSON Queries
JSON Queries use GJSON to point select specific elements. Here too, whitespace nesting is used for iterations - parent lines will run as iterators, making their children lines run for each occurance of the parent selector. If you don't like GJSON, you can always use Pipes extract your data in other ways, for example with `jq`. See more examples below.

When using pipes with to send data to program that return valid JSON, Pipet will parse the JSON and embed it in its final output.

## Playwright Queries
Playwright Queries are different and do not use whitespace nesting. Instead, queries here are simply JavaScript code that will be evaluated after the webpage loaded. If the JavaScript code returns something that can be serialized as JSON, it will be included in Pipet's output. Otherwise, you can write JavaScript that will click, scroll or perform any othe action you might want.

## Unix Pipes
Sometimes CSS Selectors and GJSON aren't enough, or perhaps you just prefer using something you already know. This is why unix pipes are first class citizen in Pipet.

```
curl https://news.ycombinator.com/
span.yclinks a
  body
  body | htmlq --attribute href a
  body | htmlq --attribute href a | wc -c

curl http://localhost:8000/some.json 
people
  name
people | jq keys
@this | jq '[.products[].name]'
```

## Next page nav
