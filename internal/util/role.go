package util

const (
	ADMIN  = "admin"
	MEMBER = "member"
)

func IsSupportedRole(role string) bool {
	switch role {
	case ADMIN, MEMBER:
		return true
	}
	return false
}
