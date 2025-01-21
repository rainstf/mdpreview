# MD Preview

Neovim plugin for hot-reload rendering of markdown in the browser.

## Features

## Dependencies

- Plenary
- MDPreview binary

## Getting Started

Using [Lazy.nvim](https://github.com/folke/lazy.nvim)

```lua
return {
    "nolibc/mdpreview",

    config = function()
        require("mdpreview").setup({
            on_event = "InsertLeave",  -- reload trigger event
            scrolling = false,  -- off by default
	})
    end
}
```

Also ensure that the binary is included in your $PATH ([How This Works](#how-this-works)).

## Usage

Calling `:MdPrev` renders your current buffer to HTML, which is displayed in your browser.

> [!NOTE]
> Synchronized scrolling is experimental and may not be accurate.

## Configuration

You can specify any of neovim's events for the reload method (see `:h events`).

To update the page content on each keystroke, set on_event to "TextChangedI".

## How This Works

This plugin consists of two main components: a Lua script and a Go binary.

There is no pre-built binary provided with this plugin, so you will have to build it from source. The binary is included in this repo (this is not a particularly good way to handle this, I know). 

Once both components are set up, calling the `:MdPrev` command from neovim will load any configuration changes and call the binary. The text written in the current buffer in neovim will then be sent (on the specified event) to the binary via a [nvim channel](https://neovim.io/doc/user/channel.html). The binary then converts the text (hopefully markdown) into HTML and sends it to the client (your browser) using [Server Sent Events (SSE)](https://html.spec.whatwg.org/multipage/server-sent-events.html).
