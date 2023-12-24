# llm

Small CLI application to interact with [Gemini][0].

## Install

```
$ go install github.com/danobi/llm@master
```

You'll also need an API key (see link above). Provide it through environment
variable `API_KEY` or put it on a single line in `~/.config/llm/key`. API keys
are currently free.

## Usage

```
$ llm "Give me a 3 word rhyme"
Cat in hat

$ echo "Give me a 3 word rhyme" | llm -
Criss cross sauce

$ llm
Reading from stdin...
^C to cancel, ^D to send
Give me a 3 word rhyme
Cat in hat

$ llm "What is this file?" - "Give me a 1 sentence answer" < /proc/mounts
This is a list of filesystems currently mounted on the system.
```

[0]: https://ai.google.dev/

## Doesn't this duplicate a lot of other projects?

Yeah, sure. I didn't even look though b/c this literally took me 30 minutes
to write. Also Gemini is free and I don't want to spend money on ChatGPT
anymore.
