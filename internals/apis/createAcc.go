package apis

import (
	"alien/internals/logger"
	"bytes"
	"errors"
	"io"
	"strings"

	http "github.com/bogdanfinn/fhttp"
)

func (c *Client) getAccount() (string, error) {
	logger.DebugLogger.Println("Fetching Cookies & Ingredients!")
	req, err := http.NewRequest("GET", "https://uk.alienwarearena.com/account/register", nil)
	if err != nil {
		logger.DebugLogger.Printf("\x1b[31mHTTP_ERR: %s\x1b[0m\n", err)
		return "", err
	}
	req.Header.Set("Host", "www.alienwarearena.com")
	req.Header.Set("Sec-Ch-Ua", `"Not(A:Brand";v="24", "Chromium";v="122"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Windows"`)
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.6261.95 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Referer", "https://www.alienwarearena.com/login?return=%2Fucf%2Fshow%2F2170237%2Fboards%2Fcontest-and-giveaways-global%2Fone-month-of-discord-nitro-exclusive-key-giveaway")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Priority", "u=0, i")
	(*c.HTTPClient).SetFollowRedirect(true)
	resp, err := (*c.HTTPClient).Do(req)
	if err != nil {
		logger.DebugLogger.Printf("\x1b[31mRESP_ERR: %s\x1b[0m\n", err)
		return "", err
	}
	defer func() {
		resp.Body.Close()
		(*c.HTTPClient).SetFollowRedirect(false)
	}()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.ErrorLogger.Printf("\x1b[31mREAD_ERR: %s\x1b[0m\n", err)
		return "", err
	}
	if resp.StatusCode != 200 {
		logger.ErrorLogger.Printf("\x1b[31m[%d] RESP_ERR: %s", resp.StatusCode, bodyText)
		return "", errors.New("unknown_response")
	}
	if bytes.Contains(bodyText, []byte(`We have detected`)) {
		logger.ErrorLogger.Printf("\x1b[31mIP BlackListed!\x1b[0m\n")
		return "", errors.New("ip_blacklist")
	}
	registration_token := string(bytes.Split(bytes.Split(bodyText, []byte(`<input type="hidden" id="user_registration__token" name="user_registration[_token]" value="`))[1], []byte(`"`))[0])
	logger.InfoLogger.Println("Registration_Token: ", registration_token)
	return registration_token, nil
}

func (c *Client) getVerification(email, registration_token string) error {
	logger.DebugLogger.Println("Creating Account Using: ", email)
	var data = strings.NewReader(`user_registration%5Bemail%5D%5Bfirst%5D=` + email + `&user_registration%5Bemail%5D%5Bsecond%5D=` + email + `&user_registration%5Bbirthdate%5D%5Bmonth%5D=5&user_registration%5Bbirthdate%5D%5Bday%5D=15&user_registration%5Bbirthdate%5D%5Byear%5D=1985&user_registration%5BtermsAccepted%5D=1&user_registration%5B_token%5D=` + registration_token + `&user_registration%5BsteamId%5D=&user_registration%5BbattlenetOauthProfileId%5D=&user_registration%5Btimezone%5D=America%2FDenver&user_registration%5BsourceInfo%5D=null&user_registration%5BreferralCode%5D=&user_registration%5Brecaptcha3%5D=`)
	req, err := http.NewRequest("POST", "https://www.alienwarearena.com/account/register", data)
	if err != nil {
		logger.DebugLogger.Printf("\x1b[31mHTTP_ERR: %s\x1b[0m\n", err)
		return err
	}
	req.Header.Set("Host", "www.alienwarearena.com")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Sec-Ch-Ua", `"Not(A:Brand";v="24", "Chromium";v="122"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Windows"`)
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Origin", "https://www.alienwarearena.com")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.6261.95 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Referer", "https://www.alienwarearena.com/account/register")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Priority", "u=0, i")
	resp, err := (*c.HTTPClient).Do(req)
	if err != nil {
		logger.DebugLogger.Printf("\x1b[31mRESP_ERR: %s\x1b[0m\n", err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 302 {
		return nil
	} else {
		bodyText, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.ErrorLogger.Printf("\x1b[31mREAD_ERR: %s\x1b[0m\n", err)
			return err
		}
		logger.ErrorLogger.Printf("\x1b[31m[%d] RESP_ERR: %s", resp.StatusCode, bodyText)
		return errors.New("unknown_response")
	}
}

func (c *Client) verifyRedirection(vlink string) (string, string, error) {
	logger.DebugLogger.Println("Redirecting Verification Link To Get VerifyToken")
	req, err := http.NewRequest("GET", vlink, nil)
	if err != nil {
		logger.DebugLogger.Printf("\x1b[31mHTTP_ERR: %s\x1b[0m\n", err)
		return "", "", err
	}
	req.Header.Set("Host", "mandrillapp.com")
	req.Header.Set("Sec-Ch-Ua", `"Not(A:Brand";v="24", "Chromium";v="122"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Windows"`)
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.6261.95 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Priority", "u=0, i")
	req.Header.Set("Connection", "close")
	(*c.HTTPClient).SetFollowRedirect(true)
	resp, err := (*c.HTTPClient).Do(req)
	if err != nil {
		logger.DebugLogger.Printf("\x1b[31mRESP_ERR: %s\x1b[0m\n", err)
		return "", "", err
	}
	defer func() {
		resp.Body.Close()
		(*c.HTTPClient).SetFollowRedirect(false)
	}()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.ErrorLogger.Printf("\x1b[31mREAD_ERR: %s\x1b[0m\n", err)
		return "", "", err
	}
	if resp.StatusCode != 200 {
		logger.ErrorLogger.Printf("\x1b[31m[%d] RESP_ERR: %s", resp.StatusCode, bodyText)
		return "", "", errors.New("unknown_response")
	}
	if bytes.Contains(bodyText, []byte(`We have detected`)) {
		logger.ErrorLogger.Printf("\x1b[31mIP BlackListed!\x1b[0m\n")
		return "", "", errors.New("ip_blacklist")
	}
	confirm_url := string(bytes.Split(bytes.Split(bodyText, []byte(`<form name="platformd_user_confirm_registration" method="post" action="/register/confirm/`))[1], []byte(`">`))[0])
	verification_token := string(bytes.Split(bytes.Split(bodyText, []byte(`<input type="hidden" id="platformd_user_confirm_registration__token" name="platformd_user_confirm_registration[_token]" value="`))[1], []byte(`"`))[0])
	logger.InfoLogger.Println("confirm_url: ", confirm_url)
	logger.InfoLogger.Println("verification_token: ", verification_token)
	return confirm_url, verification_token, nil
}

