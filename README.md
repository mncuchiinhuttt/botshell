# BotShell

> Ultra-lightweight Cross-Platform Terminal & TUI remote controller via Telegram, powered by Golang & tmux.

BotShell turns Telegram into an interactive, zero-latency dashboard for remote terminal sessions. Whether you are running coding harnesses, long-running builds, Neovim, or AI agents on **macOS**, **Linux**, or **Windows**, BotShell gives you full control and real-time observability from your mobile phone with zero visual degradation.

---

## Key Highlights

- **Cross-Platform Power Management (Keep Awake):**
  - **macOS:** `/usr/bin/caffeinate -dism -w <PID>`
  - **Linux:** Native `systemd-inhibit --what=idle:sleep:handle-lid-switch`
  - **Windows:** Native Win32 `SetThreadExecutionState` (ES_CONTINUOUS | ES_SYSTEM_REQUIRED | ES_DISPLAY_REQUIRED)
  - Automatically restores normal system sleep state when the daemon terminates or exits.
- **Pure In-Memory ANSI Rasterizer:** Converts 140–160 column terminal outputs into high-resolution PNG images directly in memory (`golang.org/x/image`). TUI frames, syntax colors, and status bars render crisp on mobile without word-wrapping artifacts or terminal escape garbage.
- **Bi-directional tmux Sync:** Seamlessly connects to a live `tmux` session. Work from your machine at your desk, then walk away and continue controlling the exact same screen from Telegram.
- **Interactive Terminal TUI Dashboard (Bubbletea):**
  - Live activity logs view.
  - Multi-window screen previewer with window switching.
  - One-key Keep-Awake toggling (`a`).
  - Background daemonize support (`b`).
- **Dynamic Slash Commands:** Auto-registers Telegram's native command menu (`/screen`, `/text`, `/tabs`, `/newtab`, `/keys`) for direct autocomplete in the chat UI.
- **Interactive Control Keyboard:** Custom inline keyboard providing instant key presses (`Esc`, `Ctrl+O`, `Ctrl+C`, `Tab`, `Enter`, arrow navigation, `y`/`n`) and in-place screen refresh.

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

## Platform Support

| Operating System | Terminal Engine | Sleep Guard Mechanism | Support Status |
| :--- | :--- | :--- | :--- |
| **macOS** (Apple Silicon & Intel) | tmux (`brew install tmux`) | `caffeinate` | Tier 1 (Full) |
| **Linux** (Ubuntu, Debian, Fedora, Arch) | tmux (`apt`, `dnf`, `pacman`) | `systemd-inhibit` | Tier 1 (Full) |
| **Windows** (WSL2 / MSYS2) | tmux | Win32 `SetThreadExecutionState` | Supported |

---

## Quick Start

### 1. Interactive Setup Wizard
Run the setup wizard to configure your Telegram Bot Token and authorized Chat IDs:

```bash
./botshell --setup
```

The interactive wizard prompts for:
1. **Telegram Bot Token:** Obtained from [@BotFather](https://t.me/botfather).
2. **Allowed Chat/User IDs:** Comma-separated list of Telegram User/Chat IDs authorized to interact with the bot.
3. **Default tmux session name:** Session identifier to attach/create (default: `botshell`).
4. **Always-Awake:** Toggle background system power assertions.

Configuration is securely written to `~/.config/botshell/config.json` with restricted permissions (`0600`).

### 2. Start the Daemon
```bash
./botshell
```

Running `./botshell` launches the interactive **Bubbletea TUI Dashboard**:
- `Tab`: Switch between **Logs View** and **Screens View**.
- `a`: Toggle Always-Awake state on the fly.
- `Left` / `Right`: Cycle through tmux windows.
- `s`: Select active window.
- `b`: Run in background (detach daemon).
- `q`: Quit and cleanly shutdown.

To run silently without TUI:
```bash
./botshell --daemon
```

---

## Telegram Commands & Interactions

| Command | Description |
| :--- | :--- |
| `/screen`, `/s` | Capture current tmux pane and deliver as a high-resolution PNG image with inline control keys |
| `/text`, `/t` | Capture current terminal buffer as clean monospace text in a Markdown code block |
| `/tabs` | Display all tmux windows with instant switch buttons |
| `/newtab <name>` | Spawn a new window in the active session |
| `/keys <key>` | Send a named key combination (e.g., `/keys Escape`, `/keys C-c`, `/keys C-o`) |
| `/start` | View daemon uptime, session status, and command cheat sheet |

---

## Building from Source

```bash
git clone https://github.com/mncuchiinhuttt/botshell.git
cd botshell

# macOS / Linux native build
go build -o botshell .

# Cross-compile for Linux (from Mac/Windows)
GOOS=linux GOARCH=amd64 go build -o botshell-linux-amd64 .

# Cross-compile for Windows
GOOS=windows GOARCH=amd64 go build -o botshell.exe .
```

---

## License

MIT License. See [LICENSE](LICENSE) for details.
