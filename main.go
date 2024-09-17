package main

import (
	"ebiten-rpg/animations"
	"ebiten-rpg/entities"
	"ebiten-rpg/spritesheet"
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func CheckCollisionVertical(sprite *entities.Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(image.Rect(
			int(sprite.X),
			int(sprite.Y),
			int(sprite.X)+16.0,
			int(sprite.Y)+16.0,
		)) {
			//fmt.Println("Player colliding with colider")
			if sprite.Dy > 0.0 {
				sprite.Y = float64(collider.Min.Y) - 16.0
			} else if sprite.Dy < 0.0 {
				sprite.Y = float64(collider.Max.Y)

			}
		}
	}
}

func CheckCollisionHorizontal(sprite *entities.Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(image.Rect(
			int(sprite.X),
			int(sprite.Y),
			int(sprite.X)+16.0,
			int(sprite.Y)+16.0,
		)) {
			//fmt.Println("Player colliding with colider")
			if sprite.Dx > 0.0 {
				sprite.X = float64(collider.Min.X) - 16.0
			} else if sprite.Dx < 0.0 {
				sprite.X = float64(collider.Max.X)

			}
		}
	}

}

type Game struct {
	player            *entities.Player
	playerSpriteSheet *spritesheet.SpriteSheet
	camera            *Camera
	colliders         []image.Rectangle
	enemies           []*entities.Enemy
	items             []*entities.Item
	tilemapJSON       *TilemapJSON
	tilesets          []Tileset
	tilemapImg        *ebiten.Image
}

