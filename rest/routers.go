package rest

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type route struct {
	Name        string
	Method      []string
	Pattern     string
	HandlerFunc http.HandlerFunc
	CheckAuth   bool
}

//Routes Rotas
type routes []route

//NewRouter router principal
func NewRouter(appParams AppParams) *mux.Router {
	appRoutes := loadRoutes(appParams)
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range appRoutes {
		var handler http.Handler
		handler = route.HandlerFunc
		if route.CheckAuth {
			handler = Auth(handler)
		}
		handler = Cors(handler)
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method...).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func httpMethods(method ...string) []string {
	method = append(method, "OPTIONS")
	return method
}

func loadRoutes(params AppParams) routes {
	var appRoutes = routes{

		route{
			"ProductList",
			httpMethods(strings.ToUpper("Get"), strings.ToUpper("Post")),
			"/product",
			makeProductsHandler(params),
			true,
		},
	}
	return appRoutes
}
