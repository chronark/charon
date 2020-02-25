package osm

import (
	"context"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/chronark/charon/pkg/log"
	"github.com/chronark/charon/service/tiles/proto/tiles"
	"github.com/micro/go-micro/v2/client"
)

func urlConstructor(parameters string) *url.URL {
	url, _ := url.Parse("http://server/?" + parameters)
	return url
}

func TestHandler_parseCoordinates(t *testing.T) {

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
				Logger: log.NewEmptyLogger(),
				Client: tiles.NewTilesService("test.srv.tiles", client.DefaultClient),
			},
			args: args{
				ctx: context.TODO(),
				r: &http.Request{
					URL: urlConstructor("x=1&y=2&z=3"),
				},
			},
			want: &tiles.Request{
				X: 1,
				Y: 2,
				Z: 3,
			},
			wantErr: false,
		},
		{
			name: "missing x",
			fields: fields{
				Logger: log.NewEmptyLogger(),
				Client: tiles.NewTilesService("test.srv.tiles", client.DefaultClient),
			},
			args: args{
				ctx: context.TODO(),
				r: &http.Request{
					URL: urlConstructor("y=2&z=3"),
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "missing y",
			fields: fields{
				Logger: log.NewEmptyLogger(),
				Client: tiles.NewTilesService("test.srv.tiles", client.DefaultClient),
			},
			args: args{
				ctx: context.TODO(),
				r: &http.Request{
					URL: urlConstructor("x=2&z=3"),
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "missing z",
			fields: fields{
				Logger: log.NewEmptyLogger(),
				Client: tiles.NewTilesService("test.srv.tiles", client.DefaultClient),
			},
			args: args{
				ctx: context.TODO(),
				r: &http.Request{
					URL: urlConstructor("y=2&x=3"),
				},
			},
			want:    nil,
			wantErr: true,
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
