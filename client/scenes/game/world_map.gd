extends Node2D

const TILE_SIZE := 16
const CHUNK_SIZE := 16
const LOAD_RADIUS := 2

enum PacketType { HANDSHAKE, WORLD_CHUNK, PLAYER_INPUT, PLAYER_STATE, CHUNK_REQUEST, BLOCK_PLACE, BLOCK_UPDATE }

var seed: int = 0
var loaded_chunks := {}
var last_chunk := Vector2i(99999, 99999)
var placed_blocks := {}
var _atlas_source_id := -1

@onready var terrain: TileMap = $Terrain
@onready var blocks_node: Node2D = $Blocks
@onready var block_texture: Texture2D = preload("res://assets/block.png")

var _network_ref = null
var _local_player_ref = null


func _ready() -> void:
	set_process_input(true)
	_setup_tileset()


func _setup_tileset() -> void:
	var tileset := TileSet.new()
	var source := TileSetAtlasSource.new()
	var texture := load("res://assets/tilemap.png")
	source.texture = texture
	source.texture_region_size = Vector2i(TILE_SIZE, TILE_SIZE)

	source.create_tile(Vector2i(0, 0))
	source.create_tile(Vector2i(1, 0))

	_atlas_source_id = tileset.add_source(source)
	terrain.tile_set = tileset


func set_references(network_node, local_player_node) -> void:
	_network_ref = network_node
	_local_player_ref = local_player_node


func set_seed(s: int) -> void:
	seed = s


func _process(_delta: float) -> void:
	if not _local_player_ref:
		return

	var tile_x := int(_local_player_ref.position.x / TILE_SIZE)
	var tile_y := int(_local_player_ref.position.y / TILE_SIZE)
	var chunk := Vector2i(_tile_to_chunk(tile_x), _tile_to_chunk(tile_y))

	if chunk != last_chunk:
		last_chunk = chunk
		_request_missing_chunks(chunk.x, chunk.y)


func _tile_to_chunk(tile: int) -> int:
	if tile >= 0:
		return tile / CHUNK_SIZE
	return (tile + 1) / CHUNK_SIZE - 1


func _request_missing_chunks(cx: int, cy: int) -> void:
	for dy in range(-LOAD_RADIUS, LOAD_RADIUS + 1):
		for dx in range(-LOAD_RADIUS, LOAD_RADIUS + 1):
			var key := Vector2i(cx + dx, cy + dy)
			if not loaded_chunks.has(key):
				loaded_chunks[key] = false
				if _network_ref and _network_ref.has_method("send_chunk_request"):
					_network_ref.send_chunk_request(key.x, key.y)


func apply_chunk(cx: int, cy: int, tiles: PackedByteArray) -> void:
	var key := Vector2i(cx, cy)
	loaded_chunks[key] = true

	var base_x := cx * CHUNK_SIZE
	var base_y := cy * CHUNK_SIZE

	for y in range(CHUNK_SIZE):
		for x in range(CHUNK_SIZE):
			var tidx := y * CHUNK_SIZE + x
			if tidx >= tiles.size():
				continue
			var tile_id := tiles[tidx]
			if tile_id >= 2:
				continue
			var atlas := Vector2i(0, 0) if tile_id == 0 else Vector2i(1, 0)
			terrain.set_cell(0, Vector2i(base_x + x, base_y + y), _atlas_source_id, atlas)


func apply_block(tile_x: int, tile_y: int, _block_type: int) -> void:
	var key := Vector2i(tile_x, tile_y)
	if placed_blocks.has(key):
		return

	placed_blocks[key] = true
	var sprite := Sprite2D.new()
	sprite.texture = block_texture
	sprite.position = Vector2(tile_x * TILE_SIZE + TILE_SIZE / 2, tile_y * TILE_SIZE + TILE_SIZE / 2)
	sprite.centered = true
	blocks_node.add_child(sprite)


func _input(event: InputEvent) -> void:
	if not event is InputEventMouseButton:
		return
	if event.button_index != MOUSE_BUTTON_LEFT or not event.pressed:
		return
	if not _network_ref:
		return

	var mouse_pos := get_global_mouse_position()
	var tile_x := int(mouse_pos.x / TILE_SIZE)
	var tile_y := int(mouse_pos.y / TILE_SIZE)

	_network_ref.send_block_place(tile_x, tile_y)
