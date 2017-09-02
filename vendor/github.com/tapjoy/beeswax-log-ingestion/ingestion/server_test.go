package ingestion

import (
	"fmt"
	"testing"

	newrelic "github.com/newrelic/go-agent"
)

var (
	hostStub        = "http://host"
	cpmStub         = "http://cpm"
	newrelicStub, _ = newrelic.NewApplication(newrelic.NewConfig("test", "abcdefghijklmnopqrstuvwxyzabcdefghijklmn"))
)

func TestNewServer(t *testing.T) {
	cases := []struct {
		tag       string
		host, cpm string
		opts      []Option
		verify    func(*server) error
	}{
		{"default", hostStub, cpmStub, []Option{}, validDefaults},
		{"options: timeouts", hostStub, cpmStub, []Option{ReadTimeout(2), WriteTimeout(3), IdleTimeout(4)}, validTimeouts},
		{"options: newrelic", hostStub, cpmStub, []Option{NewRelic(&newrelicStub)}, validNewRelic},
	}

	for _, c := range cases {
		t.Run(c.tag, func(t *testing.T) {
			s, err := NewServer(c.host, c.cpm, c.opts...)
			if err != nil {
				t.Fatal(err)
			}
			if err := c.verify(s); err != nil {
				t.Error(err)
			}
		})
	}
}

func validDefaults(s *server) error {
	if s.httpServer.Addr != hostStub {
		return fmt.Errorf("got %s, expected %s", s.httpServer.Addr, hostStub)
	}
	if s.cpmAddr != cpmStub {
		return fmt.Errorf("got %s, expected %s", s.cpmAddr, cpmStub)
	}

	// TODO: Timeouts

	if s.newrelicApp != nil {
		return fmt.Errorf("expected no new relic client, got %s", s.newrelicApp)
	}

	return nil
}

func validTimeouts(s *server) error {
	if s.httpServer.Addr != hostStub {
		return fmt.Errorf("got %s, expected %s", s.httpServer.Addr, hostStub)
	}
	if s.cpmAddr != cpmStub {
		return fmt.Errorf("got %s, expected %s", s.cpmAddr, cpmStub)
	}

	// TODO: Timeouts

	if s.newrelicApp != nil {
		return fmt.Errorf("expected no new relic client, got %s", s.newrelicApp)
	}

	return nil
}

func validNewRelic(s *server) error {
	if s.httpServer.Addr != hostStub {
		return fmt.Errorf("got %s, expected %s", s.httpServer.Addr, hostStub)
	}
	if s.cpmAddr != cpmStub {
		return fmt.Errorf("got %s, expected %s", s.cpmAddr, cpmStub)
	}

	// TODO: Timeouts

	if s.newrelicApp != newrelicStub {
		return fmt.Errorf("error setting new relic client, got %s, expected %s", s.newrelicApp, newrelicStub)
	}

	return nil
}
