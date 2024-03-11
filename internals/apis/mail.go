package apis

import (
	"alien/internals/logger"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"time"

	http "github.com/bogdanfinn/fhttp"
)

func (c *Client) createMail() (string, error) {
	req, err := http.NewRequest("GET", "https://api.tempmail.lol/v2/inbox/create", nil)
	if err != nil {
		logger.DebugLogger.Printf("\x1b[31mHTTP_ERR: %s\x1b[0m\n", err)
		return "", err
	}
	req.Header.Set("Host", "api.tempmail.lol")
	req.Header.Set("Sec-Ch-Ua", `"Not(A:Brand";v="24", "Chromium";v="122"`)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.6261.95 Safari/537.36")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Windows"`)
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://tempmail.lol/")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("If-None-Match", `"k3lfrcl8ff2u"`)
	req.Header.Set("Priority", "u=4, i")
	req.Header.Set("Connection", "close")
	resp, err := (*c.HTTPClient).Do(req)
	if err != nil {
		logger.DebugLogger.Printf("\x1b[31mRESP_ERR: %s\x1b[0m\n", err)
		return "", err
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.ErrorLogger.Printf("\x1b[31mREAD_ERR: %s\x1b[0m\n", err)
		return "", err
	}
	if resp.StatusCode != 201 {
		logger.ErrorLogger.Printf("\x1b[31m[%d] RESP_ERR: %s", resp.StatusCode, bodyText)
		return "", errors.New("unknown_response")
	}
	var email struct {
		Address string `json:"address"`
		Token   string `json:"token"`
	}
	if err := json.Unmarshal(bodyText, &email); err != nil {
		logger.ErrorLogger.Printf("\x1b[31mJSON_ERR: %s\x1b[0m\n", err)
		return "", err
	}
	c.Mail_Token = email.Token
	logger.DebugLogger.Println("Mail_Token: ", c.Mail_Token)
	logger.DebugLogger.Println("Mail: ", email)
	return email.Address, nil
}

func (c *Client) vLink() (string, error) {
	for i := 0; i < 5; i++ {
		req, err := http.NewRequest("GET", "https://api.tempmail.lol/v2/inbox?token="+c.Mail_Token, nil)
		if err != nil {
			logger.DebugLogger.Printf("\x1b[31mHTTP_ERR: %s\x1b[0m\n", err)
			return "", err
		}
		req.Header.Set("Host", "api.tempmail.lol")
		req.Header.Set("Sec-Ch-Ua", `"Not(A:Brand";v="24", "Chromium";v="122"`)
		req.Header.Set("Accept", "application/json, text/plain, */*")
		req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.6261.95 Safari/537.36")
		req.Header.Set("Sec-Ch-Ua-Platform", `"Windows"`)
		req.Header.Set("Sec-Fetch-Site", "same-origin")
		req.Header.Set("Sec-Fetch-Mode", "cors")
		req.Header.Set("Sec-Fetch-Dest", "empty")
		req.Header.Set("Referer", "https://tempmail.lol/")
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")
		req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
		req.Header.Set("If-None-Match", `"6ei8we0v63t"`)
		req.Header.Set("Priority", "u=4, i")
		req.Header.Set("Connection", "close")
		resp, err := (*c.HTTPClient).Do(req)
		if err != nil {
			logger.DebugLogger.Printf("\x1b[31mRESP_ERR: %s\x1b[0m\n", err)
			return "", err
		}
		defer resp.Body.Close()
		bodyText, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.ErrorLogger.Printf("\x1b[31mREAD_ERR: %s\x1b[0m\n", err)
			return "", err
		}
		var verfMail struct {
			Emails []struct {
				Body string `json:"body"`
			} `json:"emails"`
		}
		if err := json.Unmarshal(bodyText, &verfMail); err != nil {
			logger.ErrorLogger.Printf("\x1b[31mJSON_ERR: %s\x1b[0m\n", err)
			return "", err
		}
		if len(verfMail.Emails) == 0 {
			time.Sleep(10 * time.Second)
			continue
		}
		verification_link := strings.Split(strings.Split(verfMail.Emails[0].Body, "Activate Account\n[")[1], "]")[0]
		logger.DebugLogger.Println("Verification_link: ", verification_link)
		return verification_link, nil
	}
	logger.ErrorLogger.Println("\x1b[31mtimedout fetchin' mail\x1b[0m")
	return "", errors.New("timed_out_fetching_mail")
}
