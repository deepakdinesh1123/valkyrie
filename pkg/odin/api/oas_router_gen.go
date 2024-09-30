// Code generated by ogen, DO NOT EDIT.

package api

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/ogen-go/ogen/uri"
)

func (s *Server) cutPrefix(path string) (string, bool) {
	prefix := s.cfg.Prefix
	if prefix == "" {
		return path, true
	}
	if !strings.HasPrefix(path, prefix) {
		// Prefix doesn't match.
		return "", false
	}
	// Cut prefix from the path.
	return strings.TrimPrefix(path, prefix), true
}

// ServeHTTP serves http request as defined by OpenAPI v3 specification,
// calling handler that matches the path or returning not found error.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	elem := r.URL.Path
	elemIsEscaped := false
	if rawPath := r.URL.RawPath; rawPath != "" {
		if normalized, ok := uri.NormalizeEscapedPath(rawPath); ok {
			elem = normalized
			elemIsEscaped = strings.ContainsRune(elem, '%')
		}
	}

	elem, ok := s.cutPrefix(elem)
	if !ok || len(elem) == 0 {
		s.notFound(w, r)
		return
	}
	args := [1]string{}

	// Static code generated router with unwrapped path search.
	switch {
	default:
		if len(elem) == 0 {
			break
		}
		switch elem[0] {
		case '/': // Prefix: "/"
			origElem := elem
			if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
				elem = elem[l:]
			} else {
				break
			}

			if len(elem) == 0 {
				break
			}
			switch elem[0] {
			case 'e': // Prefix: "execution"
				origElem := elem
				if l := len("execution"); len(elem) >= l && elem[0:l] == "execution" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case '/': // Prefix: "/config"
					origElem := elem
					if l := len("/config"); len(elem) >= l && elem[0:l] == "/config" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						// Leaf node.
						switch r.Method {
						case "GET":
							s.handleGetExecutionConfigRequest([0]string{}, elemIsEscaped, w, r)
						default:
							s.notAllowed(w, r, "GET")
						}

						return
					}

					elem = origElem
				case 's': // Prefix: "s"
					origElem := elem
					if l := len("s"); len(elem) >= l && elem[0:l] == "s" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						switch r.Method {
						case "GET":
							s.handleGetAllExecutionsRequest([0]string{}, elemIsEscaped, w, r)
						default:
							s.notAllowed(w, r, "GET")
						}

						return
					}
					switch elem[0] {
					case '/': // Prefix: "/"
						origElem := elem
						if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							break
						}
						switch elem[0] {
						case 'e': // Prefix: "execute"
							origElem := elem
							if l := len("execute"); len(elem) >= l && elem[0:l] == "execute" {
								elem = elem[l:]
							} else {
								break
							}

							if len(elem) == 0 {
								// Leaf node.
								switch r.Method {
								case "POST":
									s.handleExecuteRequest([0]string{}, elemIsEscaped, w, r)
								default:
									s.notAllowed(w, r, "POST")
								}

								return
							}

							elem = origElem
						case 'j': // Prefix: "jobs/"
							origElem := elem
							if l := len("jobs/"); len(elem) >= l && elem[0:l] == "jobs/" {
								elem = elem[l:]
							} else {
								break
							}

							// Param: "JobId"
							// Leaf parameter
							args[0] = elem
							elem = ""

							if len(elem) == 0 {
								// Leaf node.
								switch r.Method {
								case "DELETE":
									s.handleDeleteExecutionJobRequest([1]string{
										args[0],
									}, elemIsEscaped, w, r)
								case "GET":
									s.handleGetExecutionJobByIdRequest([1]string{
										args[0],
									}, elemIsEscaped, w, r)
								case "PUT":
									s.handleCancelExecutionJobRequest([1]string{
										args[0],
									}, elemIsEscaped, w, r)
								default:
									s.notAllowed(w, r, "DELETE,GET,PUT")
								}

								return
							}

							elem = origElem
						case 'w': // Prefix: "workers"
							origElem := elem
							if l := len("workers"); len(elem) >= l && elem[0:l] == "workers" {
								elem = elem[l:]
							} else {
								break
							}

							if len(elem) == 0 {
								switch r.Method {
								case "GET":
									s.handleGetExecutionWorkersRequest([0]string{}, elemIsEscaped, w, r)
								default:
									s.notAllowed(w, r, "GET")
								}

								return
							}
							switch elem[0] {
							case '/': // Prefix: "/"
								origElem := elem
								if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
									elem = elem[l:]
								} else {
									break
								}

								// Param: "workerId"
								// Leaf parameter
								args[0] = elem
								elem = ""

								if len(elem) == 0 {
									// Leaf node.
									switch r.Method {
									case "DELETE":
										s.handleDeleteExecutionWorkerRequest([1]string{
											args[0],
										}, elemIsEscaped, w, r)
									default:
										s.notAllowed(w, r, "DELETE")
									}

									return
								}

								elem = origElem
							}

							elem = origElem
						}
						// Param: "execId"
						// Leaf parameter
						args[0] = elem
						elem = ""

						if len(elem) == 0 {
							// Leaf node.
							switch r.Method {
							case "GET":
								s.handleGetExecutionResultByIdRequest([1]string{
									args[0],
								}, elemIsEscaped, w, r)
							default:
								s.notAllowed(w, r, "GET")
							}

							return
						}

						elem = origElem
					}

					elem = origElem
				}

				elem = origElem
			case 'j': // Prefix: "jobs/"
				origElem := elem
				if l := len("jobs/"); len(elem) >= l && elem[0:l] == "jobs/" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case 'e': // Prefix: "execution"
					origElem := elem
					if l := len("execution"); len(elem) >= l && elem[0:l] == "execution" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						// Leaf node.
						switch r.Method {
						case "GET":
							s.handleGetAllExecutionJobsRequest([0]string{}, elemIsEscaped, w, r)
						default:
							s.notAllowed(w, r, "GET")
						}

						return
					}

					elem = origElem
				}
				// Param: "JobId"
				// Match until "/"
				idx := strings.IndexByte(elem, '/')
				if idx < 0 {
					idx = len(elem)
				}
				args[0] = elem[:idx]
				elem = elem[idx:]

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case '/': // Prefix: "/executions"
					origElem := elem
					if l := len("/executions"); len(elem) >= l && elem[0:l] == "/executions" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						// Leaf node.
						switch r.Method {
						case "GET":
							s.handleGetExecutionsForJobRequest([1]string{
								args[0],
							}, elemIsEscaped, w, r)
						default:
							s.notAllowed(w, r, "GET")
						}

						return
					}

					elem = origElem
				}

				elem = origElem
			case 'l': // Prefix: "languages"
				origElem := elem
				if l := len("languages"); len(elem) >= l && elem[0:l] == "languages" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					// Leaf node.
					switch r.Method {
					case "GET":
						s.handleGetAllLanguagesRequest([0]string{}, elemIsEscaped, w, r)
					default:
						s.notAllowed(w, r, "GET")
					}

					return
				}

				elem = origElem
			case 'v': // Prefix: "version"
				origElem := elem
				if l := len("version"); len(elem) >= l && elem[0:l] == "version" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					// Leaf node.
					switch r.Method {
					case "GET":
						s.handleGetVersionRequest([0]string{}, elemIsEscaped, w, r)
					default:
						s.notAllowed(w, r, "GET")
					}

					return
				}

				elem = origElem
			}

			elem = origElem
		}
	}
	s.notFound(w, r)
}

