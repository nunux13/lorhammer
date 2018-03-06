package testtype

import (
	"encoding/json"
	"errors"
	"lorhammer/src/model"
	"lorhammer/src/orchestrator/command"
	"lorhammer/src/tools"
	"testing"
	"time"
)

func TestFake(t *testing.T) {
	command.NewLorhammer(model.NewLorhammer{CallbackTopic: "topic1"})
	err := Start(Test{testType: Type("Fake")}, []model.Init{{}}, nil)

	if err == nil {
		t.Fatal("Fake test should return unknown testType error")
	}
}

func TestNone(t *testing.T) {
	command.NewLorhammer(model.NewLorhammer{CallbackTopic: "topic1"})
	err := Start(Test{testType: typeNone}, []model.Init{{}}, nil)

	if err != nil {
		t.Fatalf("None test should not error : %s", err)
	}
}

func TestNewTester(t *testing.T) {
	callMeMaybe := make(chan error)
	testers["other"] = func(test Test, _ []model.Init, _ tools.Mqtt) {
		if test.repeatTime != time.Duration(1*time.Minute) {
			callMeMaybe <- errors.New("Test in json was 1m must be equal to diration 1 minute")
		} else if test.rampTime != time.Duration(1*time.Minute) {
			callMeMaybe <- errors.New("Test in json was 1m must be equal to diration 1 minute")
		} else {
			callMeMaybe <- nil
		}
	}

	var test Test
	if err := json.Unmarshal([]byte(`{"type": "other", "rampTime": "1m", "repeatTime": "1m"}`), &test); err != nil {
		t.Fatal("Unmarshalling test must work")
	}

	command.NewLorhammer(model.NewLorhammer{CallbackTopic: "topic1"})
	Start(test, []model.Init{{}}, nil)

	select {
	case res := <-callMeMaybe:
		if res != nil {
			t.Fatalf("Error in other func : %s", res)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("callMeMaybe must be called")
	}
}
