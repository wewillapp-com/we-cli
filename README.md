![Wewillapp](https://wewillapp.com/images/wewillapp_logo.png)

# We CLI

a tool for scaffoling go_standard project in wewillapp

## Install
easy install with go install `go install github.com/wewillapp-com/we-cli@latest`

## Usage
`we-cli [command] [flags]`

available commands:
  - `init` initialize new project *(incomming)*
  - `create` create a resource
    - flags:
      - `--name, -n` resource name
      - `--type, -t` resource type *(model, response, form, resource)*
      - `--path, -p` resource path *(no need if create resource)*

avialable resources:
  - `model` - create a model
  - `response` - create a response
  - `form` - create a form
  - `resource` create all(model, response, form)


License:  [Gnu3](./LICENSE)
