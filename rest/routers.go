package rest

import (
	"net/http"

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
	appRoutes := routes{} //loadRoutes()
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

// func loadRoutes() routes {
// 	var appRoutes = routes{

// 		route{
// 			"Clientes",
// 			httpMethods(strings.ToUpper("Get"), strings.ToUpper("Post")),
// 			"/v1/api/clientes",
// 			mkHandlerListaCliente(clienteDAO),
// 			true,
// 		},
// 		route{
// 			"NovoCliente",
// 			httpMethods(strings.ToUpper("Post")),
// 			"/v1/api/cliente/novo",
// 			mkHandlerNovoCliente(clienteDAO),
// 			true,
// 		},
// 	}
// 	return appRoutes
// }
