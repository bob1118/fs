package eventsocket

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
)

// EventHeader represents events as a pair of key:value.
type EventHeader map[string]interface{}

// Event represents a FreeSWITCH event.
type Event struct {
	Header EventHeader // Event headers, key:val
	Body   string      // Raw body, available in some events
}

func (r *Event) String() string {
	if r.Body == "" {
		return fmt.Sprintf("%s", r.Header)
	} else {
		return fmt.Sprintf("%s\n%s", r.Header, r.Body)
	}
}

// Get returns an Event value, or "" if the key doesn't exist.
func (r *Event) Get(key string) string {
	val, ok := r.Header[key]
	if !ok || val == nil {
		return ""
	}
	if s, ok := val.([]string); ok {
		return strings.Join(s, ", ")
	}
	return val.(string)
}

// GetInt returns an Event value converted to int, or an error if conversion
// is not possible.
func (r *Event) GetInt(key string) (int, error) {
	n, err := strconv.Atoi(r.Header[key].(string))
	if err != nil {
		return 0, err
	}
	return n, nil
}

// PrettyPrint prints Event headers and body to the standard output.
func (r *Event) PrettyPrint() {
	var keys []string
	for k := range r.Header {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("%s: %#v\n", k, r.Header[k])
	}
	if r.Body != "" {
		fmt.Printf("BODY: %#v\n", r.Body)
	}
}

func (r *Event) LogPrint() {
	var keys []string
	for k := range r.Header {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		log.Printf("%s: %#v\n", k, r.Header[k])
	}
	if r.Body != "" {
		log.Printf("BODY: %#v\n", r.Body)
	}
	log.Println("-----------------------------", r.Get("Event-Name"), "-----------------------------")
}
