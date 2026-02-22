package main

import (
	"reflect"
	"testing"
)

func TestSendMessageReqStructTags(t *testing.T) {
	typ := reflect.TypeOf(SendMessageReq{})
	f, ok := typ.FieldByName("Message")
	if !ok {
		t.Fatalf("Message field not found")
	}
	if f.Tag.Get("json") != "message" {
		t.Fatalf("unexpected json tag: %q", f.Tag.Get("json"))
	}
	if f.Tag.Get("binding") != "required" {
		t.Fatalf("unexpected binding tag: %q", f.Tag.Get("binding"))
	}
}
