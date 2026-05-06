extends Node2D

@onready var network := $Network


func _ready() -> void:
	network.connect_to_server(GameSession.host, GameSession.port)


func send_player_input(keys: int) -> void:
	network.send_player_input(keys)
