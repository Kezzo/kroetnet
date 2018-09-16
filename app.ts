import PubNub = require('pubnub')
import fs = require('fs')
import { Player } from './src/player'
import { Game } from './src/game'

const pubnubConfigs = JSON.parse(fs.readFileSync('./pubnub-config.json', 'utf-8'))
const pubnub = new PubNub({
  subscribeKey: pubnubConfigs.subscribeKey,
  publishKey: pubnubConfigs.publishKey,
  secretKey: pubnubConfigs.secretKey,
  uuid: 'server',
  ssl: true
})

pubnub.subscribe({
    channels: ['world'],
});


function publishMsg(msg :any){
  pubnub.publish(
    {
      message: JSON.stringify(msg),
      channel: 'world',
      sendByPost: false, // true to send via post
      storeInHistory: false, //override default storage options
      meta: {
        "cool": "meta"
      }   // publish extra meta with the request
    },
    function (status, response) {
      if (status.error) {
        console.log(status)
      } else {
        console.log("message Published w/ timetoken", response.timetoken)
      }
    }
  );
}


//////////////////
//  RUN SERVER  //
//////////////////

let game = new Game()

pubnub.addListener({
  status: function(statusEvent) {
    if (statusEvent.category === "PNConnectedCategory") {
      var payload = {
        my: 'payload'
      };
      pubnub.publish(
        {
          message: payload
        },
        function (status) {
          console.log('Status: ', status)
        }
      );
    }
  },
  message: function(message) {
    console.log('MSG', message)
    if(message.message.command === 'JOIN'){
      let player = new Player(0,0, message.message.playerId)
      game.addPlayer(player)
    } else if(['DOWN', 'UP','LEFT', 'RIGHT'].indexOf(message.message.command) > -1){
      const newPosition = game.move(
        message.message.playerId,
        message.message.command,
        message.message.sequence
      )
      newPosition.sequence = message.message.sequence
      const msg = { "playerPositionUpdates": [ newPosition ] }
      publishMsg(msg)
    }
    console.log('Array of Players', game.players)
  },
  presence: function(presenceEvent) {
    console.log('Event: ', presenceEvent)
  }
})

