application=
{
	name="proxypublish",
	description="Application for forwarding streams to another RTMP server",
	protocol="dynamiclinklibrary",
	acceptors =
	{
		{
			ip="0.0.0.0",
			port=6665,
			protocol="inboundLiveFlv"
		},
	},
	abortOnConnectError=true,
	targetServers = 
	{
		--[[{
			targetUri="rtmp://x.xxxxxxx.fme.ustream.tv/ustreamVideo/xxxxxxx",
			targetStreamName="xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			localStreamName="gigi",
			emulateUserAgent="FMLE/3.0 (compatible; FMSc/1.0 http://www.rtmpd.com)"
		}]]--,
		{
			targetUri="rtmp://gigi:spaima@localhost/vod",
			targetStreamType="live", -- (live, record or append)
			emulateUserAgent="My user agent",
			localStreamName="stream1",
			keepAlive=true
		},
	},
	--[[externalStreams =
	{
		{
			uri="rtsp://fms20.mediadirect.ro/live2/realitatea/realitatea",
			localStreamName="stream1",
			forceTcp=true,
			keepAlive=true
		},
	},]]--
	--validateHandshake=true,
	--default=true,
}
