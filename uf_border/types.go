package uf_border

import "fmt"

// dummy type so Sprintf won't fall into recursion
type _SendBlueIn SendBlueIn

func (in SendBlueIn) String() string {
	s := fmt.Sprintf("%v", _SendBlueIn(in))
	return s
}

// smsIn comes from sendblue.co. Comments are taken from
// https://sendblue.co/docs/inbound
type SendBlueIn struct {
	// Only received if this is a group message
	GroupId string `json:"group_id"`

	// Associated account email
	AccountEmail string `json:"accountEmail"`

	// Message content
	Content string `json:"content"`

	// A CDN link to the image that was sent to your sendblue number from an
	// end-user. This link expires after 30 days
	MediaUrl string `json:"media_url"`

	// True if message is sent, false if message is received
	IsOutbound bool `json:"is_outbound"`

	// The current status of the message
	Status string `json:"status"`

	// error code (null if no error)
	ErrorCode int `json:"error_code"`

	// descriptive error message (null if no error)
	ErrorMessage string `json:"error_message"`

	// Sendblue message handle
	MessageHandler string `json:"message_handler"`

	// ISO 8601 formatted date string of the date this message was created
	DateSent string `json:"date_sent"`

	// ISO 8601 formatted date string of the date this message was last updated
	DateUpdated string `json:"date_updated"`

	// E.164 formatted phone number string of the message dispatcher
	FromNumber string `json:"from_number"`

	// E.164 formatted phone number string of your end-user (not the
	// Sendblue-provided phone number)
	Number string `json:"number"`

	// E.164 formatted phone number string of the message recipient
	ToNumber string `json:"to_number"`

	// true if the end user does not support iMessage, false otherwise
	WasDowngraded bool `json:"was_downgraded"`

	// Value of the Sendblue account plan
	Plan string `json:"plan"`
}

// smsOut comes from sendblue.co. Comments are taken from
// https://sendblue.co/docs/outbound
// Endgpoint: POST https://api.sendblue.co/api/send-message
// Headers:
//"sb-api-key-id": "d607fd40ff448cf9219ceec7c41f6dc5",
//"sb-api-secret-key": "a10df202e09c763cbaf291cfa2b361b2",
//"content-type": "application/json"
type SmsOut struct {
	// The number of the recipient of the message
	Number string `json:"number"`

	// The content of the message
	Content string `json:"content"`

	// The style of delivery of the message (see expressive messages)
	SendStyle string `json:"send_style"`

	// A CDN link to a file which is publicly accessible which will be
	// downloaded and sent to the group on our end
	MediaUrl string `json:"media_url"`

	// The URL where you want to receive the status updates of the message
	StatusCallback string `json:"statusCallback"`
}

// smsGroupOut comes from sendblue.co. Comments are taken from
// https://sendblue.co/docs/groups
// Endpoint:  POST https://api.sendblue.co/api/send-group-message
// Headers:
//"sb-api-key-id": "d607fd40ff448cf9219ceec7c41f6dc5",
//"sb-api-secret-key": "a10df202e09c763cbaf291cfa2b361b2",
//"content-type": "application/json"
type SmsGroupOut struct {
	// numbers is an array of strings which contain the E.164-formatted phone
	// numbers of the desired recipients in a group chat. The maximum number of
	// people allowed in a group chat is 25.
	Numbers []string `json:"numbers"`

	// The group_id field is a uuid with which you can message groups that you
	// have already created. This is the same as passing the same list of
	// numbers as was passed in the initial request
	GroupId string `json:"group_id"`

	// The content of the message
	Content string `json:"content"`

	// The style of delivery of the message (see expressive messages)
	SendStyle string `json:"send_style"`

	// A CDN link to a file which is publicly accessible which will be
	// downloaded and sent to the group on our end
	MediaUrl string `json:"media_url"`

	// The URL where you want to receive the status updates of the message
	StatusCallback string `json:"statusCallback"`
}

// in BETA, could be deprecated
// https://sendblue.co/docs/groups#example-1
type SmsGroupOutResponse struct {
	//"accountEmail": "nikita.jerschow@gmail.com",
	//"content": "please",
	//"is_outbound": true,
	//"status": "QUEUED",
	//"error_code": null,
	//"error_message": null,
	//"message_handle": "073c1408-a6d9-48e2-ae8c-01f026443833",
	//"date_sent": "2021-05-19T23:07:23.371Z",
	//"date_updated": "2021-05-19T23:07:23.371Z",
	//"from_number": "+13322175641",
	//"number": [  // will become 'numbers'
	//"+19173599290",
	//"+18582430295"
	//],
	//"to_number": [  // will become 'to_numbers'
	//"+19173599290",
	//"+18582430295"
	//],
	//"was_downgraded": null,
	//"plan": "blue",
	//"media_url": "https://source.unsplash.com/random.png",
	//"message_type": "group",
	//"group_id": "66e3b90d-4447-43c6-9439-15a694408ac2"
}

// https://sendblue.co/docs/outbound/#status-callback
type StatusCallbackResponse struct {
	//accountEmail	string	Associated account email
	//content	string	Message content
	//is_outbound	boolean	True if message is sent, false if message is received
	//media_url	string	A CDN link to the image that you sent our servers
	//status	string	The current status of the message
	//error_code	int	error code (null if no error)
	//error_message	string	descriptive error message (null if no error)
	//message_handle	string	Sendblue message handle
	//date_sent	string	ISO 8601 formatted date string of the date this message was created
	//date_updated	string	ISO 8601 formatted date string of the date this message was last updated
	//from_number	string	E.164 formatted phone number string of the message dispatcher
	//number	string	E.164 formatted phone number string of your end-user (not the Sendblue-provided phone number)
	//to_number	string	E.164 formatted phone number string of the message recipient
	//was_downgraded	boolean	true if the end user does not support iMessage, false otherwise
	//plan	string	Value of the Sendblue account plan
}
