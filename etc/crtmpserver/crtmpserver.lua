

confdir = "/etc/crtmpserver"
libdir = "/usr/lib/crtmpserver"

enabled_apps = confdir.."/enabled_applications.conf"
app_lib_path = libdir.."/applications"
logger_config_path = confdir.."/log_appenders.lua"

-- Print anything - including nested tables
function table_print (tt, indent, done)
	done = done or {}
	indent = indent or 0
	if type(tt) == "table" then
		for key, value in pairs (tt) do
			io.write(string.rep (" ", indent)) -- indent it
			if type (value) == "table" and not done [value] then
				done [value] = true
				io.write(string.format("[%s] => table\n", tostring (key)));
				io.write(string.rep (" ", indent+4)) -- indent it
				io.write("(\n");
				table_print (value, indent + 7, done)
				io.write(string.rep (" ", indent+4)) -- indent it
				io.write(")\n");
			else
				io.write(string.format("[%s] => %s\n",
				tostring (key), tostring(value)))
			end
		end
	else
		io.write(tt .. "\n")
	end
end

function exists(fname)
	local f = io.open(fname, "r")
	if (f) then
		return true
	else
		return false
	end
end

-- Function generate "logAppenders" section for crtmpserver
function read_logappenders()
	result = {}
	dofile(logger_config_path)
	result = logAppenders
	return result
end

-- Function generate "applications" section for crtmpserver
-- Reads apps Lua script to main configuration section
function read_applications()
	result = {}
	local app_config
	-- Must specify whre application libs can be found
	result.rootDirectory = app_lib_path

	-- Loads applications configuration listed in file "enabled_apps"
	for app in io.lines(enabled_apps) do
		app = (app:gsub("^%s*(.-)%s*$", "%1"))
		application = nil
		if ( app:match("^#.*$") or app:match("^$") or app:match("^\s+$") ) then
			--print("string '"..app.."' is unneeded, skip it")
		else
			app_config = confdir.."/applications/"..app..".lua"
			if (not exists(app_config)) then
				print("Application configuration file '"..app_config.."' not found")
				os.exit()
			end
			dofile(app_config)
			if (application == nil) then
				print("Configuration file '"..app_config.."' does not contain variable 'application'")
				os.exit()
			else
				table.insert(result, application)
			end
		end
	end
	return result
end

-- Check if logger configuration exists
if (not exists(logger_config_path)) then
	print("Logger configuration file '"..logger_config_path.."' not found");
	return
end

-- Check if list of applications exists
if (not exists(enabled_apps)) then
	print("Applications list file '"..enabled_apps.."' not found");
	return
end

-- Main section of configuration.
-- This variable reads by crtmpserver as main configuration section
-- It must be always defined
configuration =
{
	daemon = true,
	pathSeparator = "/",
	logAppenders = read_logappenders(),
	applications = read_applications()
}
-- print("__________________________")
-- table_print(configuration)
-- print("__________________________")
