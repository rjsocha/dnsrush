all: dnsrush

dnsrush: dnsrush.go
	go build dnsrush.go
	strip -x dnsrush
	upx dnsrush
