package cache

import (
	"context"
	"project/internal/model"
	"testing"
)

func Test_radisLayer_AddCache(t *testing.T) {
	type args struct {
		ctx     context.Context
		jid     uint
		jobdata model.Job
	}
	tests := []struct {
		name    string
		r       *radisLayer
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.AddCache(tt.args.ctx, tt.args.jid, tt.args.jobdata); (err != nil) != tt.wantErr {
				t.Errorf("radisLayer.AddCache() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
