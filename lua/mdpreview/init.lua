local async = require "plenary.async"
local Job = require'plenary.job'
local libuv = vim.loop

vim.api.nvim_create_user_command('MdPrev', 'lua Entry()', {})

local M = {}

local defaults = {
	on_event = "InsertLeave",
    scrolling = false,
}

M.config = defaults

function M.setup(user_options)
    user_options = user_options or {}
    for key, value in pairs(user_options) do
        if M.config[key] ~= nil then
            M.config[key] = value
        else
            error("Option " .. key .. " is not valid for this plugin.")
        end
    end
end

function BufContent(buf)
	local lines = vim.api.nvim_buf_get_lines(buf, 0, -1, false)
	local content = table.concat(lines, "\n")

	return content
end

function create_json_obj(content, evtype) -- TODO: rename to data
	local object = {
		["content" ] = content,
		["event" ] = evtype
	}
	if evtype == "scroll" or evtype == "reload" or evtype == "init" then
		return vim.json.encode(object)
	end
end

function send_initial_content(chan_id)
	local json = create_json_obj(BufContent(0), "init")
	SendServer(chan_id, json)
end

function serveInit()

    local chan_id = vim.fn.sockconnect("tcp", "127.0.0.1:8080")

	send_initial_content(chan_id)

	local event = M.config.on_event
	vim.api.nvim_create_autocmd(event, {
		callback = function()
			local json = create_json_obj(BufContent(0), "reload")
			SendServer(chan_id, json)
		end
	})
	vim.api.nvim_create_autocmd("TextChanged", {
		callback = function()
			local json = create_json_obj(BufContent(0), "reload")
			SendServer(chan_id, json)
		end
	})
	local last_pos = 0
	if M.config.scrolling == true then
		vim.api.nvim_create_autocmd({ "CursorMoved", "CursorMovedI" }, {
			callback = function()
				if last_pos ~= ShowCursorPos() then
					last_pos = ShowCursorPos()
					local scroll_offset = tostring(ShowCursorPos())
					local json = create_json_obj(scroll_offset, "scroll")
					SendServer(chan_id, json)
				end
			end
		})
	end
end

function SendServer(chan_id, object)
	local payload = object .. "\n"
	vim.fn.chansend(chan_id, payload)
end

function ShowCursorPos()
	local last_line = vim.fn.line("$")
	-- local current_line = vim.api.nvim_win_get_cursor(0)[1]
	local current_line = vim.fn.line("w0")
	local position = current_line/last_line

	return position
end


vim.api.nvim_create_autocmd({"VimLeavePre"}, {
	-- plenary has yet to address this
	-- https://github.com/nvim-lua/plenary.nvim/issues/156
	callback = function()
		vim.loop.kill(job.pid, 15)
	end
})

vim.api.nvim_create_autocmd("BufDelete", {
	buffer = vim.api.nvim_get_current_buf(),
	callback = function()
		print("[MDPreview] Buffer closed, stopping preview")
		vim.loop.kill(job.pid, 15)
	end,
})

function Entry()
	if job then
		print("[MDPreview] Failed! mdpreview is already running")
		return
	end

	job = Job:new({
		command = 'MDPreview',
		on_stdout = function(_, signal)
			if signal:find('sig_start') then
				vim.schedule(function()
					serveInit()
				end)
			end
		end
	})
	job:start()
end

return M
