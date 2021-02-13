package api

type UploadCodexRequest struct {
	// required
	// the codex category that the uploaded codex should belong to
	CodexCategoryId string `json:"codexCategoryId"`

	// required
	// the files that comprise the codex bundle
	Files []FileRef

	// optional
	// if specified, replace this codex instead of uploading a new codex
	CodexId string `json:"codexId,omitempty"`
}

type UploadCodexResponse struct {
	CodexId string `json:"codexId"`
}

func (c *Client) UploadCodex(r *UploadCodexRequest) (*UploadCodexResponse, error) {
	fields := make(map[string]string)
	fields["codexCategoryId"] = r.CodexCategoryId
	if r.CodexId != "" {
		fields["replaceCodexId"] = r.CodexId
	}
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
