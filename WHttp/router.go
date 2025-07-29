package WHttp

type HandlerFunc func(w ResponseWriter, r *Request)

type Router struct {
	routes map[string]map[string]HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]map[string]HandlerFunc),
	}
}

func (r *Router) GET(path string, handler HandlerFunc) {
	r.addRoute("GET", path, handler)
}

func (r *Router) POST(path string, handler HandlerFunc) {
	r.addRoute("POST", path, handler)
}

// addRoute 添加一个路由
func (r *Router) addRoute(method, path string, handler HandlerFunc) {
	if r.routes[method] == nil {
		r.routes[method] = make(map[string]HandlerFunc)
	}
	r.routes[method][path] = handler
}

func (r *Router) ServeHTTP(w ResponseWriter, req *Request) {
	// 获取请求路径
	path := req.Url.Path
	if path == "" {
		path = "/"
	}

	if methodRoutes, ok := r.routes[req.Method]; ok {
		if handler, ok := methodRoutes[path]; ok {
			handler(w, req)
			return
		}
	}

	w.SetStatus(404)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("404 Not Found: " + req.Method + " " + path))
}
