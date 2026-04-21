// Copyright (c) 2024 Tulir Asokan
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package whatsmeow implements a WhatsApp web multidevice client.
package whatsmeow

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

"
)

// EventHandler is a function that can handle events from WhatsApp.
type EventHandler func(evt interface{})

// Client is the main WhatsApp client struct.
type Client struct {
	Store   *store.Device
	Log     log.Logger
	RecipientDeviceCache     map[types.JID][]types.JID
	recipientDeviceCacheLock sync.Mutex

	// Event handlers registered via AddEventHandler
	eventHandlers     []wrappedEventHandler
	eventHandlersLock sync.RWMutex

	// Connection state
	isConnected atomic.Bool
	connectLock sync.Mutex

	// Context for managing goroutine lifecycles
	ctx    context.Context
	cancel context.CancelFunc

	// Unique handler ID counter
	lastHandlerID uint32
}

type wrappedEventHandler struct {
	fn EventHandler
	id uint32
}

// NewClient creates a new WhatsApp client with the given device store and logger.
func NewClient(deviceStore *store.Device, log log.Logger) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		Store:                deviceStore,
		Log:                  log,
		RecipientDeviceCache: make(map[types.JID][]types.JID),
		ctx:                  ctx,
		cancel:               cancel,
	}
}

// AddEventHandler registers a new event handler function and returns a unique ID
// that can be used to remove it later via RemoveEventHandler.
func (cli *Client) AddEventHandler(handler EventHandler) uint32 {
	cli.eventHandlersLock.Lock()
	defer cli.eventHandlersLock.Unlock()
	id := atomic.AddUint32(&cli.lastHandlerID, 1)
	cli.eventHandlers = append(cli.eventHandlers, wrappedEventHandler{fn: handler, id: id})
	return id
}

// RemoveEventHandler removes a previously registered event handler by its ID.
// Returns true if the handler was found and removed.
func (cli *Client) RemoveEventHandler(id uint32) bool {
	cli.eventHandlersLock.Lock()
	defer cli.eventHandlersLock.Unlock()
	for i, h := range cli.eventHandlers {
		if h.id == id {
			cli.eventHandlers = append(cli.eventHandlers[:i], cli.eventHandlers[i+1:]...)
			return true
		}
	}
	return false
}

// RemoveAllEventHandlers removes all registered event handlers.
func (cli *Client) RemoveAllEventHandlers() {
	cli.eventHandlersLock.Lock()
	defer cli.eventHandlersLock.Unlock()
	cli.eventHandlers = nil
}

// dispatchEvent sends the given event to all registered handlers.
func (cli *Client) dispatchEvent(evt interface{}) {
	cli.eventHandlersLock.RLock()
	handlers := make([]wrappedEventHandler, len(cli.eventHandlers))
	copy(handlers, cli.eventHandlers)
	cli.eventHandlersLock.RUnlock()
	for _, h := range handlers {
		h.fn(evt)
	}
}

// IsConnected returns true if the client currently has an active connection to WhatsApp.
func (cli *Client) IsConnected() bool {
	return cli.isConnected.Load()
}

// IsLoggedIn returns true if the client has valid credentials stored.
func (cli *Client) IsLoggedIn() bool {
	return cli.Store != nil && cli.Store.ID != nil
}

// Disconnect closes the active connection to WhatsApp.
func (cli *Client) Disconnect() {
	if !cli.isConnected.Load() {
		return
	}
	cli.Log.Infof("Disconnecting from WhatsApp")
	cli.cancel()
	cli.isConnected.Store(false)
	cli.dispatchEvent(&events.Disconnected{})
}

// WaitForConnection blocks until the client is connected or the timeout is reached.
func (cli *Client) WaitForConnection(timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if cli.IsConnected() {
			return true
		}
		time.Sleep(100 * time.Millisecond)
	}
	return cli.IsConnected()
}
