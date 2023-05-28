-- module represents a lua module for the plugin
local M = {}

M.askbard = function(question, bardcli_path, bardcli_config)
    local cmd = "sh -c 'NO_COLOR=true " .. bardcli_path .. " -c " .. bardcli_config .. " " .. question .. " 2>&1'"
    local handler = io.popen(cmd)
    local result = handler:read("*all")
    local succeeded, error_msg, retcode  = handler:close()

    -- Create a new vertical split
    vim.cmd('vsplit')
    -- Switch to the new split
    vim.cmd('wincmd l')
    -- Set the new split to read-only
    vim.bo.readonly = true
    -- Move to the beginning of the buffer
    vim.cmd('enew')
    -- Switch to the new buffer
    vim.cmd('buffer')

    -- Set the buffer's filetype to markdown
    vim.bo.filetype = 'markdown'
    vim.cmd('normal! gg')
    -- Clear the buffer
    vim.cmd('normal! dG')

    -- Insert the output into the buffer
    for line in result:gmatch("[^\r\n]+") do
        vim.api.nvim_put({line}, "l", true, true)
    end
    -- Move back to the original split
    vim.cmd('wincmd h')

    return true
end

return M
