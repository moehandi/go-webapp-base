package route

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
	//"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/moehandi/go-webapp-base/app/controller"
	"github.com/moehandi/go-webapp-base/route/middleware/acl"
	"github.com/moehandi/go-webapp-base/route/middleware/logrequest"
	//"github.com/moehandi/go-webapp-base/route/middleware/pprofhandler"
	"github.com/moehandi/go-webapp-base/helper/session"
	"github.com/gorilla/mux"
	"github.com/moehandi/go-webapp-base/route/middleware/pprofhandler"
)

// Load returns the routes and middleware
func Load() http.Handler {
	return middleware(routes())
}

// LoadHTTPS returns the HTTP routes and middleware
func LoadHTTPS() http.Handler {
	return middleware(routes())
}

// LoadHTTP returns the HTTPS routes and middleware
func LoadHTTP() http.Handler {
	return middleware(routes())

	// Uncomment this and comment out the line above to always redirect to HTTPS
	//return http.HandlerFunc(redirectToHTTPS)
}

// Optional method to make it easy to redirect from HTTP to HTTPS
func redirectToHTTPS(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "https://"+req.Host, http.StatusMovedPermanently)
}

// *****************************************************************************
// Routes
// *****************************************************************************

func routes() *mux.Router {
	r := mux.NewRouter()

	// TODO Custom error Page
	r.NotFoundHandler = http.HandlerFunc(controller.Error404)

	// Serve static files, no directory browsing
	r.PathPrefix("/static/").Handler(alice.New().ThenFunc(controller.Static))

	// Home Page
	r.Handle("/", alice.New().ThenFunc(controller.IndexGET)).Methods("GET")

	// Login
	r.Handle("/login", alice.New(acl.DisallowAuth).ThenFunc(controller.LoginGET)).Methods("GET")
	r.Handle("/login", alice.New(acl.DisallowAuth).ThenFunc(controller.LoginPOST)).Methods("POST")
	r.Handle("/logout", alice.New().ThenFunc(controller.LogoutGET))

	// Register
	r.Handle("/register", alice.New(acl.DisallowAuth).ThenFunc(controller.RegisterGET)).Methods("GET")
	r.Handle("/register", alice.New(acl.DisallowAuth).ThenFunc(controller.RegisterPOST)).Methods("POST")

	// About
	r.Handle("/about", alice.New().ThenFunc(controller.AboutGET)).Methods("GET")

	// Notepad API
	r.Handle("/api/notepad", alice.New(acl.DisallowAnon).ThenFunc(controller.ApiGetNote)).Methods("GET")
	r.Handle("/api/notepad/{id}", alice.New(acl.DisallowAnon).ThenFunc(controller.ApiNoteGetById)).Methods("GET")

	// Notepad
	r.Handle("/notepad", alice.New(acl.DisallowAnon).ThenFunc(controller.NotepadReadGET)).Methods("GET")
	r.Handle("/notepad/create", alice.New(acl.DisallowAnon).ThenFunc(controller.NotepadCreateGET)).Methods("GET")
	r.Handle("/notepad/create", alice.New(acl.DisallowAnon).ThenFunc(controller.NotepadCreatePOST)).Methods("POST")
	r.Handle("/notepad/update/{id}", alice.New(acl.DisallowAnon).ThenFunc(controller.NotepadUpdateGET)).Methods("GET")
	r.Handle("/notepad/update/{id}", alice.New(acl.DisallowAnon).ThenFunc(controller.NotepadUpdatePOST)).Methods("POST")
	r.Handle("/notepad/delete/{id}", alice.New(acl.DisallowAnon).ThenFunc(controller.NotepadDeleteGET)).Methods("GET")

	// Enable Pprof
	r.Handle("/debug/pprof/*pprof", alice.New(acl.DisallowAnon).ThenFunc(pprofhandler.Handler))

	return r
}

// *****************************************************************************
// Middleware
// *****************************************************************************

func middleware(h http.Handler) http.Handler {
	// Prevents CSRF and Double Submits
	cs := csrfbanana.New(h, session.Store, session.Name)
	cs.FailureHandler(http.HandlerFunc(controller.InvalidToken))
	cs.ClearAfterUsage(true)
	cs.ExcludeRegexPaths([]string{"/static(.*)"})
	csrfbanana.TokenLength = 32
	csrfbanana.TokenName = "token"
	csrfbanana.SingleToken = false
	h = cs

	// Log every request
	h = logrequest.Handler(h)

	// Clear handler for Gorilla Context
	h = context.ClearHandler(h)

	return h
}
