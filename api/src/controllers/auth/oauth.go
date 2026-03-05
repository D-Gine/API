/*
** D&GINE Project, 2026
** Backend
** File description:
** auth/oauth.go
 */

package auth

// type InternalData struct {
// 	Name  string `json:"name"`
// 	Email string `json:"email"`
// }

// type OauthArgs struct {
// 	Code string `json:"code"`
// 	Name string `json:"service_name"`
// }

// func httpRequest(method string, uri string, headers map[string]string, body map[string]string, debug bool) (map[string]any, error) {

// 	if debug {
// 		fmt.Printf(`
// 			=========== HTTP Request ==========
// 			Method : %s
// 			Uri : %s
// 			Header : %v
// 			Body : %v
// 			=========== HTTP Request ==========
// 		`, method, uri, headers, body)
// 	}
// 	data := url.Values{}
// 	for key, value := range body {
// 		data.Set(key, value)
// 	}

// 	req, err := http.NewRequest(method, uri, strings.NewReader(data.Encode()))
// 	if err != nil {
// 		return nil, err
// 	}
// 	for key, value := range headers {
// 		req.Header.Add(key, value)
// 	}

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	respBytes, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var result map[string]any
// 	if err := json.Unmarshal(respBytes, &result); err != nil {
// 		return nil, err
// 	}

// 	if debug {
// 		fmt.Printf(`
// 			=========== HTTP Response ==========
// 			Code : %d
// 			Body : %v
// 			=========== HTTP Response ==========
// 		`, resp.StatusCode, result)
// 	}

// 	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
// 		return nil, err
// 	}

// 	return result, nil
// }

// func linkExternalUser(c *gin.Context, creds InternalData) (UserData, error) {
// 	_, err := database.Db.Exec("INSERT INTO accounts.users (email, name) VALUES ($1, $2) ON CONFLICT (email) DO NOTHING", creds.Email, creds.Name)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database insertion error"})
// 		return UserData{}, err
// 	}

// 	var id, username, email, role sql.NullString
// 	err = database.Db.QueryRow("SELECT id, name, email, role FROM accounts.users WHERE email=$1", creds.Email).Scan(&id, &username, &email, &role)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database selection error"})
// 		return UserData{}, err
// 	}

// 	return UserData{id.String, email.String, username.String, role.String}, nil
// }

// func oauthGoogle(c *gin.Context, args OauthArgs) {
// 	// TOKEN
// 	tokenResponse, err := httpRequest("POST", "https://oauth2.googleapis.com/token",
// 		map[string]string{
// 			"Content-Type": "application/x-www-form-urlencoded",
// 		},
// 		map[string]string{
// 			"grant_type":    "authorization_code",
// 			"code":          args.Code,
// 			"redirect_uri":  "http://localhost:8081/google-callback",
// 			"client_id":     os.Getenv("GOOGLE_CLIENT_ID"),
// 			"client_secret": os.Getenv("GOOGLE_CLIENT_SECRET"),
// 		},
// 		false,
// 	)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, err)
// 		return
// 	}
// 	token, ok := tokenResponse["access_token"].(string)
// 	if !ok {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "no access token in response"})
// 		return
// 	}

// 	// USER
// 	userResponse, err := httpRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo",
// 		map[string]string{
// 			"Authorization": "Bearer " + token,
// 			"Accept":        "application/json",
// 		},
// 		map[string]string{},
// 		false,
// 	)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, err)
// 		return
// 	}
// 	username, ok := userResponse["name"].(string)
// 	if !ok {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "no username in response"})
// 		return
// 	}
// 	email, ok := userResponse["email"].(string)
// 	if !ok {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "no username in response"})
// 		return
// 	}

// 	UserData, err := linkExternalUser(c, InternalData{username, email})
// 	if err == nil {
// 		var result PostLoginResponse
// 		result.Id = UserData.Id
// 		result.Token, err = BuildToken(c, UserData)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, err)
// 		} else {
// 			c.JSON(http.StatusOK, result)
// 		}
// 	} else {
// 		c.JSON(http.StatusInternalServerError, err)
// 	}
// }

// func oauthDiscord(c *gin.Context, args OauthArgs) {
// 	// TOKEN
// 	tokenResponse, err := httpRequest("POST", "https://discord.com/api/oauth2/token",
// 		map[string]string{
// 			"Content-Type": "application/x-www-form-urlencoded",
// 		},
// 		map[string]string{
// 			"grant_type":    "authorization_code",
// 			"code":          args.Code,
// 			"redirect_uri":  "http://localhost:8081/discord-callback",
// 			"client_id":     os.Getenv("DISCORD_CLIENT_ID"),
// 			"client_secret": os.Getenv("DISCORD_CLIENT_SECRET"),
// 		},
// 		false,
// 	)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, err)
// 		return
// 	}
// 	token, ok := tokenResponse["access_token"].(string)
// 	if !ok {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "no access token in response"})
// 		return
// 	}

