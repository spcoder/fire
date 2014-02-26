package fire

var (
	_settings *Settings
)

const (
	Development = "development"
	QA          = "qa"
	UAT         = "uat"
	Staging     = "staging"
	Production  = "production"
)

type Settings struct {
	environment      string
	templateFileBase string
	staticFileBase   string
}

func SetEnvironment(s string) {
	_settings.environment = s
}

func Environment() string {
	return _settings.environment
}

func IsDevelopment() bool {
	return _settings.environment == Development
}

func IsProduction() bool {
	return _settings.environment == Production
}

func SetTemplateFileBase(b string) {
	_settings.templateFileBase = b
}

func TemplateFileBase() string {
	return _settings.templateFileBase
}

func SetStaticFileBase(b string) {
	_settings.staticFileBase = b
}

func StaticFileBase() string {
	return _settings.staticFileBase
}

func initSettings() {
	_settings = &Settings{
		templateFileBase: "templates",
		staticFileBase:   "www",
		environment:      Development,
	}
}
