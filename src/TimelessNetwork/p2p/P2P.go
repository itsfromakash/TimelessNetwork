package main

import "C"

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pion/webrtc/v4"
)

//export StartP2PEngine
func StartP2PEngine() {
	fmt.Println("[Go P2P Engine] Initializing STUN/TURN configuration...")

	startTime := time.Now()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

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

	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		fmt.Printf("[Go Error] %v\n", err)
		return
	}

	defer func() {
		fmt.Println("\n[Go P2P Engine] Closing PeerConnection safely...")
		if err := peerConnection.Close(); err != nil {
			fmt.Printf("[Go Error] Closing PeerConnection: %v\n", err)
		} else {
			fmt.Println("[Go P2P Engine] PeerConnection closed successfully.")
		}
	}()

	peerConnection.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c != nil {
			elapsed := time.Since(startTime)
			fmt.Printf("[Go P2P Engine] [%s] New ICE candidate discovered: %s\n", elapsed.Round(time.Millisecond), c.String())
		} else {
			totalTime := time.Since(startTime)
			fmt.Printf("\n[Go P2P Engine] All ICE candidates gathered successfully. Total time: %s\n", totalTime.Round(time.Millisecond))
		}
	})

	_, err = peerConnection.CreateDataChannel("chat", nil)
	if err != nil {
		fmt.Printf("[Go Error] %v\n", err)
		return
	}

	fmt.Println("[Go P2P Engine] Starting ICE gathering process...")
	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		fmt.Printf("[Go Error] %v\n", err)
		return
	}

	err = peerConnection.SetLocalDescription(offer)
	if err != nil {
		fmt.Printf("[Go Error] %v\n", err)
		return
	}

	fmt.Println("[Go P2P Engine] Engine is running in background. Waiting for OS signals...")
	<-ctx.Done()
	fmt.Println("\n[Go P2P Engine] Shutdown signal received!")
}

func main() {}