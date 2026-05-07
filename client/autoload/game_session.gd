extends Node

var host := "127.0.0.1"
var port := 42069
var server_pid := -1

func find_free_port() -> int:
	var temp_server := TCPServer.new()
	var err := temp_server.listen(0, "127.0.0.1")
	if err != OK:
		push_error("Failed to find a free port")
		return -1
	var free_port := temp_server.get_local_port()
	temp_server.stop()
	return free_port

func start_local_server(server_path: String) -> bool:
	if server_pid != -1:
		return true

	if not FileAccess.file_exists(server_path):
		print("Local server binary not found: ", server_path)
		return false

	for attempt in range(3):
		var free_port := find_free_port()
		if free_port == -1:
			continue

		port = free_port
		var args := PackedStringArray(["--local", "--port=" + str(free_port)])
		server_pid = OS.create_process(server_path, args)
		if server_pid == -1:
			print("Failed to start local server (attempt ", attempt + 1, "): ", server_path)
			continue

		return true

	return false

func stop_local_server() -> void:
	if server_pid != -1:
		OS.kill(server_pid)
		server_pid = -1

func _exit_tree() -> void:
	stop_local_server()
