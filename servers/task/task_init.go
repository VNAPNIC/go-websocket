/**
* Created by GoLand.
* User: nankai
* Date: 2019-07-29
* Time: 14:15
 */

package task

import "time"

type TimerFunc func(interface{}) bool

/**
 * Scheduled call
 * @delay first delay
 * @tick interval
 * @fun timed execution function
 * @param fun parameter
 */
func Timer(delay, tick time.Duration, fun TimerFunc, param interface{}, funcDefer TimerFunc, paramDefer interface{}) {
	go func() {
		defer func() {
			if funcDefer != nil {
				funcDefer(paramDefer)
			}
		}()

		if fun == nil {
			return
		}

		t := time.NewTimer(delay)
		defer t.Stop()

		for {
			select {
			case <-t.C:
				if fun(param) == false {
					return
				}
				t.Reset(tick)
			}
		}
	}()
}
