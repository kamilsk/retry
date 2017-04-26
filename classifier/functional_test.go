package classifier

import (
	"encoding/json"
	"errors"
	"net"
	"testing"
)

func TestFunctionalClassifier_Classify(t *testing.T) {
	var (
		errClassified       = &json.SyntaxError{}
		errNotClassified    = errors.New("is unknown error")
		jsonErrorClassifier = FunctionalClassifier(func(err error) Action {
			if err == nil {
				return Succeed
			}

			if _, is := err.(*json.SyntaxError); is {
				return Retry
			}

			return Unknown
		})
	)

	if jsonErrorClassifier.Classify(nil) != Succeed {
		t.Error("succeed is expected")
	}

	if jsonErrorClassifier.Classify(errClassified) != Retry {
		t.Error("retry is expected")
	}

	if jsonErrorClassifier.Classify(errNotClassified) != Unknown {
		t.Error("unknown is expected")
	}
}

func TestFunctionalClassifier_NetworkErrorClassifier_Classify(t *testing.T) {
	var (
		errNetworkTimeout = &net.DNSError{IsTimeout: true}
		errNetworkOther   = &net.DNSError{}
		errOther          = errors.New("is not network error")
	)

	if NetworkErrorClassifier.Classify(nil) != Succeed {
		t.Error("succeed is expected")
	}

	if NetworkErrorClassifier.Classify(errNetworkTimeout) != Retry {
		t.Error("retry is expected")
	}

	if NetworkErrorClassifier.Classify(errNetworkOther) != Fail {
		t.Error("fail is expected")
	}

	if NetworkErrorClassifier.Classify(errOther) != Unknown {
		t.Error("unknown is expected")
	}
}
