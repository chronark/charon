package osm

import (
	"context"
	"testing"

	"github.com/chronark/charon/pkg/log"
	"github.com/chronark/charon/service/tiles/proto/tiles"
	"github.com/micro/go-micro/v2/client"
)

func TestHandler_Get(t *testing.T) {
	type fields struct {
		Client client.Client
		Logger log.Factory
	}
	type args struct {
		ctx context.Context
		req *tiles.Request
		res *tiles.Response
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Successful test",
			fields: fields{
				Client: client.DefaultClient,
				Logger: log.NewDefaultLogger("service"),
			},
			args: args{
				ctx: context.TODO(),
				req: &tiles.Request{
					X: 1,
					Y: 2,
					Z: 3,
				},
				res: &tiles.Response{
					File: []byte("some bytes"),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				Client: tt.fields.Client,
				Logger: tt.fields.Logger,
			}
			if err := h.Get(tt.args.ctx, tt.args.req, tt.args.res); (err != nil) != tt.wantErr {
				t.Errorf("Handler.Get() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}


