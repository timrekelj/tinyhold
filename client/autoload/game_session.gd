extends Node

var host := "127.0.0.1"
var port := 42069
var server_pid := -1

func start_local_server(server_path: String) -> bool:
	if server_pid != -1:
		return true

	if not FileAccess.file_exists(server_path):
		print("Local server binary not found: ", server_path)
		return false

	server_pid = OS.create_process(server_path, [])
	if server_pid == -1:
		print("Failed to start local server: ", server_path)
		return false

	return true

func stop_local_server() -> void:
	if server_pid != -1:
		OS.kill(server_pid)
		server_pid = -1

func _exit_tree() -> void:
	stop_local_server()
