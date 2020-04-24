package gotwilio

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// MediaRegion is the locations of Twilio's
// TURN servers
type MediaRegion string

const (
	Australia   MediaRegion = "au1"
	Brazil      MediaRegion = "br1"
	Germany     MediaRegion = "de1"
	Ireland     MediaRegion = "ie1"
	India       MediaRegion = "in1"
	Japan       MediaRegion = "jp1"
	Singapore   MediaRegion = "sg1"
	USEastCoast MediaRegion = "us1"
	USWestCoast MediaRegion = "us2"
)

// VideoStatus is the status of a video room
type VideoStatus string

const (
	InProgress VideoStatus = "in-progress"
	Failed     VideoStatus = "failed"
	Completed  VideoStatus = "completed"
)

// VideoRoomType is how the participants connect
// to each other, whether peer-to-peer of routed
// through a TURN server.
type VideoRoomType string

const (
	PeerToPeer VideoRoomType = "peer-to-peer"
	Group      VideoRoomType = "group"
	GroupSmall VideoRoomType = "group-small"
)

// VideoCodecs are the supported codecs when
// publishing a track to the room.
type VideoCodecs string

const (
	VP8  VideoCodecs = "VP8"
	H264 VideoCodecs = "H264"
)

// ListVideoResponse is returned when listing rooms
type ListVideoReponse struct {
	Rooms []*VideoResponse `json:"rooms"`
	Meta  struct {
		Page            int64  `json:"page"`
		PageSize        int64  `json:"page_size"`
		FirstPageUrl    string `json:"first_page_url"`
		PreviousPageUrl string `json:"previous_page_url"`
		NextPageUrl     string `json:"next_page_url"`
		Url             string `json:"url"`
		Key             string `json:"key"`
	} `json:"meta"`
}

// VideoResponse is returned for a single room
type VideoResponse struct {
	AccountSid                  string        `json:"account_sid"`
	DateCreated                 time.Time     `json:"date_created"`
	DateUpdated                 time.Time     `json:"date_updated"`
	Duration                    time.Duration `json:"duration"`
	EnableTurn                  bool          `json:"enable_turn"`
	EndTime                     time.Time     `json:"end_time"`
	MaxParticipants             int64         `json:"max_participants"`
	MediaRegion                 MediaRegion   `json:"media_region"`
	RecordParticipantsOnConnect bool          `json:"record_participants_on_connect"`
	Sid                         string        `json:"sid"`
	Status                      VideoStatus   `json:"status"`
	StatusCallback              string        `json:"status_callback"`
	StatusCallbackMethod        string        `json:"status_callback_method"`
	Type                        VideoRoomType `json:"type"`
	UniqueName                  string        `json:"unique_name"`
	URL                         string        `json:"url"`
}

type createRoomOptions struct {
	EnableTurn                  bool          `json:"EnableTurn"`
	MaxParticipants             int64         `json:"MaxParticipants"`
	MediaRegion                 MediaRegion   `json:"MediaRegion"`
	RecordParticipantsOnConnect bool          `json:"RecordParticipantsOnConnect"`
	StatusCallback              string        `json:"StatusCallback"`
	StatusCallbackMethod        string        `json:"StatusCallbackMethod"`
	Type                        VideoRoomType `json:"Type"`
	UniqueName                  string        `json:"UniqueName"`
	VideoCodecs                 []VideoCodecs `json:"VideoCodecs"`
}

// DefaultVideoRoomOptions are the default options
// for creating a room.
var DefaultVideoRoomOptions = &createRoomOptions{
	EnableTurn:                  true,
	MaxParticipants:             10,
	MediaRegion:                 USEastCoast,
	RecordParticipantsOnConnect: false,
	StatusCallback:              "",
	StatusCallbackMethod:        http.MethodPost,
	Type:                        Group,
	UniqueName:                  "",
	VideoCodecs:                 []VideoCodecs{H264},
}

// ListVideoRoomOptions are the options to query
// for a list of video rooms.
type ListVideoRoomOptions struct {
	DateCreatedAfter  time.Time   `json:"DateCreatedAfter"`
	DateCreatedBefore time.Time   `json:"DateCreatedBefore"`
	Status            VideoStatus `json:"Status"`
	UniqueName        string      `json:"EnableUniqueNameTurn"`
}

// CreateVideoRoom creates a video communication session
// for participants to connect to.
// See https://www.twilio.com/docs/video/api/rooms-resource
// for more information.
func (twilio *Twilio) CreateVideoRoom(options *createRoomOptions) (videoResponse *VideoResponse, exception *Exception, err error) {
	twilioUrl := twilio.VideoUrl + "/v1/Rooms"
	formValues := createRoomOptionsToFormValues(options)

	res, err := twilio.post(formValues, twilioUrl)
	if err != nil {
		return videoResponse, exception, err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return videoResponse, exception, err
	}

	if res.StatusCode != http.StatusCreated {
		exception = new(Exception)
		err = json.Unmarshal(responseBody, exception)

		// We aren't checking the error because we don't actually care.
		// It's going to be passed to the client either way.
		return videoResponse, exception, err
	}

	videoResponse = new(VideoResponse)
	err = json.Unmarshal(responseBody, videoResponse)
	return videoResponse, exception, err
}

