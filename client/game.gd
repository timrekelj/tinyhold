extends Node2D

@onready var network := $Network


func _ready() -> void:
	network.connect_to_server(GameSession.host, GameSession.port)
