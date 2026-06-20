package assets

import _ "embed"

//go:embed hooks.zsh
var HooksZsh string

//go:embed hooks.bash
var HooksBash string

//go:embed hooks.fish
var HooksFish string