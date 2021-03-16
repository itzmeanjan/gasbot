package bot

import (
	"errors"
	"strconv"
	"strings"
)

// parseSubscriptionPayload - Attempts to parse payload of `/subscribe` command
func parseSubscriptionPayload(payload string) (string, string, float64, error) {

	splitted := strings.Split(payload, " ")
	if !(len(splitted) == 3) {
		return "", "", 0.0, errors.New("bad payload received")
	}

	threshold, err := strconv.ParseFloat(splitted[2], 64)
	if err != nil {
		return "", "", 0.0, errors.New("bad threshold received")
	}

	return splitted[0], splitted[1], threshold, nil

}
