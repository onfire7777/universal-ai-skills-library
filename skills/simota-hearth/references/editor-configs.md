# Editor Configuration Patterns

Purpose: Read this when configuring `neovim`, `vim`, or `Zed`, or when you need plugin layout, `lazy.nvim`, keymap design, `blink.cmp`, or Neovim 0.10+ guidance.

## Contents

- [Neovim configuration](#neovim-configuration)
- [Vim configuration](#vim-configuration)
- [Keymap design principles](#keymap-design-principles)
- [Minimal vs full configuration](#minimal-vs-full-configuration)
- [Neovim 0.10 plus features](#neovim-010-plus-features)
- [blink.cmp](#blinkcmp)
- [Zed configuration](#zed-configuration)

## Neovim Configuration

### Directory Structure

```text
~/.config/nvim/
├── init.lua
├── lazy-lock.json
└── lua/
    ├── config/
    │   ├── autocmds.lua
    │   ├── keymaps.lua
    │   ├── lazy.lua
    │   └── options.lua
    └── plugins/
        ├── colorscheme.lua
        ├── completion.lua
        ├── editor.lua
        ├── lsp.lua
        ├── telescope.lua
```

### `init.lua` Entry Point

```lua
require("config.options")
require("config.lazy")

vim.api.nvim_create_autocmd("User", {
  pattern = "VeryLazy",
  callback = function()
    require("config.autocmds")
    require("config.keymaps")
  end,
})
```

### Core Options

```lua
local opt = vim.opt

opt.number = true
opt.relativenumber = true
opt.tabstop = 2
opt.shiftwidth = 2
opt.expandtab = true
opt.smartindent = true
opt.ignorecase = true
```

### `lazy.nvim` Bootstrap

```lua
local lazypath = vim.fn.stdpath("data") .. "/lazy/lazy.nvim"
if not vim.loop.fs_stat(lazypath) then
  vim.fn.system({
    "git", "clone", "--filter=blob:none",
    "https://github.com/folke/lazy.nvim.git",
    "--branch=stable", lazypath,
  })
end
vim.opt.rtp:prepend(lazypath)

require("lazy").setup("plugins", {
  defaults = { lazy = true },
  install = { colorscheme = { "catppuccin" } },
  checker = { enabled = true, notify = false },
})
```

### Key Mapping Design

```lua
local map = vim.keymap.set

vim.g.mapleader = " "
vim.g.maplocalleader = " "

map("n", "<C-h>", "<C-w>h", { desc = "Go to left window" })
map("n", "<C-j>", "<C-w>j", { desc = "Go to lower window" })
map("n", "<C-k>", "<C-w>k", { desc = "Go to upper window" })
map("n", "<C-l>", "<C-w>l", { desc = "Go to right window" })
```

### LSP Configuration

```lua
return {
  {
    "neovim/nvim-lspconfig",
    event = { "BufReadPre", "BufNewFile" },
    dependencies = {
      "mason.nvim",
      "mason-lspconfig.nvim",
    },
    config = function()
      vim.diagnostic.config({
        virtual_text = { spacing = 4, prefix = "●" },
        signs = true,
        underline = true,
      })
    end,
  },
}
```

### Essential Plugin Specs

```lua
return {
  "nvim-treesitter/nvim-treesitter",
  build = ":TSUpdate",
  event = { "BufReadPost", "BufNewFile" },
  opts = {
    ensure_installed = {
      "lua", "vim", "vimdoc", "query",
      "typescript", "javascript", "tsx",
      "python", "rust", "go",
      "json", "yaml", "toml", "markdown",
      "bash", "html", "css",
    },
    highlight = { enable = true },
    indent = { enable = true },
  },
}
```

## Vim Configuration

### Minimal `.vimrc`

```vim
set nocompatible
filetype plugin indent on
syntax enable

set number relativenumber
set cursorline
set signcolumn=yes
set scrolloff=8
set nowrap
set termguicolors

set tabstop=2 shiftwidth=2 expandtab
```

### `vim-plug` Setup

```vim
let data_dir = has('nvim') ? stdpath('data') . '/site' : '~/.vim'
if empty(glob(data_dir . '/autoload/plug.vim'))
  silent execute '!curl -fLo '.data_dir.'/autoload/plug.vim --create-dirs
    \ https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim'
  autocmd VimEnter * PlugInstall --sync | source $MYVIMRC
endif

call plug#begin()
Plug 'catppuccin/vim', { 'as': 'catppuccin' }
Plug 'tpope/vim-surround'
Plug 'tpope/vim-commentary'
Plug 'junegunn/fzf', { 'do': { -> fzf#install() } }
Plug 'junegunn/fzf.vim'
call plug#end()
```

## Keymap Design Principles

1. Use mnemonic prefixes such as `<leader>f`, `<leader>g`, and `<leader>l`.
2. Keep the same key for the same action across plugins when possible.
3. Prefer discoverability over density; use `which-key` or equivalent when available.
4. Check for conflicts before rebinding defaults.
5. Put advanced actions behind leader prefixes rather than overriding core motions.

### Common Leader Namespace

| Prefix | Domain | Examples |
|--------|--------|----------|
| `<leader>f` | Find / files | `ff`, `fg`, `fb` |
| `<leader>g` | Git | `gs`, `gb`, `gd` |
| `<leader>l` | LSP | `lr`, `la`, `ld` |
| `<leader>b` | Buffer | `bd`, `bp`, `bn` |
| `<leader>w` | Window | `wv`, `ws`, `wq` |
| `<leader>t` | Toggle | `tt`, `te`, `tz` |

## Minimal vs Full Configuration

| Aspect | Minimal | Standard | Full |
|--------|---------|----------|------|
| Plugin count | `0-5` | `10-20` | `30-50+` |
| LSP | None | Primary language | Multi-language |
| Completion | Built-in | `nvim-cmp` or `blink.cmp` | Completion plus AI |
| File finder | Built-in | `telescope` | `telescope` + extensions |
| Startup time | `<50ms` | `<100ms` | `<200ms` |
| Use case | Quick edits | Daily driver | IDE replacement |

## Neovim 0.10 Plus Features

### Native Snippets (`vim.snippet`)

```lua
vim.snippet.expand("function ${1:name}(${2:params})\n\t${0}\nend")

vim.keymap.set({ "i", "s" }, "<Tab>", function()
  if vim.snippet.active({ direction = 1 }) then
    return "<cmd>lua vim.snippet.jump(1)<cr>"
  end
  return "<Tab>"
end, { expr = true })
```

### OSC 52 Clipboard

```lua
vim.opt.clipboard = "unnamedplus"
```

This works over SSH and tmux on modern terminals without extra clipboard plugins.

### Built-in Comment Toggling

```text
gc{motion}  toggle comment
gcc         toggle current line
gc          toggle selected lines
```

Neovim 0.10+ can replace dedicated comment plugins for many setups.

## `blink.cmp`

### Comparison with `nvim-cmp`

| Aspect | `nvim-cmp` | `blink.cmp` |
|--------|------------|-------------|
| Architecture | Lua-based | Rust core |
| Performance | Good | Faster under many-source loads |
| Configuration | More verbose | Simpler |
| Snippet support | External | Built-in `vim.snippet` integration |
| Recommendation | Existing setups | New setups |

### Basic `blink.cmp` Configuration

```lua
return {
  "saghen/blink.cmp",
  version = "*",
  event = "InsertEnter",
  dependencies = {
    "rafamadriz/friendly-snippets",
  },
  opts = {
    keymap = {
      preset = "default",
      ["<C-space>"] = { "show", "show_documentation", "hide_documentation" },
      ["<C-e>"] = { "cancel", "fallback" },
      ["<CR>"] = { "accept", "fallback" },
      ["<Tab>"] = { "snippet_forward", "select_next", "fallback" },
    },
  },
}
```

Choose `blink.cmp` for new setups. Keep `nvim-cmp` when the existing setup depends on sources not yet replaced.

## Zed Configuration

Zed is supported as a lightweight alternative. Hearth still treats Neovim as the primary recommendation for power users.

### `settings.json`

```jsonc
{
  "theme": {
    "mode": "dark",
    "dark": "Catppuccin Mocha",
    "light": "Catppuccin Latte"
  },
  "buffer_font_family": "JetBrains Mono",
  "buffer_font_size": 14,
  "ui_font_size": 14
}
```

### `keymap.json`

```jsonc
[
  {
    "context": "Editor && VimNormal",
    "bindings": {
      "space f f": "file_finder::Toggle",
      "space f g": "search::SelectAllMatches",
      "space e": "workspace::ToggleLeftDock"
    }
  }
]
```

### Built-in LSP

```jsonc
{
  "lsp": {
    "rust-analyzer": {
      "initialization_options": {
        "check": { "command": "clippy" }
      }
    }
  }
}
```
