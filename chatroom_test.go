package chat

import (
	"bytes"
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
func TestChatroomGetHandler(t *testing.T) {
	Db2 = &mockChatroomDb{}
	result := getHandlerTester(handleChatroomGet, "/chatrooms")
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

func postHandlerTester(handler func(http.ResponseWriter, *http.Request), address string, data string) http.Response {
	req := httptest.NewRequest(http.MethodPost, address, bytes.NewBufferString(data))
	rec := httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	handler(rec, req)
	return *rec.Result()
}
func TestChatroomCreateHandler(t *testing.T) {
	Db2 := &mockChatroomDb{}
	result := postHandlerTester(handleChatroomCreate, "/chatrooms", `{"name":"test3","id":"3"}`)
	body, _ := ioutil.ReadAll(result.Body)
	fmt.Println(string(body))
	var chatrooms []Chatroom
	json.Unmarshal([]byte(body), &chatrooms)
	fmt.Println(chatrooms)
	if chatrooms[2].Id == 3 && chatrooms[2].Name != "test3" {
		t.Errorf("chatroom Id:3 should have the name test3 but instead had %s", chatrooms[1].Name)
	}
}
