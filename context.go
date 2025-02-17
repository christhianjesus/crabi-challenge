package main

import "os"

type Context struct {
	jwtKey   string
	httpPort string
	mongoURL string
	pldURL   string
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

func (c *Context) GetMongoURL() string {
	return c.mongoURL
}

func GetContext() *Context {
	return &Context{
		jwtKey:   os.Getenv("JWT_KEY"),
		httpPort: os.Getenv("HTTP_PORT"),
		mongoURL: os.Getenv("MONGODB_URL"),
		pldURL:   os.Getenv("PLD_URL"),
	}
}
