package klog

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type duration int

type Duration interface {
	InMinutes() int
	Plus(Duration) Duration
	Minus(Duration) Duration
	ToString() string
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func NewDuration(amountHours int, amountMinutes int) Duration {
	return duration(amountHours*60) + duration(amountMinutes)
}

func (d duration) InMinutes() int {
	return int(d)
}

func (d duration) Plus(additional Duration) Duration {
	return NewDuration(0, d.InMinutes()+additional.InMinutes())
}

func (d duration) Minus(deductible Duration) Duration {
	return NewDuration(0, d.InMinutes()-deductible.InMinutes())
}

func (d duration) ToString() string {
	if d == 0 {
		return "0m"
	}
	hours := abs(int(d) / 60)
	minutes := abs(int(d) % 60)
	result := ""
	if int(d) < 0 {
		result += "-"
	}
	if hours > 0 {
		result += fmt.Sprintf("%dh", hours)
	}
	if minutes > 0 {
		result += fmt.Sprintf("%dm", minutes)
	}
	return result
}

var durationPattern = regexp.MustCompile(`^(-|\+)?((\d+)h)? ?((\d+)m)?$`)

func NewDurationFromString(hhmm string) (Duration, error) {
	match := durationPattern.FindStringSubmatch(hhmm)
	if match == nil {
		return nil, errors.New("MALFORMED_DURATION")
	}
	sign := 1
	if match[1] == "-" {
		sign = -1
	}
	if match[3] == "" && match[5] == "" {
		return nil, errors.New("MALFORMED_DURATION")
	}
	hours, _ := strconv.Atoi(match[3])
	minutes, _ := strconv.Atoi(match[5])
	if minutes > 60 {
		return nil, errors.New("UNREPRESENTABLE_DURATION")
	}
	return NewDuration(sign*hours, sign*minutes), nil
}
