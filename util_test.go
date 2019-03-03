package graceful

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseEndpoint(t *testing.T) {
	tests := []struct {
		endpoint         string
		expectError      bool
		expectedProtocol string
		expectedAddr     string
	}{
		{
			endpoint:         "unix:///tmp/s12.sock",
			expectedProtocol: "unix",
			expectedAddr:     "/tmp/s12.sock",
		},
		{
			endpoint:         "tcp://localhost:15880",
			expectedProtocol: "tcp",
			expectedAddr:     "localhost:15880",
		},
		{
			endpoint:         "npipe://./pipe/mypipe",
			expectedProtocol: "npipe",
			expectError:      true,
		},
		{
			endpoint:         "tcp1://abc",
			expectedProtocol: "tcp1",
			expectError:      true,
		},
		{
			endpoint:    "a b c",
			expectError: true,
		},
	}

	for _, test := range tests {
		protocol, addr, err := parseEndpoint(test.endpoint)
		assert.Equal(t, test.expectedProtocol, protocol)
		if test.expectError {
			assert.NotNil(t, err, "Expect error during parsing %q", test.endpoint)
			continue
		}
		assert.Nil(t, err, "Expect no error during parsing %q", test.endpoint)
		assert.Equal(t, test.expectedAddr, addr)
	}

	for _, test := range tests {
		l, f, err := CreateListenerFile(test.endpoint)
		if test.expectError {
			assert.NotNil(t, err, "Expect error during parsing %q", test.endpoint)
			continue
		}
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(test.endpoint)
		fmt.Println(f.Name())
		l.Close()
	}

}
