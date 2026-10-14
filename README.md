# BotShell

> Ultra-lightweight Cross-Platform Terminal & TUI remote controller via Telegram, powered by Golang & tmux.

BotShell turns Telegram into an interactive, zero-latency dashboard for remote terminal sessions. Whether you are running coding harnesses, long-running builds, Neovim, or AI agents on **macOS**, **Linux**, or **Windows**, BotShell gives you full control and real-time observability from your mobile phone with zero visual degradation.

---

## One-Line Quick Install

### macOS & Linux
Cài đặt tự động bản binary mới nhất vào `/usr/local/bin` chỉ bằng một câu lệnh:

```bash
curl -fsSL https://raw.githubusercontent.com/mncuchiinhuttt/botshell/main/install.sh | bash
```

Hoặc tải script về chạy trực tiếp:
```bash
wget -qO- https://raw.githubusercontent.com/mncuchiinhuttt/botshell/main/install.sh | bash
```

### Windows
Tải file zip `botshell_windows_amd64.zip` từ [Releases](https://github.com/mncuchiinhuttt/botshell/releases/latest), giải nén và chạy `botshell.exe --setup`.

---

## Key Highlights

- **Cross-Platform Power Management (Keep Awake):**
  - **macOS:** `/usr/bin/caffeinate -dism -w <PID>`
  - **Linux:** Native `systemd-inhibit --what=idle:sleep:handle-lid-switch`
  - **Windows:** Native Win32 `SetThreadExecutionState` (ES_CONTINUOUS | ES_SYSTEM_REQUIRED | ES_DISPLAY_REQUIRED)
  - Tự động khôi phục chế độ sleep mặc định khi tiến trình daemon dừng.
- **Pure In-Memory ANSI Rasterizer:** Chuyển đổi khung viền terminal rộng 140–160 cột thành ảnh PNG siêu nét trực tiếp trong bộ nhớ (`golang.org/x/image`), không làm vỡ giao diện TUI trên điện thoại.
- **Bi-directional tmux Sync:** Kết nối 2 chiều với session `tmux`. Ngồi tại máy làm việc bình thường, ra ngoài mở điện thoại tương tác tiếp.
- **Interactive Terminal TUI Dashboard (Bubbletea & Lipgloss):**
  - Giao diện trực quan hiển thị Activity Logs theo thời gian thực.
  - Chuyển tab xem preview nội dung từng screen / window trong tmux.
  - Phím tắt `a` bật/tắt Always-Awake tức thì, `u` check update, `b` chuyển chạy ngầm background.
- **Dynamic Slash Commands:** Tự động đăng ký menu lệnh Telegram (`/screen`, `/text`, `/tabs`, `/newtab`, `/keys`, `/update`, `/start`).
- **Interactive Control Keyboard:** Bàn phím nút bấm điều khiển tức thời (`Esc`, `Ctrl+O`, `Ctrl+C`, `Tab`, `Enter`, mũi tên điều hướng, `y`/`n`).

---

## Quick Start

### 1. Cấu hình ban đầu
```bash
botshell --setup
```

Wizard sẽ hỏi:
1. **Telegram Bot Token:** Lấy từ [@BotFather](https://t.me/botfather).
2. **Allowed Chat/User IDs:** ID Telegram được phép truy cập (ngăn chặn người lạ).
3. **Default tmux session name:** Tên session (mặc định: `botshell`).
4. **Always-Awake:** Bật chế độ giữ máy không sleep khi chạy.

Cấu hình được lưu tại: `~/.config/botshell/config.json`.

### 2. Khởi chạy
```bash
botshell
```

Chạy ở chế độ chạy ngầm (không mở TUI):
```bash
botshell --daemon
```

---

## Telegram Commands

| Lệnh | Chức năng |
| :--- | :--- |
| `/screen`, `/s` | Chụp màn hình terminal dạng ảnh PNG sắc nét kèm bàn phím điều khiển |
| `/text`, `/t` | Chụp terminal dưới dạng text code block |
| `/tabs` | Danh sách các cửa sổ tmux và nút chuyển tab |
| `/newtab <tên>` | Mở window mới trong session |
| `/keys <key>` | Gửi phím đặc biệt (ví dụ `/keys Escape`, `/keys C-c`) |
| `/update` | Kiểm tra phiên bản mới trên GitHub |
| `/start` | Xem thông tin trạng thái bot & trợ giúp |

---

## Pre-built Bundles & Releases

Tất cả các bản dựng bundle (`.tar.gz`, `.zip`) kèm file checksum SHA256 đều được phát hành tại [GitHub Releases](https://github.com/mncuchiinhuttt/botshell/releases/latest):

- `botshell_darwin_arm64.tar.gz` (Apple Silicon M1/M2/M3/M4)
- `botshell_darwin_amd64.tar.gz` (Intel Mac)
- `botshell_linux_amd64.tar.gz` (Linux x86_64)
- `botshell_linux_arm64.tar.gz` (Linux aarch64 / Raspberry Pi / ARM VPS)
- `botshell_windows_amd64.zip` (Windows 64-bit)

---

## License

MIT License. Copyright (c) 2026 @mncuchiinhuttt.
