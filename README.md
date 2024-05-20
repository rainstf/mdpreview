# MD Preview

MD Preview is a plugin for neovim that enables hot-reload rendering of markdown in the browser.

## Installation

using Lazy:

```lua
return {
	"abql/mdpreview",
}
```

Also ensure that the binary is included in your $PATH (see [How It Works](#how-it-works)).

## Usage

Calling `:MdPrev` will attempt to render the current buffer to HTML and display it in your $BROWSER

Please note that the synchronized scrolling feature is experimental and may not be 100% accurate.

## Configuration

There are relatively few options:

```lua
return {
	"abql/mdpreview",

	config = function()
		require("mdpreview").setup({
			on_event = "LeaveInsert",  -- reload trigger event
			scrolling = false,  -- off by default	
		})
	end
}
```
You can specify any of neovim's events for the reload method (see `:h events`).

To update the page content on each keystroke, set on_event to "TextChangedI".

## How This Works

This plugin consists of two main components: a Lua script and a Go binary.

There is no pre-built binary provided with this plugin, so you will have to build it from source.

Once both components are set up, calling the `:MdPrev` command from neovim will load any configuration changes and call the binary. The text written in the current buffer in neovim will then be sent (on the specified event) to the binary via a [nvim channel](https://neovim.io/doc/user/channel.html). The binary then converts the text (hopefully markdown) into HTML and sends it to the client (your browser) using [Server Sent Events (SSE)](https://html.spec.whatwg.org/multipage/server-sent-events.html).
