package config

type Config struct {
	ServiceName string `fig:"ServiceName,required"`
	ModulePath  string
	Database    struct {
		URI        string `fig:"uri,default:mongodb://mongoadmin:secret@localhost:27017`
		DB         string `fig:"db,default:spanner"`
		Collection string `fig:"collection,default:models"`
	}
}
