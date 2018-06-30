package dropbox

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	errors "github.com/joaosoft/errors"
	manager "github.com/joaosoft/manager"
)

type writeMode string

const (
	writeModeAdd       writeMode = "add"
	writeModeOverwrite           = "overwrite"
)

type File struct {
	client manager.IGateway
	config *DropboxConfig
}

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

func (f *File) Upload(path string, file []byte) (*uploadFileResponse, *errors.ErrorData) {
	var err error
	var bodyArgs []byte
	args := uploadFileRequest{
		Path:       path,
		Mode:       writeModeOverwrite,
		AutoRename: true,
		Mute:       false,
	}

	if bodyArgs, err = json.Marshal(args); err != nil {
		newErr := errors.NewError(err)
		log.Error("errors converting upload input arguments").ToErrorData(newErr)
		return nil, newErr
	}

	headers := manager.Headers{
		"Authorization":   {fmt.Sprintf("%s %s", f.config.Authorization.Access, f.config.Authorization.Token)},
		"Content-Type":    {"application/octet-stream"},
		"Dropbox-API-Arg": {string(bodyArgs)},
	}

	dropboxResponse := &uploadFileResponse{}
	if err != nil {
		newErr := errors.NewError(err)
		log.Error("errors marshal arguments").ToErrorData(newErr)
		return nil, newErr
	}
	if status, response, err := f.client.Request(http.MethodPost, f.config.Hosts.Content, "/files/upload", headers, file); err != nil {
		newErr := errors.NewError(err)
		log.WithField("response", response).Error("errors uploading File").ToErrorData(newErr)
		return nil, newErr
	} else if status != http.StatusOK {
		var err error
		log.WithField("response", response).Errorf("response status %d instead of %d", status, http.StatusOK).ToError(&err)
		return nil, errors.NewError(err)
	} else if response == nil {
		var err error
		log.Error("errors uploading File").ToError(&err)
		return nil, errors.NewError(err)
	} else {
		if err := json.Unmarshal(response, dropboxResponse); err != nil {
			newErr := errors.NewError(err)
			log.Error("errors converting Dropbox response data").ToErrorData(newErr)
			return nil, newErr
		}
		return dropboxResponse, nil
	}

	return nil, nil
}

type downloadFileRequest struct {
	Path string `json:"path"`
}

func (f *File) Download(path string) ([]byte, *errors.ErrorData) {
	var err error
	var bodyArgs []byte
	args := downloadFileRequest{
		Path: path,
	}

	if bodyArgs, err = json.Marshal(args); err != nil {
		newErr := errors.NewError(err)
		log.Error("errors converting download input arguments").ToErrorData(newErr)
		return nil, newErr
	}

	headers := manager.Headers{
		"Authorization":   {fmt.Sprintf("%s %s", f.config.Authorization.Access, f.config.Authorization.Token)},
		"Dropbox-API-Arg": {string(bodyArgs)},
	}

	if status, response, err := f.client.Request(http.MethodPost, f.config.Hosts.Content, "/files/download", headers, []byte("")); err != nil {
		newErr := errors.NewError(err)
		log.WithField("response", response).Error("errors downloading File").ToErrorData(newErr)
		return nil, newErr
	} else if status != http.StatusOK {
		var err error
		log.WithField("response", response).WithFields(map[string]interface{}{"response": string(response)}).Errorf("response status %d instead of %d", status, http.StatusOK).ToError(&err)
		return nil, errors.NewError(err)
	} else if response == nil {
		var err error
		log.Error("errors downloading File").ToError(&err)
		return nil, errors.NewError(err)
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

func (f *File) Delete(path string) (*deleteFileResponse, *errors.ErrorData) {
	if path == "/" {
		path = ""
	}
	body, err := json.Marshal(deleteFileRequest{
		Path: path,
	})
	if err != nil {
		newErr := errors.NewError(err)
		log.Error("errors marshal arguments").ToErrorData(newErr)
		return nil, newErr
	}

	headers := manager.Headers{
		"Authorization": {fmt.Sprintf("%s %s", f.config.Authorization.Access, f.config.Authorization.Token)},
		"Content-Type":  {"application/json"},
	}

	dropboxResponse := &deleteFileResponse{}
	if status, response, err := f.client.Request(http.MethodPost, f.config.Hosts.Api, "/files/delete_v2", headers, body); err != nil {
		newErr := errors.NewError(err)
		log.WithField("response", response).Error("errors deleting File").ToErrorData(newErr)
		return nil, newErr
	} else if status != http.StatusOK {
		var err error
		log.WithField("response", response).Errorf("response status %d instead of %d", status, http.StatusOK).ToError(&err)
		return nil, errors.NewError(err)
	} else if response == nil {
		var err error
		log.Error("errors deleting File").ToError(&err)
		return nil, errors.NewError(err)
	} else {
		if err := json.Unmarshal(response, dropboxResponse); err != nil {
			newErr := errors.NewError(err)
			log.Error("errors converting Dropbox response data").ToErrorData(newErr)
			return nil, newErr
		}
		return dropboxResponse, nil
	}

	return nil, nil
}