// DateCreatedAfter  time.Time   `json:"DateCreatedAfter"`
// DateCreatedBefore time.Time   `json:"DateCreatedBefore"`
// Status            VideoStatus `json:"Status"`
// UniqueName        string      `json:"EnableUniqueNameTurn"`

// ListVideoRooms returns a list of all video rooms.
// See https://www.twilio.com/docs/video/api/rooms-resource
// for more information.
func (twilio *Twilio) ListVideoRooms(options *ListVideoRoomOptions) (videoResponse *ListVideoReponse, exception *Exception, err error) {
	q := &url.Values{}
	if !options.DateCreatedAfter.Equal(time.Time{}) {
		q.Set("DateCreatedAfter", options.DateCreatedAfter.Format(time.RFC3339))
	}
	if !options.DateCreatedBefore.Equal(time.Time{}) {
		q.Set("DateCreatedBefore", options.DateCreatedBefore.Format(time.RFC3339))
	}
	if options.Status != "" {
		q.Set("Status", fmt.Sprintf("%s", options.Status))
	}
	if options.UniqueName != "" {
		q.Set("UniqueName", options.UniqueName)
	}

	twilioUrl := twilio.VideoUrl + "/v1/Rooms?" + q.Encode()

	res, err := twilio.get(twilioUrl)
	if err != nil {
		return videoResponse, exception, err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return videoResponse, exception, err
	}

	if res.StatusCode != http.StatusOK {
		exception = new(Exception)
		err = json.Unmarshal(responseBody, exception)

		// We aren't checking the error because we don't actually care.
		// It's going to be passed to the client either way.
		return videoResponse, exception, err
	}

	videoResponse = new(ListVideoReponse)
	err = json.Unmarshal(responseBody, videoResponse)
	return videoResponse, exception, err
}

// GetVideoRoom retrievs a single video session
// by name or by Sid.
// See https://www.twilio.com/docs/video/api/rooms-resource
// for more information.
func (twilio *Twilio) GetVideoRoom(nameOrSid string) (videoResponse *VideoResponse, exception *Exception, err error) {
	twilioUrl := twilio.VideoUrl + "/v1/Rooms/" + nameOrSid

	res, err := twilio.get(twilioUrl)
	if err != nil {
		return videoResponse, exception, err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return videoResponse, exception, err
	}

	if res.StatusCode != http.StatusOK {
		exception = new(Exception)
		err = json.Unmarshal(responseBody, exception)

		// We aren't checking the error because we don't actually care.
		// It's going to be passed to the client either way.
		return videoResponse, exception, err
	}

	videoResponse = new(VideoResponse)
	err = json.Unmarshal(responseBody, videoResponse)
	return videoResponse, exception, err
}

// EndVideoRoom stops a single video session by name
// or by Sid, and disconnects all participants.
// See https://www.twilio.com/docs/video/api/rooms-resource
// for more information.
func (twilio *Twilio) EndVideoRoom(nameOrSid string) (videoResponse *VideoResponse, exception *Exception, err error) {
	twilioUrl := twilio.VideoUrl + "/v1/Rooms/" + nameOrSid
	formValues := url.Values{}
	formValues.Set("Status", fmt.Sprintf("%s", Completed))

	res, err := twilio.post(formValues, twilioUrl)
	if err != nil {
		return videoResponse, exception, err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return videoResponse, exception, err
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		exception = new(Exception)
		err = json.Unmarshal(responseBody, exception)

		// We aren't checking the error because we don't actually care.
		// It's going to be passed to the client either way.
		return videoResponse, exception, err
	}

	videoResponse = new(VideoResponse)
	err = json.Unmarshal(responseBody, videoResponse)
	return videoResponse, exception, err
}

func createRoomOptionsToFormValues(options *createRoomOptions) url.Values {
	formValues := url.Values{}
	formValues.Set("EnableTurn", fmt.Sprintf("%t", options.EnableTurn))
	formValues.Set("MaxParticipants", fmt.Sprintf("%d", options.MaxParticipants))
	formValues.Set("MediaRegion", fmt.Sprintf("%s", options.MediaRegion))
	formValues.Set("RecordParticipantsOnConnect", fmt.Sprintf("%t", options.RecordParticipantsOnConnect))
	formValues.Set("StatusCallback", options.StatusCallback)
	formValues.Set("StatusCallbackMethod", options.StatusCallbackMethod)
	formValues.Set("Type", fmt.Sprintf("%s", options.Type))
	formValues.Set("UniqueName", options.UniqueName)
	formValues.Del("VideoCodecs")
	for _, codec := range options.VideoCodecs {
		formValues.Add("VideoCodecs", fmt.Sprintf("%v", codec))
	}
	return formValues
}
