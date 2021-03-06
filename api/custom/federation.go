package custom

import (
	"encoding/json"
	"github.com/getsentry/sentry-go"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/turt2live/matrix-media-repo/api"
	"github.com/turt2live/matrix-media-repo/common/rcontext"
	"github.com/turt2live/matrix-media-repo/matrix"
)

func GetFederationInfo(r *http.Request, rctx rcontext.RequestContext, user api.UserInfo) interface{} {
	params := mux.Vars(r)

	serverName := params["serverName"]

	rctx = rctx.LogWithFields(logrus.Fields{
		"serverName": serverName,
	})

	url, hostname, err := matrix.GetServerApiUrl(serverName)
	if err != nil {
		rctx.Log.Error(err)
		sentry.CaptureException(err)
		return api.InternalServerError(err.Error())
	}

	versionUrl := url + "/_matrix/federation/v1/version"
	versionResponse, err := matrix.FederatedGet(versionUrl, hostname, rctx)
	if err != nil {
		rctx.Log.Error(err)
		sentry.CaptureException(err)
		return api.InternalServerError(err.Error())
	}

	c, err := ioutil.ReadAll(versionResponse.Body)
	if err != nil {
		rctx.Log.Error(err)
		sentry.CaptureException(err)
		return api.InternalServerError(err.Error())
	}

	out := make(map[string]interface{})
	err = json.Unmarshal(c, &out)
	if err != nil {
		rctx.Log.Error(err)
		sentry.CaptureException(err)
		return api.InternalServerError(err.Error())
	}

	resp := make(map[string]interface{})
	resp["base_url"] = url
	resp["hostname"] = hostname
	resp["versions_response"] = out
	return &api.DoNotCacheResponse{Payload: resp}
}
