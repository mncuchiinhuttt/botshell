# BotShell

> Ultra-lightweight Cross-Platform Terminal & TUI remote controller via Telegram, powered by Golang & tmux.

BotShell turns Telegram into an interactive, zero-latency dashboard for remote terminal sessions. Whether you are running coding harnesses, long-running builds, Neovim, or AI agents on **macOS**, **Linux**, or **Windows**, BotShell provides full control and real-time observability from your mobile device with zero visual degradation.

---

## One-Line Quick Install

### macOS & Linux
Automatically install the latest release binary into `/usr/local/bin`:

```bash
curl -fsSL https://raw.githubusercontent.com/mncuchiinhuttt/botshell/main/install.sh | bash
```

Or using `wget`:
```bash
wget -qO- https://raw.githubusercontent.com/mncuchiinhuttt/botshell/main/install.sh | bash
```

### Windows
Download `botshell_windows_amd64.zip` from [Releases](https://github.com/mncuchiinhuttt/botshell/releases/latest), extract the archive, and run:
```powershell
.\botshell.exe --setup
```

---

## Key Highlights

- **Cross-Platform Sleep Guard (Keep Awake):**
  - **macOS:** Native `/usr/bin/caffeinate -dism -w <PID>` bound to the daemon PID.
  - **Linux:** Native `systemd-inhibit --what=idle:sleep:handle-lid-switch`.
  - **Windows:** Native Win32 `SetThreadExecutionState` (`ES_CONTINUOUS | ES_SYSTEM_REQUIRED | ES_DISPLAY_REQUIRED`).
  - Automatically restores default system power policies when the process terminates.
- **Pure In-Memory ANSI Rasterizer:** Converts 140–160 column terminal outputs into sharp PNG images directly in memory (`golang.org/x/image`), preventing word-wrapping artifacts and rendering clean on mobile screens.
- **Bi-directional tmux Synchronization:** Connects seamlessly to any local `tmux` session. Work locally at your desk, then continue controlling the exact same screen from Telegram on the go.
- **Interactive Terminal TUI Dashboard (Bubbletea & Lipgloss):**
  - Live activity logs view.
  - Multi-window screen previewer with window switching.
  - One-key Keep-Awake toggling (`a`).
  - Background daemonize mode (`b`).
  - Manual update checker (`u`).
- **Seamless Slash Command Routing:**
  - Auto-registered native Telegram slash command autocomplete.
  - Intelligent pass-through: any command not claimed by BotShell (such as `/exit`, `/help`, `/model`, `/compact`) is automatically forwarded directly to `omp` or your active terminal process.
  - Escape prefix: prefix any command with `//` (e.g. `//exit`) to explicitly force transmission to the terminal.
- **Interactive Control Keyboard:** Custom inline keyboard offering instant key presses (`Esc`, `Ctrl+O`, `Ctrl+C`, `Tab`, `Enter`, arrow navigation, `y`/`n`) with in-place image refresh.

---

## Architecture

```
                  +-----------------------------------+
                  |           Telegram Client         |
                  |  - Slash command autocomplete     |
                  |  - High-res PNG screen feed       |
                  |  - Inline control keyboard        |
                  +-----------------+-----------------+
                                    | Telegram Bot API
                                    v
+-----------------------------------------------------------------------+
|                       BotShell Daemon (Go)                            |
|                                                                       |
|  [ Auth & Security ]       [ OS Power Guard ]     [ ANSI Rasterizer ] |
|  Restricts chat/user IDs   macOS: caffeinate      Parses ANSI colors  |
|  via config whitelist      Linux: systemd-inhibit Renders 2D PNG grid |
|                            Windows: Win32 API                         |
+-----------------------------------+-----------------------------------+
                                    | Unix / CLI
                                    v
+-----------------------------------------------------------------------+
|                    tmux multiplexer (macOS / Linux)                   |
|                                                                       |
|  [Window 0: Agent Harness]   [Window 1: nvim]   [Window 2: build]     |
+-----------------------------------------------------------------------+
```

---

## Quick Start

### 1. Initial Configuration
Run the setup wizard:
```bash
botshell --setup
```

The interactive wizard prompts for:
1. **Telegram Bot Token:** Obtained from [@BotFather](https://t.me/botfather).
2. **Allowed Chat/User IDs:** Comma-separated list of authorized Telegram chat/user IDs.
3. **Default tmux session name:** Session identifier to attach/create (default: `botshell`).
4. **Always-Awake:** Toggle background system power assertions.

Configuration is securely written to `~/.config/botshell/config.json` with restricted permissions (`0600`).

### 2. Launching BotShell

Run with the interactive TUI Dashboard:
```bash
botshell
```

Run in silent background daemon mode:
```bash
botshell --daemon
```

---

## Telegram Commands

| Command | Description |
| :--- | :--- |
| `/screen`, `/s` | Capture current tmux pane and deliver as a high-resolution PNG image with inline control keys |
| `/text`, `/t` | Capture current terminal buffer as clean monospace text in a Markdown code block |
| `/tabs` | Display all tmux windows with instant switch buttons |
| `/newtab <name>` | Spawn a new window in the active session |
| `/keys <key>` | Send a named key combination (e.g., `/keys Escape`, `/keys C-c`) |
| `/update` | Check for newer releases on GitHub |
| `/start` | View daemon uptime, session status, and command cheat sheet |
| `//<cmd>` | Explicitly escape and forward any slash command directly to the terminal (e.g. `//exit`) |

---

## TUI Keyboard Shortcuts

| Key | Action |
| :--- | :--- |
| `Tab` | Switch between **Logs View** and **Screens View** |
| `a` | Toggle **Always-Awake** mode immediately |
| `u` | Check for GitHub release updates |
| `Left` / `Right` | Cycle through tmux windows |
| `s` | Select active tmux window |
| `b` | Detach and run daemon in background |
| `q` / `Ctrl+C` | Clean shutdown |

---

## Pre-built Bundles & Releases

Pre-compiled binary packages and SHA256 checksums are available on [GitHub Releases](https://github.com/mncuchiinhuttt/botshell/releases/latest):

- `botshell_darwin_arm64.tar.gz` (Apple Silicon M1/M2/M3/M4)
- `botshell_darwin_amd64.tar.gz` (Intel Mac)
- `botshell_linux_amd64.tar.gz` (Linux x86_64)
- `botshell_linux_arm64.tar.gz` (Linux aarch64 / Raspberry Pi / ARM VPS)
- `botshell_windows_amd64.zip` (Windows 64-bit)

---

## License

MIT License. Copyright (c) 2026 @mncuchiinhuttt.
