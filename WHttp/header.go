package WHttp

type Header map[string]string

func (h Header) Set(key, value string) {
	h[key] = value
}

func (h Header) Get(key string) string {
	return h[key]
}
