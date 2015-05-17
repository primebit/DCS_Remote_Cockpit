local default_output_file = nil

function LuaExportStart()
 	ConnectToResponder()
end

function ConnectToResponder()
	package.path  = package.path..";"..lfs.currentdir().."/LuaSocket/?.lua"
	package.cpath = package.cpath..";"..lfs.currentdir().."/LuaSocket/?.dll"
	socket = require("socket")
	host = "localhost"
	port = 9514
	c = socket.connect(host, port) -- connect to the listener socket
	if c == nil then
		return
	end
	c:setoption("tcp-nodelay",true) -- set immediate transmission mode
end

function LuaExportBeforeNextFrame()
end

function LuaExportAfterNextFrame()
end

function LuaExportStop()
	socket.try(c:send("quit"))
	c:close()
end

function LuaExportActivityNextEvent(t)
	local tNext = t

	if c == nil then
		return
	else
		local t = LoGetModelTime()
		local name = LoGetPilotName()
		local altBar = LoGetAltitudeAboveSeaLevel()
		local altRad = LoGetAltitudeAboveGroundLevel()
		local speedTrue = LoGetTrueAirSpeed()
		local speedInstr = LoGetIndicatedAirSpeed()
		local verticalSpeed = LoGetVerticalVelocity()

		if name == nil then 
			name = ""
		end
		if altBar == nil then
			altBar = 0
		end
		if altRad == nil then
			altRad = 0
		end
		if speedTrue == nil then
			speedTrue = 0
		end
		if speedInstr == nil then
			speedInstr = 0
		end
		if verticalSpeed == nil then
			verticalSpeed = 0           	
		end

		local plainInfo = LoGetSelfData()
		if plainInfo == nil then
			plainInfo = {}
			plainInfo.LatLongAlt = {}
			plainInfo.LatLongAlt.Lat = 0
			plainInfo.LatLongAlt.Long = 0
			plainInfo.Heading = 0
		end

		local payload = "{"
		payload = payload..string.format("\"time\":\"%.2f\",", t)
		payload = payload..string.format("\"name\":\"%s\",", name)
		payload = payload..string.format("\"altBar\":\"%.2f\",", altBar)
		payload = payload..string.format("\"alrRad\":\"%.2f\",", altRad)
		payload = payload..string.format("\"speedTrue\":\"%.2f\",", speedTrue)
		payload = payload..string.format("\"speedInstrumental\":\"%.2f\",", speedInstr)
		payload = payload..string.format("\"vspeed\":\"%.2f\",", verticalSpeed)

		payload = payload.."\"navigation\": {"
		payload = payload..string.format("\"heading\":\"%f\",", plainInfo.Heading)
		payload = payload..string.format("\"lat\":\"%f\",", plainInfo.LatLongAlt.Lat)
		payload = payload..string.format("\"long\":\"%f\"", plainInfo.LatLongAlt.Long)
		payload = payload.."}"

		payload = payload.."}\n"

		local sentResult = c:send(payload)
		if sentResult == nil then
			c:close()
			c = nil
		end
	end
	tNext = tNext + 0.3

	return tNext
end
