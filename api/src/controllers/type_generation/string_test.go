package typegeneration

import (
	"api/src/controllers"
	"encoding/json"
	"testing"
)

func Test_generateString(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		restriction controllers.StringRestrictions
		want        json.RawMessage
		want2       int
		wantErr     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got2, gotErr := generateString(tt.restriction)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("generateString() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("generateString() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("generateString() = %v, want %v", got, tt.want)
			}
			if true {
				t.Errorf("generateString() = %v, want %v", got2, tt.want2)
			}
		})
	}
}
