package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
)

type Game struct {
	playerX, playerY float64
	playerSize       float64
	playerVelY       float64
	isJumping        bool
	platforms        []Platform
}

type Platform struct {
	x, y, width, height float64
}

const (
	gravity       = 0.5
	jumpVelocity  = -10
	playerStartY  = 500 // Starting Y position of the player (ground level)
	playerGroundY = 500 // Y position of the ground
)

func rectCollide(x1, y1, w1, h1, x2, y2, w2, h2 float64) bool {
	return x1 < x2+w2 && x1+w1 > x2 && y1 < y2+h2 && y1+h1 > y2
}

func (g *Game) Update() error {
	// Handle player movement
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.playerX -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.playerX += 2
	}

	// Handle jumping
	if ebiten.IsKeyPressed(ebiten.KeySpace) && !g.isJumping {
		g.playerVelY = jumpVelocity
		g.isJumping = true
	}

	// Update the player's vertical position based on velocity
	g.playerY += g.playerVelY
	g.playerVelY += gravity

	// Check collisions with platforms
	for _, platform := range g.platforms {
		if rectCollide(g.playerX, g.playerY, g.playerSize, g.playerSize, platform.x, platform.y, platform.width, platform.height) {
			// Collision detected
			if g.playerY+g.playerSize <= platform.y+g.playerVelY {
				// Collision from the top
				g.playerY = platform.y - g.playerSize
				g.playerVelY = 0
				g.isJumping = false
			} else if g.playerY >= platform.y+platform.height-g.playerVelY {
				// Collision from the bottom
				g.playerY = platform.y + platform.height
				g.playerVelY = 0
			} else if g.playerX < platform.x {
				// Collision from the left
				g.playerX = platform.x - g.playerSize
			} else if g.playerX+g.playerSize > platform.x+platform.width {
				// Collision from the right
				g.playerX = platform.x + platform.width
			}
		}
	}

	// Check if the player has landed on the ground
	if g.playerY >= playerGroundY {
		g.playerY = playerGroundY
		g.playerVelY = 0
		g.isJumping = false
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw the ground
	groundColor := color.RGBA{255, 255, 255, 255}
	ebitenutil.DrawRect(screen, 0, playerGroundY+50, 800, 15, groundColor)

	// Draw the player as a square
	playerColor := color.RGBA{255, 255, 255, 255}
	ebitenutil.DrawRect(screen, g.playerX, g.playerY, g.playerSize, g.playerSize, playerColor)

	// Draw platforms
	platformColor := color.RGBA{255, 255, 255, 255}
	for _, platform := range g.platforms {
		ebitenutil.DrawRect(screen, platform.x, platform.y, platform.width, platform.height, platformColor)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 800
}

func main() {
	ebiten.SetWindowSize(800, 800)
	ebiten.SetWindowTitle("Platformer")
	game := &Game{
		playerX:    300,          // Initial X position of the player
		playerY:    playerStartY, // Initial Y position of the player
		playerSize: 50,           // Size of the player square
		playerVelY: 0,
		isJumping:  false,
		platforms: []Platform{
			{100, 450, 200, 20}, // Example platform
			{350, 350, 200, 20}, // Another platform
		},
	}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
