# Groupie Trackers

Groupie Trackers is a web application that allows users to track musical artists, their albums, and concerts. The application provides an intuitive interface to browse artist information, tour dates, and other relevant details.

## Features

- **Artist List**: View a complete list of artists with their images and basic information.
- **Artist Details**: Access detailed information about each artist, including their albums, band members, and tour dates.
- **Search**: Search for artists by : name, creation date, first album date, member name, location
- **Filters**: Search for artists by filters : Range of creation date, range of first album date, number of members, location
- **Concert Information**: View concert dates and locations for each artist.
- **Concert Geolocalization**: View geolocalization of concert dates and locations on a map
- **Error Handling**: Custom error pages to handle common errors (404, 500, etc.).

## Technologies Used

- **Backend**: [Go](https://golang.org/)
- **Frontend**: HTML, CSS, [Bootstrap](https://getbootstrap.com/)
- **Templates**: Go HTML templates
- **Data Management**: Go structs for representing artist and concert data.

## Installation

### Prerequisites

- Go (version 1.16 or higher)
- A modern web browser

### Installation Steps

1. Clone the repository:

    ```sh
    git clone https://zone01normandie.org/git/rsavary/groupie-tracker
    cd groupie-trackers
    ```

2. Install Go dependencies (if any):

    ```sh
    go mod tidy
    ```

3. Run the server:

    ```sh
    go run main.go
    ```

4. Open your web browser and navigate to:

    ```
    http://localhost:8080/home
    ```

## Project Structure

- **/cmd**: Contains files to start the application.
- **/models**: Definitions of the data structures used (artists, concerts, etc.).
- **/handlers**: HTTP route handlers.
- **/templates**: HTML templates for rendering pages.
- **/static**: Static files (CSS, JavaScript, images).

## Usage

1. **Home**: Visit the homepage to see a list of artists.
2. **Artist Details**: Click on an artist to view more details, including albums and tour dates.
3. **Navigation**: Use the navigation menu to return to the homepage or access other sections.

## Customization

### Modify the CSS

You can customize the style of the application by modifying the CSS files in the `/static` folder.

### Add Artists

Artist data can be added or modified in the data models defined in `/models`.

## Error Handling

The application handles common errors such as 404 (page not found) and 500 (internal server error) with custom error pages.

## Contributions

Contributions are welcome! Please submit a pull request or open an issue to discuss any changes you would like to make.
