package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var testUser = UserData{
	Id:    "42",
	Email: "john@example.com",
	Name:  "john",
	Image: "https://example.com/john.png",
	Role:  "admin",
}

// jwtKey fixe pour les tests (a discuter)
func setTestKey(t *testing.T) {
	t.Helper()
	old := jwtKey
	jwtKey = []byte("test-secret-key")
	t.Cleanup(func() { jwtKey = old })
}

func validClaims() jwt.MapClaims {
	return jwt.MapClaims{
		"id":       testUser.Id,
		"email":    testUser.Email,
		"username": testUser.Name,
		"image":    testUser.Image,
		"role":     testUser.Role,
		"exp":      time.Now().Add(time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}
}

func signClaims(t *testing.T, claims jwt.MapClaims, key []byte) string {
	t.Helper()
	s, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(key)
	if err != nil {
		t.Fatalf("failed to sign test token: %v", err)
	}
	return s
}

func newTestContext(req *http.Request) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	return c, w
}

// CREATETOKEN

func TestCreateToken(t *testing.T) {
	setTestKey(t)

	tokenString, err := CreateToken(testUser)
	if err != nil {
		t.Fatalf("CreateToken returned error: %v", err)
	}
	if tokenString == "" {
		t.Fatal("CreateToken returned an empty token")
	}

	parsed, err := jwt.Parse(tokenString, func(*jwt.Token) (any, error) { return jwtKey, nil })
	if err != nil {
		t.Fatalf("token cannot be parsed: %v", err)
	}
	claims := parsed.Claims.(jwt.MapClaims)

	want := map[string]string{
		"id":       testUser.Id,
		"email":    testUser.Email,
		"username": testUser.Name,
		"image":    testUser.Image,
		"role":     testUser.Role,
	}
	for key, value := range want {
		if claims[key] != value {
			t.Errorf("claim %q = %v, want %q", key, claims[key], value)
		}
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		t.Fatal("exp claim missing or not a number")
	}
	remaining := time.Until(time.Unix(int64(exp), 0))
	if remaining < 59*time.Minute || remaining > 61*time.Minute {
		t.Errorf("token should expire in ~1h, expires in %v", remaining)
	}
	if _, ok := claims["iat"].(float64); !ok {
		t.Error("iat claim missing or not a number")
	}
}

// --- ExtractFromToken ------------------------------------------------------

func Sub_ExtractFromToken(t *testing.T) {
	setTestKey(t)

	tokenString, err := CreateToken(testUser)
	if err != nil {
		t.Fatalf("CreateToken returned error: %v", err)
	}

	got, err := ExtractFromToken(tokenString)
	if err != nil {
		t.Fatalf("ExtractFromToken returned error: %v", err)
	}
	if got != testUser {
		t.Errorf("ExtractFromToken = %+v, want %+v", got, testUser)
	}
}

func Sub_ExtractFromToken_WithBearerPrefix(t *testing.T) {
	setTestKey(t)

	tokenString, _ := CreateToken(testUser)
	got, err := ExtractFromToken("Bearer " + tokenString)
	if err != nil {
		t.Fatalf("ExtractFromToken returned error: %v", err)
	}
	if got != testUser {
		t.Errorf("ExtractFromToken = %+v, want %+v", got, testUser)
	}
}

func Sub_ExtractFromToken_Invalid(t *testing.T) {
	setTestKey(t)

	expired := validClaims()
	expired["exp"] = time.Now().Add(-time.Hour).Unix()

	tests := []struct {
		name  string
		token string
	}{
		{"empty", ""},
		{"garbage", "azertyuiop"},
		{"wrong signature", signClaims(t, validClaims(), []byte("poiuytreza"))},
		{"expired", signClaims(t, expired, jwtKey)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := ExtractFromToken(tt.token); err == nil {
				t.Error("expected an error, got nil")
			}
		})
	}
}