// 	// USER
// 	userResponse, err := httpRequest("GET", "https://discord.com/api/users/@me",
// 		map[string]string{
// 			"Authorization": "Bearer " + token,
// 		},
// 		map[string]string{},
// 		false,
// 	)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, err)
// 		return
// 	}
// 	username, ok := userResponse["username"].(string)
// 	if !ok {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "no username in response"})
// 		return
// 	}
// 	email, ok := userResponse["email"].(string)
// 	if !ok {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "no email in response"})
// 		return
// 	}

// 	UserData, err := linkExternalUser(c, InternalData{username, email})
// 	if err == nil {
// 		var result PostLoginResponse
// 		result.Id = UserData.Id
// 		result.Token, err = BuildToken(c, UserData)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, err)
// 		} else {
// 			c.JSON(http.StatusOK, result)
// 		}
// 	} else {
// 		c.JSON(http.StatusInternalServerError, err)
// 	}
// }

// func oauthGithub(c *gin.Context, args OauthArgs) {
// 	// TOKEN
// 	tokenApiResponse, err := httpRequest("POST", "https://github.com/login/oauth/access_token",
// 		map[string]string{
// 			"Content-type": "application/x-www-form-urlencoded",
// 			"Accept":       "application/json",
// 		},
// 		map[string]string{
// 			"client_id":     os.Getenv("GITHUB_CLIENT_ID"),
// 			"client_secret": os.Getenv("GITHUB_CLIENT_SECRET"),
// 			"code":          args.Code,
// 			"redirect_uri":  "http://localhost:8081/github-callback",
// 		},
// 		false,
// 	)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, err)
// 		return
// 	}
// 	token, ok := tokenApiResponse["access_token"].(string)
// 	if !ok {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "no access token in response"})
// 		return
// 	}

// 	// USER
// 	userApiResponse, err := httpRequest("GET", "https://api.github.com/user",
// 		map[string]string{
// 			"Authorization": "Bearer " + token,
// 			"Accept":        "application/json",
// 		},
// 		map[string]string{},
// 		false,
// 	)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, err)
// 		return
// 	}
// 	username, ok2 := userApiResponse["login"].(string)
// 	if !ok2 {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "no user name in response"})
// 		return
// 	}

// 	// EMAIL
// 	req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
// 		return
// 	}
// 	req.Header.Add("Authorization", "Bearer "+token)
// 	req.Header.Add("Accept", "application/json")

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call API"})
// 		return
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "GitHub emails API returned error"})
// 		return
// 	}

// 	// api qui renvoie un tableau d'emails il faut trouver le bon
// 	var emails []map[string]any
// 	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse emails response"})
// 		return
// 	}

// 	var email string
// 	for _, e := range emails {
// 		if primary, _ := e["primary"].(bool); primary {
// 			if verified, _ := e["verified"].(bool); verified {
// 				email, _ = e["email"].(string)
// 				break
// 			}
// 		}
// 	}

// 	if email == "" {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "no verified email found"})
// 		return
// 	}

// 	UserData, err := linkExternalUser(c, InternalData{username, email})
// 	if err == nil {
// 		var result PostLoginResponse
// 		result.Id = UserData.Id
// 		result.Token, err = BuildToken(c, UserData)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, err)
// 		} else {
// 			c.JSON(http.StatusOK, result)
// 		}
// 	} else {
// 		c.JSON(http.StatusInternalServerError, err)
// 	}
// }

// // @BasePath /api/auth/oauth
// // Auth godoc
// // @Summary Connection of a user
// // @Schemes
// // @Description Connection of a user using an OAuth <br>Token will be set in cookies if the arguments combination is valid
// // @Tags auth
// // @Accept json
// // @Produce json
// // @Param creds body OauthArgs true "Code + Service combination of the corresponding OAuth"
// // @Success 200 {object} structs.PostResponse
// // @Router /api/auth/oauth [post]
// func OAuthLogin(c *gin.Context) {
// 	var args OauthArgs

// 	// Body Json parsing
// 	err := c.ShouldBindJSON(&args)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
// 		return
// 	}

// 	switch args.Name {
// 	case "google":
// 		oauthGoogle(c, args)
// 		return
// 	case "discord":
// 		oauthDiscord(c, args)
// 		return
// 	case "github":
// 		oauthGithub(c, args)
// 		return
// 	}
// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
// }
