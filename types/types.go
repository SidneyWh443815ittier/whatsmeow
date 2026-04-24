// Copyright (c) 2024 Tulir Asokan
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package types contains various types used by the whatsmeow library.
package types

import (
	"fmt"
	"strings"
	"time"
)

// JID represents a WhatsApp JID (Jabber ID), which is used to identify
// users, groups, and other entities on the WhatsApp network.
type JID struct {
	User   string
	Agent  uint8
	Device uint8
	Server string
	AD     bool
}

const (
	// DefaultUserServer is the standard WhatsApp user server.
	DefaultUserServer = "s.whatsapp.net"
	// GroupServer is the server used for group JIDs.
	GroupServer = "g.us"
	// LegacyUserServer is the legacy WhatsApp user server.
	LegacyUserServer = "c.us"
	// BroadcastServer is the server used for broadcast lists.
	BroadcastServer = "broadcast"
	// HiddenUserServer is used for certain hidden/system accounts.
	HiddenUserServer = "lid"
)

// NewJID creates a new JID with the given user and server.
func NewJID(user, server string) JID {
	return JID{User: user, Server: server}
}

// IsEmpty returns true if the JID has no user or server set.
func (jid JID) IsEmpty() bool {
	return jid.User == "" && jid.Server == ""
}

// IsGroup returns true if the JID refers to a group chat.
func (jid JID) IsGroup() bool {
	return jid.Server == GroupServer
}

// IsBroadcast returns true if the JID refers to a broadcast list.
func (jid JID) IsBroadcast() bool {
	return jid.Server == BroadcastServer
}

// ToNonAD returns a copy of the JID without the agent/device fields.
func (jid JID) ToNonAD() JID {
	if jid.AD {
		return JID{
			User:   jid.User,
			Server: DefaultUserServer,
		}
	}
	return jid
}

// String returns the string representation of the JID.
func (jid JID) String() string {
	if jid.AD {
		return fmt.Sprintf("%s.%d:%d@%s", jid.User, jid.Agent, jid.Device, jid.Server)
	}
	if jid.User == "" {
		return "@" + jid.Server
	}
	return jid.User + "@" + jid.Server
}

// ParseJID parses a JID string into a JID struct.
func ParseJID(jid string) (JID, error) {
	if len(jid) == 0 {
		return JID{}, fmt.Errorf("empty JID")
	}
	parts := strings.SplitN(jid, "@", 2)
	if len(parts) == 1 {
		return NewJID("", parts[0]), nil
	}
	return NewJID(parts[0], parts[1]), nil
}

// MessageID is the unique identifier for a WhatsApp message.
type MessageID = string

// MessageInfo contains metadata about a received or sent message.
type MessageInfo struct {
	// ID is the unique message identifier.
	ID MessageID
	// MessageSource contains info about who sent the message and where.
	MessageSource
	// Timestamp is when the message was sent.
	Timestamp time.Time
	// PushName is the push name of the sender at the time of the message.
	PushName string
	// Broadcast indicates if the message was sent via a broadcast list.
	Broadcast bool
}

// MessageSource contains the source/routing information for a message.
type MessageSource struct {
	// Chat is the JID of the chat where the message was sent.
	Chat JID
	// Sender is the JID of the user who sent the message.
	Sender JID
	// IsFromMe indicates if the message was sent by the current user.
	IsFromMe bool
	// IsGroup indicates if the message was sent in a group chat.
	IsGroup bool
}

// VerifiedName contains verified business name information.
type VerifiedName struct {
	// Certificate is the raw verified name certificate.
	Certificate []byte
	// Details contains the parsed details from the certificate.
	Details *VerifiedNameDetails
}

// VerifiedNameDetails contains the parsed details from a verified name certificate.
type VerifiedNameDetails struct {
	Serial  uint64
	Issuer  string
	Name    string
	LocalizedNames []LocalizedName
}

// LocalizedName represents a localized version of a verified business name.
type LocalizedName struct {
	Lg   string
	Lc   string
	Name string
}
