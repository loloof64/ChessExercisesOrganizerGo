<template>
  <v-container fluid class="px-0">
    <v-row justify="center">
      <v-dialog v-model="errorDialog" persistent max-width="300">
        <v-card>
          <v-card-title class="headline">{{errorDialogTitle}}</v-card-title>
          <v-card-text>{{errorDialogText}}</v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="green darken-1" text @click="errorDialog = false">{{okButtonText}}</v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
  </v-row>
    <loloof64-chessboard 
      size="600"
      promotion_dialog_title="Select the promotion"
      white_player_human="true"
      black_player_human="false"
      @checkmate="notifyCheckmate"
      @stalemate="notifyStalemate"
      @perpetual_draw="notifyPerpetualDraw"
      @missing_material_draw="notifyMissingMaterialDraw"
      @fifty_moves_draw="notifyFiftyMovesDraw"
      @waiting_manual_move="makeComputerMove"
    ></loloof64-chessboard>

    <v-snackbar
      v-model="gameEndStatusSnackbar"
      :timeout="1500"
    >
      {{ gameEndStatusMsg }}
    </v-snackbar>
  </v-container>
</template>

<script>
  export default {
    data () {
      return {
        message: " ",
        raised: true,
        gameEndStatusSnackbar: false,
        gameEndStatusMsg: '',
        errorDialog: false,
        errorDialogTitle: '',
        errorDialogText: '',
        okButtonText: 'Ok',
      }
    },
    methods: {
      notifyCheckmate: function(event) {
        const whiteCheckmated = event.detail.whiteTurnBeforeMove;
        this.gameEndStatusMsg = (whiteCheckmated ? 'White' : 'Black') + ' won by checkmate';
        this.gameEndStatusSnackbar = true;
      },
      notifyStalemate: function(event) {
        this.gameEndStatusMsg = 'Stalemate';
        this.gameEndStatusSnackbar = true;
      },
      notifyPerpetualDraw: function(event) {
        this.gameEndStatusMsg = 'Perpetual draw';
        this.gameEndStatusSnackbar = true;
      },
      notifyMissingMaterialDraw: function(event) {
        this.gameEndStatusMsg = 'Missing material draw';
        this.gameEndStatusSnackbar = true;
      },
      notifyFiftyMovesDraw: function(event) {
        this.gameEndStatusMsg = 'Fifty moves draw';
        this.gameEndStatusSnackbar = true;
      },
      makeComputerMove: function() {
        const board = document.querySelector('loloof64-chessboard');
        const currentPosition = board.getCurrentPosition();
        window.backend.UciEngine.PlayPosition(currentPosition).then(bestMove => {
          if (bestMove === '#EngineNotSet'){
            this.errorDialogTitle = 'Could not make then engine play';
            this.errorDialogText = 'You forgot to set up an UCI engine first.';
            this.errorDialog = true;
          }
          else if (bestMove === '#ComputationError') {
            this.errorDialogTitle = 'Could not make then engine play';
            this.errorDialogText = 'An misc. error has occured : is the given position legal and did you really select an uci engine ?';
            this.errorDialog = true;
          }
          else {
            const moveObject = this.convertMoveStringToObject(bestMove);
          }
        });
      },
      convertMoveStringToObject: function(moveString) {
        const start = this.convertAlgebraicCellToCoordinates(moveString.substring(0,2));
        const end = this.convertAlgebraicCellToCoordinates(moveString.substring(2,4));

        const board = document.querySelector('loloof64-chessboard');
          
        if (moveString.length >= 5) {
            const promotion = moveString.substring(4,5);
            board.playMove({
              startFile: start.file,
              startRank: start.rank,
              endFile: end.file,
              endRank: end.rank,
              promotion,
            });
        }
        else {
            board.playMove({
              startFile: start.file,
              startRank: start.rank,
              endFile: end.file,
              endRank: end.rank,
            });
        }
      },
      convertAlgebraicCellToCoordinates: function(cellStr) {
        const file = cellStr.charCodeAt(0) - 'a'.charCodeAt(0);
        const rank = cellStr.charCodeAt(1) - '1'.charCodeAt(0);

        return {file, rank};
      }
    },
  }
</script>

<style lang="css" scoped>
.v-btn {
  z-index: 10;
}
</style>
