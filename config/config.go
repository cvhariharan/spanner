package config

type Config struct {
	ServiceName string
	ModulePath  string
	Database    struct {
		URI        string `default:"mongodb://mongoadmin:secret@localhost:27017"`
		DB         string `default:"spanner"`
		Collection string `default:"models"`
	}
	OAuth struct {
		Enable       bool
		ClientId     string
		ClientSecret string
		RedirectUrl  string
		ConfigUrl    string
	}
	Port string `default:"5000"`
}
