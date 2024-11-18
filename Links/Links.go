package links

type Link struct {
	Name string
	Url  string
}

func GetMenuLinks() []Link {
	Links := []Link{
		{
			Name: "Startseite",
			Url:  "/",
		},
		{
			Name: "Impressum",
			Url:  "/Impressum",
		},
		{
			Name: "Datenschutz",
			Url:  "/Datenschutz",
		},
	}

	// All the Links

	return Links
}
