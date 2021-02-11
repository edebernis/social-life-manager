package models

// Location model. Describes a physical place in the world.
type Location struct {
	// Location ID. Must be unique.
	ID ID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" extensions:"x-order=1"`
	// Short descriptive name of the location, like "Home" or "Work".
	Name string `json:"name" example:"Home" extensions:"x-order=2"`
	// Full address of the location. Should contains at least street, postal code and city.
	Address string `json:"address" example:"1 rue de la Poste, 75001 Paris" extensions:"x-order=3"`
	// Location category foreign key.
	Category ID `json:"category_id" example:"550e8400-e29b-41d4-a716-446655440000" extensions:"x-order=4"`
	// User ID. Owner of the location.
	User ID `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" extensions:"x-order=5"`
}

// Locations is an array of locations
type Locations []*Location

// CreateLocation validates user input to create a new location
type CreateLocation struct {
	// Short descriptive name of the location, like "Home" or "Work".
	Name string `json:"name" example:"Home" binding:"required" extensions:"x-order=1"`
	// Full address of the location. Should contains at least street, postal code and city.
	Address string `json:"address" example:"1 rue de la Poste, 75001 Paris" binding:"required" extensions:"x-order=2"`
	// Location category foreign key.
	Category ID `json:"category_id" example:"550e8400-e29b-41d4-a716-446655440000" binding:"required" extensions:"x-order=3"`
}

// GetLocations validates user input to get locations
type GetLocations struct {
	// Location category foreign key.
	Category string `form:"category_id" example:"550e8400-e29b-41d4-a716-446655440000" binding:"omitempty,uuid" extensions:"x-order=1"`
}

// GetLocationByID validates user input to get a specific location using its id
// TODO: replace string by ID type and let the binding engine do the parsing work.
// Open issue: https://github.com/gin-gonic/gin/issues/2631
type GetLocationByID struct {
	// ID of the location.
	ID string `uri:"id" example:"550e8400-e29b-41d4-a716-446655440000" binding:"required,uuid" extensions:"x-order=1"`
}

// UpdateLocationQuery validates query user input to update an existing location
type UpdateLocationQuery struct {
	// ID of the location.
	ID string `uri:"id" example:"550e8400-e29b-41d4-a716-446655440000" binding:"required,uuid" extensions:"x-order=1"`
}

// UpdateLocation validates body user input to update an existing location
type UpdateLocation struct {
	// Short descriptive name of the location, like "Home" or "Work".
	Name string `json:"name" example:"Home" binding:"required_without_all" extensions:"x-order=1"`
	// Full address of the location. Should contains at least street, postal code and city.
	Address string `json:"address" example:"1 rue de la Poste, 75001 Paris" binding:"required_without_all" extensions:"x-order=2"`
	// Location category foreign key.
	Category ID `json:"category_id" example:"550e8400-e29b-41d4-a716-446655440000" binding:"required_without_all" extensions:"x-order=3"`
}

// DeleteLocation validates user input to delete an existing location
type DeleteLocation struct {
	// ID of the location.
	ID string `uri:"id" example:"550e8400-e29b-41d4-a716-446655440000" binding:"required,uuid" extensions:"x-order=1"`
}

// NewLocation creates a new user location
func NewLocation(id ID, name, address string, category ID, user ID) *Location {
	return &Location{
		id,
		name,
		address,
		category,
		user,
	}
}

// Category model. Allows to describe what the location is used for, such as sport, work, living, etc.
type Category struct {
	// Category ID. Must be unique.
	ID ID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" extensions:"x-order=1"`
	// Short descriptive name of the category. Like "Homes" or "Tennis Center".
	Name string `json:"name" example:"Homes" extensions:"x-order=2"`
}

// Categories is an array of categories
type Categories []*Category

// CreateCategory validates user input to create a new category
type CreateCategory struct {
	// Short descriptive name of the category. Like "Homes" or "Tennis Center".
	Name string `json:"name" example:"Homes" binding:"required" extensions:"x-order=1"`
}

// GetCategoryByID validates user input to get a specific category using its id
type GetCategoryByID struct {
	// ID of the category.
	ID string `uri:"id" example:"550e8400-e29b-41d4-a716-446655440000" binding:"required,uuid" extensions:"x-order=1"`
}

// UpdateCategoryQuery validates query user input to update an existing category
type UpdateCategoryQuery struct {
	// ID of the category.
	ID string `uri:"id" example:"550e8400-e29b-41d4-a716-446655440000" binding:"required,uuid" extensions:"x-order=1"`
}

// UpdateCategory validates body user input to update an existing category
type UpdateCategory struct {
	// Short descriptive name of the category, like "Homes" or "Sport".
	Name string `json:"name" example:"Homes" binding:"required_without_all" extensions:"x-order=1"`
}

// DeleteCategory validates user input to delete an existing category
type DeleteCategory struct {
	// ID of the category.
	ID string `uri:"id" example:"550e8400-e29b-41d4-a716-446655440000" binding:"required,uuid" extensions:"x-order=1"`
}

// NewCategory creates a new location category
func NewCategory(id ID, name string) *Category {
	return &Category{
		id,
		name,
	}
}
