package main

import (
	"ebiten-rpg/entities"
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	player      *entities.Player
	camera      *Camera
	enemies     []*entities.Enemy
	items       []*entities.Item
	tilemapJSON *TilemapJSON
	tilesets    []Tileset
	tilemapImg  *ebiten.Image
}

func (g *Game) Update() error {

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.X += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.X -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.player.Y -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.player.Y += 2
	}

	// enemies behaviour
	for _, sprite := range g.enemies {
		if sprite.FollowsPlayer {

			if sprite.X < g.player.X {
				sprite.X += 1
			} else if sprite.X > g.player.X {
				sprite.X -= 1
			}

			if sprite.Y < g.player.Y {
				sprite.Y += 1
			} else if sprite.Y > g.player.Y {
				sprite.Y -= 1
			}
		}

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

	//draw player
	screen.DrawImage(g.player.Img.SubImage(
		image.Rect(0, 0, 16, 16)).(*ebiten.Image),
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

	game := Game{
		player: &entities.Player{
			Sprite: &entities.Sprite{
				Img: playerImg,
				X:   120.0,
				Y:   120.0,
			},
			Health: 100,
		},

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
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
