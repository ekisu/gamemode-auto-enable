package gamemode_client

// #cgo CFLAGS: -I${SRCDIR}
// #include "gamemode_client_wrapper.h"
//
// int gamemode_request_start_for_wrapper(int pid);
// const char* gamemode_error_string_wrapper();
import "C"

import (
	"fmt"
	"os/exec"
	"strings"
)

func Toggle(pid int32) error {
	cmd := exec.Command("gamemoded", fmt.Sprintf("--request=%d", pid))
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("toggling gamemode for pid %d: %v", pid, err)
	}

	return nil
}

func IsActive(pid int32) (bool, error) {
	cmd := exec.Command("gamemoded", fmt.Sprintf("--status=%d", pid))
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("checking gamemode status: %v", err)
	}

	outputStr := string(output)

	if strings.Contains(outputStr, "gamemode is active") {
		return true, nil
	} else if strings.Contains(outputStr, "gamemode is inactive") {
		return false, nil
	} else {
		return false, fmt.Errorf("unexpected gamemode status output: %s", outputStr)
	}
}
