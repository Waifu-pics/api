package api

import (
	"encoding/json"
	"net/http"

	"waifu.pics/util/file"
	"waifu.pics/util/web"
)

// ListFile : listing all unverified files
func (api API) ListFile(w http.ResponseWriter, r *http.Request) {
	endpoint := r.Header.Get("endpoint")

	var res struct {
		Endpoint string `json:"endpoint"`
		Verified bool   `json:"verified"`
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&res)
	if err != nil {
		web.WriteResp(w, 500, "Error")
		return
	}

	if !findInSlice(api.Config.ENDPOINTS, endpoint) {
		web.WriteResp(w, 400, "Invalid type!")
		return
	}

	files, err := api.Database.GetFilesAdmin(res.Endpoint, res.Verified)
	if err != nil {
		web.WriteResp(w, 500, "Error")
		return
	}

	var response struct {
		Verified bool     `json:"verified"`
		Files    []string `json:"files"`
	}

	response.Verified = res.Verified
	response.Files = files

	jsonres, err := json.Marshal(response)
	if err != nil {
		web.WriteResp(w, 500, "Error")
		return
	}

	web.WriteResp(w, 200, string(jsonres))

	defer r.Body.Close()
}

// VerifyFile : Verifying user uploads
func (api API) VerifyFile(w http.ResponseWriter, r *http.Request) {
	var res struct {
		Verified bool   `json:"verified"`
		File     string `json:"file"`
	}

	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		web.WriteResp(w, 400, "Invalid JSON!")
		return
	}

	defer r.Body.Close()

	if res.Verified {
		err = api.Database.VerifyFile(res.File)
		if err != nil {
			web.WriteResp(w, 400, "File could not be verified!")
			return
		}
		web.WriteResp(w, 200, "File has been verified!")
		return
	} else {
		err = api.Database.DeleteFile(res.File)
		if err != nil {
			web.WriteResp(w, 400, "File could not be deleted!")
			return
		}
		if err := file.DeleteFile(res.File, api.Config); err != nil {
			web.WriteResp(w, 400, "File could not be deleted!")
			return
		}
		web.WriteResp(w, 400, "File has been Deleted")
		return
	}
}
