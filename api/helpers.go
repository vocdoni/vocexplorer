package api

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	types "gitlab.com/vocdoni/vocexplorer/api/dbtypes"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/logger"
)

//PingServer pings the web server
func PingServer() bool {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := c.Get("/ping")
	if err != nil || resp == nil {
		return false
	}
	return true
}

// For requests where we don't want to ReadAll the response body
func requestBody(url string) (io.ReadCloser, bool) {
	c := &http.Client{
		Timeout: 15 * time.Second,
	}
	resp, err := c.Get(url)
	if resp == nil {
		return nil, false
	}
	if err != nil {
		logger.Error(err)
		return resp.Body, false
	}
	if resp.StatusCode != http.StatusOK {
		return resp.Body, false
	}
	return resp.Body, true
}

func getHeight(url string) (int64, bool) {
	body, ok := requestBody(url)
	if body != nil {
		defer body.Close()
	}
	if !ok {
		return 0, false
	}
	var height types.Height
	err := json.NewDecoder(body).Decode(&height)
	if err != nil {
		logger.Error(err)
	}
	return height.Height, true
}

func getHeightMap(url string) (map[string]int64, bool) {
	body, ok := requestBody(url)
	if body != nil {
		defer body.Close()
	}
	if !ok {
		return map[string]int64{}, false
	}
	var heightMap types.HeightMap
	err := json.NewDecoder(body).Decode(&heightMap)
	if err != nil {
		logger.Error(err)
	}
	return heightMap.Heights, true
}

func getItemList(itemType interface{}, url string) (interface{}, bool) {
	body, ok := requestBody(url)
	if body != nil {
		defer body.Close()
	}
	if !ok {
		return nil, false
	}
	var rawItemList types.ItemList
	err := json.NewDecoder(body).Decode(&rawItemList)
	if err != nil {
		logger.Error(err)
		return nil, false
	}
	switch itemType.(type) {
	case *types.StoreBlock:
		itemList := [config.ListSize]*types.StoreBlock{}
		for i, rawItem := range rawItemList.Items {
			if len(rawItem) > 0 {
				var item types.StoreBlock
				err = json.Unmarshal(rawItem, &item)
				itemList[i] = &item
				if err != nil {
					logger.Error(err)
				}
			}
		}
		return itemList, true
	case *types.Transaction:
		itemList := [config.ListSize]*types.Transaction{}
		for i, rawItem := range rawItemList.Items {
			if len(rawItem) > 0 {
				var item types.Transaction
				err = json.Unmarshal(rawItem, &item)
				itemList[i] = &item
				if err != nil {
					logger.Error(err)
				}
			}
		}
		return itemList, true
	case *types.Envelope:
		itemList := [config.ListSize]*types.Envelope{}
		for i, rawItem := range rawItemList.Items {
			if len(rawItem) > 0 {
				var item types.Envelope
				err = json.Unmarshal(rawItem, &item)
				itemList[i] = &item
				if err != nil {
					logger.Error(err)
				}
			}
		}
		return itemList, true
	case *types.Process:
		itemList := [config.ListSize]*types.Process{}
		logger.Info("unmarshalling processes")
		for i, rawItem := range rawItemList.Items {
			logger.Info("unmarshalling process")
			if len(rawItem) > 0 {
				var item types.Process
				err = json.Unmarshal(rawItem, &item)
				if err != nil {
					logger.Error(err)
				}
				itemList[i] = &item
			}
		}
		return itemList, true
	case string:
		logger.Info("unmarshalling strings")
		itemList := [config.ListSize]string{}
		for i, rawItem := range rawItemList.Items {
			if len(rawItem) > 0 {
				item := string(rawItem)
				itemList[i] = item
			}
		}
		return itemList, true
	case *types.Validator:
		itemList := [config.ListSize]*types.Validator{}
		for i, rawItem := range rawItemList.Items {
			if len(rawItem) > 0 {
				var item types.Validator
				err = json.Unmarshal(rawItem, &item)
				itemList[i] = &item
				if err != nil {
					logger.Error(err)
				}
			}
		}
		return itemList, true
	}
	return nil, false
}
