package log

import (
	"context"
	"reflect"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewFactory(t *testing.T) {
	type args struct {
		logger *zap.Logger
	}
	tests := []struct {
		name string
		args args
		want Factory
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFactory(tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFactory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFactory_Bg(t *testing.T) {
	type fields struct {
		logger *zap.Logger
	}
	tests := []struct {
		name   string
		fields fields
		want   Logger
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Factory{
				logger: tt.fields.logger,
			}
			if got := b.Bg(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Factory.Bg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFactory_For(t *testing.T) {
	type fields struct {
		logger *zap.Logger
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Logger
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Factory{
				logger: tt.fields.logger,
			}
			if got := b.For(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Factory.For() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFactory_With(t *testing.T) {
	type fields struct {
		logger *zap.Logger
	}
	type args struct {
		fields []zapcore.Field
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Factory
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Factory{
				logger: tt.fields.logger,
			}
			if got := b.With(tt.args.fields...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Factory.With() = %v, want %v", got, tt.want)
			}
		})
	}
}
