package pd

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pingcap/log"
	"go.uber.org/zap"
)

const (
	schedulerPrefix = "pd/api/v1/schedulers"
)

var (
	// Schedulers is the random disturbance schedulers.
	Schedulers = []string{
		"shuffle-leader-scheduler",
		"shuffle-region-scheduler",
		"shuffle-hot-region-scheduler",
	}
)

// ClientExt is a mixin with some PD functions for http.Client.
type ClientExt struct {
	*http.Client
}

// DefaultClient is the default PD client.
var DefaultClient = ClientExt{http.DefaultClient}

func isOk(status int) bool {
	return status/100 == 2
}

// EnableScheduler enable the specficated schedulers.
func (c ClientExt) EnableScheduler(pdAddrs []string, schedulers ...string) error {
	for _, pd := range pdAddrs {
		for _, scheduler := range schedulers {
			body := bytes.NewBuffer([]byte(fmt.Sprintf(`{"name": "%s"}`, strings.ReplaceAll(scheduler, `"`, `\"`))))
			url := fmt.Sprintf("http://%s/%s", pd, schedulerPrefix)
			req, err := http.NewRequest(http.MethodPost, url, body)
			if err != nil {
				return err
			}
			result, err := c.Do(req)
			if err != nil {
				return err
			}
			if !isOk(result.StatusCode) {
				failBody, _ := ioutil.ReadAll(result.Body)
				_ = result.Body.Close()
				log.Error("failed to add scheduler",
					zap.String("status", result.Status),
					zap.String("response", string(failBody)),
					zap.String("address", url))
				return errors.New("pd response with non-success code")
			}
		}
	}
	return nil
}
