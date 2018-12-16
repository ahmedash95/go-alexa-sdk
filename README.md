# Alexa SDK

The purpose of this SDK is to make it easy to build Alexa skills without too much code.
it should contain all heavy logic and provide easy to use methods to create skills

what I have reached for now is somthing like below 

```Go
alexa := alexa.New()
alexa.AddIntentResponse("INTNET-NAME", HandlerFunc)
// run the http server on specific port
alexa.Run("3000")
```

# Dependinces
I used Gin framework to run the http server, probably I will use native net/http library to avoid un-used features of any framework we don't need

# How it works
alot of web/video tutorials using only aws-lambda which is an easy way to build serverless apps, but I wanted to build it on my own server.

So this SDK handles all what you need to serve alexa skill, it will handle the operations to deal with alexa and you should only fouce on your skill.

## Create Alexa instances

```GO
// alexa instance
alexa := alexa.New()

// run on port
alexa.Run("3000")
```

## Create Intent response handler

the handler accepts alexa.Request object and should return instance of StringPayload (more response types to be supported)

```Go
alexa.AddIntentResponse("WhatIsYourNameIntent",whatIsYourNameAlexaHandler)

func whatIsYourNameAlexaHandler(request alexa.Request) alexa.StringPayload {
    return alexa.StringPayload{
		Title: "ask for user name",
		Text:  "What is your name ?",
	}
}
```

## Handle Slots
you can access your intent slots using the alexa request parameter

```Go
func ResponseToAnswerWhatIsYourName(request alexa.Request) alexa.StringPayload {
    slots := request.Body.Intent.Slots
	name := slots["name"]

    return alexa.StringPayload{
		Title: "hello name",
		Text:  fmt.Sprintf("Hello %s", name.Value),
	}
}
```

## Request Object
Alexa request object contains a lot of helpfull structs to handle almost what you need for the logic

```Go
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
```


# Todo
- Implement alexa authentication
- Implement images for Echo (Alexa) with screen

# Contributing
feel free to submit PR for
