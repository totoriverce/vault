pid_file = "./pidfile"

auto_auth {
	method {
		type = "aws"
		wrap_ttl = 300
		config = {
			role = "foobar"
		}
		max_backoff = "2m"
        initial_backoff = "5s"
	}

	sink {
		type = "file"
		config = {
			path = "/tmp/file-foo"
		}
	}
}
