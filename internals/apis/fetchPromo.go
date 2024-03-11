package apis

import (
	"alien/internals/logger"
	"bytes"
	"errors"
	"io"
	"log"

	http "github.com/bogdanfinn/fhttp"
)

func (c *Client) getAccountDetails() error {
	req, err := http.NewRequest("GET", "https://www.alienwarearena.com/incomplete", nil)
	if err != nil {
		logger.DebugLogger.Printf("\x1b[31mHTTP_ERR: %s\x1b[0m\n", err)
		return err
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
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Referer", "https://www.alienwarearena.com/ucf/show/2170237/boards/contest-and-giveaways-global/one-month-of-discord-nitro-exclusive-key-giveaway")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Priority", "u=0, i")
	resp, err := (*c.HTTPClient).Do(req)
	if err != nil {
		logger.DebugLogger.Printf("\x1b[31mRESP_ERR: %s\x1b[0m\n", err)
		return err
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.ErrorLogger.Printf("\x1b[31mREAD_ERR: %s\x1b[0m\n", err)
		return err
	}
	if resp.StatusCode != 200 {
		logger.ErrorLogger.Printf("\x1b[31m[%d] RESP_ERR: %s", resp.StatusCode, bodyText)
		return errors.New("unknown_response")
	}
	if bytes.Contains(bodyText, []byte(`We have detected`)) {
		logger.ErrorLogger.Printf("\x1b[31mIP BlackListed!\x1b[0m\n")
		return errors.New("ip_blacklist")
	}
	c.User_UUID = string(bytes.Split(bytes.Split(bodyText, []byte(`var user_uuid         = "`))[1], []byte(`";`))[0])
	c.User_Id = string(bytes.Split(bytes.Split(bodyText, []byte(`var user_id           = `))[1], []byte(`;`))[0])
	c.Login_Id = string(bytes.Split(bytes.Split(bodyText, []byte(`var login_id          = `))[1], []byte(`;`))[0])
	c.User_Country = string(bytes.Split(bytes.Split(bodyText, []byte(`var user_country      = "`))[1], []byte(`";`))[0])
	logger.DebugLogger.Println("User_UUID: ", c.User_UUID)
	logger.DebugLogger.Println("User_ID: ", c.User_Id)
	logger.DebugLogger.Println("Login_ID: ", c.Login_Id)
	logger.DebugLogger.Println("User_Country: ", c.User_Country)
	return nil
}

func (c *Client) generatePromo(rp_token string) error {
	req, err := http.NewRequest("GET", "https://giveawayapi.alienwarearena.com/production/key/get?giveaway_uuid=df863897-304c-4985-830c-56414830ade7&api_key=a75eb2f0-3f7a-4742-96c7-202977acb4cf&user_uuid="+c.User_UUID+"&extra_info=%7B%22siteId%22%3A1%2C%22siteGroupId%22%3A1%2C%22loginId%22%3A"+c.Login_Id+"%2C%22countryCode%22%3A%22"+c.User_Country+"%22%2C%22userId%22%3A"+c.User_Id+"%7D&recaptcha_token="+rp_token, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Host", "giveawayapi.alienwarearena.com")
	req.Header.Set("Sec-Ch-Ua", `"Not(A:Brand";v="24", "Chromium";v="122"`)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.6261.95 Safari/537.36")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Windows"`)
	req.Header.Set("Origin", "https://www.alienwarearena.com")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://www.alienwarearena.com/")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Priority", "u=1, i")
	resp, err := (*c.HTTPClient).Do(req)
	if err != nil {
		logger.DebugLogger.Printf("\x1b[31mRESP_ERR: %s\x1b[0m\n", err)
		return err
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.ErrorLogger.Printf("\x1b[31mREAD_ERR: %s\x1b[0m\n", err)
		return err
	}
	if resp.StatusCode != 200 || !bytes.Contains(bodyText, []byte("successMessage")) {
		logger.ErrorLogger.Printf("\x1b[31m[%d] RESP_ERR: %s", resp.StatusCode, bodyText)
		return errors.New("unknown_response")
	}
	return nil
}

func (c *Client) fetchPromo() error {
	req, err := http.NewRequest("GET", "https://www.alienwarearena.com/giveaways/keys", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Host", "www.alienwarearena.com")
	req.Header.Set("Sec-Ch-Ua", `"Not(A:Brand";v="24", "Chromium";v="122"`)
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.6261.95 Safari/537.36")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Windows"`)
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://www.alienwarearena.com/ucf/show/2170237/boards/contest-and-giveaways-global/one-month-of-discord-nitro-exclusive-key-giveaway")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Priority", "u=1, i")
	resp, err := (*c.HTTPClient).Do(req)
	if err != nil {
		logger.DebugLogger.Printf("\x1b[31mRESP_ERR: %s\x1b[0m\n", err)
		return err
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.ErrorLogger.Printf("\x1b[31mREAD_ERR: %s\x1b[0m\n", err)
		return err
	}
	if resp.StatusCode != 200 {
		logger.ErrorLogger.Printf("\x1b[31m[%d] RESP_ERR: %s", resp.StatusCode, bodyText)
		return errors.New("unknown_response")
	}
	if bytes.Contains(bodyText, []byte(`We have detected`)) {
		logger.ErrorLogger.Printf("\x1b[31mIP BlackListed!\x1b[0m\n")
		return errors.New("ip_blacklist")
	}
	c.Promo = "https://discord.com/billing/promotions/" + string(bytes.Split(bytes.Split(bodyText, []byte(`[{"value":"`))[1], []byte(`",`))[0])
	logger.InfoLogger.Println("\x1b[38;5;120mRetrieved Promo: " + c.Promo[:46] + "**-*****-*****-*****\x1b[0m")
	return nil
}

func (c *Client) GetPromo() error {
	if err := c.getAccountDetails(); err != nil {
		return err
	}
	rp_token, err := c.SolveCaptcha()
	if err != nil {
		return err
	}
	if err := c.generatePromo(rp_token); err != nil {
		return err
	}
	if err := c.fetchPromo(); err != nil {
		return err
	}
	return nil
}
