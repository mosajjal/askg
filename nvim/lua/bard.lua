-- main module file
local module = require("askg.module")

local M = {}
M.config = {
  -- default config
  askg_path = "$HOME/go/bin/askg",
  askg_config_path = "$HOME/.askg.yaml",
}

-- setup is the public method to setup your plugin
M.setup = function(args)
  -- you can define your setup function here. Usually configurations can be merged, accepting outside params and
  -- you can also put some validation here for those.
  M.config = vim.tbl_deep_extend("force", M.config, args or {})
end

-- "ask" is a public method for the plugin
M.ask = function(opts)
  module.askg(opts.args, M.config.askg_path, M.config.askg_config_path)
end

return M
