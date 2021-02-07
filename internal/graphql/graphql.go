package graphql

type resp struct {
	Viewer struct {
		User struct {
			Name string `json:"name"`
		} `json:"user"`
	} `json:"viewer"`
}
