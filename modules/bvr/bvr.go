package bvr

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/angelfluffyookami/247BVR/modules/common/global"
	"github.com/angelfluffyookami/247BVR/modules/common/utils/logger"
)

var logging = logger.Log{}

func InitCache() {

}

func syncCache() error {

	return nil
}

func syncKills() error {

	return nil
}

func syncUsers() error {

	return nil
}

func syncDeaths() error {

	req, err := http.Get(global.Config.APIEndpoint + "/deaths")

	if err != nil {
		logging.Err().Panic().Message("Sanity check: " + err.Error()).Add()
	}

	var Kills []KillStruct
	body, err := io.ReadAll(req.Body)
	if err != nil {
		logging.Err().Panic().Message("Sanity check: " + err.Error()).Add()
	}

	err = json.Unmarshal(body, &Kills)
	if err != nil {
		logging.Err().Panic().Message("Sanity check: " + err.Error()).Add()
	}
	return nil
}

func syncOnline() error {

	return nil
}
