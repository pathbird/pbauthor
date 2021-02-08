package api

type UploadCodexRequest struct {
	CodexCategoryId string `json:"codexCategoryId"`
	Files           []FileRef
}

type UploadCodexResponse struct {
	CodexId string `json:"codexId"`
}

func (c *Client) UploadCodex(r *UploadCodexRequest) (*UploadCodexResponse, error) {
	fields := make(map[string]string)
	fields["codexCategoryId"] = r.CodexCategoryId
	res, err := c.postMultipart(&multipartRequest{
		route:  "author/upload-codex",
		fields: fields,
		files:  r.Files,
	})
	if err != nil {
		return nil, err
	}
	if err := res.StatusError(); err != nil {
		return nil, err
	}
	resp := &UploadCodexResponse{}
	err = res.UnmarshalJson(resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
