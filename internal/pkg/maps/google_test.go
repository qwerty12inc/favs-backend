package maps

import (
	"gitlab.com/v.rianov/favs-backend/internal/models"
	"testing"
)

func TestLocationLinkResolverImpl_ResolveLink(t *testing.T) {
	type fields struct {
	}
	type args struct {
		link string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.Coordinates
		wantErr bool
	}{
		{
			name: "valid link",
			args: args{
				link: "https://www.google.com/maps/@-6.123456,106.123456",
			},
		},
		{
			name: "valid link",
			args: args{
				link: "https://www.google.com/maps/place/Santa+Barbara+Complex/@34.7067873,33.0876892,15.75z/data=!4m15!1m5!3m4!2zMzTCsDQyJzE4LjciTiAzM8KwMDUnNDQuMiJF!8m2!3d34.705187!4d33.095623!3m8!1s0x14e0cb6916c3bfd7:0xb774504ab68aa1aa!5m2!4m1!1i2!8m2!3d34.7049327!4d33.1053062!16s%2Fg%2F1pp2t_h2m?entry=ttu",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := LocationLinkResolverImpl{}
			got, err := l.ResolveLink(tt.args.link)
			if (err != nil) != tt.wantErr {
				t.Errorf("LocationLinkResolverImpl.ResolveLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Latitude == 0 || got.Longitude == 0 {
				t.Errorf("LocationLinkResolverImpl.ResolveLink() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocationLinkResolverImpl_ResolveLink_InvalidLink(t *testing.T) {
	l := LocationLinkResolverImpl{}
	_, err := l.ResolveLink("https://www.google.com/maps")
	if err == nil {
		t.Error("LocationLinkResolverImpl.ResolveLink() = nil, want error")
	}
}
