package apis

import (
	"alien/internals/logger"
	"bytes"
	"encoding/json"
	"errors"
	"time"

	"github.com/valyala/fasthttp"
)

func (c *Client) CapBal() (float64, error) {
	req := fasthttp.AcquireRequest()
	req.Header.SetRequestURI("https://api.capsolver.com/getBalance")
	req.Header.Set("Content-Type", "application/json")
	req.Header.SetMethod("POST")
	req.SetBodyRaw([]byte(`{"clientKey": "` + c.CapKey + `"}`))
	resp := fasthttp.AcquireResponse()
	if err := fasthttp.Do(req, resp); err != nil {
		return 0, err
	}
	if resp.StatusCode() != 200 {
		logger.ErrorLogger.Printf("\x1b[31mInvalid Response: [%d] : %s\x1b[0m\n", resp.StatusCode(), resp.Body())
		return 0, errors.New("invalid response")
	}
	var response struct {
		Balance float64 `json:"balance"`
	}
	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		return 0, err
	}
	return response.Balance, nil
}

func (c *Client) createTask() (string, error) {
	req := fasthttp.AcquireRequest()
	req.Header.SetRequestURI("https://api.capsolver.com/createTask")
	req.Header.Set("Content-Type", "application/json")
	req.Header.SetMethod("POST")
	req.SetBodyRaw([]byte(`{"clientKey":"` + c.CapKey + `","task":{"type":"ReCaptchaV3EnterpriseTaskProxyLess","websiteURL":"https://www.alienwarearena.com/account/register","websiteKey":"6LfRnbwaAAAAAPYycaGDRhoUqR-T0HyVwVkGEnmC","pageAction": "getkey"}}`))
	resp := fasthttp.AcquireResponse()
	if err := fasthttp.Do(req, resp); err != nil {
		return "", err
	}
	// logger.DebugLogger.Printf("%d -> resp: %s\n", resp.StatusCode(), resp.Body())
	if resp.StatusCode() != 200 {
		logger.ErrorLogger.Printf("%d -> resp: %s\n", resp.StatusCode(), resp.Body())
		return "", errors.New("error occurred while creating task")
	}
	var response struct {
		ErrorID          int    `json:"errorId"`
		ErrorCode        string `json:"errorCode"`
		ErrorDescription string `json:"errorDescription"`
		TaskID           string `json:"taskId"`
	}
	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		logger.ErrorLogger.Printf("%d -> resp: %s\n", resp.StatusCode(), resp.Body())
		return "", err
	}
	return response.TaskID, nil
}

func (c *Client) solve(taskID string) (string, error) {
	for {
		req := fasthttp.AcquireRequest()
		req.Header.SetRequestURI("https://api.capsolver.com/getTaskResult")
		req.Header.Set("Content-Type", "application/json")
		req.Header.SetMethod("POST")
		req.SetBodyRaw([]byte(`{"clientKey": "` + c.CapKey + `","taskId": "` + taskID + `"}`))
		resp := fasthttp.AcquireResponse()
		if err := fasthttp.Do(req, resp); err != nil {
			return "", err
		}
		logger.DebugLogger.Printf("%d -> resp: %s\n", resp.StatusCode(), resp.Body())
		if bytes.Contains(resp.Body(), []byte(`ERROR_UNKNOWN_QUESTION`)) {
			return "", errors.New("ERROR_UNKNOWN_QUESTION")
		}
		if bytes.Contains(resp.Body(), []byte(`ready`)) {
			var response struct {
				Solution struct {
					UserAgent          string `json:"userAgent"`
					ExpireTime         int64  `json:"expireTime"`
					Timestamp          int64  `json:"timestamp"`
					CaptchaKey         string `json:"captchaKey"`
					GRecaptchaResponse string `json:"gRecaptchaResponse"`
				} `json:"solution"`
			}
			if err := json.Unmarshal(resp.Body(), &response); err != nil {
				return "", err
			}
			return response.Solution.GRecaptchaResponse, nil
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func (c *Client) SolveCaptcha() (string, error) {
	taskID, err := c.createTask()
	if err != nil {
		return "", err
	}
	rq_token, err := c.solve(taskID)
	if err != nil {
		return "", err
	}
	return rq_token, nil
}
