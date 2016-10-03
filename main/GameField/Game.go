package GameField

import (
	tl "github.com/JoelOtter/termloop"
	//"math"
	"strconv"
)

type GameField struct{
	Field [8][8]int
	CurrentPlayer *Player
	player1_scores int
	player2_scores int

	player1_step bool
	velocities [8]velocity
	paramx int
	paramy int
	paramas int
	text *tl.Text
}

type velocity struct {
	x, y int
}
func New() GameField{
	player := Player{
		Entity: tl.NewEntity(1, 1, 1, 1),
	}

	player.Entity.SetCell(0, 0, &tl.Cell{Fg: tl.ColorGreen, Ch: '■'})
	player.playerNumber = 1
	game := GameField{
		CurrentPlayer: &player,
	}
	game.Field[3][3] = 1
	game.Field[4][4] = 1
	game.Field[4][3] = 2
	game.Field[3][4] = 2
	game.velocities = [8]velocity{
		velocity{x:0, y:-1,},
		velocity{x:1, y:-1,},
		velocity{x:1, y:0,},
		velocity{x:1, y:1,},
		velocity{x:0, y:1,},
		velocity{x:-1, y:1,},
		velocity{x:-1, y:0,},
		velocity{x:-1, y:-1,},
	}
	game.text = tl.NewText(13,2,"", tl.ColorWhite, 0)
	return game
}

func (field *GameField) Draw(screen *tl.Screen){
	for i:=0; i<len(field.Field); i++ {
		for j:=0; j<len(field.Field[0]); j++ {

			if field.Field[i][j] == 0{
				screen.RenderCell(i, j, &tl.Cell{Fg: tl.ColorRed, Ch: '+'})
			}
			if field.Field[i][j] == 1{
				screen.RenderCell(i, j, &tl.Cell{
					Fg: tl.ColorGreen,
					Ch: '█',
				})
			}
			if field.Field[i][j] == 2{
				screen.RenderCell(i, j, &tl.Cell{
					Fg: tl.ColorWhite,
					Ch: '█',
				})
			}

		}
	}
	field.CurrentPlayer.Entity.Draw(screen)
	field.text.SetText("green scores: " + strconv.Itoa(field.player1_scores)+ " white scores:" + strconv.Itoa(field.player2_scores))
	field.text.Draw(screen)
}



func (field *GameField) Tick(event tl.Event){
	if event.Type == tl.EventKey { // Is it a keyboard event?
		x, y := field.CurrentPlayer.Entity.Position()
		switch event.Key { // If so, switch on the pressed key.
		case tl.KeyArrowRight:
			if field.isInField(x+1, y) {
				field.CurrentPlayer.Entity.SetPosition(x + 1, y)
			}
		case tl.KeyArrowLeft:
			if field.isInField(x-1, y) {
				field.CurrentPlayer.Entity.SetPosition(x - 1, y)
			}
		case tl.KeyArrowUp:
			if field.isInField(x, y-1) {
				field.CurrentPlayer.Entity.SetPosition(x, y - 1)
			}
		case tl.KeyArrowDown:
			if field.isInField(x, y+1) {
				field.CurrentPlayer.Entity.SetPosition(x, y + 1)
			}
		case tl.KeyEnter:
			field.paramx = x
			field.paramy = y
			score := field.makeStep(x,y)
			if score!=0{
				field.Field[x][y] = field.CurrentPlayer.playerNumber
				field.recalculation()
				field.changeStep()
			}
		}
	}
}

func (field *GameField) changeStep(){
	if field.player1_step{
		field.CurrentPlayer.Entity.SetCell(0, 0, &tl.Cell{Fg: tl.ColorGreen})
		field.CurrentPlayer.playerNumber = 1
	}else{
		field.CurrentPlayer.Entity.SetCell(0, 0, &tl.Cell{Fg: tl.ColorWhite})
		field.CurrentPlayer.playerNumber = 2
	}
	field.player1_step = !field.player1_step
}

func (field *GameField) recalculation(){
	field.player1_scores = 0
	field.player2_scores = 0
	for i:=0; i<len(field.Field); i++ {
		for j:=0; j<len(field.Field[0]); j++{
			if field.Field[i][j] == 1 {
				field.player1_scores++
			}
			if field.Field[i][j] == 2{
				field.player2_scores++
			}
		}
	}
}

func (field *GameField) isInField(x, y int) bool{
	if x<0 || x>=len(field.Field) || y<0 || y>=len(field.Field[0]){
		return false
	}
	return true
}

func (field *GameField) isStepIsValid(x, y, player int){

}
func (field *GameField) isAround(x, y, player int) bool {
	if field.Field[x+1][y+1] == player || field.Field[x-1][y-1] == player {
		return true
	}
	return  false
}

func (field *GameField) makeStep(x, y int) int{

	playerNumb := field.CurrentPlayer.playerNumber
	enemyNumb := 2
	if playerNumb == 1{
		enemyNumb = 2
	}else{
		enemyNumb = 1
	}
	scores := 0
	current_x := x
	current_y := y
	isLine := false
	if field.Field[x][y] == 1 || field.Field[x][y] == 2{
		return scores
	}
	for cnt:=0; cnt<len(field.velocities); 	cnt++ {
		current_x = x
		current_y = y
		for i:=0; i<8 ;i++ {
			current_x += field.velocities[cnt].x
			current_y += field.velocities[cnt].y
			if !field.isInField(current_x, current_y){
				break
			}
			if field.Field[current_x][current_y] == 0 {
				field.paramas++
				break
			}
			if field.Field[current_x][current_y] == enemyNumb {
				isLine = true
				continue
			}
			if field.Field[current_x][current_y] == playerNumb && !isLine {
				break
			} else if isLine {
				for i > 0 {
					current_x -= field.velocities[cnt].x
					current_y -= field.velocities[cnt].y
					field.paramx=field.velocities[cnt].x
					field.paramy=field.velocities[cnt].y
					scores += 1
					field.Field[current_x][current_y] = playerNumb
					i--
				}
				break
			}
		}

	}
	return scores

	//for cnt:=0; cnt<len(field.velocities); cnt++ {
	//	var length int
	//	if field.velocities[cnt].x!=0{
	//		length = len(field.Field) + (field.velocities[cnt].x*x)-1
	//	}
	//	if field.velocities[cnt].y!=0{
	//		length = len(field.Field[0]) + (field.velocities[cnt].y*y)-1
	//	}
	//	current_x := x
	//	current_y := y
	//	isLine := false
	//	for i:=0; i<length; i++ {
	//		current_x += field.velocities[cnt].x
	//		current_y += field.velocities[cnt].y
	//		if !field.isInField(current_x, current_y) || field.Field[current_x][current_y]==0 {
	//			break
	//		}
	//
	//		if field.Field[current_x][current_y] == playerNumb && !isLine {
	//			break
	//		}else if isLine{
	//			for j:=0; j<int(math.Abs(float64(current_x-x))); j++ {
	//				current_x-=field.velocities[cnt].x
	//				current_y-=field.velocities[cnt].y
	//				field.Field[current_x][current_y] = playerNumb
	//
	//				scores+=1
	//			}
	//		}
	//		if field.Field[current_x][current_y] == enemyNumb {
	//			isLine=true
	//		}
	//
	//	}
	//}
	//return scores
}

func (field *GameField) isOutOfRange(tmp int) bool{
	if tmp<0 || tmp>=len(field.Field){
		return true
	}
	return false
}