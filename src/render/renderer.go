package render

import (
	"github.com/CloudyKit/jet/v6"
	zerolog "github.com/rs/zerolog/log"
	"net/http"
)

var Views = jet.NewSet(
	jet.NewOSFileSystemLoader("/templates"),
	jet.InDevelopmentMode(),
)

func RendererPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {
	view, err := Views.GetTemplate(tmpl)
	if err != nil {
		zerolog.Error().Msg(err.Error() + "; error in render page")
		return err
	}

	err = view.Execute(w, data, nil)
	if err != nil {
		zerolog.Error().Msg(err.Error() + "; error in render page")
		return err
	}

	return nil
}
