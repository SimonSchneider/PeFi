package postgres

type (
	Config struct {
		DbName                    string
		User                      string
		Password                  string
		Host                      string
		Port                      string
		Sslmode                   string
		Fallback_application_name string
		Connect_timeout           string
		Sslcert                   string
		Sslkey                    string
		Sslrootcert               string
	}
)

func getConnectionString(c *Config) string {
	s := ""
	if c.DbName != "" {
		s += " dbname=" + c.DbName
	}
	if c.User != "" {
		s += " user=" + c.User
	}
	if c.Password != "" {
		s += " password=" + c.Password
	}
	if c.Host != "" {
		s += " host=" + c.Host
	}
	if c.Port != "" {
		s += " port=" + c.Port
	}
	if c.Sslmode != "" {
		s += " sslmode=" + c.Sslmode
	}
	return s
}
