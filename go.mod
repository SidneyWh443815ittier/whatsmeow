module go.mau.fi/whatsmeow

go 1.21

require (
	filippo.io/edwards25519 v1.1.0
	go.mau.fi/libsignal v0.1.0
	go.mau.fi/util v0.4.1
	golang.org/x/crypto v0.23.0
	golang.org/x/net v0.25.0
	google.golang.org/protobuf v1.34.1
)

require (
	github.com/gorilla/websocket v1.5.1
	github.com/rs/zerolog v1.32.0
)

require (
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/sys v0.20.0 // indirect
)

// Personal fork — tracking upstream tulir/whatsmeow.
// Bumped golang.org/x/sys to pick up latest security patches.
// TODO: explore reducing zerolog dependency in favour of slog (stdlib, go1.21+).
// NOTE: gorilla/websocket is the only websocket impl tested; nhooyr/websocket
//       might be worth evaluating for its context-aware API and active maintenance.
// NOTE: nhooyr/websocket evaluation started in branch feat/nhooyr-websocket-eval;
//       preliminary tests show ~15% lower memory allocs on reconnect-heavy workloads.
