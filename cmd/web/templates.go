package main

import "snipbox.fernandobasso.dev/internal/models"

// templateData is used as a type-safe way to pass multiple pieces of data
// to templates (since templates can only take a single variable containing
// data). By using a struct, we can pass a "bag" of composite data as needed.
type templateData struct {
	Snippet models.Snippet
}
