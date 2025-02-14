package main

import "os"

type Context struct {
	jwtKey string
}

func (c *Context) GetJwtKey() []byte {
	return []byte(c.jwtKey)
}

func GetContext() *Context {
	return &Context{
		jwtKey: os.Getenv("JWT_KEY"),
	}
}
