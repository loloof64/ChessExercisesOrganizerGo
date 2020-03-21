<template>
  <v-app id="inspire">
    <v-app-bar app fixed clipped-left>
      <v-toolbar-title>{{$t('app.title')}}</v-toolbar-title>
      <ToolbarButton :text="$t('menu.newGame.tooltip')" :action="newGameRequest" :darkMode="false">
        <v-icon>mdi-restart</v-icon>
      </ToolbarButton>
      <ToolbarButton :text="$t('menu.toggleSide.tooltip')" :action="toggleSide">
        <v-icon>mdi-arrow-up-down</v-icon>
      </ToolbarButton>
      <ToolbarButton :text="$t('menu.stopGame.tooltip')" :action="stopGameRequest">
        <v-icon>mdi-stop-circle</v-icon>
      </ToolbarButton>
      <ToolbarButton :text="$t('menu.save.tooltip')" :action="showSaveDialog">
        <v-icon>mdi-content-save</v-icon>
      </ToolbarButton>
      <ToolbarButton :text="$t('menu.settings.tooltip')" :action="showSettingsDialog">
        <v-icon>mdi-settings</v-icon>
      </ToolbarButton>
    </v-app-bar>

    <v-content>
      <v-container fluid class="px-0">
        <v-layout justify-center align-center class="px-0">
          <game-page
            ref="gameZone"
            @snackbar="showSnackbar"
            :boardBackground="boardBackgroundColor"
            :coordinates_color="boardCoordintatesColor"
            :move_highlight_color="boardLastMoveArrowColor"
            :white_cell_color="boardWhiteCellsColor"
            :black_cell_color="boardBlackCellsColor"
            :origin_cell_color="boardDndStartColor"
            :target_cell_color="boardDndEndColor"
            :dnd_cross_color="boardDndCrossColor"
            :coordinates_visible="boardCoordinatesVisible"
            :last_move_visible="boardLastMoveArrowVisible"
          ></game-page>
        </v-layout>

        <SettingsDialog ref="settingsDialog" @configurationUpdated="updateBoardAndEngine"></SettingsDialog>

        <SimpleModalDialog ref="errorDialog" :title="errorDialogTitle">
          <v-card-text>{{errorDialogText}}</v-card-text>
        </SimpleModalDialog>

        <SimpleModalDialog ref="saveSuccessfulDialog" :title="$t('modals.saveSuccess.title')">
          <v-card-text>{{$t('modals.saveSuccess.text')}}</v-card-text>
        </SimpleModalDialog>

        <SimpleModalDialog
          ref="newGameConfirmation"
          :title="$t('modals.newGame.title')"
          :confirmAction="doStartNewGame"
          cancelButton
        >
          <v-card-text>{{$t('modals.newGame.text')}}</v-card-text>
        </SimpleModalDialog>

        <SimpleModalDialog
          ref="stopGameConfirmation"
          :title="$t('modals.stopGame.title')"
          :confirmAction="doStopGame"
          cancelButton
        >
          <v-card-text>{{$t('modals.stopGame.text')}}</v-card-text>
        </SimpleModalDialog>

        <SimpleModalDialog
          ref="pgnSelectionDialog"
          id="pgnSelectionDialog"
          :title="$t('modals.pgnSelection.title')"
          :confirmAction="loadSelectedPgn"
          cancelButton
        >
          <loloof64-chessboard
            ref="previewBoard"
            size="300"
            :white_player_human="false"
            :black_player_human="false"

            :left="10"
            :top="10"

            :reversed="previewBoardReversed"
            :background="boardBackgroundColor"
            :coordinates_color="coordinates_color"
            :white_cell_color="white_cell_color"
            :black_cell_color="black_cell_color"

            :coordinates_visible="true"
          ></loloof64-chessboard>

          <div id="previewPanel">
            <div id="previewControls">
              <v-tooltip bottom>
                <template v-slot:activator="{ on }">
                    <v-btn v-on="on" @click="goPreviousPreview()" class="blue red--text"><v-icon>mdi-chevron-left</v-icon></v-btn>
                </template>
                <span>{{ $t('pgn.buttons.previous') }}</span>
              </v-tooltip>

              <v-tooltip bottom>
                <template v-slot:activator="{ on }">
                    <v-btn v-on="on" @click="goNextPreview()" class="blue red--text"><v-icon>mdi-chevron-right</v-icon></v-btn>
                </template>
                <span>{{ $t('pgn.buttons.next') }}</span>
              </v-tooltip>

              <span>
                {{selectedPgnIndex + 1}} / {{pgnArrays.length}}
              </span>
            </div>

            <div id="previewPlayers">
              {{ previewPgnWhitePlayer }} | {{ previewPgnBlackPlayer }}
            </div>

            <div id="previewEvent">
              {{ selectedPgnSite }} {{ previewPgnDate }}
            </div>
          </div>
        </SimpleModalDialog>

        <SimpleSnackBar ref="snackbar">{{ snackBarMessage }}</SimpleSnackBar>
      </v-container>
    </v-content>

    <v-footer app fixed>
      <span style="margin-left:1em">&copy; Laurent Bernab&eacute; - 2020</span>
    </v-footer>
  </v-app>
