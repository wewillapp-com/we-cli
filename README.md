![Wewillapp](https://wewillapp.com/images/wewillapp_logo.png)

# We CLI

a tool for scaffoling go_standard project in wewillapp

## Install

easy install with go install `go install go.wewillapp.com/we-cli@latest`

## Usage

`we-cli [command] [flags]`

available commands:

- `init` initialize new project _(incomming)_
- `create` create a resource
  - flags:
    - `--name, -n` resource name
    - `--type, -t` resource type _(model, response, form, resource)_
    - `--path, -p` resource path _(no need if create resource)_

avialable resources:

- `model` - create a model
- `response` - create a response
- `form` - create a form
- `resource` create all(model, response, form)

License: [Gnu3](./LICENSE)
