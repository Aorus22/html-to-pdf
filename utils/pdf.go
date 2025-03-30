package utils

import (
	"encoding/base64"
	"fmt"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

func GeneratePDF(html string) ([]byte, error) {
	url := launcher.New().
		Bin("/usr/bin/chromium-browser").
		Headless(true).
		NoSandbox(true).
		MustLaunch()

	browser := rod.New().ControlURL(url).MustConnect()
	defer browser.MustClose()

	page := browser.MustPage()

	htmlBase64 := base64.StdEncoding.EncodeToString([]byte(html))
	dataURI := fmt.Sprintf("data:text/html;base64,%s", htmlBase64)

	if err := page.Navigate(dataURI); err != nil {
		return nil, err
	}

	page.MustWaitLoad().MustWaitIdle()
	page.MustSetViewport(1200, 1600, 1, false)

	scale := new(float64)
	*scale = 0.7

	width := new(float64)
	*width = 8.27

	height := new(float64)
	*height = 11.69

	pdfParams := &proto.PagePrintToPDF{
		PrintBackground: true,
		Scale:           scale,
		PaperWidth:      width,
		PaperHeight:     height,
	}

	pdfBytes, err := pdfParams.Call(page)
	if err != nil {
		return nil, err
	}

	return pdfBytes.Data, nil
}