</template>

<script>
import GamePage from "./components/GamePage";
import ToolbarButton from "./components/ToolbarButton";
import SimpleModalDialog from "./components/SimpleModalDialog";
import SimpleSnackBar from "./components/SimpleSnackBar";
import SettingsDialog from "./components/SettingsDialog";

import Chess from "chess.js";
import pgnReader from "./libs/pgn";

export default {
  data: () => ({
    errorDialogTitle: "",
    errorDialogText: "",
    snackBarMessage: "",
    settings: undefined,
    boardBackgroundColor: undefined,
    boardCoordintatesColor: undefined,
    boardLastMoveArrowColor: undefined,
    boardWhiteCellsColor: undefined,
    boardBlackCellsColor: undefined,
    boardDndStartColor: undefined,
    boardDndEndColor: undefined,
    boardDndCrossColor: undefined,
    boardCoordinatesVisible: undefined,
    boardLastMoveArrowVisible: undefined,
    previewBoardReversed: false,
    pgnArrays: [],
    selectedPgnIndex: -1,
    selectedPgnWhite: '',
    selectedPgnBlack: '',
    selectedPgnSite: '',
    selectedPgnDate: '',
  }),
  mounted() {
    this.$i18n.locale = navigator.language.substring(0, 2);
    window.backend.SettingsManager.Load().then(content => {
      if (content === "#ConversionToJsonError") {
        this.settings = undefined;
      } else {
        let parsedContent = content;
        while (typeof parsedContent === "string")
          parsedContent = JSON.parse(parsedContent);
        this.settings = parsedContent;
        this.updateBoardAndEngine(parsedContent);
      }
    });
  },
  methods: {
    newGameRequest: function() {
      const chessBoard = document.querySelector("loloof64-chessboard");
      if (chessBoard.gameIsInProgress()) {
        this.$refs["newGameConfirmation"].open();
      } else {
        this.doStartNewGame();
      }
    },
    doStartNewGame: function() {
      const fileName = "lesson01";
      const filePath =
        "/home/laurent-bernabe/Documents/temp/pgn/" + fileName + ".pgn";
      window.backend.TextFileManager.GetTextFileContentWithPathProviden(
        filePath
      ).then(content => {
        if (content === "#ErrorReadingFile") {
          this.errorDialogTitle = this.$i18n.t("modals.failedToReadPgn.title");
          this.errorDialogText = this.$i18n.t("modals.failedToReadPgn.text");
          this.$refs["errorDialog"].open();
          return;
        }

        const allPgns = this.splitPgn(content);

        if (allPgns.length === 0) {
          this.errorDialogTitle = this.$i18n.t("modals.noGameInPgn.title");
          this.errorDialogText = this.$i18n.t("modals.noGameInPgn.text");
          this.$refs["errorDialog"].open();
          return;
        }

        this.pgnArrays = allPgns;
        this.letUserSelectPgn();
      });
    },
    stopGameRequest: function() {
      const chessBoard = document.querySelector("loloof64-chessboard");
      if (chessBoard.gameIsInProgress()) {
        this.$refs["stopGameConfirmation"].open();
      }
    },
    doStopGame: function() {
      const chessBoard = document.querySelector("loloof64-chessboard");
      chessBoard.stop();
      this.$refs["snackbar"].open(this.$i18n.t("game.status.stopped"));
      this.$refs["gameZone"].$refs["history"].selectLastMove();
    },
    toggleSide: function() {
      this.$refs["gameZone"].toggleSide();
    },
    showSettingsDialog: function() {
      this.$refs["settingsDialog"].open(this.settings);
    },
    showSnackbar: function(event) {
      this.$refs["snackbar"].open(event);
    },
    showSaveDialog: function() {
      const chessBoard = document.querySelector("loloof64-chessboard");
      if (chessBoard.gameIsInProgress()) return;

      const historyComponent = this.$refs["gameZone"].$refs["history"];
      if (!historyComponent.hasData()) return;

      const playerHasWhite = this.$refs["gameZone"].playerHasWhite();

      const whiteName = playerHasWhite
        ? this.settings.PlayerName
        : this.settings.ComputerName;
      const blackName = playerHasWhite
        ? this.settings.ComputerName
        : this.settings.PlayerName;
      const pgn = chessBoard.gamePgn({ whiteName, blackName });
      const fileName = "PgnGenere";
      const filePath =
        "/home/laurent-bernabe/Documents/temp/pgn/" + fileName + ".pgn";

      // Production mode, use window.backend.TextFileManager.SaveTextFile()
      window.backend.TextFileManager.SaveTextFileWithPathProviden(
        filePath,
        pgn
      ).then(error => {
        if (error === "#ErrorSavingFile") {
          this.errorDialogTitle = this.$i18n.t("modals.saveError.title");
          this.errorDialogText = this.$i18n.t("modals.saveError.text");
          this.$refs["errorDialog"].open();
        } else if (error === "#ErrorClosingFile") {
          this.errorDialogTitle = this.$i18n.t("modals.saveClosingError.title");
          this.errorDialogText = this.$i18n.t("modals.saveClosingError.text");
          this.$refs["errorDialog"].open();
        } else {
          this.$refs["saveSuccessfulDialog"].open();
        }
      });
    },
    updateBoardAndEngine: function(configuration) {
      // Apply modifications on board.
      this.boardBackgroundColor = configuration.BoardBackgroundColor;
      this.boardCoordintatesColor = configuration.BoardCoordinatesColor;
      this.boardLastMoveArrowColor = configuration.BoardLastMoveArrowColor;
      this.boardWhiteCellsColor = configuration.BoardWhiteCellsColor;
      this.boardBlackCellsColor = configuration.BoardBlackCellsColor;
      this.boardDndStartColor = configuration.BoardDndStartColor;
      this.boardDndEndColor = configuration.BoardDndEndColor;
      this.boardDndCrossColor = configuration.BoardDndCrossColor;
      this.boardCoordinatesVisible = [true, "true"].includes(
        configuration.CoordinatesVisible
      );
      this.boardLastMoveArrowVisible = [true, "true"].includes(
        configuration.LastMoveArrowVisible
      );

      // Update used engine
      window.backend.UciEngine.LoadEngineWithPathProviden(
        configuration.EnginePath
      ).then(error => {
        if (error === "#ConfigEngineErr") {
          this.errorDialogTitle = this.$i18n.t(
            "modals.settings.failedToSetupEngineTitle"
          );
          this.errorDialogText = this.$i18n.t(
            "modals.settings.failedToSetupEngineText"
          );
          this.$refs["errorDialog"].open();
        }
      });
    },
    splitPgn: function(fileContent) {
      let results = [];

      const isOutOfGame = 0;
      const isInHeader = 1;
      const isInGame = 2;

      let status = isOutOfGame;
      let currentGame = "";

      const isAHeaderLine = l => l.startsWith("[");
      const isBlankLine = l => l.trim().length === 0;

      const lines = fileContent.split(/\r?\n/);
      lines.forEach(l => {
        switch (status) {
          case isOutOfGame:
            if (isAHeaderLine(l)) {
              status = isInHeader;
              currentGame += l + "\n";
            } else if ( ! isBlankLine(l) ) {
              status = isInGame;
              currentGame += l + "\n";
            }

            break;

          case isInHeader:
            if (status === isInGame) {
              results.push(currentGame);
              currentGame = "";
            }
            if (isAHeaderLine(l)) {
              currentGame += l + "\n";
            }
            else if (isBlankLine(l)) {
              status = isInGame;
              currentGame += "\n";
            } else {
              currentGame += l + "\n";
             } 

            break;
          
          case isInGame:
            if (isAHeaderLine(l)) {
              results.push(currentGame);
              currentGame = "";
              status = isInHeader;
              currentGame += l + "\n";
            }
            else if (isBlankLine(l)) {
              results.push(currentGame);
              currentGame = "";
              status = isOutOfGame;
            }
            else {
              currentGame += l + "\n";
            }
        }
      });

      if (currentGame.length > 0) results.push(currentGame);

      return results;
    },
    letUserSelectPgn: function() {
      this.selectedPgnIndex = 0;
      this.$refs['pgnSelectionDialog'].open();
      setTimeout(() => {
        this.previewPgn();
      }, 80);
    },
    previewPgn: function() {
      if (this.selectedPgnIndex > this.pgnArrays.length) return;

      const selectedPgn = this.pgnArrays[this.selectedPgnIndex];
      const loader = new pgnReader({pgn: selectedPgn});
      const result = loader.load_pgn();
      const finalPosition = result.finalPosition;

      this.selectedPgnWhite = result.headers.White || '';
      this.selectedPgnBlack = result.headers.Black || '';
      this.selectedPgnSite = result.headers.Site || '';
      this.selectedPgnDate = result.headers.Date || '';

      const previewComponent = this.$refs['previewBoard'];
      previewComponent.newGame(finalPosition);
    },
    loadSelectedPgn: function() {
      const selectedPgn = this.pgnArrays[this.selectedPgnIndex];
      const loader = new pgnReader({pgn: selectedPgn});
      const result = loader.load_pgn();
      const finalPosition = result.finalPosition;

      try {
          const playerHasWhite = finalPosition.split(' ')[1] === 'w';

          this.$refs["gameZone"].newGame(finalPosition, playerHasWhite);
          const chessBoard = document.querySelector("loloof64-chessboard");
          if (chessBoard.gameIsInProgress()) {
            this.$refs["snackbar"].open(this.$i18n.t("game.status.started"));
          } else {
            this.$refs["snackbar"].open(
              this.$i18n.t("game.status.already_finished")
            );
          }
        } catch (error) {
          console.error(error);
          this.errorDialogTitle = this.$i18n.t("modals.failedToReadPgn.title");
          this.errorDialogText = this.$i18n.t("modals.failedToReadPgn.text");
          this.$refs["errorDialog"].open();
        }
    },
    goPreviousPreview: function() {
      if (this.selectedPgnIndex > 0) {
        this.selectedPgnIndex -= 1;
        this.previewPgn();
      }
    },
    goNextPreview: function() {
      if (this.selectedPgnIndex < this.pgnArrays.length - 1) {
        this.selectedPgnIndex += 1;
        this.previewPgn();
      }
    }
  },
  computed: {
    previewPgnWhitePlayer: function() {
      return this.selectedPgnWhite.length > 0 ? this.selectedPgnWhite : '???';
    },
    previewPgnBlackPlayer: function() {
      return this.selectedPgnBlack.length > 0 ? this.selectedPgnBlack : '???';
    },
    previewPgnDate: function() {
      return this.selectedPgnDate.length > 0 ? "(" + this.selectedPgnDate + ")" : '';
    }
  },
  components: {
    GamePage,
    ToolbarButton,
    SimpleModalDialog,
    SimpleSnackBar,
    SettingsDialog
  },
  props: {
    source: String
  }
};
</script>

<style scoped>
#pgnSelectionDialog {
  position: relative;
}

#previewBoard {
  position: absolute;
  left: 10px;
}

#previewPanel {
  position: absolute;
  left: 350px;
  top: 60px;
}

#previewControls > .v-btn {
  margin-left: 10px;
  margin-right: 10px;
}

#previewPlayers {
  margin-left: 20px;
  margin-top: 20px;
}

#previewEvent {
  margin-left: 20px;
  margin-top: 20px;
}
</style>