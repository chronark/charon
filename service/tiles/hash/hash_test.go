package hash

import (
	"context"
	"testing"

	tiles "github.com/chronark/charon/service/tiles/proto/tiles"
)

func TestHashRequest(t *testing.T) {
	type args struct {
		ctx context.Context
		req *tiles.Request
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{

			name: "Random Tile",
			args: args{
				ctx: context.TODO(),
				req: &tiles.Request{
					X: 1,
					Y: 2,
					Z: 3,
				},
			},
			want: "tiles/3/1/2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HashRequest(tt.args.ctx, tt.args.req); got != tt.want {
				t.Errorf("HashRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
