## What's Changed

**BotShell v1.0.0** is the first official release of the ultra-lightweight, cross-platform terminal & TUI remote controller via Telegram, powered by Golang & `tmux`.

### Core Features

* **In-Memory ANSI Terminal Rasterizer (`internal/render`)**: Converts 140–160 column terminal screens directly to crisp PNG images in memory with full TrueColor/256-color support. Eliminates mobile word-wrapping and ANSI garbage.
* **Bi-directional `tmux` Session Synchronization**: Seamlessly connects to any active local `tmux` session. Work locally at your desk, then view and control the exact same session on the go.
* **Cross-Platform Sleep Guard (Always-Awake)**:
  * **macOS**: Native `/usr/bin/caffeinate -dism -w <PID>` bound directly to daemon process lifecycle.
  * **Linux**: Native `systemd-inhibit --what=idle:sleep:handle-lid-switch`.
  * **Windows**: Native Win32 `SetThreadExecutionState`.
* **Interactive Terminal TUI Dashboard (`bubbletea` & `lipgloss`)**:
  * Real-time activity logs view.
  * Multi-window screen previewer with instant window switching (`Left`/`Right`, `s`).
  * One-touch Always-Awake toggle (`a`).
  * Manual release update checker (`u`).
  * Background daemonize trigger (`b`).
* **Rich Telegram User Experience**:
  * Auto-registered native slash command autocomplete (`/screen`, `/text`, `/tabs`, `/newtab`, `/keys`, `/update`, `/start`).
  * Dynamic interactive control keypad (`Esc`, `Ctrl+O`, `Ctrl+C`, `Tab`, `Enter`, arrows, `y`/`n`).
  * In-place message edits for smooth, flicker-free screen updates.
* **Security & Configuration**:
  * Interactive CLI setup wizard (`--setup`).
  * Strict Telegram chat ID whitelist enforcement.
  * Literal key injection (`tmux send-keys -l`) preventing shell command injection.

---

**Author**: [@mncuchiinhuttt](https://github.com/mncuchiinhuttt)  
**Full Changelog**: https://github.com/mncuchiinhuttt/botshell/commits/v1.0.0
