# BAM! - Binary Asset Manager: Full Plan

---

## Directory Structure

```
~/.bam/
  shims/                  <- IN PATH (high priority)
    node.exe
    pnpm.exe
    go.exe
  installs/               <- extracted tool versions
    node/
      22.1.0/
        bin/node.exe
        bin/npm.exe
    pnpm/
      9.5.0/
  cache/                  <- downloaded archives (kept for reinstall)
    node/
      node-v22.1.0-win-x64.zip
    pnpm/
      pnpm-v9.5.0.tar.gz
  plugins/
    builtin/              <- embedded in bam binary via go:embed
      node.toml
      pnpm.toml
      bun.toml
      deno.toml
      go.toml
    user/                 <- user-added plugins
  versions/               <- global version pins (one file per tool)
    node                  <- contains "22.1.0"
    pnpm                  <- contains "9.5.0"
  .priority               <- version resolution order
  install.log             <- full record of every touched file/PATH/dir
```

---

## Plugin TOML Format

```toml
name = "node"
aliases = ["nodejs"]
description = "Node.js JavaScript runtime"
schema = 2

[versions]
list_url = "https://nodejs.org/dist/index.json"
list_path = "$[*].version"                      # JSONPath to extract version list
strip_prefix = "v"

[download]
url = "https://nodejs.org/dist/v{version}/node-v{version}-{os}-{arch}.{ext}"
hash_url = "https://nodejs.org/dist/v{version}/SHASUMS256.txt"
hash_algo = "sha256"
hash_format = "gnu"                                                          # "gnu" = "hash  file" | "bsd" = "SHA256 (file) = hash"

[platform]
os_map = { windows = "win", linux = "linux", darwin = "darwin" }
arch_map = { amd64 = "x64", arm64 = "arm64" }
ext_map = { windows = "zip", linux = "tar.gz", darwin = "tar.gz" }

[install]
strip_components = true

[[install.bin]]
name = "node"
run = { windows = ["node.exe"], linux = ["bin/node"], darwin = ["bin/node"] }

[[install.bin]]
name = "npm"
run = { windows = [
	"node.exe",
	"node_modules/npm/bin/npm-cli.js",
], linux = [
	"bin/node",
	"lib/node_modules/npm/bin/npm-cli.js",
], darwin = [
	"bin/node",
	"lib/node_modules/npm/bin/npm-cli.js",
] }

[[install.bin]]
name = "npx"
run = { windows = [
	"node.exe",
	"node_modules/npm/bin/npx-cli.js",
], linux = [
	"bin/node",
	"lib/node_modules/npm/bin/npx-cli.js",
], darwin = [
	"bin/node",
	"lib/node_modules/npm/bin/npx-cli.js",
] }
```

---

## Shim Architecture

The generic shim is a **separate minimal Go program** (its own module), compiled per platform and embedded into the main `bam` binary via `go:embed`.

**At install time** (`bam install node@22.1.0`):

- bam extracts the embedded generic shim binary
- Copies it to `~/.bam/shims/node.exe`, `~/.bam/shims/npm.exe`, etc.
- The shim itself is stateless - no version baked in
- Resolve current OS and writes the relevant run list

```json
{
	"node": { "tool": "node", "run": ["node.exe"] },
	"npm": {
		"tool": "node",
		"run": ["node.exe", "node_modules/npm/bin/npm-cli.js"]
	},
	"npx": {
		"tool": "node",
		"run": ["node.exe", "node_modules/npm/bin/npx-cli.js"]
	}
}
```

**At invocation time** (user runs `node`):

1. Shim reads `os.Executable()` -> gets its own name (`node`)
2. Reads `~/.bam/.priority` for resolution order
3. Walks up CWD checking for version files in that order
4. Falls back to `~/.bam/versions/node` (global pin)
5. Resolves full path: `~/.bam/installs/node/22.1.0/bin/node.exe`
6. `syscall.Exec` (Unix) / `windows.CreateProcess` (Windows) - replaces itself, zero extra process overhead

---

## Version Resolution

**Formats accepted:**

- `22.1.0` - exact
- `22.x` or `22` - latest matching major
- `latest`, `lts`, `beta`, `nightly` - channel aliases defined per plugin
- `*` - absolute latest

**Global resolution:** pinned at install time, stored in `~/.bam/versions/<tool>`

**Local resolution:** at shim invocation time, walks up dirs

---

## `.priority` File

Controls which version source wins. User-editable plaintext:

```
global
.bam
package.json
.node-version
.nvmrc
.go-version
```

---

