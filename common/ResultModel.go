package common

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"
	. "github.com/logrusorgru/aurora"
)

//Implement IResult
type apiResult struct {
	debug            bool            // True if the result should log Debug level logs
	beautify_logs    bool            // True if the result should beautify the logs
	messages         []string        // And array of accumulated logs with context
	was_successful   bool            // True if the result was successful
	log_level        int             // The level to log out the messafes upon Flush
	status_code      int             // The status code that was returned from an external request
	response_message string          // The error Message returned from an external request
	parent           chan []string   // If this is a child, this will be the parent's channel
	children         []chan []string // A slice of the channels connected to any childern that are created
}

/*
	This method is for generating a apiResult
*/
func MakeAPIResult(config IConfigGetter) IResult {
	debug := false
	beautify_logs := false
	log_level := config.SafeGetCinfigVar("LOGGING_LEVEL")

	if log_level == "DEBUG" {
		debug = true
	}
	if log_level == "DEV" {
		debug = true
		beautify_logs = true
	}

	return &apiResult{
		debug:         debug,
		beautify_logs: beautify_logs,
	}
}

/*
	This sets default logging level. It is for testing.
*/
func MakeDefaultAPIRequest() IResult {
	return &apiResult{
		debug:         true,
		beautify_logs: false,
	}
}

func (this *apiResult) GetChild() IResult {
	channel := make(chan []string)
	child := &apiResult{
		debug:         this.debug,
		beautify_logs: this.beautify_logs,
		parent:        channels,
	}
	this.childern = append(this.childern, channel)

	parent_log := fmt.Sprintf("[CHILD #%s STARTED]\n", strconv.Itoa(len(this.children)))
	this.messages = append(this.messages, parent_log)
	child_log := fmt.Sprintf("[CHILD #%s OUTPUT]\n", strconv.Itoa(len(this.children)))
	child.messages = append(child.messages, child_log)

	return child
}

/*
	This function returns success flag value.
*/
func (this *apiResult) WasSuccessful() bool {
	return this.was_successful
}

/*
	This sets result success flag to true
*/
func (this *apiResult) Succeed() {
	this.was_successful = true
	if this.log_level > 1 {
		this.log_level = 1
	}
}

/*
	This function returns success flag value.
*/
func (this *apiResult) Fail() {
	this.was_successful = false
}

/*
	This implements error interface
*/
func (this *apiResult) Error() string {
	if this.messages == nil {
		return ""
	}
	return this.messages[len(this.messages)-1]
}

/*
	This is used for merging results
*/
func (this *apiResult) MergeWithResult(r IResult) {
	if r == nil {
		return
	}
	for _, v := range r.GetMessages() {
		this.messages = append(this.messages, v)
	}

	if this.log_level < r.GetLogLevel() {
		this.log_level = r.GetLogLevel()
	}

	this.children = append(this.children, r.GetChildren()...)
	this.response_message = r.GetResponseMessage()
	this.status_code = r.GetStatusCode()
}

func (this *apiResult) GetChildren() []chan []string {
	return this.children
}

/*
	This returns messages in result
*/
func (this *apiResult) GetMessages() []string {
	return this.messages
}

/*
	This returns current logging level
*/
func (this *apiResult) GetLogLevel() int {
	return this.log_level
}

/*
	Get Status Code
*/
func (this *apiResult) GetStatusCode() int {
	return this.status_code
}

/*
	Set Status Code
*/
func (this *apiResult) SetStatusCode(code int) {
	this.status_code = code
}

/*
	Get response error
*/
func (this *apiResult) GetResponseMessage() string {
	return this.response_message
}

/*
	Set response error
*/
func (this *apiResult) SetResponseMessage(msg string) {
	this.response_message = msg
}

/*
	Append the debug logs
*/
func (this *apiResult) Debugf(templates string, args ...interface{}) {
	if !this.debug {
		return
	}

	original_message := fmt.Sprintf(Blue(templates), args...)
	this.addLog("[Debug]", original_message)
}

/*
	Append the info logs
*/
func (this *apiResult) Infof(templates string, args ...interface{}) {
	if this.log_level < 1 {
		this.log_level = 1
	}
	original_message := fmt.Sprintf(Green(templates), args...)
	this.addLog("[Info]", original_message)
}

/*
	Append the error logs
*/
func (this *apiResult) Errorf(templates string, args ...interface{}) {
	if this.log_level < 2 {
		this.log_level = 2
	}
	original_message := fmt.Sprintf(Red(templates), args...)
	this.addLog("[Error]", original_message)
}

/*
	Append the warning logs
*/
func (this *apiResult) Warningf(templates string, args ...interface{}) {
	if this.log_level < 2 {
		this.log_level = 2
	}
	original_message := fmt.Sprintf(templates, args...)
	this.addLog("[Warning]", original_message)
}

/*
	Helper function to add logs
*/
func (this *apiResult) addLog(header string, org_msg string) {
	_, file, line, _ := runtime.Caller(2)

	org_msg = strings.TrimSuffix(org_msg, "\n")
	output := fmt.Sprintf(
		header+" %s %s %s::%d",
		org_msg,
		time.Now(),
		file,
		line,
	)
	if !strings.HasSuffix(output, "\n") {
		output += "\n"
	}

	this.messages = append(this.messages, output)
}

/*
	This add contextual information for log
*/
func (this *apiResult) DebugMessagef(templates string, args ...interface{}) {
	if !this.debug {
		return
	}

	if !strings.HasSuffix(templates, "\n") {
		templates += "\n"
	}

	original_message := fmt.Sprintf(templates, args...)
	this.messages = append(this.messages, fmt.Sprintf(
		"[Message] %s",
		original_message,
	))
}

/*
	This function combine all logs and flush it on screen.
	This get used at top-most level
*/
func (this *apiResult) Flush() {

	my_logs_length := len(this.messages)

	for i, child := range this.children {
		select {
		case child_output := <-child:
			this.messages = append(this.messages, child_output...)
		case <-time.After(time.Minute * 5):
			this.Errorf("CHILD %d DID NOT COME HOME!! We're flushing without them", i+1)
		}
	}

	output := ""
	if this.beautify_logs {
		for i, msg := range this.messages {
			if i < my_logs_length {
				num := strconv.Itoa(i)
				this.messages[i] = num + " " + msg
			} else {
				this.messages[i] = "-" + msg
			}
			output += this.messages[i]
		}
	} else {
		for i, msg := range this.messages {
			if i < my_logs_length {
				num := strconv.Itoa(i)
				this.messages[i] = num + ") " + msg
			} else {
				this.messages[i] = "-" + msg
			}
			output += this.messages[i]
		}
		output = strings.Replace(output, "\n", " :|: ", -1)
	}

	if this.parent != nil {
		select {
		case this.parent <- this.messages:
		case <-time.After(time.Minute * 5):
			this.Errorf("PARENT NOT LISTENING!! We'll move on without them")
			fmt.Println(output)
		}
	} else {
		fmt.Println(output)
	}

	this.messages = []string{}
}
