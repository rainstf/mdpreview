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

### Dependencies

- [Plenary](https://github.com/nvim-lua/plenary.nvim)
- [Go](https://go.dev/doc/install)

## Getting Started

Using [Lazy.nvim](https://github.com/folke/lazy.nvim)

```lua
return {
    "nolibc/mdpreview",
    dependencies = { 'nvim-lua/plenary.nvim' },
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

You may omit the lazy-loading `keys` spec and call the user commands directly instead.


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

Navigate to binary/ in the directory the package manager cloned. For lazy.nvim, cloned plugins are usually found in ~/.local/share/nvim/lazy/

2. `go mod tidy && go build .`

Ensure all dependencies the binary requires are downloaded. Then build it.

3. `sudo cp ./MDPreview /usr/bin`

Add it to a directory in PATH. This may require root privileges depending on the location.

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

> [!WARNING]
> Synchronized scrolling is experimental and may not be accurate.

- `reload_on`: Neovim event used to reload the page.
You may specify Neovim events (:h events) for the reload method.
For instance, use 'TextChangedI' to update the page content on each keystroke.

- `scrolling`: Controls automatically page scroll.

## How This Works

MDPreview consists of two components—a Lua script and a Go binary.

**Why It Matters**: Understanding the plugin and web server interaction is crucial if you encounter difficulties.

Once both components are installed, you can interact with the plugin via its Lua API. The server is a binary executable and is spawned as a user process. This is handled by the Lua component. Configuration changes are loaded upon spawning the server.

Text written in the current buffer will be sent—on the reload event—to the server over an [nvim channel](https://neovim.io/doc/user/channel.html). The server converts the text to HTML and gives it to the client (your browser) using [Server Sent Events (SSE)](https://html.spec.whatwg.org/multipage/server-sent-events.html). It is then rendered by the client.
