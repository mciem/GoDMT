package discord

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strings"

	http "github.com/bogdanfinn/fhttp"
	"github.com/mciem/GoDMT/internal/utils"

	tls_client "github.com/bogdanfinn/tls-client"
)

type Discord struct {
	SessionID string
	Client    tls_client.HttpClient
	UA        string
	Token     string
	XSup      string
	Proxy     string
	Headers   map[string]string
}

func NewDiscord(token string, proxy string) (Discord, error) {
	jar := tls_client.NewCookieJar()
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(5),
		tls_client.WithClientProfile(tls_client.Chrome_112),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithCookieJar(jar),
		tls_client.WithProxyUrl(proxy),
		tls_client.WithRandomTLSExtensionOrder(),
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		return Discord{}, err
	}

	d := Discord{
		Client: client,
		UA:     "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36",
		Token:  token,
		XSup:   "eyJvcyI6IldpbmRvd3MiLCJicm93c2VyIjoiQ2hyb21lIiwiZGV2aWNlIjoiIiwic3lzdGVtX2xvY2FsZSI6InBsLVBMIiwiYnJvd3Nlcl91c2VyX2FnZW50IjoiTW96aWxsYS81LjAgKFdpbmRvd3MgTlQgMTAuMDsgV2luNjQ7IHg2NCkgQXBwbGVXZWJLaXQvNTM3LjM2IChLSFRNTCwgbGlrZSBHZWNrbykgQ2hyb21lLzExMi4wLjAuMCBTYWZhcmkvNTM3LjM2IiwiYnJvd3Nlcl92ZXJzaW9uIjoiMTEyLjAuMC4wIiwib3NfdmVyc2lvbiI6IjEwIiwicmVmZXJyZXIiOiIiLCJyZWZlcnJpbmdfZG9tYWluIjoiIiwicmVmZXJyZXJfY3VycmVudCI6IiIsInJlZmVycmluZ19kb21haW5fY3VycmVudCI6IiIsInJlbGVhc2VfY2hhbm5lbCI6InN0YWJsZSIsImNsaWVudF9idWlsZF9udW1iZXIiOjIyNDU3NSwiY2xpZW50X2V2ZW50X3NvdXJjZSI6bnVsbH0=",
		Proxy:  proxy,
		Headers: map[string]string{
			"accept":             "*/*",
			"accept-language":    "en-US;q=0.8,en;q=0.7",
			"content-type":       "application/json",
			"connection":         "keep-alive",
			"host":               "discord.com",
			"origin":             "https://discord.com",
			"sec-ch-ua":          `"Chromium";v="112", "Google Chrome";v="112", "Not;A=Brand";v="24"`,
			"sec-ch-ua-mobile":   "?0",
			"sec-ch-ua-platform": `"Windows"`,
			"sec-fetch-dest":     "empty",
			"sec-fetch-mode":     "cors",
			"sec-fetch-site":     "same-origin",
			"user-agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36",
		},
	}

	ers := d.initClient()
	if ers != nil {
		return Discord{}, ers
	}

	return d, nil
}

func (d *Discord) initClient() error {
	req, err := http.NewRequest(http.MethodGet, "https://discord.com", nil)

	if err != nil {
		return err
	}

	req.Header = http.Header{
		http.HeaderOrderKey: headerOrder,
	}

	for k, v := range d.Headers {
		req.Header.Set(k, v)
	}

	r, ers := d.Client.Do(req)
	if ers != nil {
		return ers
	}

	cookies := ""

	c := r.Header["Set-Cookie"][0]
	s := strings.Split(c, "/,/")

	for _, cookie := range s {
		ss := strings.Split(cookie, " ")

		cookies += ss[0] + " "
	}

	d.Headers["referer"] = "https://discord.com/"
	d.Headers["cookies"] = cookies

	d.Headers["x-debug-options"] = "bugReporterEnabled"
	d.Headers["x-discord-locale"] = "pl"
	d.Headers["x-discord-timezone"] = "Europe/Warsaw"
	d.Headers["x-super-properties"] = d.XSup

	if d.Token != "" {
		d.Headers["authorization"] = d.Token

		sess, err := d.OnlineToken()
		if err != nil {
			return err
		}

		d.SessionID = sess
	}

	return nil
}

func (d *Discord) CheckInvite(invite string) (InviteData, error) {
	req, err := http.NewRequest(http.MethodPost, "https://discord.com/api/v9/invites/"+invite+"?with_counts=true&with_expiration=true", bytes.NewReader([]byte("{}")))
	if err != nil {
		return InviteData{}, err
	}

	req.Header = http.Header{
		http.HeaderOrderKey: headerOrder,
	}

	for k, v := range d.Headers {
		req.Header.Set(k, v)
	}

	r, errs := d.Client.Do(req)
	if errs != nil {
		return InviteData{}, err
	}

	body, errr := ioutil.ReadAll(r.Body)
	if errr != nil {
		return InviteData{}, err
	}

	var invD InviteData
	json.Unmarshal(body, &invD)

	return invD, nil
}

