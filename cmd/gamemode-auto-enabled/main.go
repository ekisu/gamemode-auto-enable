package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ekisu/gamemode-auto-enable/internal/gamemode_client"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/v4/process"
)

func main() {
	setupLogger()

	log.Info().Msg("Starting gamemode-auto-enabled")

	for {
		enableGamemodeIfNeeded()

		time.Sleep(5 * time.Second)
	}
}

func setupLogger() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	var err error
	level := zerolog.InfoLevel
	overrideLevel, ok := os.LookupEnv("LOG_LEVEL")
	if ok {
		level, err = zerolog.ParseLevel(overrideLevel)
		if err != nil {
			log.Warn().Err(err).Str("LOG_LEVEL", overrideLevel).Msg("invalid LOG_LEVEL, defaulting to INFO")
			level = zerolog.InfoLevel
		}
	}

	zerolog.SetGlobalLevel(level)
}

func enableGamemodeIfNeeded() {
	steamProcess, err := findSteamProcess()
	if err != nil {
		log.Warn().Err(err).Msg("couldn't find steam process")

		return
	}

	if steamProcess == nil {
		log.Trace().Msg("steam process not found")

		return
	}

	gameProcesses, err := findGameProcessesRecursively(steamProcess)
	if err != nil {
		log.Warn().Err(err).Msg("couldn't find game processes")

		return
	}

	if len(gameProcesses) == 0 {
		log.Trace().Msg("no game processes found")

		return
	}

	for _, gameProcess := range gameProcesses {
		pid := gameProcess.Pid

		status, err := gamemode_client.QueryStatus(pid)
		if err != nil {
			log.Error().Err(err).Int32("pid", pid).Msg("failed to check gamemode status for process")

			continue
		}

		if status.IsActive() {
			log.Trace().Int32("pid", pid).Int("status", int(status)).Msg("gamemode already active for process")

			continue
		}

		log.Info().Int32("pid", pid).Msg("new game process found, attempting to gamemode")

		err = gamemode_client.RequestStartFor(pid)
		if err != nil {
			log.Error().Err(err).Int32("pid", pid).Msg("failed to enable gamemode for process")
		} else {
			log.Info().Int32("pid", pid).Msg("gamemode enabled for process")
		}
	}
}

func findSteamProcess() (*process.Process, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("listing processes: %v", err)
	}

	for _, p := range processes {
		name, err := p.Name()
		if err != nil {
			continue
		}

		if name == "steam" {
			return p, nil
		}
	}

	return nil, nil
}

func findGameProcessesRecursively(p *process.Process) ([]*process.Process, error) {
	children, err := p.Children()
	if err != nil {
		return nil, fmt.Errorf("getting process children: %v", err)
	}

	log.Trace().Int("count", len(children)).Msg("got child processes")

	var gameProcesses []*process.Process
	for _, child := range children {
		exe, err := child.Exe()
		if err != nil {
			log.Trace().Err(err).Int32("pid", child.Pid).Msg("couldn't get child executable path")
			continue
		}

		log.Trace().Int32("pid", child.Pid).Str("exe", exe).Msg("checking steam child process")

		if isGameExecutable(exe) {
			gameProcesses = append(gameProcesses, child)

			// No need to go deeper, as games shouldn't have game processes as children
			continue
		}

		subGameProcesses, err := findGameProcessesRecursively(child)
		if err != nil {
			log.Trace().Err(err).Int32("pid", child.Pid).Msg("couldn't get game processes from child")
			continue
		}

		gameProcesses = append(gameProcesses, subGameProcesses...)
	}

	return gameProcesses, nil
}

func isGameExecutable(exe string) bool {
	// TODO implement a proper logic here
	return strings.Contains(exe, "/steamapps/common/")
}
