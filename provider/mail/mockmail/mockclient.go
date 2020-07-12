package mockmail

import (
	"context"
	"sync"
	"testing"
)

type CallbackFunc func(T *testing.T, ctx context.Context, to []string, from, subject, message string)

type MockClient struct {
	T *testing.T
	Context context.Context
	CallbackFunc CallbackFunc
	Group *sync.WaitGroup
}

func New(client *MockClient) *MockClient {
	m := client
	m.Group.Add(1)
	return m
}

func (m *MockClient) SendMail(to []string, from, subject, message string) error {
	defer m.Group.Done()
	m.CallbackFunc(m.T, m.Context, to, from, subject, message)

	return nil
}