func (d *Discord) JoinServer(invite string, data InviteData) (bool, string, error) {
	headers := d.Headers

	headers["referer"] = "https://discord.com/invite/" + invite
	headers["x-context-properties"] = BuildXContext(data)

	pd, er := json.Marshal(D1{
		Session_id: d.SessionID,
	})
	if er != nil {
		return false, "", er
	}

	req, err := http.NewRequest(http.MethodPost, "https://discord.com/api/v9/invites/"+invite, bytes.NewReader(pd))
	if err != nil {
		return false, "", err
	}

	req.Header = http.Header{
		http.HeaderOrderKey: headerOrder,
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	r, errs := d.Client.Do(req)
	if errs != nil {
		return false, "", err
	}

	s, x := utils.HandleStatusCode(r.StatusCode, "join")

	return s, x, nil
}

func (d *Discord) SendFriendRequest(id string, gID string, cID string) (bool, string, error) {
	headers := d.Headers

	headers["referer"] = "https://discord.com/channels/" + gID + "/" + cID
	headers["x-context-properties"] = "eyJsb2NhdGlvbiI6IlVzZXIgUHJvZmlsZSJ9"

	req, err := http.NewRequest(http.MethodPut, "https://discord.com/api/v9/users/@me/relationships/"+id, bytes.NewReader([]byte("{}")))
	if err != nil {
		return false, "", err
	}

	req.Header = http.Header{
		http.HeaderOrderKey: headerOrder,
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	r, errs := d.Client.Do(req)
	if errs != nil {
		return false, "", err
	}

	s, x := utils.HandleStatusCode(r.StatusCode, "send friend req")

	return s, x, nil

}

func (d *Discord) CreateChannel(id string, gID string, cID string) (bool, string, string, error) {
	headers := d.Headers

	headers["referer"] = "https://discord.com/channels/" + gID + "/" + cID
	headers["x-context-properties"] = "e30="

	pd, er := json.Marshal(CreateChannelPayload{
		Recipients: []string{id},
	})
	if er != nil {
		return false, "", "", er
	}

	req, err := http.NewRequest(http.MethodPost, "https://discord.com/api/v9/users/@me/channels", bytes.NewReader(pd))
	if err != nil {
		return false, "", "", err
	}

	req.Header = http.Header{
		http.HeaderOrderKey: headerOrder,
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	r, errs := d.Client.Do(req)
	if errs != nil {
		return false, "", "", errs
	}

	s, x := utils.HandleStatusCode(r.StatusCode, "send DM")

	body, er := ioutil.ReadAll(r.Body)
	if er != nil {
		return false, "", "", er
	}

	var data CreateChannelResponse
	json.Unmarshal(body, &data)

	return s, data.ID, x, nil
}

func (d *Discord) SendMessage(message string, gID string, cID string) (bool, string, error) {
	headers := d.Headers

	sx := ""
	if gID != "" {
		headers["referer"] = "https://discord.com/channels/" + gID + "/" + cID
		sx = "send message"
	} else {
		headers["referer"] = "https://discord.com/channels/" + cID
		sx = "send DM"
	}

	pd, er := json.Marshal(SendMessagePayload{
		Content: message,
		Tts:     false,
	})
	if er != nil {
		return false, "", er
	}

	req, err := http.NewRequest(http.MethodPost, "https://discord.com/api/v9/channels/"+cID+"/messages", bytes.NewReader(pd))
	if err != nil {
		return false, "", err
	}

	req.Header = http.Header{
		http.HeaderOrderKey: headerOrder,
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	r, errs := d.Client.Do(req)
	if errs != nil {
		return false, "", errs
	}

	s, x := utils.HandleStatusCode(r.StatusCode, sx)

	return s, x, nil
}

func (d *Discord) ChangeDisplayName(name string) (bool, string, error) {
	pd, er := json.Marshal(DisplayNamePayload{
		GlobalName: name,
	})
	if er != nil {
		return false, "", er
	}

	req, err := http.NewRequest(http.MethodPatch, "https://discord.com/api/v9/users/@me", bytes.NewReader(pd))
	if err != nil {
		return false, "", err
	}

	req.Header = http.Header{
		http.HeaderOrderKey: headerOrder,
	}

	for k, v := range d.Headers {
		req.Header.Set(k, v)
	}

	r, errs := d.Client.Do(req)
	if errs != nil {
		return false, "", errs
	}

	s, x := utils.HandleStatusCode(r.StatusCode, "change display name")

	return s, x, nil
}

func (d *Discord) ChangeBio(bio string) (bool, string, error) {
	pd, er := json.Marshal(BioPayload{
		Bio: bio,
	})
	if er != nil {
		return false, "", er
	}

	req, err := http.NewRequest(http.MethodPatch, "https://discord.com/api/v9/users/@me", bytes.NewReader(pd))
	if err != nil {
		return false, "", err
	}

	req.Header = http.Header{
		http.HeaderOrderKey: headerOrder,
	}

	for k, v := range d.Headers {
		req.Header.Set(k, v)
	}

	r, errs := d.Client.Do(req)
	if errs != nil {
		return false, "", errs
	}

	s, x := utils.HandleStatusCode(r.StatusCode, "change bio")

	return s, x, nil
}

func (d *Discord) CheckToken() (bool, error) {
	req, err := http.NewRequest(http.MethodPost, "https://discord.com/api/v9/users/@me/settings", bytes.NewReader([]byte{}))
	if err != nil {
		return false, err
	}

	req.Header = http.Header{
		http.HeaderOrderKey: headerOrder,
	}

	for k, v := range d.Headers {
		req.Header.Set(k, v)
	}

	r, errs := d.Client.Do(req)
	if errs != nil {
		return false, err
	}

	return r.StatusCode == 200, nil
}
