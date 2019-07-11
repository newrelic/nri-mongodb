package arguments

import (
	"testing"
)

func TestArgumentList_Validate(t *testing.T) {
	testCases := []struct {
		argumentList  ArgumentList
		expectedError bool
	}{
		{
			argumentList: ArgumentList{
				Username:              "testUser",
				Password:              "testPass",
				Host:                  "testHost",
				ClusterName:           "testClusterName",
				Port:                  "27071",
				AuthSource:            "admin",
				Ssl:                   false,
				SslInsecureSkipVerify: false,
			},
			expectedError: false,
		},
		{
			argumentList: ArgumentList{
				Password:              "testPass",
				Host:                  "testHost",
				Port:                  "27071",
				AuthSource:            "admin",
				Ssl:                   false,
				SslInsecureSkipVerify: false,
			},
			expectedError: true,
		},
		{
			argumentList: ArgumentList{
				Username:              "testUser",
				Host:                  "testHost",
				Port:                  "27071",
				AuthSource:            "admin",
				Ssl:                   false,
				SslInsecureSkipVerify: false,
			},
			expectedError: true,
		},
		{
			argumentList: ArgumentList{
				Username:              "testUser",
				Password:              "testPass",
				Port:                  "27071",
				AuthSource:            "admin",
				Ssl:                   false,
				SslInsecureSkipVerify: false,
			},
			expectedError: true,
		},
		{
			argumentList: ArgumentList{
				Username:              "testUser",
				Password:              "testPass",
				Host:                  "testHost",
				Port:                  "testPort",
				AuthSource:            "admin",
				Ssl:                   false,
				SslInsecureSkipVerify: false,
			},
			expectedError: true,
		},
		{
			argumentList: ArgumentList{
				Username:              "testUser",
				Password:              "testPass",
				ClusterName:           "testClusterName",
				Host:                  "testHost",
				Port:                  "2000",
				AuthSource:            "admin",
				Ssl:                   false,
				SslInsecureSkipVerify: true,
			},
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		err := tc.argumentList.Validate()
		if (err != nil) != tc.expectedError {
			t.Errorf("expected error: %v, got: %v for %v", tc.expectedError, err, tc.argumentList)
		}
	}
}
