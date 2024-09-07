package utils

import (
	"encoding/json"
	"fmt"
	"groupietracker/pkg/models"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// FetchArtists récupère les données des artistes depuis l'API.
func FetchArtists() ([]models.Artist, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var artists []models.Artist
	err = json.NewDecoder(resp.Body).Decode(&artists)
	if err != nil {
		return nil, err
	}
	return artists, nil
}

func FetchLocations() ([]models.Location, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/locations")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var locations []models.Location
	err = json.NewDecoder(resp.Body).Decode(&locations)
	if err != nil {
		return nil, err
	}
	return locations, nil
}

// FetchLocations récupère les données des locations (lieux de concert) depuis l'API.
func FetchLocationsById(id string) (*models.Location, error) {
	url := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/locations/%s", id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var locations models.Location
	err = json.NewDecoder(resp.Body).Decode(&locations)
	if err != nil {
		return nil, err
	}
	return &locations, nil
}

func FetchDates() ([]models.Date, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/dates")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var dates []models.Date
	err = json.NewDecoder(resp.Body).Decode(&dates)
	if err != nil {
		return nil, err
	}
	return dates, nil
}

// FetchDates récupère les données des dates de concert depuis l'API.
func FetchDatesById(id string) (*models.Date, error) {
	url := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/dates/%s", id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var dates models.Date
	err = json.NewDecoder(resp.Body).Decode(&dates)
	if err != nil {
		return nil, err
	}
	return &dates, nil
}

// FetchRelations récupère les données des relations entre les artistes, dates et lieux de concert depuis l'API.
func FetchRelationsById(id string) (*models.Relation, error) {
	url := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/relation/%s", id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var relations models.Relation
	err = json.NewDecoder(resp.Body).Decode(&relations)
	if err != nil {
		return nil, err
	}
	return &relations, nil
}

// FetchArtistDetails récupère les détails d'un artiste par ID.
func FetchArtistDetails(artistID string) (*models.Artist, error) {
	url := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/artists/%s", artistID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch artist details: status code %d", resp.StatusCode)
	}

	var artist models.Artist
	if err := json.NewDecoder(resp.Body).Decode(&artist); err != nil {
		return nil, err
	}

	return &artist, nil
}

func FitlerArtists(artists []models.Artist, creationDateMin, creationDateMax int, firstAlbumMin, firstAlbumMax, members, location string) []models.Artist {
	var filtered []models.Artist
	filtered = FilterArtistsByCreationDate(artists, creationDateMin, creationDateMax)
	filtered = FilterArtistsByFirstAlbum(filtered, firstAlbumMin, firstAlbumMax)
	filtered = FilterArtistsByMembers(filtered, members)
	if location != "" {
		var relation []models.Relation
		for _, artist := range filtered {
			relations, _ := FetchRelationsById(strconv.Itoa(artist.ID))
			relation = append(relation, *relations)
		}
		filtered = FilterArtistsByLocation(filtered, relation, location)
	}
	return filtered
}

// Filter by date ranges (both creation and first album)
func FilterArtistsByCreationDate(artists []models.Artist, creationDateMin, creationDateMax int) []models.Artist {
	var filtered []models.Artist
	for _, artist := range artists {
		if artist.CreationDate >= creationDateMin && artist.CreationDate <= creationDateMax {
			filtered = append(filtered, artist)
		}
	}
	return filtered
}

// Filter by date ranges (both creation and first album)
func FilterArtistsByFirstAlbum(artists []models.Artist, firstAlbumMin, firstAlbumMax string) []models.Artist {
	var filtered []models.Artist
	firstAlbumMinInt, _ := strconv.Atoi(firstAlbumMin)
	firstAlbumMaxInt, _ := strconv.Atoi(firstAlbumMax)
	for _, artist := range artists {
		artistFirstAlbum, _ := strconv.Atoi(getYearFromDate(artist.FirstAlbum))
		if artistFirstAlbum >= firstAlbumMinInt && artistFirstAlbum <= firstAlbumMaxInt {
			filtered = append(filtered, artist)
		}
	}
	return filtered
}

