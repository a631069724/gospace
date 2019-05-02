package conf

type Config struct {
	RunSql *RunSql
	Log    *Log
	Redis  *Redis
	Oracle *Oracle
	Tick   int
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

type RunSql struct {
	Sql []string
}
