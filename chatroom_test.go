package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockChatroomDb struct{}

func (m *mockChatroomDb) Chatrooms() []Chatroom {
	return []Chatroom{{Id: 1, Name: "test"}, {Id: 2, Name: "test2"}}
}
func TestChatroomSelectHandler(t *testing.T) {
	Db2 = &mockChatroomDb{}
	result := getterHandlerTester(handleChatroomSelect, "/chatrooms")
	body, _ := ioutil.ReadAll(result.Body)
	fmt.Println(string(body))
	var chatrooms []Chatroom
	json.Unmarshal([]byte(body), &chatrooms)
	fmt.Println(chatrooms)
	if chatrooms[1].Id == 2 && chatrooms[1].Name != "test2" {
		t.Errorf("chatroom Id:2 should have the name test2 but instead had %s", chatrooms[1].Name)
	}
}

func getHandlerTester(handler func(http.ResponseWriter, *http.Request), address string) http.Response {
	req := httptest.NewRequest(http.MethodGet, address, nil)
	rec := httptest.NewRecorder()
	handler(rec, req)

	return *rec.Result()
}
func postHandlerTester(handler func(http.ResponseWriter, *http.Request), address string) http.Response {

	req := httptest.NewRequest(http.MethodGet, address, nil)
	rec := httptest.NewRecorder()
	handler(rec, req)
	return *rec.Result()
}
func TestChatroomCreateHandler(t *testing.T) {
	Db2 := &mockChatroomDb{}
}
