extends Node2D

@onready var network := $Network
@onready var world_map := $WorldMap


func _ready() -> void:
	network.set_world_map(world_map)
	network.connect_to_server(GameSession.host, GameSession.port)


func send_player_input(keys: int) -> void:
	network.send_player_input(keys)
