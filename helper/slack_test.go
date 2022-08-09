package helper

import "testing"

func TestSendNotification(t *testing.T) {
	type args struct {
		title string
		body  string
		ctx   string
		err   error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Case 1: Success",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SendNotification(tt.args.title, tt.args.body, tt.args.ctx, tt.args.err)
		})
	}
}