// Route is route object.
type Route struct {
	name        string
	summary     string
	operationID string
	pathPattern string
	count       int
	args        [1]string
}

// Name returns ogen operation name.
//
// It is guaranteed to be unique and not empty.
func (r Route) Name() string {
	return r.name
}

// Summary returns OpenAPI summary.
func (r Route) Summary() string {
	return r.summary
}

// OperationID returns OpenAPI operationId.
func (r Route) OperationID() string {
	return r.operationID
}

// PathPattern returns OpenAPI path.
func (r Route) PathPattern() string {
	return r.pathPattern
}

// Args returns parsed arguments.
func (r Route) Args() []string {
	return r.args[:r.count]
}

// FindRoute finds Route for given method and path.
//
// Note: this method does not unescape path or handle reserved characters in path properly. Use FindPath instead.
func (s *Server) FindRoute(method, path string) (Route, bool) {
	return s.FindPath(method, &url.URL{Path: path})
}

// FindPath finds Route for given method and URL.
func (s *Server) FindPath(method string, u *url.URL) (r Route, _ bool) {
	var (
		elem = u.Path
		args = r.args
	)
	if rawPath := u.RawPath; rawPath != "" {
		if normalized, ok := uri.NormalizeEscapedPath(rawPath); ok {
			elem = normalized
		}
		defer func() {
			for i, arg := range r.args[:r.count] {
				if unescaped, err := url.PathUnescape(arg); err == nil {
					r.args[i] = unescaped
				}
			}
		}()
	}

	elem, ok := s.cutPrefix(elem)
	if !ok {
		return r, false
	}

	// Static code generated router with unwrapped path search.
	switch {
	default:
		if len(elem) == 0 {
			break
		}
		switch elem[0] {
		case '/': // Prefix: "/"
			origElem := elem
			if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
				elem = elem[l:]
			} else {
				break
			}

			if len(elem) == 0 {
				break
			}
			switch elem[0] {
			case 'e': // Prefix: "execution"
				origElem := elem
				if l := len("execution"); len(elem) >= l && elem[0:l] == "execution" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case '/': // Prefix: "/config"
					origElem := elem
					if l := len("/config"); len(elem) >= l && elem[0:l] == "/config" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						// Leaf node.
						switch method {
						case "GET":
							r.name = "GetExecutionConfig"
							r.summary = "Get execution config"
							r.operationID = "getExecutionConfig"
							r.pathPattern = "/execution/config"
							r.args = args
							r.count = 0
							return r, true
						default:
							return
						}
					}

					elem = origElem
				case 's': // Prefix: "s"
					origElem := elem
					if l := len("s"); len(elem) >= l && elem[0:l] == "s" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						switch method {
						case "GET":
							r.name = "GetAllExecutions"
							r.summary = "Get all executions"
							r.operationID = "getAllExecutions"
							r.pathPattern = "/executions"
							r.args = args
							r.count = 0
							return r, true
						default:
							return
						}
					}
					switch elem[0] {
					case '/': // Prefix: "/"
						origElem := elem
						if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							break
						}
						switch elem[0] {
						case 'e': // Prefix: "execute"
							origElem := elem
							if l := len("execute"); len(elem) >= l && elem[0:l] == "execute" {
								elem = elem[l:]
							} else {
								break
							}

							if len(elem) == 0 {
								// Leaf node.
								switch method {
								case "POST":
									r.name = "Execute"
									r.summary = "Execute a script"
									r.operationID = "execute"
									r.pathPattern = "/executions/execute"
									r.args = args
									r.count = 0
									return r, true
								default:
									return
								}
							}

							elem = origElem
						case 'j': // Prefix: "jobs/"
							origElem := elem
							if l := len("jobs/"); len(elem) >= l && elem[0:l] == "jobs/" {
								elem = elem[l:]
							} else {
								break
							}

							// Param: "JobId"
							// Leaf parameter
							args[0] = elem
							elem = ""

							if len(elem) == 0 {
								// Leaf node.
								switch method {
								case "DELETE":
									r.name = "DeleteExecutionJob"
									r.summary = "Delete execution job"
									r.operationID = "deleteExecutionJob"
									r.pathPattern = "/executions/jobs/{JobId}"
									r.args = args
									r.count = 1
									return r, true
								case "GET":
									r.name = "GetExecutionJobById"
									r.summary = "Get execution job"
									r.operationID = "getExecutionJobById"
									r.pathPattern = "/executions/jobs/{JobId}"
									r.args = args
									r.count = 1
									return r, true
								case "PUT":
									r.name = "CancelExecutionJob"
									r.summary = "Cancel Execution Job"
									r.operationID = "cancelExecutionJob"
									r.pathPattern = "/executions/jobs/{JobId}"
									r.args = args
									r.count = 1
									return r, true
								default:
									return
								}
							}

							elem = origElem
						case 'w': // Prefix: "workers"
							origElem := elem
							if l := len("workers"); len(elem) >= l && elem[0:l] == "workers" {
								elem = elem[l:]
							} else {
								break
							}

							if len(elem) == 0 {
								switch method {
								case "GET":
									r.name = "GetExecutionWorkers"
									r.summary = "Get all execution workers"
									r.operationID = "getExecutionWorkers"
									r.pathPattern = "/executions/workers"
									r.args = args
									r.count = 0
									return r, true
								default:
									return
								}
							}
							switch elem[0] {
							case '/': // Prefix: "/"
								origElem := elem
								if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
									elem = elem[l:]
								} else {
									break
								}

								// Param: "workerId"
								// Leaf parameter
								args[0] = elem
								elem = ""

								if len(elem) == 0 {
									// Leaf node.
									switch method {
									case "DELETE":
										r.name = "DeleteExecutionWorker"
										r.summary = "Delete execution worker"
										r.operationID = "deleteExecutionWorker"
										r.pathPattern = "/executions/workers/{workerId}"
										r.args = args
										r.count = 1
										return r, true
									default:
										return
									}
								}

								elem = origElem
							}

							elem = origElem
						}
						// Param: "execId"
						// Leaf parameter
						args[0] = elem
						elem = ""

						if len(elem) == 0 {
							// Leaf node.
							switch method {
							case "GET":
								r.name = "GetExecutionResultById"
								r.summary = "Get execution result by id"
								r.operationID = "getExecutionResultById"
								r.pathPattern = "/executions/{execId}"
								r.args = args
								r.count = 1
								return r, true
							default:
								return
							}
						}

						elem = origElem
					}

					elem = origElem
				}

				elem = origElem
			case 'j': // Prefix: "jobs/"
				origElem := elem
				if l := len("jobs/"); len(elem) >= l && elem[0:l] == "jobs/" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case 'e': // Prefix: "execution"
					origElem := elem
					if l := len("execution"); len(elem) >= l && elem[0:l] == "execution" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						// Leaf node.
						switch method {
						case "GET":
							r.name = "GetAllExecutionJobs"
							r.summary = "Get all execution jobs"
							r.operationID = "getAllExecutionJobs"
							r.pathPattern = "/jobs/execution"
							r.args = args
							r.count = 0
							return r, true
						default:
							return
						}
					}

					elem = origElem
				}
				// Param: "JobId"
				// Match until "/"
				idx := strings.IndexByte(elem, '/')
				if idx < 0 {
					idx = len(elem)
				}
				args[0] = elem[:idx]
				elem = elem[idx:]

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case '/': // Prefix: "/executions"
					origElem := elem
					if l := len("/executions"); len(elem) >= l && elem[0:l] == "/executions" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						// Leaf node.
						switch method {
						case "GET":
							r.name = "GetExecutionsForJob"
							r.summary = "Get executions of given job"
							r.operationID = "getExecutionsForJob"
							r.pathPattern = "/jobs/{JobId}/executions"
							r.args = args
							r.count = 1
							return r, true
						default:
							return
						}
					}

					elem = origElem
				}

				elem = origElem
			case 'l': // Prefix: "languages"
				origElem := elem
				if l := len("languages"); len(elem) >= l && elem[0:l] == "languages" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					// Leaf node.
					switch method {
					case "GET":
						r.name = "GetAllLanguages"
						r.summary = "Get all languages"
						r.operationID = "getAllLanguages"
						r.pathPattern = "/languages"
						r.args = args
						r.count = 0
						return r, true
					default:
						return
					}
				}

				elem = origElem
			case 'v': // Prefix: "version"
				origElem := elem
				if l := len("version"); len(elem) >= l && elem[0:l] == "version" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					// Leaf node.
					switch method {
					case "GET":
						r.name = "GetVersion"
						r.summary = "Get version"
						r.operationID = "getVersion"
						r.pathPattern = "/version"
						r.args = args
						r.count = 0
						return r, true
					default:
						return
					}
				}

				elem = origElem
			}

			elem = origElem
		}
	}
	return r, false
}
