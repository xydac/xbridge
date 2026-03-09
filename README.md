# xbridge

A zero-config build proxy for your Mac mini.

Single binary. No dependencies. Point it at your Xcode project. Build and test from any machine on your network.

![demo](docs/demo.gif)

---

## What it does

A small HTTP server that runs on your Mac and lets you trigger iOS builds, manage simulators, and grab screenshots — all from your Linux box, a CI script, or even `curl`.

## Install

```bash
# Homebrew (builds from source, no Gatekeeper warnings)
brew tap xydac/xbridge
brew install xbridge

# Or download the binary directly
VERSION=$(gh release view --repo xydac/xbridge --json tagName -q .tagName)
curl -fsSL "https://github.com/xydac/xbridge/releases/download/${VERSION}/xbridge_${VERSION#v}_darwin_arm64.tar.gz" | tar xz
sudo mv xbridge /usr/local/bin/

# Or build from source
go install github.com/xydac/xbridge/cmd/xbridge@latest
```

## Quick Start

### On the Mac (server)

```bash
cd ~/projects/my-app
xbridge serve
# → Server running on :7900
# → Detected MyApp.xcworkspace (scheme: MyApp)
# → Ready.
```

### On Linux (client)

```bash
# Point it at your Mac
xbridge remote set mac:7900

# Build
xbridge build
xbridge build --watch          # stream logs
xbridge build --profile staging

# Simulators
xbridge sim list
xbridge sim boot "iPhone 16"
xbridge screenshot

# Interact with the simulator
xbridge tap 220 400
xbridge swipe 220 800 220 200
xbridge text "hello world"
xbridge key 1                      # home button

# The full combo: pull → build → install → launch
xbridge run
```

### With curl

```bash
curl http://mac:7900/health
curl -X POST http://mac:7900/build
curl http://mac:7900/simulators/BOOTED/screenshot > screen.png
```

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/health` | Server status, Xcode version, disk space |
| `POST` | `/build` | Trigger build (returns job ID) |
| `POST` | `/build/clean` | Clean build folder |
| `GET` | `/build/:id` | Build status + errors |
| `GET` | `/build/:id/logs` | Build logs (SSE stream) |
| `GET` | `/build/:id/artifact` | Download .app bundle |
| `GET` | `/simulators` | List available simulators |
| `POST` | `/simulators/boot` | Boot a simulator |
| `POST` | `/simulators/shutdown` | Shutdown a simulator |
| `GET` | `/simulators/:udid/screenshot` | Take screenshot (PNG) |
| `POST` | `/simulators/:udid/install` | Install .app |
| `POST` | `/simulators/:udid/launch` | Launch app by bundle ID |
| `POST` | `/simulators/:udid/openurl` | Open a deep link |
| `POST` | `/simulators/:udid/tap` | Tap at (x, y) coordinates |
| `POST` | `/simulators/:udid/swipe` | Swipe from (x1,y1) to (x2,y2) |
| `POST` | `/simulators/:udid/text` | Type text into focused field |
| `POST` | `/simulators/:udid/key` | Send a key press |
| `GET` | `/simulators/:udid/logs` | App logs (SSE stream) |
| `GET` | `/git/status` | Branch, commit, dirty state |
| `POST` | `/git/pull` | Pull latest |
| `POST` | `/git/checkout` | Switch branch |

Use `BOOTED` as the UDID to target the currently booted simulator.

## Configuration

### Zero config (auto-detect)

```bash
xbridge serve
# Scans for .xcworkspace or .xcodeproj, picks first scheme
```

### CLI flags

```bash
xbridge serve --port 7900 --project ./MyApp.xcworkspace --scheme MyApp --key "secret"
```

### Config file (`xbridge.yaml`)

```yaml
project: ./MyApp.xcworkspace
scheme: MyApp
configuration: Debug

simulator:
  device: "iPhone 16 Pro"
  runtime: "iOS 18.2"

server:
  port: 7900
  key: "my-secret-key"

hooks:
  pre_build: "pod install"

profiles:
  staging:
    configuration: Debug
    build_args:
      - "SWIFT_ACTIVE_COMPILATION_CONDITIONS=STAGING"
    simulator:
      device: "iPhone 16 Pro Max"
    env:
      API_URL: "https://staging.example.com"
```

Precedence: CLI flags > xbridge.yaml > auto-detect.

## Multi-project

Just run multiple instances:

```bash
cd ~/projects/main-app && xbridge serve --port 7900
cd ~/projects/admin-panel && xbridge serve --port 7901
```

```bash
xbridge remote add main mac:7900
xbridge remote add admin mac:7901
xbridge build --remote main
xbridge build --remote admin
```

## Run as a service

```bash
xbridge install-service --project ~/projects/my-app --port 7900
# Creates ~/Library/LaunchAgents/com.xbridge.my-app.plist
```

## Development

```bash
make build       # build binary
make test        # run tests
make run         # go run ./cmd/xbridge serve
make build-all   # cross-compile darwin-arm64, darwin-amd64, linux-amd64
```

## Requirements

- **Xcode** — for building and simulators
- **idb** (optional) — required for `tap`, `swipe`, `text`, and `key` commands. Install with `brew install idb-companion`.

## Tech Stack

- **Go 1.22+** — single binary, no runtime deps
- **chi** — lightweight stdlib-compatible router
- **cobra** — CLI framework
- **gorilla/websocket** — WebSocket for log streaming
- 4 external modules total. Everything else is stdlib.

## License

MIT
