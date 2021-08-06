package handlers

import (
	"github.com/DapperBlondie/simple-chat-app/src/render"
	"net/http"
)

type ApplicationConfig struct {

}

var AppConf *ApplicationConfig

func (ac *ApplicationConfig) Home(w http.ResponseWriter, r *http.Request)  {
	if http.MethodGet != r.Method {
		http.Error(w, "Error in method usage.", http.StatusMethodNotAllowed)
		return
	}

	err := render.RendererPage(w, "Home.jet", nil)
	if err != nil {
		http.Error(w, "Error in rendering page", http.StatusInternalServerError)
		return
	}
	return
}