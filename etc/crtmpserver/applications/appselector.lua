application =
{
	-- The name of the application. It is mandatory and must be unique
	name = "appselector",
	-- Short description of the application. Optional
	description = "Application for selecting the rest of the applications",
	-- The type of the application. Possible values are:
	-- dynamiclinklibrary - the application is a shared library
	protocol = "dynamiclinklibrary",
	-- the complete path to the library. This is optional. If not provided,
	-- the server will try to load the library from here
	-- <rootDirectory>/<name>/lib<name>.{so|dll|dylib}
	-- library="/some/path/to/some/shared/library.so"

	-- Tells the server to validate the clien's handshake before going further.
	-- It is optional, defaulted to true
	validateHandshake = true,
	-- this is the folder from where the current application gets it's content.
	-- It is optional. If not specified, it will be defaulted to:
	-- <rootDirectory>/<name>/mediaFolder
	-- mediaFolder="/some/directory/where/media/files/are/stored"
	-- the application will also be known by that names. It is optional
	--aliases=
	--{
	--	"simpleLive",
	--	"vod",
	--	"live",
	--},
	-- This flag designates the default application. The default application
	-- is responsable of analyzing the "connect" request and distribute
	-- the future connection to the correct application.
	default = true,
	acceptors =
	{
		{
			ip = "0.0.0.0",
			port = 1935,
			protocol = "inboundRtmp"
		},
		{
			ip = "0.0.0.0",
			port = 8080,
			protocol = "inboundRtmpt"
		},
		--[[{
			ip = "0.0.0.0",
			port = 8081,
			protocol = "inboundRtmps",
			sslKey = "server.key",
			sslCert = "server.crt"
		},]]--
	}
}
