package osm

import (
	"context"
	"net/http"
	"reflect"
	"testing"
	"net/url"

	"github.com/micro/go-micro/v2/client"
	"github.com/chronark/charon/pkg/log"
	"github.com/chronark/charon/service/tiles/proto/tiles"
)

func TestHandler_parseCoordinates(t *testing.T) {
	validURL, _ := url.Parse("http://dummy/?x=1&y=2&z=3")


	type fields struct {
		Logger log.Factory
		Client tiles.TilesService
	}
	type args struct {
		ctx context.Context
		r   *http.Request
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *tiles.Request
		wantErr bool
	}{
		{
			name: "Successful parse",
			fields: fields{
				Logger: log.NewDefaultLogger("service"),
				Client: tiles.NewTilesService("test.srv.tiles", client.DefaultClient),
			},
			args: args{
				ctx: context.TODO(),
				r: &http.Request{
					URL: validURL,
				},
			},
			want: &tiles.Request{
				X: 1,
				Y: 2,
				Z: 3,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				Logger: tt.fields.Logger,
				Client: tt.fields.Client,
			}
			got, err := h.parseCoordinates(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Handler.parseCoordinates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Handler.parseCoordinates() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_Get(t *testing.T) {
	type fields struct {
		Logger log.Factory
		Client tiles.TilesService
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				Logger: tt.fields.Logger,
				Client: tt.fields.Client,
			}
			h.Get(tt.args.w, tt.args.r)
		})
	}
}
