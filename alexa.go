package alexa

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/ahmedash95/amazon-alexa-sdk/storage"
	"github.com/gin-gonic/gin"
)

/**Alexa is the main struct for this sdk,
it helps to manage all intents with slots and running the http server
*/
type Alexa struct {
	Intent          map[string]Intent
	IntentsResponse map[string]func(request Request) StringPayload `json:"intents_response"`
}

// Intent is the Alexa skill intent
type Intent struct {
	Name  string          `json:"name"`
	Slots map[string]Slot `json:"slots"`
}

// Slot is an Alexa skill slot
type Slot struct {
	Name        string      `json:"name"`
	Value       string      `json:"value"`
	Resolutions Resolutions `json:"resolutions"`
}

type Resolutions struct {
	ResolutionPerAuthority []struct {
		Values []struct {
			Value struct {
				Name string `json:"name"`
				Id   string `json:"id"`
			} `json:"value"`
		} `json:"values"`
	} `json:"resolutionsPerAuthority"`
}

type Request struct {
	Version string  `json:"version"`
	Session Session `json:"session"`
	Body    ReqBody `json:"request"`
	Context Context `json:"context"`
}

// Session represents the Alexa skill session
type Session struct {
	New         bool   `json:"new"`
	SessionID   string `json:"sessionId"`
	Application struct {
		ApplicationID string `json:"applicationId"`
	} `json:"application"`
	Attributes map[string]interface{} `json:"attributes"`
	User       struct {
		UserID      string `json:"userId"`
		AccessToken string `json:"accessToken,omitempty"`
	} `json:"user"`
}

// Context represents the Alexa skill request context
type Context struct {
	System struct {
		APIAccessToken string `json:"apiAccessToken"`
		Device         struct {
			DeviceID string `json:"deviceId,omitempty"`
		} `json:"device,omitempty"`
		Application struct {
			ApplicationID string `json:"applicationId,omitempty"`
		} `json:"application,omitempty"`
	} `json:"System,omitempty"`
}

// ReqBody is the actual request information
type ReqBody struct {
	Type        string `json:"type"`
	RequestID   string `json:"requestId"`
	Timestamp   string `json:"timestamp"`
	Locale      string `json:"locale"`
	Intent      Intent `json:"intent,omitempty"`
	Reason      string `json:"reason,omitempty"`
	DialogState string `json:"dialogState,omitempty"`
}

// Response is the response back to the Alexa speech service
type Response struct {
	Version           string                 `json:"version"`
	SessionAttributes map[string]interface{} `json:"sessionAttributes,omitempty"`
	Body              ResBody                `json:"response"`
}

// ResBody is the actual body of the response
type ResBody struct {
	OutputSpeech     *Payload     `json:"outputSpeech,omitempty"`
	Card             *Payload     `json:"card,omitempty"`
	Reprompt         *Reprompt    `json:"reprompt,omitempty"`
	Directives       []Directives `json:"directives,omitempty"`
	ShouldEndSession bool         `json:"shouldEndSession"`
}

// Reprompt is imformation
type Reprompt struct {
	OutputSpeech Payload `json:"outputSpeech,omitempty"`
}

// Directives is imformation
type Directives struct {
	Type          string         `json:"type,omitempty"`
	SlotToElicit  string         `json:"slotToElicit,omitempty"`
	UpdatedIntent *UpdatedIntent `json:"UpdatedIntent,omitempty"`
	PlayBehavior  string         `json:"playBehavior,omitempty"`
	AudioItem     struct {
		Stream struct {
			Token                string `json:"token,omitempty"`
			URL                  string `json:"url,omitempty"`
			OffsetInMilliseconds int    `json:"offsetInMilliseconds,omitempty"`
		} `json:"stream,omitempty"`
	} `json:"audioItem,omitempty"`
}

// UpdatedIntent is to update the Intent
type UpdatedIntent struct {
	Name               string                 `json:"name,omitempty"`
	ConfirmationStatus string                 `json:"confirmationStatus,omitempty"`
	Slots              map[string]interface{} `json:"slots,omitempty"`
}

// Image ...
type Image struct {
	SmallImageURL string `json:"smallImageUrl,omitempty"`
	LargeImageURL string `json:"largeImageUrl,omitempty"`
}

// Payload ...
type Payload struct {
	Type    string `json:"type,omitempty"`
	Title   string `json:"title,omitempty"`
	Text    string `json:"text,omitempty"`
	SSML    string `json:"ssml,omitempty"`
	Content string `json:"content,omitempty"`
	Image   Image  `json:"image,omitempty"`
}

type StringPayload struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type User struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

var currentRequest Request

func getResponsePayload(response StringPayload) Response {
	r := Response{
		Version: "1.0",
		Body: ResBody{
			OutputSpeech: &Payload{
				Type: "PlainText",
				Text: response.Text,
			},
			Card: &Payload{
				Type:    "Simple",
				Title:   response.Title,
				Content: response.Text,
			},
			ShouldEndSession: true,
		},
	}
	return r
}

func New() Alexa {
	var alexa Alexa
	alexa.IntentsResponse = make(map[string]func(request Request) StringPayload)

	return alexa
}

func (alexa *Alexa) AddIntentResponse(name string, callback func(request Request) StringPayload) {
	alexa.IntentsResponse[name] = callback
}

func (alexa *Alexa) Run(port string) {
	storage.Start()

	router := gin.Default()

	router.POST("/", func(c *gin.Context) {
		var alexaRequest Request
		requestBody, _ := c.GetRawData()
		json.Unmarshal(requestBody, &alexaRequest)

		currentRequest = alexaRequest

		intentName := alexaRequest.Body.Intent.Name

		var response StringPayload

		intent, Ok := alexa.IntentsResponse[intentName]
		if Ok {
			response = intent(alexaRequest)
		} else {
			response = StringPayload{Text: "Intent not found"}
		}

		c.JSON(200, getResponsePayload(response))
	})

	router.Run(":3000")
}

const AlexaStoregeUserKey = "alexaStorageUserInfo_"

func GetUser() User {
	var user User
	info, err := storage.Get(AlexaStoregeUserKey + currentRequest.Session.User.UserID)
	if err == nil {
		json.Unmarshal(info, &user)
	}
	return user
}

func SetUser(u User) (bool, error) {
	userByte, _ := json.Marshal(u)
	if strings.TrimSpace(u.Name) == "" {
		return false, errors.New("the user's name must not be empty")
	}

	storage.Set(AlexaStoregeUserKey+currentRequest.Session.User.UserID, userByte)
	return true, nil
}
