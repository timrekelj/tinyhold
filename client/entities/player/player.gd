extends CharacterBody2D

@export var player_id: int = 0
@export var is_local: bool = false

@onready var sprite: AnimatedSprite2D = $AnimatedSprite2D
@onready var camera: Camera2D = $Camera2D


func _ready() -> void:
	if camera:
		camera.enabled = is_local
	if sprite:
		sprite.play("idle")


func _physics_process(_delta: float) -> void:
	if not is_local:
		return

	var keys := 0
	if Input.is_action_pressed("up"): keys |= 1
	if Input.is_action_pressed("down"): keys |= 2
	if Input.is_action_pressed("left"): keys |= 4
	if Input.is_action_pressed("right"): keys |= 8

	var game = get_parent()
	if game and game.has_method("send_player_input"):
		game.send_player_input(keys)


func update_state(x_sub: int, y_sub: int) -> void:
	var new_pos := Vector2(x_sub / 100.0, y_sub / 100.0)
	var vel_approx := new_pos - position

	if sprite:
		if vel_approx.length() > 0.01:
			if sprite.animation != "run":
				sprite.play("run")
			if vel_approx.x < 0:
				sprite.flip_h = true
			elif vel_approx.x > 0:
				sprite.flip_h = false
		else:
			if sprite.animation != "idle":
				sprite.play("idle")

	position = new_pos
