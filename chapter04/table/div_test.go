package table_test

import (
	"errors"
	"testing"

	"github.com/PacktPublishing/Test-Driven-Development-in-Go/chapter04/table"
	"github.com/stretchr/testify/assert"
)

func TestDivide(t *testing.T) {
	tests := map[string]struct {
		x, y    int
		wantErr error
		want    string
	}{
		"pos x, pos y":   {x: 8, y: 4, want: "2.00"},
		"neg x, neg y":   {x: -4, y: -8, want: "0.50"},
		"equal x, y":     {x: 4, y: 4, want: "1.00"},
		"max x, pos y":   {x: 127, y: 2, want: "63.50"},
		"min x, pos y":   {x: -128, y: 2, want: "-64.00"},
		"zero x, pos y":  {x: 0, y: 2, want: "0.00"},
		"pos x, zero y":  {x: 10, y: 0, wantErr: errors.New("cannot divide by 0")},
		"zero x, zero y": {x: 0, y: 0, wantErr: errors.New("cannot divide by 0")},
		"max x, max y":   {x: 127, y: 127, want: "1.00"},
		"min x, min y":   {x: -128, y: -128, want: "1.00"},
	}

	for name, rtc := range tests {
		tc := rtc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			x, y := int8(tc.x), int8(tc.y)
			r, err := table.Divide(x, y)
			if tc.wantErr != nil {
				assert.Equal(t, tc.wantErr, err)
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, tc.want, *r)
		})
	}
}
