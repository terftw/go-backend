package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/terftw/go-backend/internal/api/handlers"
)

func SetupRoutes(r chi.Router, h *handlers.Handlers) {
	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	// Public routes
	r.Group(func(r chi.Router) {
		r.Route("/api/v1/auth", func(r chi.Router) {
			r.Get("/google", h.Auth.InitiateGoogleOAuth)
			r.Get("/google/callback", h.Auth.HandleGoogleCallback)
		})
	})

	// // Protected routes
	// r.Group(func(r chi.Router) {
	// 	r.Route("/api/v1", func(r chi.Router) {
	// 		// User routes
	// 		r.Route("/users", func(r chi.Router) {
	// 			r.Get("/me", h.User.GetUser)
	// 			r.Put("/me", h.User.UpdateUser)
	// 		})

	// 		// Add other resource routes here
	// 		// r.Route("/resources", func(r chi.Router) {
	// 		//     r.Get("/", h.Resource.List)
	// 		//     r.Post("/", h.Resource.Create)
	// 		// })
	// 	})
	// })
}
