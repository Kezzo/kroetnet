import { Player } from './player'

export class Game {
  players:Player[] = [];
  moves:any = []
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
  move(id: string, command: string, seq: number){
    for(var i=0; i<this.players.length;i++){
      if (id === this.players[i].id){

        if(this.players[i].lastSeq < seq){
          this.moves[this.players[i].id] = {
            "timestamp": Date.now(),
            "buffer" : [{"command": command, "sequence": seq}]
          };
        }else if (this.moves[this.players[i].id] !== undefined) {
          this.moves[this.players[i].id].buffer.sort(function(a, b) {
            return parseFloat(a.sequence) - parseFloat(b.sequence);
          });
          if(this.moves[this.players[i].id].timestamp - Date.now() > 100){
            this.moves[this.players[i].id].buffer.forEach((ele) =>{
              this.move(this.players[i].id, ele.command, ele.seq)
              this.players[i].lastSeq = ele.seq
            });
          } else {
            const ele = this.moves[this.players[i].id].buffer
            this.move(this.players[i].id, ele.command, ele.seq)
            this.players[i].lastSeq = ele.seq
          }
        }

        this.players[i].lastSeq = seq;
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

