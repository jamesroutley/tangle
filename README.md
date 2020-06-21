# Tangle

Tangle is a literate programming tool which:

1. Extracts the contents of code blocks from Markdown files
2. Concatenates them
3. Prints them, or writes them to a file

It's designed to make writing tutorials and blog posts which contain snippets of
code as simple as possible by letting you easily test and validate any code
you've written.

## Features

- A `--watch` option, which automatically re-runs Tangke when your Markdown files change
- Saving code written in different languages in different output files
- Placing code blocks in the outputted file in a different order to how they're
  defined in the Markdown file
- Rewriting code blocks

## Install

Tangle is written in Go, and can be installed with:

```sh
$ go get -u github.com/jamesroutley/tangle
```

## A simple example

Let's say we're writing an article on iteration in JavaScript, stored in a file
`iteration.md`:

````markdown
# Iteration in JavaScript

Given an array of numbers:

```javascript
let numbers = [1, 2, 3, 4, 5];
```

We can iterate over them using this syntax:

```javascript
for (let i = 0; i < numbers.length; i++) {
  let number = numbers[i];
  // Do something...
}
```
````

We can run Tangle on it with:

```sh
$ tangle iteration.md --outfile iteration.js
$ cat iteration.js
let numbers = [1, 2, 3, 4, 5];

for (let i = 0; i < numbers.length; i++) {
  let number = numbers[i];
  // Do something...
}
```

Tangle has extracted the JavaScript code, and written it to a file, where we can
check it executes correctly, or run a linter or tests against it.

Tangle also has a `--watch` option, which automatically re-extracts code when
changes are made to your Markdown file.

## Named code blocks

Let's say you're writing a blog post about optimising some code

By default, Tangle will concatenate all code blocks. Sometimes however, you want
a later code block to replace an earlier one. You can do this by naming your
code blocks by adding a name after the language definition:

````markdown
```javascript greeting
console.log("hello");
```

```javascript greeting
console.log("hello world");
```
````

If two blocks have the same name, the later one will replace the earlier one. If
we run Tangle on the file above:

```sh
$ tangle greeting.md --outfile greeting.js
$ cat greeting.js
console.log("hello world")
```

Code block names can't contain spaces.

## Config file

The simple usage above should cater for most use cases, but Tangle also supports
more complex projects which involve multiple files, multiple outputs and/or
multiple languages. To start using these features, we'll need to write a config
file which tells Tangle what you want it to do.

Config files have the following schema:

```javascript
{
  "targets": [
    {
      // The files that Tangle extracts code from
      "sources": ["example.md"],

      // Name of the file that Tangle writes the extracted code to
      "outfile": "example.go",

      // Filters instruct Tangle to only include certain code blocks
      "filters": {
        // The language filter only includes code blocks of this language. The
        // language must be explicitly written next to the code block opening
        // tag, Tangle can't infer the language from the code itself
        "language": "go",
        // The regex filter only includes code blocks where the code matches
        // this regex.
        "regex": "^// example.go"
      },

      // Order specifies the order in which Tangle will write the code blocks
      // to the output file. To use this feature, you must name your code
      // blocks. You don't need to include all code blocks in this list, so you
      // can also use order to selectively include blocks.
      "order": ["block_1", "block_2"]
    }
  ]
}
```

You can use this config to handle a wide range of use cases.

## More complex examples

### Multiple languages

Let's say you're writing an article on styling HTML lists. Your article will
contain some HTML code, and some CSS code. You can use Tangle's language filter
to split this code into two different files:

```json
{
  "targets": [
    {
      "sources": ["list_styling.md"],
      "outfile": "index.html",
      "filters": {
        "language": "html"
      }
    },
    {
      "sources": ["list_styling.md"],
      "outfile": "styles.css",
      "filters": {
        "language": "css"
      }
    }
  ]
}
```

### Multiple steps

Let's say you're writing a 'get started' tutorial which has multiple steps.
Later steps build on code from eariler steps, but replace certain parts. We'd
like Tangle to output a different file for each step, so we can check things are
working. First, we need to name our code blocks, then use the `order` param to
select the blocks we want at each step:

```json
{
  "targets": [
    {
      "sources": ["get-started.md"],
      "outfile": "step_1.py",
      "order": ["initialise_game", "game_setup"]
    },
    {
      "sources": ["get-started.md"],
      "outfile": "step_2.py",
      "order": ["initialise_game", "game_setup", "game_loop"]
    },
    {
      "sources": ["get-started.md"],
      "outfile": "step_3.py",
      "order": ["initialise_game", "game_setup", "game_loop_with_physics"]
    }
  ]
}
```

### Building code from all your blog posts

Let's say you write a programming blog, and you've got multiple posts which
contain code snippets. You can use a single Tangle config file to build all of
these:

```json
{
  "targets": [
    {
      "sources": ["write-a-binary-tree.md"],
      "outfile": "../source_code/binary_tree/binary_tree.py"
    },
    {
      "sources": ["bulid-a-dns-server.md"],
      "outfile": "../source_code/dns-server/main.go"
    },
    {
      "sources": ["get-started-with-p5.md"],
      "outfile": "../source_code/p5-get-started/index.html",
      "filters": {
        "language": "html"
      }
    },
    {
      "sources": ["get-started-with-p5.md"],
      "outfile": "../source_code/p5-get-started/sketch.js",
      "filters": {
        "language": "javascript"
      }
    }
  ]
}
```
