extends Node

enum PacketType { HANDSHAKE, WORLD_CHUNK, PLAYER_INPUT, PLAYER_STATE, CHUNK_REQUEST, BLOCK_PLACE, BLOCK_UPDATE }

var socket := StreamPeerTCP.new()
var connected := false
var player_id := 0
var seed := 0

var player_scene = preload("res://entities/player/player.tscn")
var players := {}

var world_map = null


func connect_to_server(host: String, port: int = 42069) -> void:
	var err := socket.connect_to_host(host, port)
	if err != OK:
		print("connect_to_host failed: ", err)
		return
	print("Connecting to %s:%d" % [host, port])


func set_world_map(wm) -> void:
	world_map = wm


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
		while socket.get_available_bytes() > 0:
			var prev_pending = _pending_type
			handle_packet()
			if _pending_type == prev_pending and _pending_type != -1:
				break


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


func send_chunk_request(cx: int, cy: int) -> void:
	var packet := PackedByteArray()
	packet.append(PacketType.CHUNK_REQUEST)
	packet.append(8)
	packet.append(0)
	packet.append(cx & 0xFF)
	packet.append((cx >> 8) & 0xFF)
	packet.append((cx >> 16) & 0xFF)
	packet.append((cx >> 24) & 0xFF)
	packet.append(cy & 0xFF)
	packet.append((cy >> 8) & 0xFF)
	packet.append((cy >> 16) & 0xFF)
	packet.append((cy >> 24) & 0xFF)
	socket.put_data(packet)


func send_block_place(tile_x: int, tile_y: int) -> void:
	var packet := PackedByteArray()
	packet.append(PacketType.BLOCK_PLACE)
	packet.append(8)
	packet.append(0)
	packet.append(tile_x & 0xFF)
	packet.append((tile_x >> 8) & 0xFF)
	packet.append((tile_x >> 16) & 0xFF)
	packet.append((tile_x >> 24) & 0xFF)
	packet.append(tile_y & 0xFF)
	packet.append((tile_y >> 8) & 0xFF)
	packet.append((tile_y >> 16) & 0xFF)
	packet.append((tile_y >> 24) & 0xFF)
	socket.put_data(packet)


var _pending_type := -1
var _pending_length := 0


func handle_packet():
	var available := socket.get_available_bytes()

	if _pending_type == -1:
		if available < 3:
			return
		var result = socket.get_data(3)
		var err = result[0]
		var header = result[1]
		if err != OK or header.size() < 3:
			return
		_pending_type = header[0]
		_pending_length = header[1] | (header[2] << 8)
		return

	if available < _pending_length:
		return

	var result = socket.get_data(_pending_length)
	var err = result[0]
	var payload = result[1]
	if err != OK:
		_pending_type = -1
		return

	var ptype = _pending_type
	var length = _pending_length
	_pending_type = -1

	match ptype:
		PacketType.HANDSHAKE:
			var version = payload[0]
			player_id = payload[1] | (payload[2] << 8)
			seed = payload.decode_s64(3)
			print("Handshake response: version=%d, player_id=%d, seed=%d" % [version, player_id, seed])
			if world_map:
				world_map.set_seed(seed)
			spawn_player(player_id, true)
			if world_map:
				world_map.set_references(self, players[player_id])
		PacketType.PLAYER_STATE:
			handle_player_state(payload)
		PacketType.WORLD_CHUNK:
			handle_world_chunk(payload)
		PacketType.BLOCK_UPDATE:
			handle_block_update(payload)
		_:
			print("Unknown packet type: %d" % ptype)


func handle_world_chunk(payload: PackedByteArray) -> void:
	if payload.size() < 8 + 256:
		print("WorldChunk payload too small: %d" % payload.size())
		return

	var cx := payload.decode_s32(0)
	var cy := payload.decode_s32(4)
	var tiles := payload.slice(8)

	if world_map:
		world_map.apply_chunk(cx, cy, tiles)


func handle_block_update(payload: PackedByteArray) -> void:
	if payload.size() < 9:
		return

	var tile_x := payload.decode_s32(0)
	var tile_y := payload.decode_s32(4)
	var block_type := payload[8]

	if world_map:
		world_map.apply_block(tile_x, tile_y, block_type)


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
