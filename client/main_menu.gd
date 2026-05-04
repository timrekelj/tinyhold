extends Control

@onready var singleplayer_button: Button = $VBoxContainer/SingleplayerButton
@onready var multiplayer_button: Button = $VBoxContainer/MultiplayerButton
@onready var ip_line_edit: LineEdit = $VBoxContainer/ServerIpLineEdit

func _ready() -> void:
	singleplayer_button.pressed.connect(_on_singleplayer_pressed)
	multiplayer_button.pressed.connect(_on_multiplayer_pressed)

func _on_singleplayer_pressed() -> void:
	GameSession.host = "127.0.0.1"
	GameSession.port = 42069

	var server_path := ProjectSettings.globalize_path("res://../server/build/tinyhold-server")
	if not GameSession.start_local_server(server_path):
		return

	await get_tree().create_timer(0.5).timeout
	get_tree().change_scene_to_file("res://game.tscn")

func _on_multiplayer_pressed() -> void:
	var host := ip_line_edit.text.strip_edges()
	if host.is_empty():
		print("Enter an IP address")
		return

	GameSession.host = host
	GameSession.port = 42069

	get_tree().change_scene_to_file("res://game.tscn")
