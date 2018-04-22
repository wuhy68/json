package godropbox

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/joaosoft/go-error/service"
	"github.com/joaosoft/go-manager/service"
)

type file struct {
	client gomanager.IGateway
	config *DropboxConfig
}

type writeMode string

const (
	writeModeAdd       writeMode = "add"
	writeModeOverwrite           = "overwrite"
)

type uploadFileRequest struct {
	Path       string    `json:"path"`
	Mode       writeMode `json:"mode"`
	AutoRename bool      `json:"autorename"`
	Mute       bool      `json:"mute"`
}

type uploadFileResponse struct {
	Name           string    `json:"name"`
	ID             string    `json:"id"`
	ClientModified time.Time `json:"client_modified"`
	ServerModified time.Time `json:"server_modified"`
	Rev            string    `json:"rev"`
	Size           int       `json:"size"`
	PathLower      string    `json:"path_lower"`
	PathDisplay    string    `json:"path_display"`
	SharingInfo    struct {
		ReadOnly             bool   `json:"read_only"`
		ParentSharedFolderID string `json:"parent_shared_folder_id"`
		ModifiedBy           string `json:"modified_by"`
	} `json:"sharing_info"`
	PropertyGroups []struct {
		TemplateID string `json:"template_id"`
		Fields     []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"fields"`
	} `json:"property_groups"`
	HasExplicitSharedMembers bool   `json:"has_explicit_shared_members"`
	ContentHash              string `json:"content_hash"`
}

func (f *file) Upload(path string, file []byte) (*uploadFileResponse, *goerror.ErrorData) {
	var err error
	var body []byte
	args := uploadFileRequest{
		Path:       path,
		Mode:       writeModeOverwrite,
		AutoRename: true,
		Mute:       false,
	}

	if body, err = json.Marshal(args); err != nil {
		newErr := goerror.NewError(err)
		log.Error("error converting upload input body").ToErrorData(newErr)
		return nil, newErr
	}

	headers := gomanager.Headers{
		"Authorization":   {fmt.Sprintf("%s %s", f.config.Authorization.Access, f.config.Authorization.Token)},
		"Content-Type":    {"application/octet-stream"},
		"Dropbox-API-Arg": {string(body)},
	}

	dropboxResponse := &uploadFileResponse{}
	if status, response, err := f.client.Request(http.MethodPost, f.config.Hosts.Content, "/files/upload", headers, string(file)); err != nil {
		newErr := goerror.NewError(err)
		log.WithField("response", response).Error("error uploading file").ToErrorData(newErr)
		return nil, newErr
	} else if status != http.StatusOK {
		var err error
		log.WithField("response", response).Errorf("response status %d instead of %d", status, http.StatusOK).ToError(&err)
		return nil, goerror.NewError(err)
	} else if response == nil {
		var err error
		log.Error("error uploading file").ToError(&err)
		return nil, goerror.NewError(err)
	} else {
		if err := json.Unmarshal(response, dropboxResponse); err != nil {
			newErr := goerror.NewError(err)
			log.Error("error converting Dropbox response data").ToErrorData(newErr)
			return nil, newErr
		}
		return dropboxResponse, nil
	}

	return nil, nil
}

type downloadFileRequest struct {
	Path string `json:"path"`
}

func (f *file) Download(path string) ([]byte, *goerror.ErrorData) {
	var err error
	var body []byte
	args := downloadFileRequest{
		Path: path,
	}

	if body, err = json.Marshal(args); err != nil {
		newErr := goerror.NewError(err)
		log.Error("error converting download input body").ToErrorData(newErr)
		return nil, newErr
	}

	headers := gomanager.Headers{
		"Authorization":   {fmt.Sprintf("%s %s", f.config.Authorization.Access, f.config.Authorization.Token)},
		"Dropbox-API-Arg": {string(body)},
	}

	if status, response, err := f.client.Request(http.MethodPost, f.config.Hosts.Content, "/files/download", headers, nil); err != nil {
		newErr := goerror.NewError(err)
		log.WithField("response", response).Error("error downloading file").ToErrorData(newErr)
		return nil, newErr
	} else if status != http.StatusOK {
		var err error
		log.WithField("response", response).WithFields(map[string]interface{}{"response": string(response)}).Errorf("response status %d instead of %d", status, http.StatusOK).ToError(&err)
		return nil, goerror.NewError(err)
	} else if response == nil {
		var err error
		log.Error("error downloading file").ToError(&err)
		return nil, goerror.NewError(err)
	} else {
		return response, nil
	}

	return nil, nil
}

type deleteFileRequest struct {
	Path string `json:"path"`
}

type deleteFileResponse struct {
	Metadata struct {
		Tag            string    `json:".tag"`
		Name           string    `json:"name"`
		ID             string    `json:"id"`
		ClientModified time.Time `json:"client_modified"`
		ServerModified time.Time `json:"server_modified"`
		Rev            string    `json:"rev"`
		Size           int       `json:"size"`
		PathLower      string    `json:"path_lower"`
		PathDisplay    string    `json:"path_display"`
		SharingInfo    struct {
			ReadOnly             bool   `json:"read_only"`
			ParentSharedFolderID string `json:"parent_shared_folder_id"`
			ModifiedBy           string `json:"modified_by"`
		} `json:"sharing_info"`
		PropertyGroups []struct {
			TemplateID string `json:"template_id"`
			Fields     []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			} `json:"fields"`
		} `json:"property_groups"`
		HasExplicitSharedMembers bool   `json:"has_explicit_shared_members"`
		ContentHash              string `json:"content_hash"`
	} `json:"metadata"`
}

func (f *file) Delete(path string) (*deleteFileResponse, *goerror.ErrorData) {
	if path == "/" {
		path = ""
	}
	body := deleteFileRequest{
		Path: path,
	}

	headers := gomanager.Headers{
		"Authorization": {fmt.Sprintf("%s %s", f.config.Authorization.Access, f.config.Authorization.Token)},
		"Content-Type":  {"application/json"},
	}

	dropboxResponse := &deleteFileResponse{}
	if status, response, err := f.client.Request(http.MethodPost, f.config.Hosts.Api, "/files/delete_v2", headers, body); err != nil {
		newErr := goerror.NewError(err)
		log.WithField("response", response).Error("error deleting file").ToErrorData(newErr)
		return nil, newErr
	} else if status != http.StatusOK {
		var err error
		log.WithField("response", response).Errorf("response status %d instead of %d", status, http.StatusOK).ToError(&err)
		return nil, goerror.NewError(err)
	} else if response == nil {
		var err error
		log.Error("error deleting file").ToError(&err)
		return nil, goerror.NewError(err)
	} else {
		if err := json.Unmarshal(response, dropboxResponse); err != nil {
			newErr := goerror.NewError(err)
			log.Error("error converting Dropbox response data").ToErrorData(newErr)
			return nil, newErr
		}
		return dropboxResponse, nil
	}

	return nil, nil
}
