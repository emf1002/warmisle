package service

func validPriority(p string) bool {
	return p == "normal" || p == "important" || p == "urgent"
}
