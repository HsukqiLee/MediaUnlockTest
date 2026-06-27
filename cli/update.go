package main

import (
	core "MediaUnlockTest/pkg/core"
)

func checkUpdate(force bool) {
	cfg := core.UpdateConfig{
		AppName:         "unlock-test",
		VersionURL:      "https://unlock.icmp.ing/test/latest/version",
		BinaryURLPrefix: "https://unlock.icmp.ing/test/latest/unlock-test",
		Silent:          false,
		ForceUpdate:     force,
	}
	core.CheckUpdate(cfg)
}
