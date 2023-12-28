package work

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gen2brain/beeep"
)

type selectedWorkTask struct {
	task string
}

func selectTask(selectedTask string) tea.Cmd {
	return func() tea.Msg {
		return selectedWorkTask{task: selectedTask}
	}
}

type selectedWorkTimeFromList struct {
	time time.Duration
}

func selectTimeFromList(selectedTime time.Duration) tea.Cmd {
	return func() tea.Msg {
		return selectedWorkTimeFromList{time: selectedTime}
	}
}

type selectedWorkTimeFromInput struct {
	time time.Duration
}

func selectTimeFromInput(selectedTime time.Duration) tea.Cmd {
	return func() tea.Msg {
		return selectedWorkTimeFromInput{time: selectedTime}
	}
}

type selectedCustomTime struct {
	time  time.Duration
	error error
}

func selectCustomTime(selectedTime string) tea.Cmd {
	return func() tea.Msg {
		value, err := time.ParseDuration(selectedTime)
		if err != nil {
			return selectedCustomTime{error: err}
		}

		return selectedCustomTime{time: value}
	}
}

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type workFinished struct {
	error error
}

func finishWork(time time.Duration) tea.Cmd {
	return func() tea.Msg {
		var finalError error
		if err := beeep.Beep(44000, 10000); err != nil {
			finalError = fmt.Errorf("failed to notify with a beep that work finished: %w", err)
		}

		if err := beeep.Notify(
			"Work Pilot: Work Finished!",
			fmt.Sprintf("Congratulations! You've worked for %d minute(s).", time),
			"",
		); err != nil {
			finalError = fmt.Errorf("%w: failed to notify with a notification that work finished: %w", finalError, err)
		}

		return workFinished{error: finalError}
	}
}
