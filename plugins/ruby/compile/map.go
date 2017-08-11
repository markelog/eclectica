package compile

// In order to not guess "-p..." suffixes with remote server
// It's easier just to map latest versions
var remoteVersions = map[string]string{
	"1.8.5": "ruby-1.8.5-p115.tar.gz",
	"1.8.6": "ruby-1.8.6-p420.tar.gz",
	"1.8.7": "ruby-1.8.7-p358.tar.gz",
	"1.9.0": "ruby-1.9-stable.tar.gz",
	"1.9.2": "ruby-1.9.2-p330.tar.gz",
	"1.9.3": "ruby-1.9.3-p551.tar.gz",
	"2.0.0": "ruby-2.0.0-p648.tar.gz",
}
