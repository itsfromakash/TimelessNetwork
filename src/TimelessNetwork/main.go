package main

import (
	"fmt"
	"github.com/pion/webrtc/v4"
)

func main() {
	// 1. STUN Server එක Configure කිරීම
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				// Google හි නොමිලේ දෙන STUN Server එකක්
				URLs: []string{"stun:stun.l.google.com:19302"}, 
			},
		},
	}

	// 2. Peer Connection එකක් නිර්මාණය කිරීම
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}
	defer peerConnection.Close()

	// 3. ICE Candidate කෙනෙක් (Public IP/Port) ලැබුණු විට ක්‍රියාත්මක වන Event එක
	peerConnection.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c != nil {
			// මෙතැනදී ඔබට ඔබේ Public IP සහ Port එක බලාගත හැක
			fmt.Printf("නව ICE Candidate (Public IP) හමුවුණා: %s\n", c.String())
		}
	})

	// 4. Data Channel එකක් සෑදීම (පණිවිඩ හුවමාරු කර ගැනීමට)
	dataChannel, err := peerConnection.CreateDataChannel("chat", nil)
	if err != nil {
		panic(err)
	}

	dataChannel.OnOpen(func() {
		fmt.Println("P2P සම්බන්ධතාවය සාර්ථකව ආරම්භ වුණා!")
		dataChannel.SendText("ဟලෝ (Hello) අනෙක් පස සිටින යාළුවා!")
	})

	// (වැදගත්) මෙතැන් සිට ඉදිරියට යාමට SDP Offer/Answer එකක් සාදා Signaling Server එකක් හරහා අනෙක් පරිගණකය සමඟ හුවමාරු කරගත යුතුය.
	select {} // වැඩසටහන දිගටම Run වී තිබීමට
}
