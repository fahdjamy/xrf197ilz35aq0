package handlers

import (
	"net/http"
	"xrf197ilz35aq0/core/service"
	"xrf197ilz35aq0/internal"
	"xrf197ilz35aq0/server/http/middleware"
)

// SetupHandlers creates route handlers and calls their RoutesHandler#RegisterRoutes methods for each handler to
// individually set up their own routes and call the necessary handlers
func SetupHandlers(serveMux *http.ServeMux, logger internal.Logger, services service.Services) http.Handler {
	routes := make([]RoutesHandler, 0)

	// create handlers
	healthRoutesHandler := NewHealthRoutes(logger)
	authRoutesHandler := NewAuthHandler(logger, services)
	orgRoutesHandler := NewOrgHandler(logger, services.OrgService)
	userRoutesHandler := NewUserHandler(logger, services.UserService)
	permissionRoutesHandler := NewPermHandler(logger, services.PermissionService)

	routes = append(routes, orgRoutesHandler)
	routes = append(routes, userRoutesHandler)
	routes = append(routes, authRoutesHandler)
	routes = append(routes, healthRoutesHandler)
	routes = append(routes, permissionRoutesHandler)

	// register routes To the http mux server
	for _, handler := range routes {
		handler.RegisterRoutes(serveMux)
	}

	loggerMiddleware := middleware.NewLoggerHandler(logger)

	// wrap entire mux with middleware
	authMiddleware := middleware.NewAuthenticationMiddleware(logger, services.AuthService)
	wrappedServer := loggerMiddleware.Handler(serveMux)
	wrappedServer = authMiddleware.Handle(wrappedServer)

	return wrappedServer
}

type RoutesHandler interface {
	RegisterRoutes(serveMux *http.ServeMux)
}
