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

	// 1. OS Signals (Ctrl+C, SIGTERM) අල්ලා ගැනීමට Context එකක් සාදා ගැනීම
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop() // Function එක අවසන් වන විට signal notification අක්‍රිය කරයි

	// 2. STUN සහ TURN Servers Configuration
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

	// 3. Peer Connection එක සෑදීම
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}

	// Graceful Cleanup: Program එක අවසන් වන විට Connection එක නිසි පරිදි ක්ලෝස් වේ
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

	// 5. Data Channel එක නිර්මාණය
	_, err = peerConnection.CreateDataChannel("chat", nil)
	if err != nil {
		panic(err)
	}

	// 6. Offer එක සාදා Local Description එක set කිරීම
	fmt.Println("Starting ICE gathering process (creating local offer)...")
	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		panic(err)
	}

	err = peerConnection.SetLocalDescription(offer)
	if err != nil {
		panic(err)
	}

	// 7. select {} වෙනුවට OS Signal එකක් එනතෙක් Block වී සිටීම
	fmt.Println("Application is running. Press Ctrl+C to stop safely...")

	<-ctx.Done() // Ctrl+C (SIGINT) හෝ SIGTERM ලැබෙන තෙක් මෙතන නතර වී පවතී

	fmt.Println("\nShutdown signal received! Cleaning up resources...")
}
