-- Configuration for LogAppender subsystem of crtmpserver

logAppenders=
{
	{
		-- name of the appender. Not too important, but is mandatory
		name="console appender",
		-- type of the appender. We can have the following values:
		-- console, coloredConsole and file
		-- NOTE: console appenders will be ignored if we run the server
		-- as a daemon
		type="coloredConsole",
		-- the level of logging. 6 is the FINEST message, 0 is FATAL message.
		-- The appender will "catch" all the messages below or equal to this level
		-- bigger the level, more messages are recorded
		level=6
	},
	{
		name="file appender",
		type="file",
		level=6,
		-- the file where the log messages are going to land
		fileName="/var/log/crtmpserver/main.log",
                --newLineCharacters="\r\n",
                --fileHistorySize=10,
                --fileLength=1024*256,
                --singleLine=true
	}
}

