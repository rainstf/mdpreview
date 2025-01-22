local async = require'plenary.async'
local Job = require'plenary.job'

local SIGTERM = 15
local M = {}

local defaults = {
	server = {
		conn_timeout  = 500,
		conn_attempts = 10,
	},
	events = {
		reload_on = 'InsertLeave',
		scrolling = false,
	},
}
M.config = defaults

vim.api.nvim_create_user_command('MdPreviewStart', 'lua Entry()', {})
vim.api.nvim_create_user_command('MdPreviewStop', 'lua Stop()', {})

function M.setup(user_options)
	local function merge(defaults, overrides)
		for k, v in pairs(overrides) do
			if type(defaults[k]) == "table" and type(v) == "table" then
				merge(defaults[k], v)
			else
				defaults[k] = v
			end
		end
	end
	user_options = user_options or {}
	merge(M.config, user_options)
end

function BufContent(buf)
	local lines = vim.api.nvim_buf_get_lines(buf, 0, -1, false)
	local content = table.concat(lines, '\n')

	return content
end

function create_json_obj(content, evtype)
	local object = {
		['content' ] = content,
		['event' ] = evtype
	}
	if evtype == 'scroll' or evtype == 'reload' or evtype == 'init' then
		return vim.json.encode(object)
	end
end

function send_initial_content(chan_id)
	local json = create_json_obj(BufContent(0), 'init')
	server_send(chan_id, json)
end

function is_markdown()
	return vim.filetype.match({ buf = 0 }) == 'markdown'
end

function server_connect(address)
	local chan_id = vim.fn.sockconnect('tcp', address)
	send_initial_content(chan_id)

	-- transport events 
	vim.api.nvim_create_autocmd(M.config.events.reload_on, {
		callback = function()
			if not is_markdown() then
				print('MDPreview Failed! Reason: not a markdown file')
				return
			end
			local json = create_json_obj(BufContent(0), 'reload')
			server_send(chan_id, json)
		end
	})
	vim.api.nvim_create_autocmd('TextChanged', {
		callback = function()
			if not is_markdown() then
				print('MDPreview Failed! Reason: not a markdown file')
				return
			end
			local json = create_json_obj(BufContent(0), 'reload')
			server_send(chan_id, json)
		end
	})
	local last_pos = 0
	if M.config.events.scrolling == true then
		vim.api.nvim_create_autocmd({ 'CursorMoved', 'CursorMovedI' }, {
			callback = function()
				if last_pos ~= fetch_cursor_pos() then
					last_pos = fetch_cursor_pos()
					local scroll_offset = tostring(fetch_cursor_pos())
					local json = create_json_obj(scroll_offset, 'scroll')
					server_send(chan_id, json)
				end
			end
		})
	end
end

function server_send(chan_id, object)
	local payload = object .. '\n'
	vim.fn.chansend(chan_id, payload)
end

function fetch_cursor_pos()
	local last_line = vim.fn.line('$')
	local current_line = vim.fn.line('w0')
	local position = current_line/last_line
	return position
end

function Entry()
	if job then
		print('MDPreview Failed! Reason: mdpreview is already running')
		return
	end
	if not is_markdown() then
		print('MDPreview Failed! Reason: not a markdown file')
		return
	end

	job = Job:new({
		command = 'MDPreview',
		on_stdout = function(_, signal)
			if signal:find('sig_start') then
				vim.schedule(function()
					server_connect('localhost:8080')
				end)
			end
		end
	})
	job:start()

	vim.api.nvim_create_autocmd({'VimLeavePre'}, {
		-- plenary has yet to address this
		-- https://github.com/nvim-lua/plenary.nvim/issues/156
		callback = function()
			vim.loop.kill(job.pid, SIGTERM)
		end
	})
	-- kill server when buffer used to spawn it dies
	vim.api.nvim_create_autocmd('BufDelete', {
		buffer = vim.api.nvim_get_current_buf(),
		callback = function()
			print('MDPreview Shutdown: Buffer was closed')
			vim.loop.kill(job.pid, SIGTERM)
		end,
	})
end

function Stop()
    if job then
        vim.loop.kill(job.pid, SIGTERM)
        print("MDPreview stopped")
        job = nil
    else
        print("MDPreview is not running.")
    end
end

return M
