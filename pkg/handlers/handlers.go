package handlers

import (
	"encoding/json"
	"fmt"
	"groupietracker/pkg/models"
	"groupietracker/pkg/utils"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	// Récupérer les paramètres de la requête
	creationDateMinStr := r.URL.Query().Get("creationDateMin")
	creationDateMaxStr := r.URL.Query().Get("creationDateMax")
	albumDateMin := r.URL.Query().Get("albumDateMin")
	albumDateMax := r.URL.Query().Get("albumDateMax")
	members := r.URL.Query().Get("members")
	location := r.URL.Query().Get("location")

	// Convertir les valeurs des dates
	creationDateMin, _ := strconv.Atoi(creationDateMinStr)
	creationDateMax, _ := strconv.Atoi(creationDateMaxStr)

	// Récupérer tous les artistes
	allArtists, err := utils.FetchArtists()
	if err != nil {
		ErrorPageHandler(w, r, "Unable to fetch artists", http.StatusInternalServerError)
		return
	}

	var filteredArtists []models.Artist

	// Si des filtres sont appliqués
	if creationDateMinStr != "" || creationDateMaxStr != "" || albumDateMin != "" || albumDateMax != "" || members != "" || location != "" {
		filteredArtists = utils.FitlerArtists(allArtists, creationDateMin, creationDateMax, albumDateMin, albumDateMax, members, location)
	} else {
		filteredArtists = allArtists
	}

	// Afficher les artistes filtrés
	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	tmpl.Execute(w, filteredArtists)
}

