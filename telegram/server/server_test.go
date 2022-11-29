package server

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	telegram "inoutlog/telegram/client"
)

type mockTelegramClient struct {
	t testing.T
}

func (m *mockTelegramClient) SendMessage(id string, text string) error {
	m.t.Logf("test log: SendMessage(%v,%v)", id, text)
	return nil
}

func (m *mockTelegramClient) GetUpdates(timeout int) ([]telegram.Update, error) {
	m.t.Fatal("GetUpdates called when shouldn't")
	return nil, nil
}

func TestServerGetDate(t *testing.T) {
	s := Server{
		TelegramClient: &mockTelegramClient{},
	}

	zeroTs := time.Date(1988, time.October, 7, 0, 0, 0, 0, time.UTC)
	ts := zeroTs.Add(17*time.Hour + 30*time.Minute + 12*time.Second + 123*time.Nanosecond)
	tests := []struct {
		text      string
		expected  time.Time
		shouldErr bool
	}{
		{text: "/in 13:40", expected: zeroTs.Add(13*time.Hour + 40*time.Minute)},
		{text: "/in 14:20", expected: zeroTs.Add(14*time.Hour + 20*time.Minute)},
		{text: "/in", expected: ts},
		{text: "", expected: ts},
		{text: "/in 25:20", shouldErr: true},
		{text: "/in 4:20 but late", shouldErr: true},
	}
	for _, tt := range tests {
		got, err := s.getDate(tt.text, "6")
		if tt.shouldErr {
			require.Error(t, err, "text :%s should err", tt.text)
			continue
		}
		require.True(t, got.Equal(tt.expected), "got: %s, expected: %s", got.String(), tt.expected.String())
	}
}
