# Gemini (Bard) CLI

Use Gemini in your CLI!

## Installation

`askg` has now a client/server architecture. `askgd` is the daemon that runs and always listens on `localhost:12345` for RPC requests. `askg` is the client that sends the requests to the daemon.

### Install the daemon

```sh
go install github.com/mosajjal/askg/cmd/askgd@latest
```

in order to run the daemon, first run the `browser` subcommand to open a browser and login to `gemini.google.com`. currently only Chrome is supported. Make sure that you already open the browser and login to `gemini.google.com` before running the `browser` subcommand.

```sh
$ askgd browser
```

then run the daemon using the following command

```sh
$ askgd run
```

the daemon will now listen on `localhost:12345` for RPC requests, it also rotates and commits cookies automatically. 


### Install the client

```sh
go install github.com/mosajjal/askg/cmd/askg@latest
```

```md
$ askg "what is the meaning of life?"
The meaning of life is a question that has been asked by philosophers and theologians for centuries. There is no one answer that will satisfy everyone, but some possible answers include:

* To find happiness and fulfillment.
* To make a difference in the world.
* To learn and grow as a person.
* To connect with others and build relationships.
* To experience the beauty of the world.
* To leave a legacy behind.

Ultimately, the meaning of life is up to each individual to decide. There is no right or wrong answer, and what matters most is that you find something that gives your life meaning.

Here are some additional thoughts on the meaning of life:

* The meaning of life is not something that is given to us, but something that we create.
* The meaning of life is not something that we find once and for all, but something that we discover and rediscover throughout our lives.
* The meaning of life is not something that is the same for everyone, but something that is unique to each individual.

If you are searching for the meaning of life, I encourage you to explore your own values, beliefs, and experiences. What is important to you? What makes you happy? What do you want to achieve in your life? The answers to these questions may help you to find your own meaning in life.
```

## Use a Proxy

To use a HTTP(s) or SOCKS4/5 proxy to access Google, set `HTTP_PROXY` as well as `HTTPS_PROXY` environment variables before running `askg`. 


## Use as a neovim plugin

in your neovim config, add the following plugin

```lua
Plug 'mosajjal/askg', {'rtp': 'nvim'}
```

by default, the plugin looks for `askg` in `$HOME/go/bin/askg` and the configuration file at `$HOME/.askg.yaml`

to change that, run the setup function of Gemini using the following

```lua
lua require('askg').setup({askg_path="$HOME/go/bin/askg", askg_config_path="$HOME/.askg.yaml"})
```

## Usage

`
:Askg "write hello world in javascript"
`

the above will open a new vsplit and return the results in the new buffer

![interactive](static/neovim.png)
