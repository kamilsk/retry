package classifier_test

import (
	"errors"
	"testing"

	. "github.com/kamilsk/retry/v3/classifier"
)

func TestWhitelistClassifier_Classify(t *testing.T) {
	var (
		errInWhitelist    = errors.New("is in blacklist")
		errNotInWhitelist = errors.New("is not in blacklist")
	)
	list := WhitelistClassifier([]error{errInWhitelist})

	if list.Classify(nil) != Succeed {
		t.Error("succeed is expected")
	}

	if list.Classify(errNotInWhitelist) != Fail {
		t.Error("fail is expected")
	}

	if list.Classify(errInWhitelist) != Retry {
		t.Error("retry is expected")
	}
}