// ErrorPageHandler gère les erreurs et affiche une page d'erreur personnalisée.
func ErrorPageHandler(w http.ResponseWriter, r *http.Request, message string, status int) {
	w.Header().Set("Content-Type", "text/html")
	//w.WriteHeader(status)

	data := struct {
		StatusCode int
		Message    string
	}{
		StatusCode: status,
		Message:    message,
	}

	tmpl, err := template.ParseFiles("templates/error.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// NotFoundHandler gère les erreurs 404.
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	ErrorPageHandler(w, r, "Not found", http.StatusNotFound)
}

func ArtistDetailsPage(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/artist/"):]

	// Fetch artist details
	artist, err := utils.FetchArtistDetails(id)
	if err != nil || artist.ID == 0 {
		ErrorPageHandler(w, r, "Unable to fetch artist details", http.StatusInternalServerError)
		return
	}

	// Fetch artist relations (concert locations)
	relations, err := utils.FetchRelationsById(id)
	if err != nil {
		fmt.Println(err)
		ErrorPageHandler(w, r, "Unable to fetch relations", http.StatusInternalServerError)
		return
	}

	// Clean the locations within the relations and fetch coordinates
	locationsWithCoordinates := []models.LocationWithCoordinates{}
	for location, dates := range relations.DatesLocations {
		// Séparer ville et pays
		city, country, err := utils.SplitCityAndCountry(location)
		if err != nil {
			fmt.Println("Error splitting location:", err)
			continue
		}

		// Fetch coordinates using geocoding
		coords, err := utils.Geocode(city, country)
		if err != nil {
			fmt.Println("Error fetching coordinates for", city, country, err)
			continue
		}

		// Append location with coordinates and dates
		locationWithCoords := models.LocationWithCoordinates{
			City:      utils.CleanLocationName(location),
			Latitude:  coords.Latitude,
			Longitude: coords.Longitude,
			Dates:     dates,
		}
		locationsWithCoordinates = append(locationsWithCoordinates, locationWithCoords)
	}

	// Convert locationsWithCoordinates to JSON
	locationsJSON, err := json.Marshal(locationsWithCoordinates)
	if err != nil {
		ErrorPageHandler(w, r, "Error generating location data", http.StatusInternalServerError)
		return
	}

	// Prepare the data to be passed to the template
	data := struct {
		Artist        models.Artist
		LocationsJSON string
	}{
		Artist:        *artist,
		LocationsJSON: string(locationsJSON), // Pass the JSON as a string
	}

	// Parse and execute the template
	tmpl, err := template.ParseFiles("templates/artist.html")
	if err != nil {
		ErrorPageHandler(w, r, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		ErrorPageHandler(w, r, "Internal Server Error", http.StatusInternalServerError)
	}
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := strings.ToLower(r.URL.Query().Get("q"))
	artists, err := utils.FetchArtists()
	if err != nil {
		ErrorPageHandler(w, r, "Unable to fetch data", http.StatusInternalServerError)
		return
	}

	// Determine the type of search (artist/band, member, creation date, first album date, or concert date)
	cleanedQuery, searchType := cleanSearchQuery(query)

	var searchResults []models.Artist

	// Search based on the type
	for _, artist := range artists {
		// Fetch relations for each artist to include locations in the search
		relations, err := utils.FetchRelationsById(strconv.Itoa(artist.ID))
		if err != nil {
			log.Printf("Error fetching relations for artist %d: %v", artist.ID, err)
			continue
		}

		// Clean the location names in relations
		cleanedDatesLocations := make(map[string][]string)
		for location, dates := range relations.DatesLocations {
			cleanedLocation := utils.CleanLocationName(location)
			cleanedDatesLocations[cleanedLocation] = dates
		}
		relations.DatesLocations = cleanedDatesLocations

		switch searchType {
		case "artist/band":
			// Search by artist name
			if strings.EqualFold(artist.Name, cleanedQuery) {
				searchResults = append(searchResults, artist)
			}
		case "member":
			// Search by member name
			for _, member := range artist.Members {
				if strings.EqualFold(member, cleanedQuery) {
					searchResults = append(searchResults, artist)
					break
				}
			}
		case "creation date":
			// Search by creation date
			if strings.EqualFold(artist.Name, cleanedQuery) {
				searchResults = append(searchResults, artist)
			}
		case "first album date":
			// Search by first album date
			if strings.EqualFold(artist.Name, cleanedQuery) {
				searchResults = append(searchResults, artist)
			}
		case "concert date":
			// Search by concert date in relations
			if strings.EqualFold(artist.Name, cleanedQuery) {
				searchResults = append(searchResults, artist)
			}
		default:
			// General search across artist name, album, creation date, members, and locations
			if strings.Contains(strings.ToLower(artist.Name), cleanedQuery) ||
				strings.Contains(strings.ToLower(artist.FirstAlbum), cleanedQuery) ||
				strings.Contains(strings.ToLower(strconv.Itoa(artist.CreationDate)), cleanedQuery) {
				searchResults = append(searchResults, artist)
			}

			// Search in members
			for _, member := range artist.Members {
				if strings.Contains(strings.ToLower(member), cleanedQuery) {
					searchResults = append(searchResults, artist)
					break
				}
			}

			// Search in locations and concert dates
			for location, dates := range relations.DatesLocations {
				if strings.Contains(strings.ToLower(location), cleanedQuery) {
					searchResults = append(searchResults, artist)
					break
				}
				for _, date := range dates {
					if strings.Contains(strings.ToLower(date), cleanedQuery) {
						searchResults = append(searchResults, artist)
						break
					}
				}
			}
		}
	}

	// Remove duplicates if necessary
	uniqueResults := make(map[int]models.Artist)
	for _, artist := range searchResults {
		uniqueResults[artist.ID] = artist
	}

	// Convert the map to a slice
	finalResults := make([]models.Artist, 0, len(uniqueResults))
	for _, artist := range uniqueResults {
		finalResults = append(finalResults, artist)
	}

	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		ErrorPageHandler(w, r, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, finalResults)
}

// cleanSearchQuery removes any category suffix from the search query and returns the cleaned query and type.
func cleanSearchQuery(query string) (string, string) {
	// Define possible suffixes related to different search types
	suffixes := []string{
		" - artist/band",
		" - member",
		" - location",
		" - creation date",
		" - first album date",
		" - concert date",
	}

	for _, suffix := range suffixes {
		if strings.HasSuffix(query, suffix) {
			return strings.TrimSuffix(query, suffix), strings.TrimPrefix(suffix, " - ")
		}
	}

	// If no specific suffix is found, return the query as is
	return query, ""
}

func SuggestionHandler(w http.ResponseWriter, r *http.Request) {
	query := strings.ToLower(r.URL.Query().Get("q"))
	artists, err := utils.FetchArtists()
	if err != nil {
		http.Error(w, "Unable to fetch data", http.StatusInternalServerError)
		return
	}

	var suggestions []string

	// Check if the query is a date (YYYY-MM-DD) or a year (YYYY)
	isDate := regexp.MustCompile(`^\d{2}-\d{2}-\d{4}$`).MatchString(query)
	isYear := regexp.MustCompile(`^\d{4}$`).MatchString(query)

	for _, artist := range artists {
		// Fetch relations for each artist to include locations in the suggestions
		relations, err := utils.FetchRelationsById(strconv.Itoa(artist.ID))
		if err != nil {
			log.Printf("Error fetching relations for artist %d: %v", artist.ID, err)
			continue
		}

		// Clean the location names in relations
		cleanedDatesLocations := make(map[string][]string)
		for location, dates := range relations.DatesLocations {
			cleanedLocation := utils.CleanLocationName(location)
			cleanedDatesLocations[cleanedLocation] = dates
		}
		relations.DatesLocations = cleanedDatesLocations

		// Generate suggestions based on the artist's name
		artistSuggestion := fmt.Sprintf("%s - artist/band", artist.Name)
		if strings.Contains(strings.ToLower(artist.Name), query) && !contains(suggestions, artistSuggestion) {
			suggestions = append(suggestions, artistSuggestion)
		}

		// Generate suggestions based on the members' names
		for _, member := range artist.Members {
			memberSuggestion := fmt.Sprintf("%s - member", member)
			if strings.Contains(strings.ToLower(member), query) && !contains(suggestions, memberSuggestion) {
				suggestions = append(suggestions, memberSuggestion)
			}
		}

		// Generate suggestions based on the locations
		for location := range relations.DatesLocations {
			locationSuggestion := fmt.Sprintf("%s - location", location)
			if strings.Contains(strings.ToLower(location), query) && !contains(suggestions, locationSuggestion) {
				suggestions = append(suggestions, locationSuggestion)
			}
		}

		// Generate suggestions based on the date
		if isDate {
			for _, dates := range relations.DatesLocations {
				for _, date := range dates {
					if date == query {
						concertSuggestion := fmt.Sprintf("%s - concert date", artist.Name)
						if !contains(suggestions, concertSuggestion) {
							suggestions = append(suggestions, concertSuggestion)
						}
					}
				}
			}

			if artist.FirstAlbum == query {
				albumSuggestion := fmt.Sprintf("%s - first album date", artist.Name)
				if !contains(suggestions, albumSuggestion) {
					suggestions = append(suggestions, albumSuggestion)
				}
			}
		}

		// Generate suggestions based on the year
		if isYear {
			year, _ := strconv.Atoi(query)
			if artist.CreationDate == year {
				creationSuggestion := fmt.Sprintf("%s - creation date", artist.Name)
				if !contains(suggestions, creationSuggestion) {
					suggestions = append(suggestions, creationSuggestion)
				}
			}
		}
	}

	// Send the suggestions as JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(suggestions)
}

// Helper function to check if a slice contains a particular string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
