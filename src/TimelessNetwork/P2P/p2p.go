package P2P

import "C" // Enables Cgo so external languages (C, C++, Swift) can call this Go code

import (
	"context"      // Controls program lifecycle and cancellations
	"fmt"          // Used for printing output to the terminal
	"os"           // Interfaces with the Operating System
	"os/signal"    // Listens for system signals like Ctrl+C
	"syscall"      // Low-level system signal definitions
	"time"         // Handles time measurement and timestamps
	"github.com/pion/webrtc/v4" // WebRTC library for P2P networking
)

//export StartP2PEngine
// Exposes this function to external languages when compiled as a library
func StartP2PEngine() {
	fmt.Println("[P2P Engine] Initializing STUN/TURN configuration...")

	startTime := time.Now() // Save current time to measure connection setup duration

	// Create a context that automatically cancels when the user presses Ctrl+C or terminates the app
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop() // Clean up signal listeners when this function ends

	// Configure STUN (finds public IP) and TURN (relays data if firewalls block P2P)
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"}, // Free public STUN server
			},
			{
				URLs:           []string{"turn:relay.expressturn.com:3478"}, // TURN server credentials
				Username:       "000000002105517497",
				Credential:     "JOtrJ/AaPmxxIswWVo3hPPCPFMg=",
				CredentialType: webrtc.ICECredentialTypePassword,
			},
		},
	}

	// Initialize the WebRTC peer connection object
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		fmt.Printf("[Error] %v\n", err)
		return
	}

	// Guarantee that peerConnection is closed safely when this function finishes
	defer func() {
		fmt.Println("\n[P2P Engine] Closing PeerConnection safely...")
		if err := peerConnection.Close(); err != nil {
			fmt.Printf("[Error] Closing PeerConnection: %v\n", err)
		} else {
			fmt.Println("[P2P Engine] PeerConnection closed successfully.")
		}
	}()

	// Event listener: Triggers every time a new network route (ICE candidate) is discovered
	peerConnection.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c != nil {
			elapsed := time.Since(startTime)
			fmt.Printf("[P2P Engine] [%s] New ICE candidate discovered: %s\n", elapsed.Round(time.Millisecond), c.String())
		} else {
			// A nil candidate signals that candidate gathering is complete
			totalTime := time.Since(startTime)
			fmt.Printf("\n[P2P Engine] All ICE candidates gathered successfully. Total time: %s\n", totalTime.Round(time.Millisecond))
		}
	})

	// Open a WebRTC DataChannel named "chat" for bidirectional messaging
	_, err = peerConnection.CreateDataChannel("chat", nil)
	if err != nil {
		fmt.Printf("[Error] %v\n", err)
		return
	}

	fmt.Println("[P2P Engine] Starting ICE gathering process...")
	
	// Create the SDP offer (contains connection details and encryption parameters)
	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		fmt.Printf("[Error] %v\n", err)
		return
	}

	// Apply local description to start background network path gathering
	err = peerConnection.SetLocalDescription(offer)
	if err != nil {
		fmt.Printf("[Error] %v\n", err)
		return
	}

	// Pause and keep the background tasks alive until a shutdown signal is received
	fmt.Println("[P2P Engine] Engine is running in background. Waiting for OS signals...")
	<-ctx.Done()
	fmt.Println("\n[P2P Engine] Shutdown signal received!")
}
