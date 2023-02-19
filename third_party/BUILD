
go_module(
    name = "warnings",
    module = "gopkg.in/warnings.v0",
    version = "v0.1.2",
    visibility = ["PUBLIC"],
)

go_module(
  name = "gcfg",
  install = [
    ".",
    "scanner",
    "token",
    "types",
  ],
  module = "github.com/please-build/gcfg",
  version = "v1.6.0",
  visibility = ["PUBLIC"],
  deps = [":warnings"],
)
