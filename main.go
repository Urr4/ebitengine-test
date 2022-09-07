package main

import (
	"fmt"
	"github.com/gen2brain/raylib-go/raylib"
	"os"
	"strconv"
	"strings"
)

const (
	screenWidth  = 1000
	screenHeight = 480
)

var (
	running = true

	backgroundColor = rl.NewColor(147, 211, 196, 255)
	grassSprite     rl.Texture2D
	fenceSprite     rl.Texture2D
	tex             rl.Texture2D

	playerSprite                                  rl.Texture2D
	playerSrc                                     rl.Rectangle
	playerDest                                    rl.Rectangle
	playerSpeed                                   float32 = 3
	playerMoving                                  bool
	playerDir                                     int
	playerUp, playerDown, playerLeft, playerRight bool
	playerFrame                                   int

	frameCount int

	tileDest            rl.Rectangle
	tileSrc             rl.Rectangle
	tileMap             []int
	srcMap              []string
	mapWidth, mapHeight int

	musicPaused bool
	music       rl.Music

	camera rl.Camera2D
)

func drawScene() {

	for i := 0; i < len(tileMap); i++ {
		if tileMap[i] != 0 {
			tileDest.X = tileDest.Width * float32(i%mapWidth)
			tileDest.Y = tileDest.Height * float32(i/mapWidth)

			switch srcMap[i] {
			case "g":
				tex = grassSprite
			case "f":
				tex = fenceSprite
			default:
				tex = grassSprite
			}

			tileSrc.X = tileSrc.Width * float32((tileMap[i]-1)%int(tex.Width/int32(tileSrc.Width)))
			tileSrc.Y = tileSrc.Height * float32((tileMap[i]-1)/int(tex.Width/int32(tileSrc.Width)))

			rl.DrawTexturePro(tex, tileSrc, tileDest, rl.NewVector2(tileDest.Width, tileDest.Height), 0, rl.White)
		}
	}

	rl.DrawTexturePro(playerSprite, playerSrc, playerDest, rl.NewVector2(playerDest.Width, playerDest.Height), 0, rl.White)
}

func input() {
	if rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyUp) {
		playerMoving = true
		playerUp = true
		playerDir = 1
	}
	if rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyDown) {
		playerMoving = true
		playerDown = true
		playerDir = 0
	}
	if rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft) {
		playerMoving = true
		playerLeft = true
		playerDir = 2
	}
	if rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight) {
		playerMoving = true
		playerRight = true
		playerDir = 3
	}
	if rl.IsKeyPressed(rl.KeyQ) {
		musicPaused = !musicPaused
	}

}

func update() {
	running = !rl.WindowShouldClose()

	playerSrc.X = playerSrc.Width * float32(playerFrame)

	if playerMoving {
		if playerUp {
			playerDest.Y -= playerSpeed
		}
		if playerDown {
			playerDest.Y += playerSpeed
		}
		if playerLeft {
			playerDest.X -= playerSpeed
		}
		if playerRight {
			playerDest.X += playerSpeed
		}
		if frameCount%8 == 0 {
			playerFrame++
		}
	} else if frameCount%45 == 1 {
		playerFrame++
	}
	frameCount++
	if playerFrame > 3 {
		playerFrame = 0
	}
	if !playerMoving && playerFrame > 1 {
		playerFrame = 0
	}

	playerSrc.Y = playerSrc.Height * float32(playerDir)

	rl.UpdateMusicStream(music)
	if musicPaused {
		rl.PauseMusicStream(music)
	} else {
		rl.ResumeMusicStream(music)
	}

	camera.Target = rl.NewVector2(playerDest.X-(playerDest.Width/2), playerDest.Y-(playerDest.Height/2))

	playerMoving = false
	playerUp, playerDown, playerLeft, playerRight = false, false, false, false
}

func loadMap(mapFile string) {
	file, err := os.ReadFile(mapFile)

	if err != nil {
		fmt.Println("Couldn't load file", mapFile)
		os.Exit(1)
	}

	remNewLines := strings.Replace(string(file), "\n", " ", -1)
	sliced := strings.Split(remNewLines, " ")

	mapWidth = -1
	mapHeight = -1
	for i := 0; i < len(sliced); i++ {
		s, _ := strconv.ParseInt(sliced[i], 10, 64)
		m := int(s)
		if mapWidth == -1 {
			mapWidth = m
		} else if mapHeight == -1 {
			mapHeight = m
		} else if i < mapWidth*mapHeight+2 {
			tileMap = append(tileMap, m)
		} else {
			srcMap = append(srcMap, sliced[i])
		}
	}
	if len(tileMap) > mapWidth*mapHeight {
		tileMap = tileMap[:len(tileMap)-1]
	}
}

func render() {
	rl.BeginDrawing()
	rl.ClearBackground(backgroundColor)
	rl.BeginMode2D(camera)

	drawScene()

	rl.EndMode2D()
	rl.EndDrawing()
}

// Go calls init funcs before the main func
func init() {
	rl.InitWindow(screenWidth, screenHeight, "Funny Game Stuff")
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)

	grassSprite = rl.LoadTexture("assets/Tilesets/Grass.png")
	fenceSprite = rl.LoadTexture("assets/Tilesets/Fence.png")
	tileDest = rl.NewRectangle(0, 0, 16, 16)
	tileSrc = rl.NewRectangle(0, 0, 16, 16)
	playerSprite = rl.LoadTexture("assets/Characters/BasicCharakterSpritesheet.png")

	playerSrc = rl.NewRectangle(0, 0, 48, 48)
	playerDest = rl.NewRectangle(200, 200, 100, 100)

	camera = rl.NewCamera2D(rl.NewVector2(float32(screenWidth/2), float32(screenHeight/2)), rl.NewVector2(playerDest.X, playerDest.Y), 0.0, 1.0)

	rl.InitAudioDevice()
	music = rl.LoadMusicStream("assets/Sound/Farm.mp3")
	musicPaused = false
	rl.PlayMusicStream(music)

	loadMap("static/world.map")
}

func quit() {
	rl.UnloadTexture(grassSprite)
	rl.UnloadTexture(playerSprite)
	rl.UnloadMusicStream(music)
	rl.CloseAudioDevice()
	rl.CloseWindow()
}

func main() {
	for running {
		input()
		update()
		render()
	}
	quit()
}
