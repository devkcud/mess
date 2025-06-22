# 🧹 mess

**mess** is your fast, friendly command-line helper for whipping up directories and files. It keeps your project neat-ish without arcane syntax or ceremony.

Compared to older tools like `mk`, mess is more intuitive, flexible, and built for actual humans (also fixed a ton of old bugs mk had).

## 🚀 Usage

```sh
mess [-flags] <..|dir/|dir/file|file>...
```

### 📐 Behavior Rules

- `dir/` → Creates a directory and adds it to the stack. Everything created afterward goes inside it.
- `file` → Creates a file in the current stack location.
- `dir/file` → Creates the specified directory and file, but does not push the directory to the stack.
- `..` → Pops the last directory off the stack. Back up one level like a well-behaved script.

### 🧩 Flags

- `-h` or `--help`: The "what does this flag do?" menu.
- `-b <dir>` or `--base <dir>`: Set the base working directory (default: your current pwd).
- `-d` or `--dry`: Dry run mode. No files harmed, just simulated structure.
- `-e` or `--echo`: Print out shell commands instead of creating anything. Similar to dry run, but less pretty.
- `--loglevel <0-4>`: How chatty should it be?
  - `0`: 😶 Error only
  - `1`: ⚠️ Warnings
  - `2`: ℹ️ Info
  - `3`: 🐛 Debug
  - `4`: 🧵 Trace everything. Yes, everything. Almost.

## 🛠️ Examples

### 📄 Create a file

```sh
mess hello.txt
```

Drops `hello.txt` in the base directory. Easy.

### 🗂️ Nested directory + file

```sh
mess src/lib/components/Button.svelte
```

Creates `src/lib/components/` if it doesn't exist and create the file `Button.svelte`.

### ⬅️ Back up with ..

```sh
mess project/ docs/ README.md .. src/ index.js
```

Creates `project/docs/README.md`, goes up one level, and `project/src/index.js`.

Alternatively, you could simply do:

```sh
mess project/ docs/README.md src/index.js
```

Note: `dir/file` creates the file but does not push the directory to the stack.

### 🫥 Dry run

```sh
~ $ mess -d notes/ day-1.md day-2.md day-3.md
```

Sends you:

```
/home/<user>/notes/
├── day-1.md
├── day-2.md
└── day-3.md
```

### 🎭 Echo mode

```sh
~ $ mess -e cli/ cmd/goon/main.go internal/ modules/download.go testing/framework.go .. pkg/utils/commands.go
```

Should spit out:

```sh
mkdir -p /home/<user>/cli/cmd/goon
mkdir -p /home/<user>/cli/internal/modules
mkdir -p /home/<user>/cli/internal/testing
mkdir -p /home/<user>/cli/pkg/utils
touch /home/<user>/cli/cmd/goon/main.go
touch /home/<user>/cli/internal/modules/download.go
touch /home/<user>/cli/internal/testing/framework.go
touch /home/<user>/cli/pkg/utils/commands.go
```

## ✨ Why mess?

Because file and folder creation should be fast, flexible, and slightly entertaining. **mess** helps you build structure without building a headache.

## 📜 License

[The Unlicense](LICENSE): use it, break it, improve it. No strings attached.
