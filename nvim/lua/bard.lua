-- main module file
local module = require("bard.module")

local M = {}
M.config = {
  -- default config
  bardcli_path = "$HOME/go/bin/bard-cli",
  bardcli_config_path = "$HOME/.bardcli.yaml",
}

-- setup is the public method to setup your plugin
M.setup = function(args)
  -- you can define your setup function here. Usually configurations can be merged, accepting outside params and
  -- you can also put some validation here for those.
  M.config = vim.tbl_deep_extend("force", M.config, args or {})
end

-- "ask" is a public method for the plugin
M.ask = function(opts)
  module.askbard(opts.args, M.config.bardcli_path, M.config.bardcli_config_path)
end

return M
