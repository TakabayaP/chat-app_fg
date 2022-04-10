package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndexHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	index(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("a%v", err)
	}
	if string(data) != "chat server running!" {
		t.Errorf("expected chat server running! got %v", string(data))
	}
}

type mockChatroomDb struct{}

func (m *mockChatroomDb) Chatrooms() []Chatroom {
	return []Chatroom{{Id: 1, Name: "test"}, {Id: 2, Name: "test2"}}
}
func TestChatroomHandler(t *testing.T) {
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

func getterHandlerTester(handler func(http.ResponseWriter, *http.Request), address string) http.Response {
	req := httptest.NewRequest(http.MethodGet, address, nil)
	rec := httptest.NewRecorder()
	handler(rec, req)

	return *rec.Result()
}
