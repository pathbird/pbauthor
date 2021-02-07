package graphql

type CodexCategory struct {
	ID string `json:"id"`
	Name string `json:"name"`
}

type Course struct {
	ID string `json:"id"`
	Name string `json:"name"`
	CodexCategories []CodexCategory `json:"codexCategories"`
}

type CourseEdge struct {
	Course Course `json:"course"`
	Role string `json:"role"`
}

type User struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
}