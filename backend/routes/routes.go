package routes

import (
	"backend/handlers"
	"backend/middleware"
	"database/sql"
	"net/http"

	"github.com/pkg/sftp"
)

func SetupRoutes(database *sql.DB, sftp *sftp.Client) {
	// Route de redirection intelligente pour les portails RH et admin
	http.HandleFunc("/gestion", middleware.RequireJWT(database, func(w http.ResponseWriter, r *http.Request) {
		// Empêcher le cache du navigateur pour cette route
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // pour forcer le navigateur à recharger la page à chaque fois
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		claims := middleware.GetClaims(r)
		switch claims.Role {
		case "admin":
			http.ServeFile(w, r, "/usr/share/nginx/html/pages/admin.html")
		case "rh":
			http.ServeFile(w, r, "/usr/share/nginx/html/pages/rh.html")
		default:
			http.Error(w, "Accès refusé", http.StatusForbidden)
		}
	}))

	//  Routes API
	http.HandleFunc("/api/forms", middleware.RequireRole(database, []string{"rh", "admin"}, handlers.GetAllForms(database)))
	http.HandleFunc("/api/form", middleware.RequireRole(database, []string{"rh", "admin"}, handlers.GetForm(database)))
	http.HandleFunc("/api/form/add", middleware.RequireRole(database, []string{"rh", "admin"}, handlers.AjoutForm(database, sftp)))

	// Routes Admin
	http.HandleFunc("/api/form/delete", middleware.RequireRole(database, []string{"admin"}, handlers.SuppForm(database)))
	http.HandleFunc("/api/form/modify", middleware.RequireRole(database, []string{"admin"}, handlers.ModifForm(database)))
	http.HandleFunc("/api/users", middleware.RequireRole(database, []string{"admin"}, handlers.GetAllUsersHandler(database)))

	// déconnexion (le client n'a pas accès directement au token en JS)

	http.HandleFunc("/api/logout", handlers.LogoutHandler(database))

	// Routes publiques (pas de JWT)
	http.HandleFunc("/api/user/add", handlers.AjoutUser(database))
	http.HandleFunc("/api/login", handlers.LoginHandler(database))
}