## `.bam` File Format

```
node=22.x
pnpm=latest
go=1.22.x
```

---

## `install.log` (Uninstall Hygiene)

JSON record of everything bam has ever touched:

```json
{
	"path_modifications": [
		{ "file": "~/.bashrc", "added": "export PATH=~/.bam/shims:$PATH", "timestamp": "..." }
	],
	"shims_created": ["~/.bam/shims/node.exe"],
	"dirs_created": ["~/.bam/installs/node/22.1.0"],
	"installed_tools": [{ "tool": "node", "version": "22.1.0", "installed_at": "..." }]
}
```

`bam uninstall --all` reads this log and reverses everything cleanly.

---

## CLI Commands

```
bam install <tool>[@<version>]         install a tool (defaults to latest)
bam install <tool>[@<version>] --use   install a tool (defaults to latest) and pins it. currently pins regardless
bam uninstall <tool>[@<version>]       remove a version
bam uninstall --all                    full nuke using install.log
bam use <tool>@<version>               set global pin
bam list                               show all installed tools + versions
bam list <tool>                        show available versions of a tool
bam which <tool>                       show resolved binary path for current dir
bam update                             update bam itself
bam plugin add <path/url>              add a user plugin
bam plugin list                        list all plugins (builtin + user)
bam env                                print bam dirs, PATH status, active versions
bam setup                              run PATH setup (also runs on first install)
```

---

## Built-in Plugins (Phase 1)

`node`, `npm` (bundled with node), `pnpm`, `yarn`, `bun`, `deno`, `go`

---

## Bootstrap Installer

A **single tiny Go binary** (or shell script fallback) that:

1. Detects OS/arch
2. Downloads the correct `bam` binary from GitHub releases
3. Places it in `~/.bam/`
4. Runs `bam setup` (creates dirs, handles PATH)

`bam setup` asks the user clearly:

> "Add `~/.bam/shims` to PATH in `~/.bashrc`? [Y/n]"
> Then also prints the manual line regardless.

---

## Known Plugin Format TODOs

Issues discovered during implementation - fix when building Component 10 (built-in plugins).

| #   | Issue                                 | Detail                                                                                                                                                                                                                                                                                                              |
| --- | ------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1   | Platform-specific `bin` paths         | Windows node has no `bin/` - needs `bin_map` like `os_map` e.g. `bin_map = { windows = ["node.exe"], linux = ["bin/node"] }`                                                                                                                                                                                        |
| 2   | Hash file format field                | Currently hardcoded to `{hash}  {filename}` (GNU format). Add `hash_format` to TOML to support other layouts                                                                                                                                                                                                        |
| 3   | `strip_prefix` on version list        | Defined in TOML but not yet used in resolver - wire it up                                                                                                                                                                                                                                                           |
| 4   | pnpm has no hash file                 | `hash_url = ""` should be handled gracefully (warn + proceed)                                                                                                                                                                                                                                                       |
| 5   | Imperative plugin support - undecided | I have not yet decided if a comprehensive declarative TOML system is enough to handle possible weirdness other tools may have. Running un-compiled go code is not feasible.                                                                                                                                         |
| 6   | Global packages support               | Right now our toml format defines a path and args to run, meaning we can define node, npm and npx for one tool. How do we handle global packages? 1. Node installs them INTO the installs folder so its a) not pristine and b) now tied to version and 2. How do we, for example, make a shim for `tsc` or similar? |

---

## Implementation Order

| #   | Component                       | Notes                                              |
| --- | ------------------------------- | -------------------------------------------------- |
| 1   | Core dirs + config bootstrap    | First-run setup, OS/arch detection                 |
| 2   | Plugin loader                   | TOML parsing, builtin + user discovery             |
| 3   | Version resolver                | latest/x patterns, fetch+cache version list        |
| 4   | Downloader + cache + hash check | Check cache before download, verify hash           |
| 5   | Extractor                       | tar.gz, tar.xz, zip, strip_components              |
| 6   | Install manager                 | Wire 3-5, write to installs/, update versions/     |
| 7   | Shim binary                     | Separate Go module, embedded into bam              |
| 8   | Shim generator                  | Copy+rename shim on install                        |
| 9   | PATH manager                    | Shell detection, profile modification, install.log |
| 10  | Built-in plugins                | node, pnpm, bun, deno, go TOMLs                    |
| 11  | Local version resolution        | .bam, package.json, .priority, dir walking         |
| 12  | Full CLI                        | All commands wired up                              |
| 13  | Bootstrap installer             | Standalone tiny binary                             |
