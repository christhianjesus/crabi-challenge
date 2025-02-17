package main

import "os"

type Context struct {
	jwtKey   string
	httpPort string
}

func (c *Context) GetJwtKey() []byte {
	return []byte(c.jwtKey)
}

func (c *Context) GetHttpPort() string {
	if c.httpPort == "" {
		return "8080"
	}

	return c.httpPort
}

func GetContext() *Context {
	return &Context{
		jwtKey:   os.Getenv("JWT_KEY"),
		httpPort: os.Getenv("HTTP_PORT"),
	}
}
