package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func main() {
	if os.Geteuid() != 0 {
		fmt.Fprintln(os.Stderr, "✗ This program must be run as root. Try: sudo reset-spotlight")
		os.Exit(1)
	}

	fmt.Println("🔍 Resetting Spotlight index...")

	fmt.Println("Disabling Spotlight indexing...")

	cmd := exec.Command("mdutil", "-a", "-i", "off")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "✗ Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("→ Waiting for mds to settle...")
	time.Sleep(10 * time.Second)

	fmt.Println("Re-enabling Spotlight indexing...")

	cmd = exec.Command("mdutil", "-a", "-i", "on")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "✗ Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Erasing and rebuilding index...")

	cmd = exec.Command("mdutil", "-a", "-E")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "✗ Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Restarting mds daemon...")

	cmd = exec.Command("launchctl", "kickstart", "-k", "system/com.apple.metadata.mds")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "✗ Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Done. Spotlight is rebuilding in the background.")
	fmt.Println("   This may take a few minutes depending on your disk size.")

}
