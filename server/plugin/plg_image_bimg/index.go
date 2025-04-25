package plg_image_golang

import (
	"io"
	"net/http"

	"github.com/davidbyttow/govips/v2/vips"
	. "github.com/mickael-kerjean/filestash/server/common"
)

const THUMB_SIZE int = 150

func init() {
	Hooks.Register.Thumbnailer("image/jpeg", thumbnailer{})
	Hooks.Register.Thumbnailer("image/png", thumbnailer{})
	Hooks.Register.Thumbnailer("image/gif", thumbnailer{})
}

type thumbnailer struct{}

func (this thumbnailer) Generate(reader io.ReadCloser, ctx *App, res *http.ResponseWriter, req *http.Request) (io.ReadCloser, error) {
	query := req.URL.Query()
	mType := GetMimeType(query.Get("path"))

	if query.Get("thumbnail") != "true" {
		return reader, nil
	} else if mType != "image/jpeg" && mType != "image/png" && mType != "image/gif" {
		return reader, nil
	}

	newImage, err := vips.NewImageFromReader(reader)
	if err != nil {
		return reader, err
	}
	err = newImage.Thumbnail(300, 300, vips.InterestingAll)
	if err != nil {
		return reader, err
	}
	param := vips.NewJpegExportParams()
	param.Quality = 85
	data, _, err := newImage.ExportJpeg(param)

	return NewReadCloserFromBytes(data), err
}
