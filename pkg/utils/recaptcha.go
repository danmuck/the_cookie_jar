package utils

import (
    "encoding/json"
    "io/ioutil"
    "net/http"
    "net/url"
)

const recaptchaServerName = "https://www.google.com/recaptcha/api/siteverify"

type RecaptchaResponse struct {
    Success     bool     `json:"success"`
    Score       float64  `json:"score"`
    Action      string   `json:"action"`
    ChallengeTS string   `json:"challenge_ts"`
    Hostname    string   `json:"hostname"`
    ErrorCodes  []string `json:"error-codes"`
}

func VerifyRecaptcha(recaptchaResponse string) (bool, error) {
    secretKey := "6LfzcZUqAAAAAGna7mMSIl9ItWtiHGxkmasRe2OM"

    resp, err := http.PostForm(recaptchaServerName,
        url.Values{
            "secret":   {secretKey},
            "response": {recaptchaResponse},
        })
    if err != nil {
        return false, err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return false, err
    }

    var result RecaptchaResponse
    err = json.Unmarshal(body, &result)
    if err != nil {
        return false, err
    }

    return result.Success, nil
}
