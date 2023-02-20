package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth = 1000
	screenHeight = 480
)
var(
	running = true
	bkgColor = rl.NewColor(147, 211, 196, 255)

	grassSprite rl.Texture2D
	hillSprite rl.Texture2D
	fenceSprite rl.Texture2D
	houseSprite rl.Texture2D
	waterSprite rl.Texture2D
	tilledSprite rl.Texture2D
	tex rl.Texture2D

	playerSprite rl.Texture2D

	playerSrc rl.Rectangle
	playerDest rl.Rectangle
	playerMoving bool
	playerDir int
	playerUp, playerDown, playerRight, playerLeft bool
	playerFrame int

	frameCount int

	tileDest rl.Rectangle
	tileSrc rl.Rectangle
	tileMap []int
	srcMap []string
	mapW, mapH int


	playerSpeed float32 = 1.4

	musicPaused bool
	music rl.Music

	cam rl.Camera2D
)

func drawScene() {
	//rl.DrawTexture(grassSprite, 100, 50, rl.White)

	for i:=0; i<len(tileMap); i++ {
		if tileMap[i] != 0 {
			tileDest.X = tileDest.Width * float32(i % mapW)
			tileDest.Y = tileDest.Height * float32(i / mapW)

			if srcMap[i] == "g" { tex = grassSprite}
			if srcMap[i] == "l" { tex = hillSprite}
			if srcMap[i] == "f" { tex = fenceSprite}
			if srcMap[i] == "h" { tex = houseSprite}
			if srcMap[i] == "w" { tex = waterSprite}
			if srcMap[i] == "t" { tex = tilledSprite}
			
			if srcMap[i] == "h" || srcMap[i] == "f" {
				tileSrc.X = 0
				rl.DrawTexturePro(grassSprite, tileSrc, tileDest, rl.NewVector2(tileDest.Width, tileDest.Height), 0, rl.White)
			}


			tileSrc.X = tileSrc.Width * float32((tileMap[i]-1) % int(tex.Width/int32(tileSrc.Width)))
			tileSrc.Y = tileSrc.Height * float32((tileMap[i]-1) / int(tex.Width/int32(tileSrc.Width)))	
			rl.DrawTexturePro(tex, tileSrc, tileDest, rl.NewVector2(tileDest.Width, tileDest.Height), 0, rl.White)
		}
	}

	rl.DrawTexturePro(playerSprite, playerSrc, playerDest, rl.NewVector2(playerDest.Width, playerDest.Height), 0, rl.White)
}

func input() {
	//bevegelse opp
	if rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyUp) {
		playerMoving = true
		playerDir = 1
		playerUp = true
	}
	//bevegelse Ned
	if rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyDown) {
		playerMoving = true
		playerDir = 0
		playerDown = true
	}
	//bevegelse venstre
	if rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft) {
		playerMoving = true
		playerDir = 2
		playerLeft = true
	}

	// bevegelse høyre
	if rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight) {
		playerMoving = true
		playerDir = 3
		playerRight = true
	}

	//skru av musikk
	if rl.IsKeyPressed(rl.KeyQ) {
		musicPaused = !musicPaused
	}

} 
func update() {
	//hvis ikke spillet kjører, lukk vindu
	running = !rl.WindowShouldClose()

	//størrelse på karakter
	playerSrc.X = playerSrc.Width * float32(playerFrame)

	if playerMoving {
		if playerUp {playerDest.Y -= playerSpeed}
		if playerDown {playerDest.Y += playerSpeed}
		if playerLeft {playerDest.X -= playerSpeed}
		if playerRight {playerDest.X += playerSpeed}
		if frameCount % 8 == 1 {playerFrame++}
	} else if frameCount % 45 == 1 {
		playerFrame++	
	}

	frameCount++
	if playerFrame > 3 {	playerFrame = 0}
	if !playerMoving && playerFrame > 1 { playerFrame = 0}


	playerSrc.X = playerSrc.Width * float32(playerFrame)
	playerSrc.Y = playerSrc.Height * float32(playerDir)

	rl.UpdateMusicStream(music)
	if musicPaused {
		rl.PauseMusicStream(music)
	}else {
		rl.ResumeMusicStream(music)
	}

	cam.Target = rl.NewVector2(float32(playerDest.X-(playerDest.Width/2)), float32(playerDest.Y-(playerDest.Height/2)))



	playerMoving = false
	playerUp, playerDown, playerRight, playerLeft = false, false, false, false
}
func render() {

	rl.BeginDrawing()	
	rl.ClearBackground(bkgColor)
rl.BeginMode2D(cam)

	
	drawScene()

	rl.EndMode2D()
	rl.EndDrawing()
}

func loadMap(mapFile string) {
	file, err := ioutil.ReadFile(mapFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	remNewLines := strings.Replace(string(file), "\r", " ", -1)
	sliced := strings.Split(remNewLines, " ")
	mapW = -1
	mapH = -1
	for i:=0; i<len(sliced); i++ {
		s, _ := strconv.ParseInt(sliced[i], 10, 64)
		m := int(s)
		if mapW == -1 {
			mapW = m
		} else if mapH == -1 {
			mapH = m
		} else if i < mapW*mapH+2{
			tileMap = append(tileMap, m)
		} else {
			srcMap = append(srcMap, sliced[i])
		}
	}
	if len(tileMap) > mapW*mapH { tileMap = tileMap[:len(tileMap)-1]}


	/*mapW = 5
	mapH = 5
	for i:=0; i<(mapW*mapH); i++ {
		tileMap = append(tileMap, 1)
	}
	*/
}

func init() {
	rl.InitWindow(screenWidth, screenHeight, "Sykeste Ever!")
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)
	grassSprite = rl.LoadTexture("res/Tilesets/groundtiles/oldtiles/Grass.png")
	hillSprite = rl.LoadTexture("res/Tilesets/groundtiles/oldtiles/Hills.png")
	fenceSprite = rl.LoadTexture("res/Tilesets/Fences.png")
	houseSprite = rl.LoadTexture("res/Tilesets/House.png")
	waterSprite = rl.LoadTexture("res/Tilesets/Water.png")
	tilledSprite = rl.LoadTexture("res/Tilesets/groundtiles/oldtiles/Tilled.png")

	tileDest = rl.NewRectangle(0,0, 16, 16)
	tileSrc = rl.NewRectangle(0, 0, 16, 16)

	playerSprite = rl.LoadTexture("res/Characters/BasicCharakterSpritesheet.png")

	playerSrc = rl.NewRectangle(0, 0, 48, 48)
	playerDest = rl.NewRectangle(200, 200, 60, 60)

	rl.InitAudioDevice()
	music = rl.LoadMusicStream("res/Music/GameSound.mp3")
	musicPaused = false
	rl.PlayMusicStream(music)

	cam = rl.NewCamera2D(rl.NewVector2(float32(screenWidth/2), float32(screenHeight/2)), rl.NewVector2(float32(playerDest.X-(playerDest.Width/2)), float32(playerDest.Y-(playerDest.Height/2))), 0.0, 1.5)

	cam.Zoom = 3

	loadMap("one.map")

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