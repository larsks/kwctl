package formatters

import "reflect"

// HeadersFromStruct extracts field names from a struct and returns them as a slice of strings.
// It checks for a 'header' struct tag first, and falls back to the field name if not present.
// Only exported (capitalized) fields are included.
func HeadersFromStruct(v any) []string {
	headers := []string{}

	// Get the type of the value
	t := reflect.TypeOf(v)

	// Handle pointer types
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Only process struct types
	if t.Kind() != reflect.Struct {
		return headers
	}

	// Iterate through all fields
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Only include exported fields (fields that start with a capital letter)
		if !field.IsExported() {
			continue
		}

		// Check for 'header' tag first
		if headerTag := field.Tag.Get("header"); headerTag != "" {
			headers = append(headers, headerTag)
		} else {
			// Fall back to field name
			headers = append(headers, field.Name)
		}
	}

	return headers
}
