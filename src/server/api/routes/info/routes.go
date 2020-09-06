package info

import (
	"fmt"

	"github.com/Riku32/waifu.pics/src/api"
	"github.com/labstack/echo"
)

// RecentFiles : recent files
func (i Route) RecentFiles(c echo.Context) error {
	data, err := i.Database.GetRecent(10)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, api.Basic{Message: api.ErrServer})
	}

	var files []FileData
	for _, v := range data {
		files = append(files, FileData{
			Uploaded: v.Uploaded,
			Name:     v.Name,
			URL:      i.Config.Web.Cdn + v.Name,
			Type:     v.Type,
			Nsfw:     v.Nsfw,
			Verified: v.Verified,
		})
	}

	return c.JSON(200, files)
}
