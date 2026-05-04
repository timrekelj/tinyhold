extends Node2D

enum PacketType { HANDSHAKE, WORLD_CHUNK, PLAYER_INPUT, PLAYER_STATE, ENTITY_UPDATE, HEARTBEAT, DISCONNECT, INVENTORY }

var socket := StreamPeerTCP.new()
var connected := false
var player_id := 0

var player_scene = preload("res://player.tscn")
var players := {}

func connect_to_server(host: String, port: int = 42069) -> void:
	var err := socket.connect_to_host(host, port)
	if err != OK:
		print("connect_to_host failed: ", err)
		return
	print("Connecting to %s:%d" % [host, port])


func _process(_delta: float) -> void:
	var err := socket.poll()
	if err != OK:
		print("Socket poll failed: ", err)
		return

	var status = socket.get_status()

	if !connected:
		match status:
			StreamPeerTCP.STATUS_CONNECTED:
				connected = true
				send_handshake()
			StreamPeerTCP.STATUS_ERROR:
				print("Failed to connect!")
			_:
				pass
	else:
		if socket.get_available_bytes() > 0:
			handle_packet()


func send_handshake():
	var packet := PackedByteArray([0, 1, 0, 1])
	socket.put_data(packet)
	print("Handshake sent!")


func send_player_input(keys: int) -> void:
	var packet := PackedByteArray()
	packet.append(PacketType.PLAYER_INPUT)
	packet.append(5)
	packet.append(0)
	packet.append(keys)
	packet.append(0)
	packet.append(0)
	packet.append(0)
	packet.append(0)
	socket.put_data(packet)


func handle_packet():
	if socket.get_available_bytes() < 3:
		return

	var result = socket.get_data(3)
	var err = result[0]
	var header = result[1]
	if err != OK or header.size() < 3:
		return

	var ptype = header[0]
	var length = header[1] | (header[2] << 8)

	if socket.get_available_bytes() < length:
		return

	result = socket.get_data(length)
	err = result[0]
	var payload = result[1]
	if err != OK:
		return

	match ptype:
		PacketType.HANDSHAKE:
			var version = payload[0]
			player_id = payload[1] | (payload[2] << 8)
			print("Handshake response: version=%d, player_id=%d" % [version, player_id])
			spawn_player(player_id, true)
		PacketType.PLAYER_STATE:
			handle_player_state(payload)
		_:
			print("Unknown packet type: %d" % ptype)


func spawn_player(id: int, is_local_player: bool):
	if players.has(id):
		var p = players[id]
		if is_local_player and not p.is_local:
			p.is_local = true
			if p.camera:
				p.camera.enabled = true
		return
	else:
		var p = player_scene.instantiate()
		p.player_id = id
		p.is_local = is_local_player
		players[id] = p
		add_child(p)
		print("Spawned player %d (local=%s)" % [id, is_local_player])


func handle_player_state(payload: PackedByteArray):
	if payload.size() < 10:
		return

	var id = payload.decode_u16(0)
	var x = payload.decode_s32(2)
	var y = payload.decode_s32(6)

	if not players.has(id):
		spawn_player(id, false)

	var p = players[id]
	p.update_state(x, y)
