# reset_spotlight

A Go program to fully reset the macOS Spotlight index when search stops working correctly — apps not appearing, stale results, or Spotlight behaving unexpectedly.

---

## Requirements

- macOS (tested on 15.x Sequoia)
- Go 1.18 or later
- `sudo` access

---

## Installation

**Clone or download `main.go`, then build:**

```bash
go build -o reset-spotlight main.go
```

Optionally move it to your PATH so you can run it from anywhere:

```bash
sudo mv reset-spotlight /usr/local/bin/
```

---

## Usage

```bash
sudo ./reset-spotlight
```

Or if installed in your PATH:

```bash
sudo reset-spotlight
```

> The program checks for root privileges at startup and exits immediately with a clear message if not run with `sudo`.

---

## What It Does

| Step | Command | Description |
|------|---------|-------------|
| 1 | `mdutil -a -i off` | Disables Spotlight indexing on all volumes |
| 2 | `sleep 10s` | Waits for `mds` to settle between state changes |
| 3 | `mdutil -a -i on` | Re-enables indexing on all volumes |
| 4 | `mdutil -a -E` | Erases and schedules a full index rebuild |
| 5 | `launchctl kickstart -k system/com.apple.metadata.mds` | Restarts the `mds` metadata daemon |

Each step streams its output directly to the terminal in real time. If any step fails, the program prints the error and exits immediately — it does not silently continue.

---

## After Running

Spotlight rebuilds its index in the background. Depending on the size of your disk and number of files, this can take anywhere from **a few minutes to an hour**.

You can monitor progress via:

- **Activity Monitor** — watch `mds_stores` CPU usage; it will settle when indexing is complete
- **Terminal:**
  ```bash
  sudo mdutil -a -s
  ```
  Look for `Indexing enabled.` on your volumes.

---

## Troubleshooting

**"Index is already changing state" error**
> `mds` is mid-transition. Wait 10–15 seconds and run the program again. If it persists, force a clean restart first:
> ```bash
> sudo killall mds
> ```

**"Spotlight server is disabled" error**
> The `mds` daemon isn't running. Restart it manually:
> ```bash
> sudo launchctl kickstart -k system/com.apple.metadata.mds
> ```
> Then re-run the program.

**Spotlight still not working after rebuild**
> Try excluding and re-adding your disk in Spotlight's privacy list:
> **System Settings → Siri & Spotlight → Spotlight Privacy**
> Add your drive, wait a moment, then remove it to force a fresh index.

---

## Notes

- This program does not touch Spotlight's privacy/exclusion list. Any folders manually excluded will remain excluded after the reset.
- On volumes with SIP (System Integrity Protection) restrictions, some operations may require booting into Recovery Mode.

---

## License

MIT — do whatever you want with it.