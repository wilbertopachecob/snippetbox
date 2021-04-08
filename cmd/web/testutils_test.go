package main

import (
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
	"time"
	"wilbertopachecob/snippetbox/pkg/models/mock"

	"github.com/golangcollege/sessions"
)

// Create a newTestApplication helper which returns an instance of our
// application struct containing mocked dependencies.
func newTestApplication(t *testing.T) *application {

	templateCache, err := newTemplateCache("./../../ui/html/")
	if err != nil {
		t.Fatal(err)
	}

	// Create a session manager instance, with the same settings as production.
	session := sessions.New([]byte("s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge"))
	session.Lifetime = 1 * time.Hour
	session.Secure = true

	return &application{
		infolog:       log.New(ioutil.Discard, "", 0),
		errorlog:      log.New(ioutil.Discard, "", 0),
		snippets:      &mock.SnippetModel{},
		users:         &mock.UserModel{},
		templateCache: templateCache,
		session:       session,
	}
}

// Define a custom testServer type which anonymously embeds a httptest.Server
// instance.
type testServer struct {
	*httptest.Server
}

// Create a newTestServer helper which initalizes and returns a new instance
// of our custom testServer type.
func newTestServer(t *testing.T, h http.Handler) *testServer {
	tls := httptest.NewTLSServer(h)

	// Initialize a new cookie jar.
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	// Add the cookie jar to the client, so that response cookies are stored
	// and then sent with subsequent requests.
	tls.Client().Jar = jar
	// Disable redirect-following for the client. Essentially this function
	// is called after a 3xx response is received by the client, and returning
	// the http.ErrUseLastResponse error forces it to immediately return the
	// received response.
	tls.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return &testServer{tls}
}

// Implement a get method on our custom testServer type. This makes a GET
// request to a given url path on the test server, and returns the response
// status code, headers and body.

func (tls *testServer) get(t *testing.T, URL string) (int, http.Header, []byte) {
	rs, err := tls.Client().Get(tls.URL + URL)
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()
	return rs.StatusCode, rs.Header, body
}

// Create a postForm method for sending POST requests to the test server.
// The final parameter to this method is a url.Values object which can contain
// any data that you want to send in the request body.
func (tls *testServer) postForm(t *testing.T, URL string, form url.Values) (int, http.Header, []byte) {
	rs, err := tls.Client().PostForm(tls.URL+URL, form)
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()
	return rs.StatusCode, rs.Header, body
}

// Define a regular expression which captures the CSRF token value from the
// HTML for our user signup page.
var csrfTokenRX = regexp.MustCompile(`<input type="hidden" name="csrf_token" value='(.*?)'/>`)

func extractCSRFToken(t *testing.T, body []byte) string {
	// Use the FindSubmatch method to extract the token from the HTML body.
	// Note that this returns an array with the entire matched pattern in the
	// first position, and the values of any captured data in the subsequent
	// positions.
	matches := csrfTokenRX.FindSubmatch(body)
	if len(matches) < 2 {
		t.Fatal("no csrf token found in body")
	}
	return html.UnescapeString(string(matches[1]))
}
