package console

var (
	white = "\x1b[97m"
	black = "\x1b[90m"

	logs = map[string]string{
		"DBG": "\x1b[94m",
		"ERR": "\x1b[91m",
		"SCC": "\x1b[92m",
		"FLD": "\x1b[93m",
		"INP": "\x1b[96m",
		"WRN": "\x1b[93m",
		"HUM": "\x1b[95m",
		"WBS": "\x1b[35m",
		"SLV": "\x1b[35m",
	}
)
