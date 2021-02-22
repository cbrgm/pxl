package pxl

import "testing"

func TestResizeImageBounds(t *testing.T) {
	type args struct {
		maxWidth  int
		maxHeight int
		dx        int
		dy        int
	}
	tests := []struct {
		name  string
		args  args
		wantX int
		wantY int
	}{
		{
			name: "Test MaxWidth lower than dx",
			args: args{
				maxWidth:  1024,
				maxHeight: 1024,
				dx:        1900,
				dy:        1024,
			},
			wantX: 1024,
			wantY: 551,
		},
		{
			name: "Test MaxHeight lower than dy",
			args: args{
				maxWidth:  1024,
				maxHeight: 900,
				dx:        1900,
				dy:        1024,
			},
			wantX: 1669,
			wantY: 900,
		},
		{
			name: "Test MaxHeight priority over MaxWidth",
			args: args{
				maxWidth:  1024,
				maxHeight: 800,
				dx:        1900,
				dy:        1024,
			},
			wantX: 1484,
			wantY: 800,
		},
		{
			name: "Test MaxWidth and MaxHeight greater than dx",
			args: args{
				maxWidth:  1024,
				maxHeight: 800,
				dx:        900,
				dy:        600,
			},
			wantX: 900,
			wantY: 600,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotX, gotY := resizeImageBounds(tt.args.maxWidth, tt.args.maxHeight, tt.args.dx, tt.args.dy)
			if gotX != tt.wantX {
				t.Errorf("resizeImageBounds() gotX = %v, want %v", gotX, tt.wantX)
			}
			if gotY != tt.wantY {
				t.Errorf("resizeImageBounds() gotY = %v, want %v", gotY, tt.wantY)
			}
		})
	}
}