func (g *Game) Update() error {

	g.player.Dx = 0.0
	g.player.Dy = 0.0

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.Dx = 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.Dx = -2
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.player.Dy = -2
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.player.Dy = 2
	}

	g.player.X += g.player.Dx

	CheckCollisionHorizontal(g.player.Sprite, g.colliders)

	g.player.Y += g.player.Dy

	CheckCollisionVertical(g.player.Sprite, g.colliders)

	//animations
	activeAnim := g.player.ActiveAnimation(int(g.player.Dx), int(g.player.Dy))
	if activeAnim != nil {
		activeAnim.Update()
	}

	// enemies behaviour
	for _, sprite := range g.enemies {
		sprite.Dx = 0.0
		sprite.Dy = 0.0

		if sprite.FollowsPlayer {

			if sprite.X < g.player.X {
				sprite.Dx += 1
			} else if sprite.X > g.player.X {
				sprite.Dx -= 1
			}

			if sprite.Y < g.player.Y {
				sprite.Dy += 1
			} else if sprite.Y > g.player.Y {
				sprite.Dy -= 1
			}
		}

		sprite.X += sprite.Dx
		CheckCollisionHorizontal(sprite.Sprite, g.colliders)

		sprite.Y += sprite.Dy
		CheckCollisionVertical(sprite.Sprite, g.colliders)

	}

	//health pick
	for _, item := range g.items {
		if g.player.X == item.X && g.player.Y == item.Y {
			g.player.Health += item.HealAmount
			fmt.Printf("Picked health!: %d \n", g.player.Health)
		}

	}

	offset := 8.0
	g.camera.FollowTarget(g.player.X+offset, g.player.Y+offset, 320, 240)
	g.camera.Constrain(float64(g.tilemapJSON.Layers[0].Width)*16.0, float64(g.tilemapJSON.Layers[0].Height)*16.0, 320, 240)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{120, 180, 255, 255})
	//ebitenutil.DebugPrint(screen, "Hello, World!")

	opts := ebiten.DrawImageOptions{}

	//loop over tile layers
	for layerIndex, layer := range g.tilemapJSON.Layers {
		// loop over tile data
		//log.Printf("LAYERDATA: : %v", layer.Data)

		for index, id := range layer.Data {

			if id == 0 {
				continue
			}
			//getting tile position
			x := index % layer.Width
			y := index / layer.Width

			//covert tile position to pixel position
			x *= 16
			y *= 16

			// log.Printf("layerIndex: %d", layerIndex)
			// log.Printf("id : %d", id)

			img := g.tilesets[layerIndex].Img(id)

			//set the drawimage options to draw the tile at x,y
			opts.GeoM.Translate(float64(x), float64(y))

			//fix position, top left
			opts.GeoM.Translate(0.0, -(float64(img.Bounds().Dy()) + 16))

			opts.GeoM.Translate(g.camera.X, g.camera.Y)

			//draw the tile
			screen.DrawImage(
				img,
				&opts,
			)

			//reset the opts for the next itaration
			opts.GeoM.Reset()
		}
	}

	// set drawimageoptions for the player
	opts.GeoM.Translate(g.player.X, g.player.Y)
	opts.GeoM.Translate(g.camera.X, g.camera.Y)

	playerFrame := 0
	activeAnim := g.player.ActiveAnimation(int(g.player.Dx), int(g.player.Dy))
	if activeAnim != nil {
		playerFrame = activeAnim.Frame()
	}

	//draw player
	screen.DrawImage(g.player.Img.SubImage(
		g.playerSpriteSheet.Rect(playerFrame),
	).(*ebiten.Image),
		&opts,
	)

	opts.GeoM.Reset()

	for _, sprite := range g.enemies {
		opts.GeoM.Translate(sprite.X, sprite.Y)
		opts.GeoM.Translate(g.camera.X, g.camera.Y)

		screen.DrawImage(sprite.Img.SubImage(image.Rect(0, 0, 16, 16)).(*ebiten.Image),
			&opts,
		)
		opts.GeoM.Reset()
	}

	opts.GeoM.Reset()

	for _, sprite := range g.items {
		opts.GeoM.Translate(sprite.X, sprite.Y)
		opts.GeoM.Translate(g.camera.X, g.camera.Y)

		screen.DrawImage(sprite.Img.SubImage(image.Rect(0, 0, 16, 16)).(*ebiten.Image),
			&opts,
		)
		opts.GeoM.Reset()
	}

	//Colliders
	for _, collider := range g.colliders {
		vector.StrokeRect(screen,
			float32(collider.Min.X)+float32(g.camera.X),
			float32(collider.Min.Y)+float32(g.camera.Y),
			float32(collider.Dx()), float32(collider.Dy()), 1.0,
			color.RGBA{255, 0, 0, 255}, true,
		)
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
	//return ebiten.WindowSize()
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Go RPG!")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	playerImg, _, err := ebitenutil.NewImageFromFile("assets/img/ninja.png")
	if err != nil {
		log.Fatal(err)
	}

	skeletonImg, _, err := ebitenutil.NewImageFromFile("assets/img/skeleton.png")
	if err != nil {
		log.Fatal(err)
	}

	sushiImg, _, err := ebitenutil.NewImageFromFile("assets/img/sushi.png")
	if err != nil {
		log.Fatal(err)
	}
	tilemapImg, _, err := ebitenutil.NewImageFromFile("assets/img/tilesetFloor.png")
	if err != nil {
		log.Fatal(err)
	}
	tilemapJSON, err := NewTilemapJSON("assets/maps/tileset/spawn.json")
	if err != nil {
		log.Fatal(err)
	}

	tilesets, err := tilemapJSON.GenTilesets()
	if err != nil {
		log.Fatal(err)
	}

	playerSpriteSheet := spritesheet.NewSpriteSheet(4, 7, 16)

	game := Game{
		player: &entities.Player{
			Sprite: &entities.Sprite{
				Img: playerImg,
				X:   120.0,
				Y:   120.0,
			},
			Health: 100,
			Animations: map[entities.PLayerState]*animations.Animation{
				entities.Up:    animations.NewAnimation(5, 13, 4, 20.0),
				entities.Down:  animations.NewAnimation(4, 12, 4, 20.0),
				entities.Left:  animations.NewAnimation(6, 14, 4, 20.0),
				entities.Right: animations.NewAnimation(7, 15, 4, 20.0),
			},
		},
		playerSpriteSheet: playerSpriteSheet,
		enemies: []*entities.Enemy{
			{
				Sprite: &entities.Sprite{
					Img: skeletonImg,
					X:   100.0,
					Y:   100.0,
				},
				FollowsPlayer: true,
			},
			{
				Sprite: &entities.Sprite{
					Img: skeletonImg,
					X:   150.0,
					Y:   50.0,
				},
				FollowsPlayer: false,
			},
			{
				Sprite: &entities.Sprite{
					Img: skeletonImg,
					X:   75.0,
					Y:   75.0,
				},
				FollowsPlayer: false,
			},
		},

		items: []*entities.Item{
			{
				Sprite: &entities.Sprite{
					Img: sushiImg,
					X:   250.0,
					Y:   95.0,
				},
				HealAmount: 5,
			},
		},

		tilemapJSON: tilemapJSON,
		tilemapImg:  tilemapImg,
		tilesets:    tilesets,
		camera:      NewCamera(0.0, 0.0),
		colliders: []image.Rectangle{
			image.Rect(100, 100, 116, 116),
		},
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
