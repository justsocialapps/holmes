package justlib

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/justsocialapps/assert"
)

func TestMain(m *testing.M) {
	flag.Parse()
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestTryZeroTimesShouldCallFunctionIndefinitelyAndEventuallyReturn(t *testing.T) {
	a := assert.NewAssert(t)
	cnt := 0
	rand.Seed(time.Now().UnixNano())
	successAfterNTries := rand.Intn(100000)
	err := Try(0, 0, func() error {
		cnt++
		if cnt == successAfterNTries {
			return nil
		}
		return errors.New("Some error")
	})
	a.Equal(err, nil, "Try returned with an error")
	a.Equal(cnt, successAfterNTries, "The callback wasn't called enough times")
}

func TestTryTwoTimesShouldCallPassingFunctionOnlyOnce(t *testing.T) {
	a := assert.NewAssert(t)
	called := 0
	err := Try(2, 0, func() error {
		called++
		return nil
	})
	a.Equal(called, 1, fmt.Sprintf("f has been called %d times", called))
	a.Equal(err, nil, "Try returned with an error")
}

func TestTryTwoTimesShouldCallFailingFunctionTwoTimes(t *testing.T) {
	a := assert.NewAssert(t)
	called := 0
	err := Try(2, 0, func() error {
		called++
		if called == 1 {
			return errors.New("Warning. Warp reactor core primary coolant failure")
		}
		return nil
	})
	a.Equal(called, 2, fmt.Sprintf("f has been called %d times", called))
	a.Equal(err, nil, "Try returned with an error")
}
