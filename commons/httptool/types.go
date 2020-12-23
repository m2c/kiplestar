package httptool

const (
	TAG_TYPE__QUERY  = "query"
	TAG_TYPE__BODY   = "body"
	TAG_TYPE__PATH   = "path"
	TAG_TYPE__HEADER = "header"
	TAG_TYPE__COOKIE = "cookie"
)

const (
	ContentTypeFormUrlencoded = "application/x-www-form-urlencoded"
	ContentTypeFormData       = "multipart/form-data"
	ContentTypeJson           = "application/json"
)

type RequestItem struct {
	Key   string
	Value string
}