func Sub_ExtractFromToken_MissingOrInvalidClaim(t *testing.T) {
	setTestKey(t)

	for _, claim := range []string{"id", "email", "username", "role", "image"} {
		t.Run("missing "+claim, func(t *testing.T) {
			claims := validClaims()
			delete(claims, claim)

			if _, err := ExtractFromToken(signClaims(t, claims, jwtKey)); err == nil {
				t.Errorf("expected an error when %q is missing", claim)
			}
		})
	}

	t.Run("id with wrong type", func(t *testing.T) {
		claims := validClaims()
		claims["id"] = 42 // nombre au lieu d'une string

		if _, err := ExtractFromToken(signClaims(t, claims, jwtKey)); err == nil {
			t.Error("expected an error when id is not a string")
		}
	})
}

func Test_ExtractFromToken(t *testing.T) {
	t.Run("extract from token", Sub_ExtractFromToken)
	t.Run("with bearer prefix", Sub_ExtractFromToken_WithBearerPrefix)
	t.Run("invalid token", Sub_ExtractFromToken_Invalid)
	t.Run("missing or invalid claim", Sub_ExtractFromToken_MissingOrInvalidClaim)
}

// VERIFYTOKEN

func TestVerifyToken(t *testing.T) {
	setTestKey(t)

	valid := signClaims(t, validClaims(), jwtKey)

	t.Run("valid token", func(t *testing.T) {
		tok, err := verifyToken(valid)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !tok.Valid {
			t.Error("token should be valid")
		}
	})

	t.Run("valid token with Bearer prefix", func(t *testing.T) {
		if _, err := verifyToken("Bearer " + valid); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("invalid tokens",
		func(t *testing.T) {
			for name, tokenString := range map[string]string{
				"empty":           "",
				"garbage":         "kyurem",
				"wrong signature": signClaims(t, validClaims(), []byte("poiuytreza")),
			} {
				if _, err := verifyToken(tokenString); err == nil {
					t.Errorf("%s: expected an error, got nil", name)
				}
			}
		})
}

// BUILDTOKEN

func TestBuildToken(t *testing.T) {
	setTestKey(t)

	c, w := newTestContext(httptest.NewRequest(http.MethodPost, "/login", nil))

	tokenString, err := BuildToken(c, testUser)
	if err != nil {
		t.Fatalf("BuildToken returned error: %v", err)
	}

	got, err := ExtractFromToken(tokenString)
	if err != nil || got != testUser {
		t.Fatalf("returned token is not valid for testUser (user=%+v, err=%v)", got, err)
	}

	var cookie *http.Cookie
	for _, ck := range w.Result().Cookies() {
		if ck.Name == "token" {
			cookie = ck
		}
	}
	if cookie == nil {
		t.Fatal("cookie 'token' was not set")
	}
	if cookie.Value != tokenString {
		t.Error("cookie value differs from the returned token")
	}
	if !cookie.HttpOnly {
		t.Error("cookie should be HttpOnly")
	}
	if cookie.MaxAge != 3600 {
		t.Errorf("cookie MaxAge = %d, want 3600", cookie.MaxAge)
	}
	if cookie.Path != "/" {
		t.Errorf("cookie Path = %q, want %q", cookie.Path, "/")
	}
}

// GETUSERFROMTOKEN

func TestGetUserFromToken(t *testing.T) {
	setTestKey(t)

	valid, _ := CreateToken(testUser)

	otherUser := UserData{Id: "643", Email: "thresh@iram.com", Name: "resh", Image: "", Role: "user"}
	otherToken, _ := CreateToken(otherUser)

	tests := []struct {
		name     string
		cookie   string
		header   string
		wantUser UserData
		wantErr  bool
	}{
		{name: "token in cookie", cookie: valid, wantUser: testUser},
		{name: "Bearer header", header: "Bearer " + valid, wantUser: testUser},
		{name: "raw header without Bearer", header: valid, wantUser: testUser},
		{name: "cookie has priority over header", cookie: valid, header: "Bearer " + otherToken, wantUser: testUser},
		{name: "no token", wantErr: true},
		{name: "invalid cookie", cookie: "garbage", wantErr: true},
		{name: "invalid header", header: "Bearer garbage", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/me", nil)
			if tt.cookie != "" {
				req.AddCookie(&http.Cookie{Name: "token", Value: tt.cookie})
			}
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}
			c, _ := newTestContext(req)

			got, err := GetUserFromToken(c)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected an error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.wantUser {
				t.Errorf("got %+v, want %+v", got, tt.wantUser)
			}
		})
	}
}
