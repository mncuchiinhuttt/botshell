# BotShell

> Ultra-lightweight macOS Terminal & TUI remote controller via Telegram, powered by Golang & tmux.

BotShell turns Telegram into an interactive, zero-latency dashboard for remote terminal sessions running on macOS. Whether you are running coding harnesses, long-running builds, Neovim, or AI agents, BotShell gives you full control and real-time observability from your mobile phone with zero visual degradation.

---

## Key Highlights

- **Pure In-Memory ANSI Rasterizer:** Converts 140–160 column terminal outputs into high-resolution PNG images directly in memory (`golang.org/x/image`). TUI frames, syntax colors, and status bars render crisp on mobile without word-wrapping artifacts or terminal escape garbage.
- **Bi-directional tmux Sync:** Seamlessly connects to a live `tmux` session. Work from your Mac at your desk, then walk away and continue controlling the exact same screen from Telegram.
- **Dynamic Slash Commands:** Auto-registers Telegram's native command menu (`/screen`, `/text`, `/tabs`, `/newtab`, `/keys`) so you get auto-complete directly inside the chat UI.
- **Interactive Control Keyboard:** Custom inline keyboard providing instant key presses (`Esc`, `Ctrl+O`, `Ctrl+C`, `Tab`, `Enter`, arrow navigation, `y`/`n`) and in-place screen refresh.
- **Automatic macOS Sleep Prevention (caffeinate):** Leverages `/usr/bin/caffeinate -dism -w <PID>` to keep your Mac awake while BotShell is active. Automatically restores default sleep policies when the process terminates or exits.
- **Zero Third-Party Binary Dependencies:** Aside from standard `tmux`, the entire ANSI parsing, rasterizing, and bot dispatch logic is self-contained in a single, compact Go binary.

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
|                          BotShell Daemon (Go)                         |
|                                                                       |
|  [ Auth & Security ]       [ Caffeinate Guard ]   [ ANSI Rasterizer ] |
|  Restricts chat/user IDs   Holds power assertion  Parses ANSI colors  |
|  via config whitelist      bound to process PID   Renders 2D PNG grid |
+-----------------------------------+-----------------------------------+
                                    | Unix CLI
                                    v
+-----------------------------------------------------------------------+
|                       tmux multiplexer (macOS)                        |
|                                                                       |
|  [Window 0: Agent Harness]   [Window 1: nvim]   [Window 2: build]     |
+-----------------------------------------------------------------------+
```

---

## Prerequisites

- **macOS:** Tested on macOS 14+ / Apple Silicon & Intel.
- **tmux:** Install via Homebrew if not already installed:
  ```bash
  brew install tmux
  ```
- **Go:** 1.24+ (only required if building from source).

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
4. **macOS Always-Awake:** Toggle background system power assertions.

Configuration is securely written to `~/.config/botshell/config.json` with restricted permissions (`0600`).

### 2. Start the Daemon
```bash
./botshell
```

*(If executed without an existing configuration file, BotShell automatically boots the setup wizard).*

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

### Direct Typing
Sending any standard text message to the bot will pipe that string directly into the active tmux pane (`send-keys -l "<input>" Enter`), followed by an immediate screen refresh.

---

## Security Model

- **Whitelist Enforcement:** Every update is checked against the configured `allowed_chat_ids`. Unlisted IDs receive no responses and actions are dropped before execution.
- **Shell Escaping:** Input strings are forwarded using `tmux send-keys -l` (literal mode), preventing command injection or unwanted shell expansion of metacharacters (`$`, `&`, `;`, `"`).
- **Process Isolation:** Power management assertions are strictly bound to the daemon PID; terminating the daemon immediately relinquishes system sleep holds.

---

## Building from Source

```bash
git clone https://github.com/mncuchiinhuttt/botshell.git
cd botshell
go build -o botshell .
```

To run test suites:
```bash
go test -v ./...
```

---

## License

MIT License. See [LICENSE](LICENSE) for details.
