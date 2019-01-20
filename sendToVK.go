package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

// https://oauth.vk.com/authorize?client_id=APP_ID&redirect_uri=https://oauth.vk.com/blank.html&display=page&scope=photos,stories,wall,offline&v=5.92&revoke=1&response_type=token

func vkGetWallUploadServer(groupID int, accessToken string) getWallUploadServer {
	unResp := getWallUploadServer{}
	resp, err := http.Get("https://api.vk.com/method/" + "photos.getWallUploadServer?" +
		url.Values{
			"access_token": {accessToken},
			"v":            {"5.92"},
			"group_id":     {strconv.Itoa(groupID)}}.Encode())
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&unResp)
	if err != nil {
		log.Fatal(err)
	}
	return unResp
}

func postFile(filename string, targetUrl string) uploadResponse {
	unResp := uploadResponse{}
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	// this step is very important
	fileWriter, err := bodyWriter.CreateFormFile("photo", filename)
	if err != nil {
		log.Fatal(err)
	}

	// open file handle
	fh, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		log.Fatal(err)
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(targetUrl, contentType, bodyBuf)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&unResp)
	if err != nil {
		log.Fatal(err)
	}
	return unResp
}

func vkSavePhoto(upResp uploadResponse, groupID int, accessToken string) savedPhoto {
	unResp := savedPhoto{}
	resp, err := http.Get("https://api.vk.com/method/" + "photos.saveWallPhoto?" +
		url.Values{
			"group_id":     {strconv.Itoa(groupID)},
			"access_token": {accessToken},
			"v":            {"5.92"},
			"server":       {strconv.Itoa(upResp.Server)},
			"hash":         {upResp.Hash},
			"photo":        {upResp.Photo}}.Encode())
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&unResp)
	if err != nil {
		log.Fatal(err)
	}
	return unResp
}

func sendMatchToVk(matchID int64, text string, isFull bool) (err error) {
	unResp := postID{}

	groupID := appconfig.VkGroupID
	accessToken := appconfig.VkAPIkey

	upServer := vkGetWallUploadServer(groupID, accessToken)
	upLink := postFile("tmp/"+strconv.FormatInt(matchID, 10)+".png", upServer.Response.UploadURL)
	upPhoto := vkSavePhoto(upLink, groupID, accessToken)

	resp, err := http.Get("https://api.vk.com/method/" + "wall.post?" +
		url.Values{"owner_id": {strconv.Itoa(-groupID)},
			"access_token": {accessToken},
			"v":            {"5.92"},
			"message":      {text},
			"attachments":  {"photo" + strconv.Itoa(upPhoto.Response[0].OwnerID) + "_" + strconv.Itoa(upPhoto.Response[0].ID) + ",https://www.opendota.com/matches/" + strconv.FormatInt(matchID, 10) + "/"},
			"from_group":   {"1"}}.Encode())
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&unResp)
	if err != nil {
		log.Fatal(err)
	}
	if isFull {
		markFull(matchID)
		addPostID(matchID, unResp.Response.PostID)
	} else {
		markShort(matchID)
		addPostID(matchID, unResp.Response.PostID)
	}
	return nil
}

func editMatchAtVk(matchID int64, post int, text string) (err error) {

	groupID := appconfig.VkGroupID
	accessToken := appconfig.VkAPIkey

	upServer := vkGetWallUploadServer(groupID, accessToken)
	upLink := postFile("tmp/"+strconv.FormatInt(matchID, 10)+".png", upServer.Response.UploadURL)
	upPhoto := vkSavePhoto(upLink, groupID, accessToken)

	resp, err := http.Get("https://api.vk.com/method/" + "wall.edit?" +
		url.Values{"owner_id": {strconv.Itoa(-groupID)},
			"access_token": {accessToken},
			"v":            {"5.92"},
			"message":      {text},
			"post_id":      {strconv.Itoa(post)},
			"attachments":  {"photo" + strconv.Itoa(upPhoto.Response[0].OwnerID) + "_" + strconv.Itoa(upPhoto.Response[0].ID) + ",https://www.opendota.com/matches/" + strconv.FormatInt(matchID, 10) + "/"},
			"from_group":   {"1"}}.Encode())
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	markFull(matchID)
	return nil
}
