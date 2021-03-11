package data

import "fmt"

// CurrentGasPrice - When `gasz` service is queried, it'll send
// response of this form back
type CurrentGasPrice struct {
	Fast    float64 `json:"fast"`
	Fastest float64 `json:"fastest"`
	SafeLow float64 `json:"safeLow"`
	Average float64 `json:"average"`
}

func (c *CurrentGasPrice) String() string {

	return fmt.Sprintf("fastest : %.2f Gwei | fast : %.2f Gwei | average : %.2f Gwei | safeLow : %.2f Gwei", c.Fastest, c.Fast, c.Average, c.SafeLow)

}
