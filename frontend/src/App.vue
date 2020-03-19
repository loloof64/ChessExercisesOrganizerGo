<template>
  <v-app id="inspire">
    <v-app-bar app fixed clipped-left>
      <v-toolbar-title>{{$t('app.title')}}</v-toolbar-title>
      <ToolbarButton :text="$t('menu.newGame.tooltip')" :action="newGameRequest" :darkMode="false"><v-icon>mdi-restart</v-icon></ToolbarButton>
      <ToolbarButton :text="$t('menu.toggleSide.tooltip')" :action="toggleSide"><v-icon>mdi-arrow-up-down</v-icon></ToolbarButton>
      <ToolbarButton :text="$t('menu.stopGame.tooltip')" :action="stopGameRequest"><v-icon>mdi-stop-circle</v-icon></ToolbarButton>
      <ToolbarButton :text="$t('menu.save.tooltip')" :action="showSaveDialog"><v-icon>mdi-content-save</v-icon></ToolbarButton>
      <ToolbarButton :text="$t('menu.settings.tooltip')" :action="showSettingsDialog"><v-icon>mdi-settings</v-icon></ToolbarButton>
    </v-app-bar>

    <v-content>
      <v-container fluid class="px-0">
        <v-layout justify-center align-center class="px-0">
          <game-page ref="gameZone" @snackbar="showSnackbar"
            :boardBackground="boardBackgroundColor"
          ></game-page>
        </v-layout>

        <SettingsDialog ref="settingsDialog"
          @configurationUpdated="updateBoardAndEngine"
        ></SettingsDialog>

        <SimpleModalDialog ref="errorDialog" :title="errorDialogTitle">
            <v-card-text>{{errorDialogText}}</v-card-text>
        </SimpleModalDialog>

        <SimpleModalDialog ref="saveSuccessfulDialog" :title="$t('modals.saveSuccess.title')">
            <v-card-text>{{$t('modals.saveSuccess.text')}}</v-card-text>
        </SimpleModalDialog>

        <SimpleModalDialog ref="newGameConfirmation" 
          :title="$t('modals.newGame.title')"
          :confirmAction="doStartNewGame"
          cancelButton
        >
            <v-card-text>{{$t('modals.newGame.text')}}</v-card-text>
        </SimpleModalDialog>

        <SimpleModalDialog ref="stopGameConfirmation"
          :title="$t('modals.stopGame.title')"
          :confirmAction="doStopGame"
          cancelButton
        >
          <v-card-text>{{$t('modals.stopGame.text')}}</v-card-text>
        </SimpleModalDialog>

        <SimpleSnackBar ref="snackbar">
          {{ snackBarMessage }}
        </SimpleSnackBar>

      </v-container>
    </v-content>

    <v-footer app fixed>
      <span style="margin-left:1em">&copy; Laurent Bernab&eacute; - 2020</span>
    </v-footer>
  </v-app>
</template>

