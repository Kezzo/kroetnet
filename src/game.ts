import { Player } from './player'

export class Game {
  players:Player[] = [];
  constructor(){}
  addPlayer(player: Player){
    // this.players.forEach( (ply) => {
    //   if (player.id === ply.id){
    //     ply.posX = 0
    //     ply.posY = 0
    //   } else {
        this.players = this.players.concat(player)
      // }
    // });
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

