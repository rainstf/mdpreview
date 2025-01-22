<h1 align="center">
MD Preview
</h1>

<p align="center">
Neovim plugin for hot-reload rendering of markdown in the browser.
</p>

---

## Features

- **CommonMark**: uses [goldmark](https://pkg.go.dev/github.com/yuin/goldmark@v1.7.8) which is compliant with CommonMark 0.31.2.
- **Images**: embed local images with base64 encoding.
- **YouTube Embeds**: just as it sounds.
- **Syntax Highlighting**: supports syntax-highlighting fenced code blocks.

## Getting Started

Using [Lazy.nvim](https://github.com/folke/lazy.nvim)

```lua
return {
    "nolibc/mdpreview",
    opts = {
        events = {
            reload_on = 'InsertLeave',
            scrolling = false
        }
    },
    keys = {
        { "<leader>ms", "<cmd>MdPreviewStart<cr>" },
        { "<leader>mf", "<cmd>MdPreviewStop<cr>" }
    }
}
```

You may ommit the lazy-loading `keys` spec and call the user commands directly instead.


### Building MDPreview Server

This plugin uses a binary executable file to communicate with Neovim and host a web server. It does not come pre-built, so you'll need a Go compiler to build it.

> [!CAUTION]
> Before running anything, please verify its contents. It's never a good idea to blindly execute scripts/commands—especially ones found on the internet.

```bash
cd $HOME/.local/share/nvim/lazy/mdpreview && bash ./install.sh
```

The above shell command will automate the build and install process. But doing it manually is probably safer.

### Build It Manually

Alternatively, you could execute the commands manually:

1. `cd $HOME/.local/share/nvim/lazy/mdpreview/binary`

In the above, we navigate to the directory the package manager cloned. For lazy.nvim, cloned plugins are usually found in ~/.local/share/nvim/lazy/

2. `go mod tidy && go build .`

For step 2, we ensure all dependencies the binary requires are downloaded. We then build the binary executable file.

3. `sudo cp ./MDPreview /usr/bin`

Step 2 gave us the binary file. We now need to add it do a directory in our PATH. This may require root privileges depending on where you want to move it.

## Dependencies

- [Plenary](https://github.com/nvim-lua/plenary.nvim)
- [Go](https://go.dev/doc/install)
- MDPreview server

## Configuration

Configuration is currently limited.

```lua
opts = {
    events = {
        reload_on = 'InsertLeave',
        scrolling = false
    }
}
```

- `reload_on`: Neovim event used to reload the page.
You can specify neovim events for the reload method (`:h events`).
For instance, use 'TextChangedI' to update the page content on each keystroke.

- `scrolling`: Controls automatically page scroll.

> [!WARNING]
> Synchronized scrolling is experimental and may not be accurate.

## How This Works

MDPreview consists of two components—a Lua script and a Go binary.

**Why It Matters**: Understanding the plugin and web server interaction is crucial if you encounter difficulties.

Once both components are installed, you can interact with the plugin via its Lua API. The server comes in the form of a binary executable and is spawned as a user process. To do this, Lua code tries to find the MDPreview binary in your system PATH. Configuration changes are loaded upon spawning the server. Once running, text written in the current buffer will then be sent— on the reload event—to the server over an [nvim channel](https://neovim.io/doc/user/channel.html). Next, the server converts the text to HTML and finally gives it to the client (your browser) using [Server Sent Events (SSE)](https://html.spec.whatwg.org/multipage/server-sent-events.html).

There isn't a pre-built server binary, so you'll need to build it from source. This requires you to have Go installed on your system.
