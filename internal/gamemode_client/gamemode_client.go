package gamemode_client

import (
	"fmt"
	"os"

	"github.com/godbus/dbus/v5"
)

func RequestStartFor(pid int32) error {
	ret, err := performRequest[int]("RegisterGameByPID", pid)
	if err != nil {
		return fmt.Errorf("performing request: %v", err)
	}

	switch *ret {
	case 0:
		// Success
		return nil
	case -1:
		return fmt.Errorf("request failed")
	case -2:
		return fmt.Errorf("request rejected")
	default:
		return fmt.Errorf("unknown return code: %d", *ret)
	}
}

type GamemodeStatus int

const (
	GamemodeInactive               GamemodeStatus = 0
	GamemodeActive                 GamemodeStatus = 1
	GamemodeActiveClientRegistered GamemodeStatus = 2
)

func (s *GamemodeStatus) IsActive() bool {
	return *s == GamemodeActive || *s == GamemodeActiveClientRegistered
}

func QueryStatus(pid int32) (GamemodeStatus, error) {
	ret, err := performRequest[int]("QueryStatusByPID", pid)
	if err != nil {
		return GamemodeInactive, fmt.Errorf("performing request: %v", err)
	}

	switch *ret {
	case 0:
		return GamemodeInactive, nil
	case 1:
		return GamemodeActive, nil
	case 2:
		return GamemodeActiveClientRegistered, nil
	case -1:
		return GamemodeInactive, fmt.Errorf("query failed")
	default:
		return GamemodeInactive, fmt.Errorf("unknown return code: %d", *ret)
	}
}

func performRequest[T any](method string, pid int32) (*T, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, fmt.Errorf("connecting to D-Bus session bus: %v", err)
	}

	currentPid := os.Getpid()

	obj := conn.Object("com.feralinteractive.GameMode", "/com/feralinteractive/GameMode")
	call := obj.Call("com.feralinteractive.GameMode."+method, 0, pid, currentPid)
	if call.Err != nil {
		return nil, fmt.Errorf("calling %s: %v", method, call.Err)
	}

	var result T
	err = call.Store(&result)
	if err != nil {
		return nil, fmt.Errorf("storing result of %s: %v", method, err)
	}

	return &result, nil
}
