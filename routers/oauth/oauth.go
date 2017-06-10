package oauth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/RangelReale/osin"
	"github.com/felipeweb/osin-mysql"
	"github.com/gin-gonic/gin"
)

var osinServer *osin.Server

var rh RequestHandler

// Initialise initialises the oauth router.
func Initialise(handler RequestHandler) error {
	store := mysql.New(handler.GetDB(), "osin_")

	rh = handler

	err := store.CreateSchemas()
	if err != nil {
		return err
	}

	config := osin.NewServerConfig()
	config.AllowClientSecretInParams = true
	config.AllowGetAccessRequest = true
	// Hmm... Wondering why we're making access tokens everlasting, and
	// disabling refresh tokens? http://telegra.ph/On-refresh-tokens-06-10
	config.AccessExpiration = 0
	osinServer = osin.NewServer(config, store)
	return nil
}

// Client is the data passed to the template of the authorization page.
type Client struct {
	ID             string
	Name           string
	CreatorID      int
	Avatar         string
	Authorizations []string
}

func clientFromOsinClient(oc osin.Client) Client {
	userDataRaw := oc.GetUserData().(string)

	var userData [3]string
	json.Unmarshal([]byte(userDataRaw), &userData)

	cid, _ := strconv.Atoi(userData[1])

	return Client{
		ID:        oc.GetId(),
		Name:      userData[0],
		CreatorID: cid,
		Avatar:    userData[2],
		// Authorizations is managed by caller
	}
}

// RequestHandler contains the functions to which Oauth will delegate.
type RequestHandler interface {
	// Checks whether the user is logged in. If they are, it returns true. Otherwise, it redirects
	// the user to the login page.
	CheckLoggedInOrRedirect(c *gin.Context) bool
	// DisplayAuthorizePage renders the page where the user can validate the request for authorization
	DisplayAuthorizePage(client Client, c *gin.Context)
	// CheckCSRF is a function that checks the POSTed parameter `csrf` is valid.
	CheckCSRF(c *gin.Context, s string) bool
	// GetDB retrieves MySQL's DB.
	GetDB() *sql.DB
	// GetUserID retrieves the ID of the currently logged in user.
	GetUserID(c *gin.Context) int
}

// Authorize handles a request for user authorization
func Authorize(c *gin.Context) {
	resp := osinServer.NewResponse()
	defer resp.Close()

	// first we let osinserver handle the authorize request
	if ar := osinServer.HandleAuthorizeRequest(resp, c.Request); ar != nil {
		// we then make sure to be logged in
		if !rh.CheckLoggedInOrRedirect(c) {
			return
		}

		// and show the authorization page
		client := clientFromOsinClient(ar.Client)
		client.Authorizations = safeScopes(strings.Split(ar.Scope, " "))
		if c.PostForm("appid") != ar.Client.GetId() || !rh.CheckCSRF(c, c.PostForm("csrf")) {
			rh.DisplayAuthorizePage(client, c)
			return
		}

		// all good, authorization succeded
		ar.Authorized = c.PostForm("approve") == "1"
		ar.UserData = strconv.Itoa(rh.GetUserID(c))
		osinServer.FinishAuthorizeRequest(resp, c.Request, ar)
	}
	if resp.IsError && resp.InternalError != nil {
		fmt.Println(resp.InternalError)
	}
	osin.OutputJSON(resp, c.Writer, c.Request)
}

// an "identify" scope is taken for granted
var scopes = [...]string{
	"read_confidential",
	"write",
}

// safeScopes removes duplicates and invalid elements from the raw scopes
func safeScopes(rawScopes []string) []string {
	retScopes := make([]string, 0, len(scopes))
	for _, el := range scopes {
		for _, rawScope := range rawScopes {
			if rawScope == el {
				retScopes = append(retScopes, rawScope)
				break
			}
		}
	}
	return retScopes
}

// Token handles a request from a client to obtain an access token.
func Token(c *gin.Context) {
	resp := osinServer.NewResponse()
	defer resp.Close()

	c.Request.ParseForm()

	if ar := osinServer.HandleAccessRequest(resp, c.Request); ar != nil {
		ar.Authorized = true
		ar.GenerateRefresh = false
		osinServer.FinishAccessRequest(resp, c.Request, ar)
	}
	if resp.IsError && resp.InternalError != nil {
		fmt.Printf("ERROR: %s\n", resp.InternalError)
	}
	delete(resp.Output, "expire_in")
	osin.OutputJSON(resp, c.Writer, c.Request)
}
