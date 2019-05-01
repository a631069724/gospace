package conf

type Config struct {
	Log    *Log
	Redis  *Redis
	Oracle *Oracle
}

type Log struct {
	Path string
}

type Redis struct {
	Uri     string //redis://:[password]@[ip]:[port]/[db]
	WriteId string
}

type Oracle struct {
	UserName string
	Pwd      string
	Source   string
}
