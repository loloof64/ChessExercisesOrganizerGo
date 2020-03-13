<template>
  <v-container fluid class="px-0">
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

    <v-row justify-center align-center class="ma-auto">
      <v-col>
        <loloof64-chessboard
          size="600"
          :promotion_dialog_title="$t('modals.promotion.title')"
          white_player_human="true"
          black_player_human="true"
          @checkmate="notifyCheckmate"
          @stalemate="notifyStalemate"
          @perpetual_draw="notifyPerpetualDraw"
          @missing_material_draw="notifyMissingMaterialDraw"
          @fifty_moves_draw="notifyFiftyMovesDraw"
          @waiting_manual_move="makeComputerMove"
          @move_done="addMoveToHistory"
        ></loloof64-chessboard>
      </v-col>

      <v-col>
        <MovesHistory
          id="history"
          :history="orderedHistory"
          @position_requested="setPosition($event)"
        ></MovesHistory>
      </v-col>
    </v-row>

  </v-container>
</template>

<script>
import MovesHistory from "./MovesHistory";

/*
    History should be something like (here simplified)
    [
        {
            moveNumber: 1,
            whiteTurn: true,
            moveFan: "e4"
        },
        {
            moveNumber: 1,
            whiteTurn: false,
            moveFan: "Nc6"
        },
        {
            moveNumber: 2,
            whiteTurn: true,
            moveFan: "Nf3"
        },
        {
            moveNumber: 2,
            whiteTurn: false,
            moveFan: "e5"
        },
    ]

    orderedHistory should be something like (here simplified)
    [
        {
          moveNumber: 1,
          white: {
            moveFan: "e4"
          },
          black: {
            moveFan: "Nc6"
          }
        },
        {
          moveNumber: 2,
          white: {
            moveFan: "Nf3"
          },
          black: {
            moveFan: "e5"
          }
        },
      ]
    */

export default {
  data() {
    return {
      message: " ",
      raised: true,
      errorDialog: false,
      errorDialogTitle: "",
      errorDialogText: "",
      okButtonText: "Ok",
      history: [],
      orderedHistory: []
    };
  },
  methods: {
    newGame: function() {
      const boardComponent = document.querySelector("loloof64-chessboard");
      boardComponent.newGame();

      this.history = [];
      this.updateOrderedHistory();
    },
    addMoveToHistory: function(event) {
      this.history = [...this.history, event.detail.moveObject];
      this.updateOrderedHistory();
    },
    notifyCheckmate: function(event) {
      const whiteCheckmated = event.detail.whiteTurnBeforeMove;
      const gameEndStatusMsg =
        this.$i18n.t('game.ended.checkmate', {
          winner: this.$i18n.t(whiteCheckmated ? 'game.side.white' : "game.side.black")});
      this.$emit('snackbar', gameEndStatusMsg);
    },
    notifyStalemate: function(event) {
      const gameEndStatusMsg = this.$i18n.t('game.ended.stalemate');
      this.$emit('snackbar', gameEndStatusMsg);
    },
    notifyPerpetualDraw: function(event) {
      const gameEndStatusMsg =  this.$i18n.t('game.ended.perpetualDraw');
      this.$emit('snackbar', gameEndStatusMsg);
    },
    notifyMissingMaterialDraw: function(event) {
      const gameEndStatusMsg =  this.$i18n.t('game.ended.missingMaterial');
      this.$emit('snackbar', gameEndStatusMsg);
    },
    notifyFiftyMovesDraw: function(event) {
      const gameEndStatusMsg =  this.$i18n.t('game.ended.fiftyMovesRule');
      this.$emit('snackbar', gameEndStatusMsg);
    },
    makeComputerMove: function() {
      const board = document.querySelector("loloof64-chessboard");
      const currentPosition = board.getCurrentPosition();
      window.backend.UciEngine.PlayPosition(currentPosition).then(bestMove => {
        if (bestMove === "#EngineNotSet") {
          this.errorDialogTitle = this.$i18n.t('modals.engineCouldNotPlay.engineNotSet.title');
          this.errorDialogText = this.$i18n.t('modals.engineCouldNotPlay.engineNotSet.text');
          this.errorDialog = true;
        } else if (bestMove === "#ComputationError") {
          this.errorDialogTitle = this.$i18n.t('modals.engineCouldNotPlay.computationError.title');
          this.errorDialogText =
            this.$i18n.t('modals.engineCouldNotPlay.computationError.text');
          this.errorDialog = true;
        } else {
          const moveObject = this.convertMoveStringToObject(bestMove);
        }
      });
    },
    convertMoveStringToObject: function(moveString) {
      const start = this.convertAlgebraicCellToCoordinates(
        moveString.substring(0, 2)
      );
      const end = this.convertAlgebraicCellToCoordinates(
        moveString.substring(2, 4)
      );

      const board = document.querySelector("loloof64-chessboard");

      if (moveString.length >= 5) {
        const promotion = moveString.substring(4, 5);
        board.playMove({
          startFile: start.file,
          startRank: start.rank,
          endFile: end.file,
          endRank: end.rank,
          promotion
        });
      } else {
        board.playMove({
          startFile: start.file,
          startRank: start.rank,
          endFile: end.file,
          endRank: end.rank
        });
      }
    },
    convertAlgebraicCellToCoordinates: function(cellStr) {
      const file = cellStr.charCodeAt(0) - "a".charCodeAt(0);
      const rank = cellStr.charCodeAt(1) - "1".charCodeAt(0);

      return { file, rank };
    },
    updateOrderedHistory() {
      let currentMoveNumber = undefined;
      let currentHistoryLine = undefined;
      let update = this.history.reduce((acc, curr) => {
        if (currentMoveNumber !== curr.moveNumber) {
          if (currentHistoryLine !== undefined) {
            acc.push(currentHistoryLine);
          }
          currentMoveNumber = curr.moveNumber;
          currentHistoryLine = { moveNumber: currentMoveNumber };
        }
        if (curr.whiteTurn) {
          currentHistoryLine["white"] = curr;
        } else {
          currentHistoryLine["black"] = curr;
        }

        return acc;
      }, []);
      // Don't forget to push the current edited element if any !
      if (currentHistoryLine !== undefined) {
        update.push(currentHistoryLine);
      }

      this.orderedHistory = update;
    },
    setPosition: function(evt) {
      const whitePlayer = evt.whitePlayer;
      const historyIndex = evt.historyIndex;
      const board = document.querySelector("loloof64-chessboard");
      const historyLine = this.orderedHistory[historyIndex];

      const positionObject = whitePlayer
        ? historyLine.white
        : historyLine.black;

      board.setPositionAndLastMove({...positionObject});
    }
  },
  components: {
    MovesHistory
  }
};
</script>

<style lang="css" scoped>
.v-btn {
  z-index: 10;
}

#history {
  overflow: scroll;
  height: 600px;
}
</style>
