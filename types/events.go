// Copyright (c) 2024 Tulir Asokan
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package types contains various types used by whatsmeow.
package types

import (
	"time"
)

// MessageSource contains the basic sender and chat information for any incoming message.
type MessageSource struct {
	// The chat where the message was sent.
	Chat JID
	// The sender of the message. This is the same as Chat for private chats.
	Sender JID
	// Whether the message was sent by the current user instead of someone else.
	IsFromMe bool
	// Whether the message is in a group chat.
	IsGroup bool
}

// Message is the event dispatched when a new WhatsApp message is received.
type Message struct {
	Info    MessageInfo
	Message interface{}
	IsEphemeral bool
	IsViewOnce  bool
}

// MessageInfo contains metadata about an incoming message.
type MessageInfo struct {
	MessageSource
	// The unique message ID.
	ID string
	// The server timestamp of the message.
	Timestamp time.Time
	// The push name (display name) of the sender.
	PushName string
	// Whether the message has been broadcast (status update).
	Broadcast bool
	// The JID of the device that sent the message (for multi-device).
	SenderDevice int
}

// Receipt is the event dispatched when a message receipt (read/delivered) is received.
type Receipt struct {
	MessageSource
	// The IDs of the messages that were read/delivered.
	MessageIDs []string
	// The type of receipt (read, delivered, etc.).
	Type ReceiptType
	// The timestamp of the receipt.
	Timestamp time.Time
}

// ReceiptType is the type of a message receipt.
type ReceiptType string

const (
	// ReceiptTypeDelivered means the message was delivered to the device.
	ReceiptTypeDelivered ReceiptType = "delivered"
	// ReceiptTypeRead means the message was read by the user.
	ReceiptTypeRead ReceiptType = "read"
	// ReceiptTypeReadSelf means the current user read the message on another device.
	ReceiptTypeReadSelf ReceiptType = "read-self"
	// ReceiptTypePlayed means a voice message was played.
	ReceiptTypePlayed ReceiptType = "played"
)

// Presence is the event dispatched when a contact's presence (online/offline) changes.
type Presence struct {
	// The JID of the contact.
	From JID
	// Whether the contact is currently available.
	Unavailable bool
	// The last time the contact was seen (only set if Unavailable is true).
	LastSeen time.Time
}

// Connected is the event dispatched when the client successfully connects to WhatsApp.
type Connected struct{}

// Disconnected is the event dispatched when the client disconnects from WhatsApp.
type Disconnected struct {
	// Whether the disconnection was intentional (e.g., called Disconnect()).
	LoggedOut bool
}

// QR is the event dispatched when a QR code is available for scanning.
type QR struct {
	// The QR code codes to display. Multiple codes may be provided as they rotate.
	Codes []string
}

// PairSuccess is the event dispatched when QR pairing succeeds.
type PairSuccess struct {
	// The JID assigned to this device.
	ID JID
	// The business name or push name of the account.
	BusinessName string
	// The platform of the device.
	Platform string
}

// LoggedOut is the event dispatched when the client is logged out by the server.
type LoggedOut struct {
	// Whether the logout was initiated on another device.
	OnConnect bool
	// The reason for the logout.
	Reason ConnectFailureReason
}

// ConnectFailureReason is the reason for a connection failure.
type ConnectFailureReason int

const (
	ConnectFailureLoggedOut      ConnectFailureReason = 401
	ConnectFailureMainDeviceGone ConnectFailureReason = 442
	ConnectFailureUnknownLogout  ConnectFailureReason = 440
)
