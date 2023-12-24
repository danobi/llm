# llm

Small CLI application to interact with [Gemini][0].

## Install

`go install` to build and install `llm`.

You'll also need an API key (see link above). Provide it through environment
variable `API_KEY` or put it on a single line in `~/.config/llm/key`.

## Usage

```
$ llm "Give me a 3 word rhyme"
Cat in hat

$ echo "Give me a 3 word rhyme" | llm -
Criss cross sauce

$ llm "What is this file?" - < /proc/mounts
The provided information appears to be a list of mounted file systems on a Linux system, along with their respective mount points, file system types, and mount options. Here's a breakdown of each file system:

1. **proc:**
   - Mount point: /proc
   - File system type: proc
   - Options: rw,nosuid,nodev,noexec,relatime
   - Description: The proc file system provides information about processes, system memory, and other kernel data structures. It is a pseudo file system that exists in memory and does not occupy any disk space.

2. **sys:**
   - Mount point: /sys
[...]
```

[0]: https://ai.google.dev/
