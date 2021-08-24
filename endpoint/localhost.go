package endpoint

type localHost struct {
	Marathon    func() string
	MesosMaster func() string
	Chronos     func() string
}

var LocalHost = &localHost{
	Marathon: func() string {
		return "0.0.0.0:8080"
	},
	MesosMaster: func() string {
		return "0.0.0.0:5050"
	},
	Chronos: func() string {
		return "0.0.0.0:4400"
	},
}