<script>
  import GamePage from "./components/GamePage";
  import ToolbarButton from './components/ToolbarButton';
  import SimpleModalDialog from './components/SimpleModalDialog';
  import SimpleSnackBar from './components/SimpleSnackBar';
  import SettingsDialog from './components/SettingsDialog';

  import Chess from 'chess.js';
  import pgnReader from './libs/pgn';

  export default {
    data: () => ({
      errorDialogTitle: '',
      errorDialogText: '',
      snackBarMessage: '',
      settings: undefined,
      boardBackgroundColor: '#124589',
    }),
    mounted() {
      this.$i18n.locale = navigator.language.substring(0, 2);
      window.backend.SettingsManager.Load().then(content => {
        if (content === '#ConversionToJsonError') {
          this.settings = undefined;
        }
        else {
          this.settings = JSON.parse(content);
        }
      });
    },
    methods: {
      newGameRequest: function() {
        const chessBoard = document.querySelector('loloof64-chessboard');
        if (chessBoard.gameIsInProgress()) {
          this.$refs['newGameConfirmation'].open();
        }
        else {
          this.doStartNewGame();
        }
      },
      doStartNewGame: function() {
        const fileName = 'PgnVierge';
        const filePath = '/home/laurent-bernabe/Documents/temp/pgn/' + fileName + '.pgn';
        window.backend.TextFileManager.GetTextFileContentWithPathProviden(filePath)
        .then(content => {
          if (content === '#ErrorReadingFile') {
            this.errorDialogTitle = this.$i18n.t('modals.failedToReadPgn.title');
            this.errorDialogText = this.$i18n.t('modals.failedToReadPgn.text');
            this.$refs['errorDialog'].open();
            return;
          }

          const loader = new pgnReader({
            pgn: content,
          });

          try {
            const result = loader.load_pgn();
            const chessInstance = new Chess(result.finalPosition);

            const resultingPosition = chessInstance.fen();
            const playerHasWhite = chessInstance.turn() === 'w';

            this.$refs['gameZone'].newGame(resultingPosition, playerHasWhite);
            const chessBoard = document.querySelector('loloof64-chessboard');
            if (chessBoard.gameIsInProgress()) {
              this.$refs['snackbar'].open(this.$i18n.t('game.status.started'));
            }
            else {
              this.$refs['snackbar'].open(this.$i18n.t('game.status.already_finished'));
            }
          }
          catch (error) {
            console.error(error);
            this.errorDialogTitle = this.$i18n.t('modals.failedToReadPgn.title');
            this.errorDialogText = this.$i18n.t('modals.failedToReadPgn.text');
            this.$refs['errorDialog'].open();
          }
        });
      },
      stopGameRequest: function() {
        const chessBoard = document.querySelector('loloof64-chessboard');
        if (chessBoard.gameIsInProgress()) {
          this.$refs['stopGameConfirmation'].open();
        }
      },
      doStopGame: function() {
        const chessBoard = document.querySelector('loloof64-chessboard');
        chessBoard.stop();
        this.$refs['snackbar'].open(this.$i18n.t('game.status.stopped'));
        this.$refs['gameZone'].$refs['history'].selectLastMove()
      },
      toggleSide: function() {
        const chessBoard = document.querySelector('loloof64-chessboard');
        chessBoard.toggleSide();
      },
      showSettingsDialog: function() {
        this.$refs['settingsDialog'].open(this.settings);
      },
      showSnackbar: function(event) {
        this.$refs['snackbar'].open(event);
      },
      showSaveDialog: function() {
        const chessBoard = document.querySelector('loloof64-chessboard');
        if (chessBoard.gameIsInProgress()) return;


        const historyComponent = this.$refs['gameZone'].$refs['history'];
        if (!historyComponent.hasData()) return;

        const pgn = chessBoard.gamePgn();
        const fileName = 'PgnGenere';
        const filePath = '/home/laurent-bernabe/Documents/temp/pgn/' + fileName + '.pgn';

        window.backend.TextFileManager.SaveTextFile(pgn).then(error => {
          if (error === '#ErrorSavingFile')
          {
            this.errorDialogTitle = this.$i18n.t('modals.saveError.title');
            this.errorDialogText = this.$i18n.t('modals.saveError.text');
            this.$refs['errorDialog'].open();
          }
          else if (error === '#ErrorClosingFile')
          {
            this.errorDialogTitle = this.$i18n.t('modals.saveClosingError.title');
            this.errorDialogText = this.$i18n.t('modals.saveClosingError.text');
            this.$refs['errorDialog'].open();
          }
          else {
            this.$refs['saveSuccessfulDialog'].open();
          }
        });
      },
      updateBoardAndEngine: function(configuration) {
        // Update used engine
        window.backend.UciEngine.LoadEngineWithPathProviden(configuration.EnginePath).then(error => {
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

        // Apply modifications on board.
        this.boardBackgroundColor = configuration.BoardBackgroundColor;        
      }
    },
    components: {
      GamePage,
      ToolbarButton,
      SimpleModalDialog,
      SimpleSnackBar,
      SettingsDialog,
    },
    props: {
      source: String
    }
  }
</script>