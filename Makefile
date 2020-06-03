all: dnsrush

dnsrush: dnsrush.go
	CGO_ENABLED=0 go build dnsrush.go
	strip -x dnsrush
	upx dnsrush
