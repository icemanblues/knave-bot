package main

import "testing"

func TestIsSlackUser(t *testing.T) {
	testCases := []struct {
		name     string
		user     string
		ok       bool
		expected string
	}{
		{"escaped", "<@UAWQFTRT7|roland.kluge>", true, "UAWQFTRT7"},
		{"canonical", "UAWQFTRT7", true, "UAWQFTRT7"},
		{"fail", "fail", false, ""},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			actual, ok := IsSlackUser(test.user)
			if ok != test.ok {
				t.Errorf("Expecting %v but actual %v\n", test.ok, ok)
			}
			if actual != test.expected {
				t.Errorf("Expecting %v but actual %v\n", test.expected, actual)
			}
		})
	}
}
