package fire

var (
	_settings *Settings
)

type Settings struct {
	templateFileBase string
	staticFileBase   string
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
	}
}
