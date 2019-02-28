package database

import (
	"testing"

	"github.com/marcospsbrito/dic/config"
)

func Test_connect(t *testing.T) {
	type args struct {
		dbURL  string
		dbName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Connect to db", args{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := connect(tt.args.dbURL, tt.args.dbName)
			if (err != nil) != tt.wantErr {
				t.Errorf("connect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Error("connect() returns nil")
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		config config.Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"should return new connection", args{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
