# Gemini (Bard) CLI Neovim plugin

## Installation

in your neovim config, add the following plugin

```lua
Plug 'mosajjal/askg', {'rtp': 'nvim'}
```

by default, the plugin looks for `askg` in `$HOME/go/bin/askg`.

to change that, run the setup function of Gemini using the following

```lua
lua require('askg').setup({askg_path="$HOME/go/bin/askg"})
```

## Usage

`
:Askg "write hello world in javascript"
`

the above will open a new vsplit and return the results in the new buffer

