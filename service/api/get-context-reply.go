package api

import (
	"github.com/gensimone/WASA-project/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// getContextReply is an example of HTTP endpoint that returns "Hello World!" as a plain text. The signature of this
// handler accepts a reqcontext.RequestContext (see httpRouterHandler).
func (rt *_router) getContextReply(w http.ResponseWriter, _ *http.Request, _ httprouter.Params, _ reqcontext.RequestContext) {
	w.Header().Set("content-type", "text/plain")
	_, _ = w.Write([]byte("Hello World!"))
}
