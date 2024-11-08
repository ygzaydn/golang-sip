# Example SIP Request

An example SIP Request
```md
INVITE sip:bob@example.com SIP/2.0
Via: SIP/2.0/UDP alicepc.example.com;branch=z9hG4bK776asdhds
Max-Forwards: 70
To: Bob <sip:bob@example.com>
From: Alice <sip:alice@example.com>;tag=1928301774
Call-ID: a84b4c76e66710@alicepc.example.com
CSeq: 314159 INVITE
Contact: <sip:alice@alicepc.example.com>
Content-Type: application/sdp
Content-Length: 142

v=0
o=Alice 2890844526 2890844526 IN IP4 alicepc.example.com
s=Session SDP
c=IN IP4 192.0.2.101
t=0 0
m=audio 49172 RTP/AVP 0
a=rtpmap:0 PCMU/8000
```

# Example SIP Response

An example SIP Response

```md
SIP/2.0 200 OK
Via: SIP/2.0/UDP alicepc.example.com;branch=z9hG4bK776asdhds
To: Bob <sip:bob@example.com>;tag=456248
From: Alice <sip:alice@example.com>;tag=1928301774
Call-ID: a84b4c76e66710@alicepc.example.com
CSeq: 314159 INVITE
Contact: <sip:bob@bobphone.example.com>
Content-Type: application/sdp
Content-Length: 147

v=0
o=Bob 2890844526 2890844526 IN IP4 bobphone.example.com
s=Session SDP
c=IN IP4 192.0.2.201
t=0 0
m=audio 49170 RTP/AVP 0
a=rtpmap:0 PCMU/8000
```

To distinguish SIP Request and Response

1- First Line Check
2- Header Analysis
3- Reqular Expressions

-   CSeq Header Check is important