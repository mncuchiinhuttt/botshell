#!/usr/bin/env bash
set -e

# BotShell One-Line Installer
# Supported: macOS, Linux (amd64, arm64)

REPO="mncuchiinhuttt/botshell"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

printf "${BLUE}==>${NC} Installing BotShell...\n"

# 1. Detect OS
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$OS" in
  darwin) OS="darwin" ;;
  linux) OS="linux" ;;
  *)
    printf "${RED}Error:${NC} Unsupported OS: %s\n" "$OS"
    exit 1
    ;;
esac

# 2. Detect Architecture
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64|amd64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *)
    printf "${RED}Error:${NC} Unsupported architecture: %s\n" "$ARCH"
    exit 1
    ;;
esac

# 3. Resolve latest release
TARBALL="botshell_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/latest/download/${TARBALL}"

printf "${BLUE}==>${NC} Detected platform: ${GREEN}%s/%s${NC}\n" "$OS" "$ARCH"
printf "${BLUE}==>${NC} Downloading %s...\n" "$URL"

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

if command -v curl >/dev/null 2>&1; then
  curl -sSL "$URL" -o "$TMP_DIR/$TARBALL"
elif command -v wget >/dev/null 2>&1; then
  wget -q "$URL" -O "$TMP_DIR/$TARBALL"
else
  printf "${RED}Error:${NC} Neither curl nor wget found.\n"
  exit 1
fi

# 4. Extract
tar -xzf "$TMP_DIR/$TARBALL" -C "$TMP_DIR"

BIN_PATH="$TMP_DIR/${OS}_${ARCH}/botshell"
if [ ! -f "$BIN_PATH" ]; then
  # Fallback if flat
  BIN_PATH="$TMP_DIR/botshell"
fi

if [ ! -f "$BIN_PATH" ]; then
  printf "${RED}Error:${NC} Extracted binary not found.\n"
  exit 1
fi

chmod +x "$BIN_PATH"

# 5. Install to destination
printf "${BLUE}==>${NC} Installing to %s...\n" "$INSTALL_DIR"
if [ -w "$INSTALL_DIR" ]; then
  mv "$BIN_PATH" "$INSTALL_DIR/botshell"
else
  printf "${BLUE}==>${NC} Superuser permissions required to install to %s:\n" "$INSTALL_DIR"
  sudo mv "$BIN_PATH" "$INSTALL_DIR/botshell"
fi

# 6. Verify tmux
if ! command -v tmux >/dev/null 2>&1; then
  printf "\n${RED}Notice:${NC} 'tmux' was not found on your system.\n"
  if [ "$OS" = "darwin" ]; then
    printf "Please install it via Homebrew: ${GREEN}brew install tmux${NC}\n"
  else
    printf "Please install it via your package manager: ${GREEN}sudo apt install tmux${NC}\n"
  fi
fi

printf "\n${GREEN}Success!${NC} BotShell has been installed to %s/botshell\n" "$INSTALL_DIR"
printf "Run ${GREEN}botshell --setup${NC} to configure your Telegram bot.\n"
