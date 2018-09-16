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
      // player in game
      if (id === this.players[i].id){
        // if there are buffered moves
        if (this.moves[this.players[i].id] &&
          this.moves[this.players[i].id].buffer.length > 0) {
          this.moves[this.players[i].id].buffer.sort(function(a, b) {
            return parseFloat(a.sequence) - parseFloat(b.sequence);
          });

          console.log('DATE DIFF', this.moves[this.players[i].id].timestamp, Date.now())
          // execute all moves older than
          if(this.moves[this.players[i].id].timestamp - Date.now() > 10000){
            this.moves[this.players[i].id].buffer.forEach((ele) =>{
              this.move(this.players[i].id, ele.command, ele.seq)
              this.moves[this.players[i].id].buffer.pop()
              return;
            });
          }else {
            // execute just the next move
            const nextMove = this.moves[this.players[i].id].buffer.forEach((ele) =>{
              console.log('CHECK ELE: ', ele.sequence, seq)
              if(ele.sequence == seq++){
                this.move(this.players[i].id, nextMove.command, nextMove.sequence)
                this.moves[this.players[i].id].buffer.pop()
                return;
              }
            })
          }
          // if there are no buffered moves check if current seq is the next one
        } else if(this.players[i].lastSeq !== seq-1){
          // add new buffer because its not the next move
          if (!this.moves[this.players[i].id]) {
            this.moves[this.players[i].id] = {
              "timestamp": Date.now(),
              "buffer" : [{"command": command, "sequence": seq}]
            };
            console.log('NEW BUFFER')
          } else {
            this.moves[this.players[i].id].buffer = 
              this.moves[this.players[i].id].buffer.concat({
                "command": command, "sequence": seq 
              })
            console.log('ADD TO BUFFER')
          }
          console.log('BUFFER ', this.moves[this.players[i].id])

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

