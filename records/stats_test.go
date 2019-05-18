package records

import "testing"

func TestGetAvg(t *testing.T) {
	cases := []struct {
		name        string
		elems       []string
		expected    int
		expectedErr bool
	}{
		{
			"Should be able to compute avg",
			[]string{
				"20",
				"20",
				"0",
				"40",
			},
			20,
			false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			avg, err := getAvg(tc.elems)
			if err != nil {
				t.Error(err)
			}
			if avg != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, avg)
			}
		})
	}
}
