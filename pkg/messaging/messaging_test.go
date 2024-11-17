package messaging

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestData_Validate(t *testing.T) {
	tests := []struct {
		name    string
		d       *Data
		wantErr error
	}{
		{
			name:    "ErrorDataNil",
			d:       nil,
			wantErr: ErrDataNil,
		},
		{
			name:    "ErrorMsgNil",
			d:       &Data{},
			wantErr: ErrMessageNil,
		},
		{
			name:    "Success",
			d:       &Data{Msg: []byte{}},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.d.Validate()
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
