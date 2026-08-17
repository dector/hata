package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"hata/internal/api"
	"hata/internal/apiui"
	"hata/internal/db"
	"hata/internal/devicestatus"
	"hata/internal/webauth"
	"hata/internal/webui"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func startServer(database db.DB, ctx context.Context) {
	// Setup chi router
	r := chi.NewRouter()

	// Add middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Create API handlers
	authHandler := api.NewAuthHandler(database.Repos())
	serverHandler := api.NewServerHandler()
	houseHandler := api.NewHouseHandler(database.Repos())
	deviceController := api.NewRealDeviceController()
	deviceHandler := api.NewDeviceHandlerWithController(database.Repos(), deviceController)
	shoppingListHandler := api.NewShoppingListHandler(database.Repos())
	devicestatus.StartPoller(ctx, database.Repos(), deviceController, 0)

	// Static assets
	r.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir("public"))))
	r.Handle("/js/*", http.StripPrefix("/js/", http.FileServer(http.Dir("public/js"))))

	// Add routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/app", http.StatusFound)
	})

	devIndexHandler := func(w http.ResponseWriter, r *http.Request) {
		homeData := webui.HomePageData{}
		auth, ok, err := webauth.TryAuthFromRequest(r, database.Repos())
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		if ok {
			homeData.IsLoggedIn = true
			homeData.DisplayName = strings.TrimSpace(auth.DisplayName)
			if homeData.DisplayName == "" {
				homeData.DisplayName = auth.Username
			}
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := webui.HomePage(homeData).Render(r.Context(), w); err != nil {
			http.Error(w, "failed to render dev index page", http.StatusInternalServerError)
			return
		}
	}
	r.Get("/dev", devIndexHandler)
	r.Get("/dev/", devIndexHandler)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	webAuthHandler := webauth.NewHandler(authHandler, database.Repos())

	// API routes
	r.Route("/api/latest", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", authHandler.Login)
		})
		r.Get("/house", houseHandler.List)
		r.Get("/house/{houseId}/device", deviceHandler.ListByHouse)
		r.Patch("/house/{houseId}/device/{deviceId}/state", deviceHandler.SetState)
		r.Patch("/house/{houseId}/device/{deviceId}/light", deviceHandler.SetLight)
		r.Get("/house/{houseId}/shopping-list", shoppingListHandler.ListByHouse)
		r.Get("/house/{houseId}/shopping-list/{listId}", shoppingListHandler.GetByHouseAndUID)
		r.Get("/house/{houseId}/shopping-list/{listId}/item", shoppingListHandler.ListItems)
		r.Post("/house/{houseId}/shopping-list/{listId}/item", shoppingListHandler.CreateItem)
		r.Patch("/house/{houseId}/shopping-list/{listId}/item/{itemId}", shoppingListHandler.UpdateItem)
		r.Patch("/house/{houseId}/shopping-list/{listId}/item/{itemId}/check", shoppingListHandler.SetItemChecked)
		r.Delete("/house/{houseId}/shopping-list/{listId}/item/{itemId}", shoppingListHandler.DeleteItem)
		r.Get("/device", deviceHandler.ListByUser)
		r.Get("/ping", serverHandler.Ping)
	})

	// Browser auth routes
	r.Route("/auth", func(r chi.Router) {
		r.Get("/login", webAuthHandler.LoginPage)
		r.Post("/login", webAuthHandler.Login)
		r.Get("/logout", webAuthHandler.LogoutPage)
		r.Post("/logout", webAuthHandler.Logout)
	})

	r.With(webauth.RequirePageAuth(database.Repos())).Get("/me", webAuthHandler.MePage)
	r.With(webauth.RequirePageAuth(database.Repos())).Get("/app", webAuthHandler.AppPage)
	r.With(webauth.RequirePageAuth(database.Repos())).Post("/app/active-house", webAuthHandler.SetActiveHouse)
	r.With(webauth.RequirePageAuth(database.Repos())).Get("/h/{houseId}/manage", webAuthHandler.HouseManagePage)
	r.With(webauth.RequirePageAuth(database.Repos())).Post("/h/{houseId}/device/{deviceId}/toggle", webAuthHandler.ToggleHouseDevice)
	r.With(webauth.RequirePageAuth(database.Repos())).Post("/h/{houseId}/device/{deviceId}/state", webAuthHandler.SetHouseDeviceState)
	r.With(webauth.RequirePageAuth(database.Repos())).Post("/h/{houseId}/device/{deviceId}/light", webAuthHandler.SetHouseDeviceLight)
	r.With(webauth.RequirePageAuth(database.Repos())).Post("/h/{houseId}/manage/devices", webAuthHandler.AddHouseDevice)
	r.With(webauth.RequirePageAuth(database.Repos())).Post("/h/{houseId}/manage/devices/{deviceId}/rename", webAuthHandler.RenameHouseDevice)
	r.With(webauth.RequirePageAuth(database.Repos())).Post("/h/{houseId}/manage/discovery-networks", webAuthHandler.AddHouseDiscoveryNetwork)
	r.With(webauth.RequirePageAuth(database.Repos())).Post("/h/{houseId}/manage/discovery-networks/{networkId}/delete", webAuthHandler.DeleteHouseDiscoveryNetwork)
	r.With(webauth.RequirePageAuth(database.Repos())).Get("/h/{houseId}/manage/devices/discover", webAuthHandler.HouseDeviceDiscovery)

	// API UI routes (dev-only)
	if os.Getenv("HATA_DEV") == "1" || os.Getenv("AIR") == "1" {
		authUIHandler := apiui.NewAuthUIHandler()
		houseUIHandler := apiui.NewHouseUIHandler()
		deviceUIHandler := apiui.NewDeviceUIHandler()
		pingUIHandler := apiui.NewPingUIHandler()
		apiUIIndexHandler := apiui.NewIndexHandler()
		r.Get("/apiui", apiUIIndexHandler.Index)
		r.Get("/apiui/", apiUIIndexHandler.Index)
		r.Get("/apiui/latest/auth/login", authUIHandler.LoginPage)
		r.Get("/apiui/latest/house", houseUIHandler.ListPage)
		r.Get("/apiui/latest/device", deviceUIHandler.ListPage)
		r.Get("/apiui/latest/house/{houseId}/device", houseUIHandler.HouseDevicePage)
		r.Get("/apiui/latest/ping", pingUIHandler.PingPage)
	}

	host := "localhost"
	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "4501"
	}
	port = strings.TrimPrefix(port, ":")
	listenAddr := fmt.Sprintf("%s:%s", host, port)

	portAccess := port
	if os.Getenv("AIR_PROXY") == "1" {
		portAccess = os.Getenv("AIR_PROXY_PORT")
		if strings.TrimSpace(portAccess) == "" {
			portAccess = "4500"
		}
	}

	// Start server
	log.Printf("Starting server on http://%s:%s\n", host, portAccess)
	if err := http.ListenAndServe(listenAddr, r); err != nil {
		log.Fatalf("failed starting server: %v", err)
	}
}
