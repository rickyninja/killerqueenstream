application=
{
	name="stresstest",
	description="Application for stressing a streaming server",
	protocol="dynamiclinklibrary",
	targetServer="localhost",
	targetApp="vod",
	active=false,
	--[[streams =
	{
		"lg00","lg01","lg02","lg03","lg04","lg05","lg06","lg07","lg08",
		"lg09","lg10","lg11","lg12","lg13","lg14","lg15","lg16","lg17",
		"lg18","lg19","lg20","lg21","lg22","lg23","lg24","lg25","lg26",
		"lg27","lg28","lg29","lg30","lg31","lg32","lg33","lg34","lg35",
		"lg36","lg37","lg38","lg39","lg40","lg41","lg42","lg43","lg44",
		"lg45","lg46","lg47","lg48","lg49"
	},]]--
	streams =
	{
		"mp4:lg.mp4"
	},
	numberOfConnections=10,
	randomAccessStreams=false
}
