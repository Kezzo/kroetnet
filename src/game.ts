import { Player } from './player'

export class Game {
  players:Player[] = [];
  constructor(){}
  addPlayer(player: Player){
    for(var i=0; i < this.players.length;i++){
      if (player.id === this.players[i].id){
        this.players[i].posX = 0
        this.players[i].posY = 0
      } else {
        this.players = this.players.concat(player)
      }
    }
  }
  move(id: string, command: string){
    for(var i=0; i<this.players.length;i++){
      if (id === this.players[i].id){
        switch(command){
          case "UP": {
            return this.players[i].up()
          }
          case "DOWN": {
            return this.players[i].down()
          };
          case "LEFT": {
            return this.players[i].left()
          };
          case "RIGHT": {
            return this.players[i].right()
          };
        };
      }
    }
  }
}

