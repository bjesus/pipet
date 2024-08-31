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

# Try it now
1. Create a new Pipet spec file containing this:
2. run `go run https://github.com/bjesus/pipet/ myfile.pipet`

# Pipet spec files
The details on where the data you require exists should all be stored in Pipet spec file. This a text file normal text file containing one or more blocks, separated with an empty line. Line beginning with `#` are ignored. The first line of every block defines where the resource is and how to get it, and the following lines are the extractors, defining exactly what you'd like to take out of the resource.

```
curl https://en.wikipedia.org/wiki/Main_Page
; read Wikipedia's "on this day" and today's feature article subject
div#mp-otd li
  body
div#mp-tfa > p > b > a

curl https://wttr.in/Alert%20Canada?format=j1
; how cold is it over there?
current_condition.0.FeelsLikeC
current_condition.0.FeelsLikeF

playwright https://github.com/bjesus/pipet
; let's get all the stars, watchers, and forks from this repro
Array.from(document.querySelectorAll('.about-margin .Link')).map(e => e.innerText.trim()).filter(t=> /^\d/.test(t) )
```

Blocks can start with either `curl` or `playwright`. Pipet doesn't just call these things `curl` because it's cool - it actually uses curl to fetch the resource. This might sound weird, but it's meant so that you can use your browser to find the request containing the information you are interested at, right click it, choose "Copy as cURL", and paste here. This ensures that your headers and cookies are all the same, making it very easy to get data which is behind a login page or is hidden from bots.

Starting a block with `playwright` will download a headless browser to your computer and start it with the specified URL.

The lines following the first line are your _queries_. There are 3 different type of queries - for HTML files, for JSON files, and for websites loaded using `playwright`.

## HTML Queries

## JSON Queries

## Playwright Queries

## Pipes
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

# Running Pipet

## Installation

### Pre-built
Download the latest release from the Releases page. `chmod +x pipet` and run `./pipet`.

### Compile
You will need to have Go installed for this installation method.
You can use Go to install Pipet using `go install https://github.com/bjesus/pipet@latest`.  Otherwise you can run it without installing using `go run`.

### Distros
Packages are currently only available for Arch Linux.

## Usage

```
USAGE:
   pipet [global options] command [command options]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --json                                   Output as JSON (default: false)
   --separator value [ --separator value ]  Separator for text output (can be used multiple times)
   --template value                         Path to template file for output
   --max-pages value                        Maximum number of pages to scrape (default: 3)
   --interval value                         Maximum number of pages to scrape (default: 3)
   --on-change value                        Path to template file for output
   --help, -h                               show help
```