func (c *Client) verifyConfirm(confirm_url, verification_token string) error {
	logger.DebugLogger.Println("Verifying Account")
	var data = strings.NewReader(`platformd_user_confirm_registration%5Bconfirm%5D=&platformd_user_confirm_registration%5B_token%5D=` + verification_token)
	req, err := http.NewRequest("POST", "https://www.alienwarearena.com/register/confirm/"+confirm_url, data)
	if err != nil {
		logger.DebugLogger.Printf("\x1b[31mHTTP_ERR: %s\x1b[0m\n", err)
		return err
	}
	req.Header.Set("Host", "www.alienwarearena.com")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Sec-Ch-Ua", `"Not(A:Brand";v="24", "Chromium";v="122"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Windows"`)
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Origin", "https://www.alienwarearena.com")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.6261.95 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Referer", "https://www.alienwarearena.com/register/confirm/"+confirm_url)
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Priority", "u=0, i")
	resp, err := (*c.HTTPClient).Do(req)
	if err != nil {
		logger.DebugLogger.Printf("\x1b[31mRESP_ERR: %s\x1b[0m\n", err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 302 {
		return nil
	} else {
		bodyText, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.ErrorLogger.Printf("\x1b[31mREAD_ERR: %s\x1b[0m\n", err)
			return err
		}
		logger.ErrorLogger.Printf("\x1b[31m[%d] RESP_ERR: %s", resp.StatusCode, bodyText)
		return errors.New("unknown_response")
	}
}

func (c *Client) CreateAccount() error {
	reg_token, err := c.getAccount()
	if err != nil {
		return err
	}
	// rnd, err := utils.RandomString(16)
	// if err != nil {
	// 	logger.ErrorLogger.Printf("\x1b[31mRND_ERR: %s\x1b[0m\n", err)
	// 	return err
	// }

	rnd, err := c.createMail()
	if err != nil {
		return err
	}
	if err := c.getVerification(rnd, reg_token); err != nil {
		return err
	}
	vlink, err := c.vLink()
	if err != nil {
		return err
	}
	confirm_url, verification_url, err := c.verifyRedirection(vlink)
	if err != nil {
		return err
	}
	if err := c.verifyConfirm(confirm_url, verification_url); err != nil {
		return err
	}
	logger.DebugLogger.Println("\x1b[38;5;120mVerified Account: " + rnd + "@zadowchi.live" + "\x1b[0m")
	return nil
}

//[{65ead0ed951b7051e8eb7db0 /accounts/65ead0b86eb6feb8e50bf534 <30904262.20240308084844.65ead0ec9239d8.62910718@mail128-14.atl41.mandrillapp.com> {noreply@alienwarearena.com Alienware Arena} [{epupqqan16zpmtbv7a80@awgarstone.com }] Activate Your Alienware Arena Account WELCOME! Thank you for registering for an Alienware Arena account. Your account will allow you access to great content across… false false false 38414 /messages/65ead0ed951b7051e8eb7db0/download 2024-03-08 08:48:44 +0000 +0000 2024-03-08 08:48:45 +0000 +0000}]
