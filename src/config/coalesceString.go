package config

func CoalesceString(args ...*string) *string {
	for _, arg := range args {
		if arg != nil && *arg != "" {
			return arg
		}
	}
	return nil
}
