package disel

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	CRLF           = "\r\n"
	SPACE          = " "
	HOST_SEP       = ": "
	UA_SEP         = ": "
	CL_SEP         = ": "
	PATH_SEP       = "/"
	TEXT_PLAIN     = "text/plain"
	SUCCESS_TEXT   = "OK"
	NOT_FOUND_TEXT = "Not found"
	HTTP_VERSION   = "HTTP/1.1"

	HTTP            = "HTTP"
	HOST            = "host"
	USER_AGENT      = "User-Agent"
	CONTENT_LENGTH  = "Content-Length"
	ACCEPT_ENCODING = "Accept-Encoding"
	CONTENT_TYPE    = "Content-Type"
)

type HttpRequest struct {
	Method         string
	Path           string
	PathParams     []string
	Host           string
	UserAgent      string
	Version        string
	Body           *strings.Reader
	ContentType    string
	ContentLength  int
	AcceptEncoding string
}

type HttpResponse struct {
	status      int
	contentType string
	body        string
}

func DeserializeRequest(req string) HttpRequest {
	reqArray := strings.Split(req, CRLF)
	var request HttpRequest

	n := len(reqArray)
	if n == 0 {
		return request
	}

	for idx, line := range reqArray {
		if strings.Contains(line, HTTP) {
			request.parseFirstLine(line)
		}
		if strings.HasPrefix(line, HOST) {
			request.parseHost(line)
		}
		if strings.HasPrefix(line, USER_AGENT) {
			request.parseUserAgent(line)
		}
		if strings.HasPrefix(line, CONTENT_TYPE) {
			request.parseContentType(line)
		}
		if strings.HasPrefix(line, CONTENT_LENGTH) {
			request.parseContentLength(line)
		}
		if strings.HasPrefix(line, ACCEPT_ENCODING) {
			request.parseAcceptEncoding(line)
		}
		if idx == n-1 {
			request.parseBody(line)
		}
	}
	return request
}
func (r *HttpRequest) parseFirstLine(firstLineInfo string) {
	firstLineArray := strings.Split(firstLineInfo, SPACE)
	if len(firstLineArray) > 0 {
		method := getIndex(0, len(firstLineArray), firstLineArray, "string").(string)
		pathStr := getIndex(1, len(firstLineArray), firstLineArray, "string").(string)
		pathParams := strings.Split(pathStr, PATH_SEP)
		path := pathParams[1]
		_ = path
		version := getIndex(2, len(firstLineArray), firstLineArray, "string").(string)

		r.Method = method
		r.Path = pathStr //
		r.Version = version
		r.PathParams = pathParams[2:]
	}
}
func (r *HttpRequest) parseHost(hostInfo string) {
	hostArray := strings.Split(hostInfo, HOST_SEP)
	if len(hostArray) >= 2 {
		r.Host = hostArray[1]
	}
}
func (r *HttpRequest) parseUserAgent(userAgentInfo string) {
	userAgentArray := strings.Split(userAgentInfo, UA_SEP)
	if len(userAgentArray) >= 2 {
		r.UserAgent = userAgentArray[1]
	}
}
func (r *HttpRequest) parseAcceptEncoding(acceptEncodingInfo string) {
	acceptEncodingArray := strings.Split(acceptEncodingInfo, UA_SEP)
	if len(acceptEncodingArray) >= 2 {
		r.AcceptEncoding = acceptEncodingArray[1]
	}
}
func (r *HttpRequest) parseContentLength(contentLengthInfo string) {
	contentLengthArray := strings.Split(contentLengthInfo, CL_SEP)
	if len(contentLengthArray) >= 2 {
		cl, err := strconv.Atoi(contentLengthArray[1])
		if err == nil {
			r.ContentLength = cl
		}
	}
}
func (r *HttpRequest) parseBody(bodyInfo string) {
	body := bodyInfo
	r.Body = strings.NewReader(body)
}
func (r *HttpRequest) parseContentType(contentTypeInfo string) {
	contentTypeArray := strings.Split(contentTypeInfo, UA_SEP)
	if len(contentTypeArray) >= 2 {
		r.ContentType = contentTypeArray[1]
	}
}

/* ---- SERIALIZE ------ */

func Serialize(status int, content string, contentType string) string {

	var text string
	if status == 200 {
		text = SUCCESS_TEXT
	} else if status == 404 {
		text = NOT_FOUND_TEXT
	}

	firstLine := fmt.Sprintf("%s %d %s", HTTP_VERSION, status, text)
	if len(content) == 0 {
		return firstLine + "\r\n" + "Content-length: 0" + "\r\n\r\n"
	}
	secondLine := fmt.Sprintf("Content-Type: %s", contentType)
	thirdLine := fmt.Sprintf("Content-Length: %d", len(content))
	return firstLine + "\r\n" + secondLine + "\r\n" + thirdLine + "\r\n\r\n" + content
}

func getIndex(index int, size int, array []string, dataType string) any {
	if index >= size {
		if dataType == "string" {
			return ""
		} else if dataType == "int" {
			return -1
		} else if dataType == "bool" {
			return false
		}
	}
	return array[index]
}
