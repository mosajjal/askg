vim.api.nvim_create_user_command("Askbard", require("bard").ask, { nargs='?' })