// Filter by members (e.g., "1", "2-3", "4-5", "6+")
func FilterArtistsByMembers(artists []models.Artist, members string) []models.Artist {
	if members == "" {
		return artists // No filter applied
	}

	var filtered []models.Artist
	memberRanges := strings.Split(members, ",") // Assuming members comes as a comma-separated string
	for _, artist := range artists {
		memberMatch := false
		for _, r := range memberRanges {
			switch r {
			case "1":
				if len(artist.Members) == 1 {
					memberMatch = true
				}
			case "2":
				if len(artist.Members) == 2 {
					memberMatch = true
				}
			case "3":
				if len(artist.Members) == 3 {
					memberMatch = true
				}
			case "4":
				if len(artist.Members) == 4 {
					memberMatch = true
				}
			case "5":
				if len(artist.Members) == 5 {
					memberMatch = true
				}
			case "6":
				if len(artist.Members) == 6 {
					memberMatch = true
				}
			case "7":
				if len(artist.Members) == 7 {
					memberMatch = true
				}
			case "8":
				if len(artist.Members) == 8 {
					memberMatch = true
				}
			}
		}
		if memberMatch {
			filtered = append(filtered, artist)
		}
	}
	return filtered
}

// FilterArtistsByLocation - filtre les artistes en fonction des lieux (en utilisant la structure Location et les IDs)
func FilterArtistsByLocation(artists []models.Artist, relations []models.Relation, location string) []models.Artist {

	location = strings.ToLower(location)
	var filteredArtists []models.Artist

	// Créer un map pour stocker les artistes associés aux locations filtrées
	artistIDs := make(map[int]bool)
	for _, relation := range relations {
		var cleanedLocation []string
		for k := range relation.DatesLocations {
			cleanedLocation = append(cleanedLocation, CleanLocationName(k))
		}
		if contains(cleanedLocation, location) {
			artistIDs[relation.ID] = true
		} else {
			artistIDs[relation.ID] = false
		}
	}

	// Filtrer les artistes en fonction des ID trouvés dans les locations
	for _, artist := range artists {
		if artistIDs[artist.ID] {
			filteredArtists = append(filteredArtists, artist)
		}
	}
	return filteredArtists
}

// CleanLocationName nettoie le nom d'une location en supprimant les underscores et tirets
// et en mettant en majuscule la première lettre de chaque mot.
func CleanLocationName(location string) string {
	// Remplacer les underscores et tirets par des espaces
	location = strings.ReplaceAll(location, "_", " ")
	location = strings.ReplaceAll(location, "-", " - ")

	// Mettre en majuscule la première lettre de chaque mot
	return strings.Title(location)
}

// Helper function to check if a slice contains a particular string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.Contains(strings.ToLower(s), item) {
			return true
		}
	}
	return false
}

func getYearFromDate(date string) string {
	// Utilise la fonction Split pour diviser la chaîne par "/"
	parts := strings.Split(date, "-")
	// L'année est la troisième partie, index 2
	if len(parts) == 3 {
		return parts[2]
	}
	return "" // Retourne une chaîne vide si la date est incorrecte
}

// Geocode fetches coordinates (latitude and longitude) for a given location using OpenStreetMap Nominatim API
func Geocode(city string, country string) (models.Coordinates, error) {
	var coords models.Coordinates
	// On inclut le pays pour des résultats plus précis
	query := fmt.Sprintf("%s, %s", city, country)
	escapedLocation := url.QueryEscape(query)
	apiURL := fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=json&limit=1", escapedLocation)

	resp, err := http.Get(apiURL)
	if err != nil {
		return coords, err
	}
	defer resp.Body.Close()

	var results []struct {
		Lat string `json:"lat"`
		Lon string `json:"lon"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return coords, err
	}

	if len(results) > 0 {
		lat, err := strconv.ParseFloat(results[0].Lat, 64)
		if err != nil {
			return coords, err
		}
		lon, err := strconv.ParseFloat(results[0].Lon, 64)
		if err != nil {
			return coords, err
		}
		coords.Latitude = lat
		coords.Longitude = lon
	}

	return coords, nil
}

func SplitCityAndCountry(location string) (string, string, error) {
	// Divise la chaîne sur la virgule et supprime les espaces autour
	parts := strings.Split(location, "-")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid location format: %s", location)
	}

	// Nettoyer les espaces autour de la ville et du pays
	city := strings.TrimSpace(parts[0])
	country := strings.TrimSpace(parts[1])

	return city, country, nil
}
