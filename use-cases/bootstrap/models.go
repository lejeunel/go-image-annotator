package bootstrap

type Request struct {
	InitialAdminEmail    string
	InitialAdminPassword string
}

type Response struct {
	Skipped bool
}
