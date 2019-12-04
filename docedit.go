package main

import (
    "io/ioutil"
	"github.com/gin-gonic/gin"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func docEdit(c *gin.Context) {
	if getContext(c).User.ID == 0 {
		resp403(c)
		return
	}
	
	if !((getContext(c).User.Privileges & 3145727) > 0){
		resp403(c)
		return
	}
	
	if c.PostForm("data") == "" {
		addMessage(c, errorMessage{T(c, "Something went wrong.")})
		getSession(c).Save()
		c.Redirect(302, "/settings/rules")
		return
	}
	
	d1 := []byte(c.PostForm("data"))
    err := ioutil.WriteFile("./website-docs/en/rules.md", d1, 0775)
    check(err)

	addMessage(c, successMessage{T(c, "Successfully edited doc")})
	getSession(c).Save()
	c.Redirect(302, "/settings/doc")
}