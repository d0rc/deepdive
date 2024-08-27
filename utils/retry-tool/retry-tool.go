package retry_tool

import "time"

var DefaultRetrySleep = time.Second * 1

func RetryCallWithCount[T any](f func() (T, error), retryCount uint64) (T, error) {
	var err error
	var res T
	for i := uint64(0); i < retryCount; i++ {
		res, err = f()
		if err == nil {
			return res, nil
		}
		time.Sleep(DefaultRetrySleep)
	}
	return res, err
}
