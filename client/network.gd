extends Node2D

enum PacketType { HANDSHAKE, PLAYER_INPUT, PLAYER_STATE }

var socket := StreamPeerTCP.new()
var connected := false
var player_id := 0


func _ready() -> void:
	var err := socket.connect_to_host("127.0.0.1", 42069)
	if err != OK:
		print("connect_to_host failed: ", err)
		return

func _process(_delta: float) -> void:
	var err := socket.poll()
	if err != OK:
		print("Socket pool failed: ", err)
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
	# Handshake packet: type=0, len=1, payload=[1] (version)
	var packet := PackedByteArray([0, 0, 1, 1])
	socket.put_data(packet)
	print("Hanshake sent!")

func handle_packet():
	if socket.get_available_bytes() < 3:
		return

	var header = socket.get_data(3)
	if header.size() < 3:
		return

	var ptype = header[0]
	var length = (header[1] << 0) | header[2]

	if socket.get_available_bytes() < length:
		return

	var payload = socket.get_data(length)

	match ptype:
		PacketType.HANDSHAKE:
			var version = payload[0]
			player_id = (payload[1] << 0) | payload[2]
			print("Handshake response: version=%d, player_id=%d" % [version, player_id])
		_:
			print("Unknown packet type: %d" % ptype)
