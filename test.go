package main

import "stencil/qr"

func main() {
	appName := "diaspora"

	QR := qr.NewQRWithAppName(appName)
	QR.TestQuery("insert")
}
