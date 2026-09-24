package p2p

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pion/webrtc/v4"
)

func F_p2p() {
	fmt.Println("Initializing STUN/TURN configuration...")

	startTime := time.Now()

	// 1. Create a context to capture OS signals (Ctrl+C, SIGTERM)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop() // Disable signal notification when the function completes

	// 2. Configure STUN and TURN Servers
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
			{
				URLs:           []string{"turn:relay.expressturn.com:3478"},
				Username:       "000000002105517497",
				Credential:     "JOtrJ/AaPmxxIswWVo3hPPCPFMg=",
				CredentialType: webrtc.ICECredentialTypePassword,
			},
		},
	}

	// 3. Create the PeerConnection
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}

	// Graceful Cleanup: Ensure the connection is properly closed upon exit
	defer func() {
		fmt.Println("\nClosing PeerConnection safely...")
		if err := peerConnection.Close(); err != nil {
			fmt.Printf("Error closing PeerConnection: %v\n", err)
		} else {
			fmt.Println("PeerConnection closed successfully.")
		}
	}()

	// 4. ICE Candidates Event Handler
	peerConnection.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c != nil {
			elapsed := time.Since(startTime)
			fmt.Printf("[%s] New ICE candidate discovered: %s\n", elapsed.Round(time.Millisecond), c.String())
		} else {
			totalTime := time.Since(startTime)
			fmt.Printf("\nAll ICE candidates gathered successfully. Total elapsed time: %s\n", totalTime.Round(time.Millisecond))
		}
	})

	// 5. Create a Data Channel
	_, err = peerConnection.CreateDataChannel("chat", nil)
	if err != nil {
		panic(err)
	}

	// 6. Create Offer and set Local Description
	fmt.Println("Starting ICE gathering process (creating local offer)...")
	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		panic(err)
	}

	err = peerConnection.SetLocalDescription(offer)
	if err != nil {
		panic(err)
	}

	// 7. Block and wait for an OS signal instead of an empty select {}
	fmt.Println("Application is running. Press Ctrl+C to stop safely...")

	<-ctx.Done() // Waits here until Ctrl+C (SIGINT) or SIGTERM is received

	fmt.Println("\nShutdown signal received! Cleaning up resources...")
}