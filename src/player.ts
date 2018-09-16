const MAXY = 24
const MAXX = 24

export class Player {
  posX: number;
  posY: number;
  id: string;
  constructor(positionX: number, positionY: number, playerId: string){
    this.posX = positionX;
    this.posY = positionY;
    this.id = playerId;
  }
  private newPositon(){
    return {
      "playerId": this.id,
      "xPosition": this.posX,
      "zPosition": this.posY,
      "counter": 0
    }
  }
  up(){
    if(this.posY + 1 <= 24){
     this.posY+= 1 
    }
    return this.newPositon()
  }
  down(){
    if(this.posY - 1 >= -24){
      this.posY-= 1
    }
    return this.newPositon()
  }
  left(){
    if(this.posX - 1 >= -24){
      this.posX-= 1
    }
    return this.newPositon()
  }
  right(){
    if(this.posX + 1 <= 24){
      this.posX+= 1
    }
    return this.newPositon()
  }
}
//
// "playerPositionUpdates": [
//     {
//       "playerId": "43563456",
//       "xPosition": "234",
//       "zPosition": "234"
//     }
//   ]
