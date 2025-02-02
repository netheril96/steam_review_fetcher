package steamreviewfetcher

type GameData struct {
	Success bool `json:"success"`
	Data    Game `json:"data"`
}

type Game struct {
	Type                string             `json:"type"`
	Name                string             `json:"name"`
	SteamAppID          int                `json:"steam_appid"`
	RequiredAge         int                `json:"required_age"`
	IsFree              bool               `json:"is_free"`
	DetailedDescription string             `json:"detailed_description"`
	AboutTheGame        string             `json:"about_the_game"`
	ShortDescription    string             `json:"short_description"`
	SupportedLanguages  string             `json:"supported_languages"`
	HeaderImage         string             `json:"header_image"`
	CapsuleImage        string             `json:"capsule_image"`
	CapsuleImageV5      string             `json:"capsule_imagev5"`
	Website             string             `json:"website"`
	Developers          []string           `json:"developers"`
	Publishers          []string           `json:"publishers"`
	PriceOverview       PriceOverview      `json:"price_overview"`
	Packages            []int              `json:"packages"`
	PackageGroups       []PackageGroup     `json:"package_groups"`
	Platforms           Platforms          `json:"platforms"`
	Categories          []Category         `json:"categories"`
	Genres              []Genre            `json:"genres"`
	Screenshots         []Screenshot       `json:"screenshots"`
	Movies              []Movie            `json:"movies"`
	Recommendations     Recommendations    `json:"recommendations"`
	ReleaseDate         ReleaseDate        `json:"release_date"`
	SupportInfo         SupportInfo        `json:"support_info"`
	Background          string             `json:"background"`
	BackgroundRaw       string             `json:"background_raw"`
	ContentDescriptors  ContentDescriptors `json:"content_descriptors"`
}

type Requirements struct {
	Minimum     string `json:"minimum"`
	Recommended string `json:"recommended"`
}

type PriceOverview struct {
	Currency         string `json:"currency"`
	Initial          int    `json:"initial"`
	Final            int    `json:"final"`
	DiscountPercent  int    `json:"discount_percent"`
	InitialFormatted string `json:"initial_formatted"`
	FinalFormatted   string `json:"final_formatted"`
}

type PackageGroup struct {
	Name                    string `json:"name"`
	Title                   string `json:"title"`
	Description             string `json:"description"`
	SelectionText           string `json:"selection_text"`
	SaveText                string `json:"save_text"`
	DisplayType             int    `json:"display_type"`
	IsRecurringSubscription string `json:"is_recurring_subscription"`
	Subs                    []Sub  `json:"subs"`
}

type Sub struct {
	PackageID                int    `json:"packageid"`
	PercentSavingsText       string `json:"percent_savings_text"`
	PercentSavings           int    `json:"percent_savings"`
	OptionText               string `json:"option_text"`
	OptionDescription        string `json:"option_description"`
	CanGetFreeLicense        string `json:"can_get_free_license"`
	IsFreeLicense            bool   `json:"is_free_license"`
	PriceInCentsWithDiscount int    `json:"price_in_cents_with_discount"`
}

type Platforms struct {
	Windows bool `json:"windows"`
	Mac     bool `json:"mac"`
	Linux   bool `json:"linux"`
}

type Category struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

type Genre struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

type Screenshot struct {
	ID            int    `json:"id"`
	PathThumbnail string `json:"path_thumbnail"`
	PathFull      string `json:"path_full"`
}

type Movie struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Thumbnail string `json:"thumbnail"`
	Webm      Webm   `json:"webm"`
	Mp4       Mp4    `json:"mp4"`
	Highlight bool   `json:"highlight"`
}

type Webm struct {
	Low  string `json:"480"`
	High string `json:"max"`
}

type Mp4 struct {
	Low  string `json:"480"`
	High string `json:"max"`
}

type Recommendations struct {
	Total int `json:"total"`
}

type ReleaseDate struct {
	ComingSoon bool   `json:"coming_soon"`
	Date       string `json:"date"` // Consider using time.Time and parsing the date
}

type SupportInfo struct {
	URL   string `json:"url"`
	Email string `json:"email"`
}

type ContentDescriptors struct {
	IDs   []int  `json:"ids"`
	Notes string `json:"notes"`
}

type Ratings struct {
	Dejus        RatingInfo `json:"dejus"`
	SteamGermany RatingInfo `json:"steam_germany"`
}

type RatingInfo struct {
	RatingGenerated string `json:"rating_generated"`
	Rating          string `json:"rating"`
	RequiredAge     string `json:"required_age"`
	Banned          string `json:"banned"`
	UseAgeGate      string `json:"use_age_gate"`
	Descriptors     string `json:"descriptors"`
}
