package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"math/rand"
	"time"
)

type SnakeSquare struct {
	square *rl.Rectangle
	next   *SnakeSquare
}

type Snake struct {
	head *SnakeSquare
}

func (s *Snake) CollisionCheck() bool {
	head := s.head.square

	current := s.head.next
	for current != nil {
		didHit := head.X == current.square.X && head.Y == current.square.Y
		if didHit {
			return true
		}
		current = current.next
	}

	return false
}

func (s *Snake) Draw() {
	current := s.head
	for current != nil {
		rl.DrawRectangleRec(*current.square, rl.Red)
		current = current.next
	}
}

func (s *Snake) Move(direction rl.Vector2) {
	currentLastPosition := rl.NewVector2(s.head.square.X, s.head.square.Y)
	s.head.square.X += direction.X * 10
	s.head.square.Y += direction.Y * 10
	current := s.head.next

	for current != nil {
		nextValue := currentLastPosition
		currentLastPosition = rl.NewVector2(current.square.X, current.square.Y)
		current.square.X = nextValue.X
		current.square.Y = nextValue.Y
		current = current.next
	}

}

func (s *Snake) Insert(direction rl.Vector2) {
	newSquare := &SnakeSquare{
		square: &rl.Rectangle{X: 50, Y: 50, Width: 10, Height: 10},
		next:   nil,
	}

	if s.head == nil {
		s.head = newSquare
		return
	}

	current := s.head
	for current != nil {
		if current.next == nil {
			break
		}
		current = current.next
	}

	newSquare.square.X = current.square.X + (direction.X * 12)
	newSquare.square.Y = current.square.Y + (direction.Y * 12)
	current.next = newSquare
}

func main() {
	rl.InitWindow(300, 300, "Snake")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	direction := rl.Vector2{X: 1, Y: 0}
	speed := float32(10)

	snake := Snake{}
	snake.Insert(direction)
	snake.Insert(direction)
	snake.Insert(direction)

	moveInterval := time.Millisecond * 100
	lastMoveTime := time.Now()
	rec := rl.NewRectangle(250, 50, 10, 10)
	hasCollided := false
	var lastDirectionChange time.Time

	for !rl.WindowShouldClose() {
		canMove := time.Since(lastMoveTime) >= moveInterval
		if canMove {
			snake.Move(direction)
			lastMoveTime = time.Now()
		}

		if time.Since(lastDirectionChange) >= moveInterval {
			if rl.IsKeyPressed(rl.KeyA) && direction.X != 1 {
				direction.X = -1
				direction.Y = 0
				lastDirectionChange = time.Now()
			} else if rl.IsKeyPressed(rl.KeyD) && direction.X != -1 {
				direction.X = 1
				direction.Y = 0
				lastDirectionChange = time.Now()
			} else if rl.IsKeyPressed(rl.KeyW) && direction.Y != 1 {
				direction.Y = -1
				direction.X = 0
				lastDirectionChange = time.Now()
			} else if rl.IsKeyPressed(rl.KeyS) && direction.Y != -1 {
				direction.Y = 1
				direction.X = 0
				lastDirectionChange = time.Now()
			}
		}

		if rl.CheckCollisionRecs(rec, *snake.head.square) && !hasCollided {
			hasCollided = true
		}

		fmt.Printf("%v\n", snake.head.square.X)
		isGameOver := snake.head.square.X+snake.head.square.Width > 300 ||
			snake.head.square.X < 0 ||
			snake.head.square.Y+snake.head.square.Height > 300 || snake.head.square.Y < 0

		if isGameOver || snake.CollisionCheck() && !hasCollided {
			panic("Game over")
		}

		if hasCollided {
			rec.X = float32(rand.Intn(280) + 10)
			rec.Y = float32(rand.Intn(280) + 10)
			snake.Insert(direction)
			hasCollided = false
			speed += 5
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		snake.Draw()
		rl.DrawRectangleRec(rec, rl.White)

		rl.EndDrawing()
	}
}
