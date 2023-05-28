# Bard CLI Neovim plugin

## Installation

in your neovim config, add the following plugin

```lua
Plug 'mosajjal/bard-cli', {'rtp': 'nvim'}
```

by default, the plugin looks for `bard-cli` in `$HOME/go/bin/bard-cli` and the configuration file at `$HOME/.bardcli.yaml`

to change that, run the setup function of bard using the following

```lua
lua require('bard').setup({bardcli_path="$HOME/go/bin/bard-cli", bardcli_config_path="$HOME/.bardcli.yaml"})
```

## Usage

`
:Askbard "write hello world in javascript"
`

the above will open a new vsplit and return the results in the new buffer

