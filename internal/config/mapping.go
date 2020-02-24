package config

func FindMappingForRecipient(user string, mappings []Mapping) *Mapping {
	for _, mapping := range mappings {
		if user == mapping.User {
			return &mapping
		}
	}
	return nil
}
