# goioc/web: Web Framework for Go, based on goioc/di
[![goioc](https://habrastorage.org/webt/ym/pu/dc/ympudccm7j7a3qex_jjroxgsiwg.png)](https://github.com/goioc)

![Go](https://github.com/goioc/web/workflows/Go/badge.svg)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/goioc/web/?tab=doc)
[![CodeFactor](https://www.codefactor.io/repository/github/goioc/web/badge)](https://www.codefactor.io/repository/github/goioc/web)
[![Go Report Card](https://goreportcard.com/badge/github.com/goioc/WEB)](https://goreportcard.com/report/github.com/goioc/web)
[![codecov](https://codecov.io/gh/goioc/web/branch/master/graph/badge.svg)](https://codecov.io/gh/goioc/web)
[![Dependabot Status](https://api.dependabot.com/badges/status?host=github&repo=goioc/web)](https://dependabot.com)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=goioc_web&metric=alert_status)](https://sonarcloud.io/dashboard?id=goioc_web)
[![DeepSource](https://static.deepsource.io/deepsource-badge-light-mini.svg)](https://deepsource.io/gh/goioc/web/?ref=repository-badge)

## How this framework is different from others?

1. First of all, `goioc/web` is working using Dependency Injection and is based on [goioc/di](https://github.com/goioc/di), which is basically the IoC Container.
2. Secondly - and this is the most exciting part - web-endpoints in `goioc/web` can have (almost) arbitrary signature! 
No more `func(w http.ResponseWriter, r *http.Request)` handlers! If your endpoint receives a string and produces a binary stream, just declare it as is:

```go
...
func (e *endpoint) Hello(name string) io.Reader {
	return bytes.NewBufferString("Hello, " + name + "!")
}
...
```

Of course, you can still directly use `http.ResponseWriter` and `*http.Request`, if you like.

## Basic concepts

The main entity in `goioc/web` is the Endpoint, which is represented by the interface of the same name. Here's the example implementation:

```go
type endpoint struct {
}

func (e endpoint) HandlerFuncName() string {
	return "Hello"
}

func (e *endpoint) Hello(name string) io.Reader {
	return bytes.NewBufferString("Hello, " + name + "!")
}
```

`Endpoint` interface has one method that returns the name of the method that will be used as an endpoint.

In order for `goioc/web` to pick up this endpoint, it should be registered in the DI Container:

```go
_, _ = di.RegisterBean("endpoint", reflect.TypeOf((*endpoint)(nil)))
```

And then the container should be initialized (please, refer to the [goioc/di](https://github.com/goioc/di) documentation for more details):

```go
_ = di.InitializeContainer()
```

Finally, the web-server can be started, either using the built-in function:

```go
_ = web.ListenAndServe(":8080")
```
... or using returned `Router`
```go
router, _ := CreateRouter()
_ = http.ListenAndServe("", router)
```

## Routing

So, how does the framework know where to bind this endpoint to? 
For the routing functionality `goioc/web` leverages [gorilla/mux](https://github.com/gorilla/mux) library.
Don't worry: you don't have to cope with this library directly: `goioc/web` provides a set of convenient wrappers around it.
The wrappers are implemented as tags in the endpoint-structure. Let's slightly update our previous example:

```go
...
type endpoint struct {
	method interface{} `web.methods:"GET"`
	path   interface{} `web.path:"/hello"`
}
...
```

Now our endpoint is bound to the `GET` requests at the `/hello` path. Yes, it's that simple! 🙂

### `goioc/web` tags

| **Tag**       | **Value**                                 | **Example**                                           |
|---------------|-------------------------------------------|-------------------------------------------------------|
| `web.methods` | List of HTTP-methods.                     | `web.methods:"POST,PATCH"`                            |
| `web.path`    | URL sub-path. Can contain path variables. | `web.path:"/articles/{category}/{id:[0-9]+}"`         |
| `web.queries` | Key-value paris of the URL query part.    | `web.queries:"foo,bar,id,{id:[0-9]+}"`                |
| `web.headers` | Key-value paris of the request headers.   | `web.headers:"Content-Type,application/octet-stream"` |
| `web.matcher` | ID of the bean of type `*mux.MatcherFunc`.| `web.matcher:"matcher"`                               |

