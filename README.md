# What?

The Choria Go daemon is embeddable into other Go applications.

This is a demonstration that runs on a Raspberry Pi, it reads weather and humidity from a DHT220 sensor.

It then starts an embedded Choria server and register itself as a reigstration data provider.

It also starts an embedded agent that would be usable from the normal ruby `mco rpc` cli.

```
root@f35711d:/usr/src/app# DH2200_PIN=GPIO_4 ./choriapi
INFO[0000] Initial servers: []choria.Server{choria.Server{Host:"demo.nats.io", Port:4222, Scheme:"nats"}}  component=server identity=f35711d
INFO[0000] Connected to nats://demo.nats.io:4222         component=server identity=f35711d
INFO[0000] Registering new agent discovery of type discovery  component=server identity=f35711d subsystem=agents
INFO[0000] Subscribing agent discovery to mcollective.broadcast.agent.discovery  component=server identity=f35711d subsystem=agents
INFO[0000] Subscribing node f35711d to mcollective.node.f35711d  component=server identity=f35711d
INFO[0000] Registering new agent dht220 of type dht220   component=server identity=f35711d subsystem=agents
INFO[0000] Subscribing agent dht220 to mcollective.broadcast.agent.dht220  component=server identity=f35711d subsystem=agents
2017/12/04 22:30:18 Starting to send data every 60 seconds
2017/12/04 22:30:18 Publishing {"temperature":18.8,"humidy":55,"time":"2017-12-04T22:30:18.610248213Z"}
INFO[0000] Sending a broadcast message to NATS target 'mcollective.broadcast.agent.temperature' for message dca08402796f4ca78b578f4a7f5570d6 type request
```

Here the Choria security is disabled so no encryption or TLS, the data this sends can be seen here:

```
% nats-sub -s nats://demo.nats.io:4222 mcollective.broadcast.agent.temperature
[#1] Received on [mcollective.broadcast.agent.temperature] : '{"protocol":"choria:transport:1","data":"eyJwcm90b2NvbCI6ImNob3JpYTpzZWN1cmU6cmVxdWVzdDoxIiwibWVzc2FnZSI6IntcInByb3RvY29sXCI6XCJjaG9yaWE6cmVxdWVzdDoxXCIsXCJtZXNzYWdlXCI6XCJleUowWlcxd1pYSmhkSFZ5WlNJNk1UZ3VPQ3dpYUhWdGFXUjVJam8xTmk0MUxDSjBhVzFsSWpvaU1qQXhOeTB4TWkwd05GUXlNam96TmpvME1pNDJNelkyTVRrNU1UTmFJbjA9XCIsXCJlbnZlbG9wZVwiOntcInJlcXVlc3RpZFwiOlwiYmJiOTVjNzNlYWRiNDFjMzlmODJhN2VlZGZmZWNlZTFcIixcInNlbmRlcmlkXCI6XCJmMzU3MTFkXCIsXCJjYWxsZXJpZFwiOlwiY2hvcmlhPWYzNTcxMWRcIixcImNvbGxlY3RpdmVcIjpcIm1jb2xsZWN0aXZlXCIsXCJhZ2VudFwiOlwidGVtcGVyYXR1cmVcIixcInR0bFwiOjYwLFwidGltZVwiOjE1MTI0MjcwMDIsXCJmaWx0ZXJcIjp7XCJmYWN0XCI6W10sXCJjZl9jbGFzc1wiOltdLFwiYWdlbnRcIjpbXSxcImlkZW50aXR5XCI6W10sXCJjb21wb3VuZFwiOltdfX19Iiwic2lnbmF0dXJlIjoiaW5zZWN1cmUiLCJwdWJjZXJ0IjoiaW5zZWN1cmUifQ==","headers":{"reply-to":"dev.null","mc_sender":"f35711d","seen-by":[["nats://demo.nats.io:4222","f35711d","nats://demo.nats.io:4222"]]}}'
```

The data is Base64 encoded and contains the onion layers leading down to the real data:

From the Transport get the Secure Request:
```
% cat|base64 -d
eyJwcm90b2NvbCI6ImNob3Jp.....LCJwdWJjZXJ0IjoiaW5zZWN1cmUifQ==
{"protocol":"choria:secure:request:1","message":"{\"protocol\":\"choria:request:1\",\"message\":\"eyJ0ZW1wZXJhdHVyZSI6MTguOCwiaHVtaWR5Ijo1Ni41LCJ0aW1lIjoiMjAxNy0xMi0wNFQyMjozNjo0Mi42MzY2MTk5MTNaIn0=\",\"envelope\":{\"requestid\":\"bbb95c73eadb41c39f82a7eedffecee1\",\"senderid\":\"f35711d\",\"callerid\":\"choria=f35711d\",\"collective\":\"mcollective\",\"agent\":\"temperature\",\"ttl\":60,\"time\":1512427002,\"filter\":{\"fact\":[],\"cf_class\":[],\"agent\":[],\"identity\":[],\"compound\":[]}}}","signature":"insecure","pubcert":"insecure"}
```

Get the data in the Request:

```
% cat|base64 -d
eyJ0ZW1wZXJhdHVyZSI6MTguOCwiaHVtaWR5Ijo1Ni41LCJ0aW1lIjoiMjAxNy0xMi0wNFQyMjozNjo0Mi42MzY2MTk5MTNaIn0=
{"temperature":18.8,"humidy":56.5,"time":"2017-12-04T22:36:42.636619913Z"}
```

In includes a DDL file for the ruby choria, if installed and configured you can do:

```
[rip@dev1]% mco rpc dht220 reading --config .mcollective.choriapi
Discovering hosts using the mc method for 2 second(s) .... 1

 * [ ============================================================> ] 1 / 1


f35711d
        humidy: 54.6
   Temperature: 18.5
          Time: 2017-12-05T08:11:59.623588181Z


Finished processing 1 / 1 hosts in 696.85 ms
```