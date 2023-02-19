# clerc

Clerc is a simple utility for [Please build system](https://please.build) which extracts the compilation database (`compile_commands.json`) for the LSP tools
like `ccls` and `clangd`.

**This tools is very much in development, PRs and improvement ideas are welcome.**

The compilation database can be used for various tools, most notably for the LSPs like `clangd` to provide accurate information to the editor.

For detailed description have read through the Resources.

### How to run clerc?

```
 clerc //dir/target:cc_target > compile_commands.json 
```

## Resources
 - [Eli Bendersky's outline of the compilation databases](https://eli.thegreenplace.net/2014/05/21/compilation-databases-for-clang-based-tools)
 - [Papin's outline of the compilation databases](https://sarcasm.github.io/notes/dev/compilation-database.html)

