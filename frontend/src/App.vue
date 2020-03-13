<template>
  <v-app id="inspire">
    <v-app-bar app fixed clipped-left>
      <v-toolbar-title>{{$t('app.title')}}</v-toolbar-title>
      <ToolbarButton :text="$t('menu.newGame.tooltip')" :action="newGameRequest"><v-icon>mdi-restart</v-icon></ToolbarButton>
      <ToolbarButton :text="$t('menu.toggleSide.tooltip')" :action="toggleSide"><v-icon>mdi-arrow-up-down</v-icon></ToolbarButton>
      <ToolbarButton :text="$t('menu.stopGame.tooltip')" :action="stopGame"><v-icon>mdi-stop-circle</v-icon></ToolbarButton>
      <ToolbarButton :text="$t('menu.settings.tooltip')" :action="showSettingsDialog"><v-icon>mdi-settings</v-icon></ToolbarButton>
    </v-app-bar>

    <v-content>
      <v-container fluid class="px-0">
        <v-layout justify-center align-center class="px-0">
          <game-page ref="gameZone"></game-page>
        </v-layout>

        <SimpleModalDialog ref="settingsDialog" :title="$t('modals.settings.title')">
          <v-btn @click="selectEngine()">{{$t('modals.settings.configureEngine')}}</v-btn>
        </SimpleModalDialog>

        <SimpleModalDialog ref="errorDialog" :title="errorDialogTitle">
            <v-card-text>{{errorDialogText}}</v-card-text>
        </SimpleModalDialog>

        <SimpleModalDialog ref="newGameConfirmation" 
          :title="$t('modals.newGame.title')"
          :confirmAction="doStartNewGame"
          cancelButton
        >
            <v-card-text>{{$t('modals.newGame.text')}}</v-card-text>
        </SimpleModalDialog>

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

  export default {
    data: () => ({
      errorDialogTitle: '',
      errorDialogText: '',
    }),
    mounted() {
      this.$i18n.locale = navigator.language.substring(0, 2);
      window.backend.UciEngine.SetEnginePathManually('/home/laurent-bernabe/Programmes/Echecs/stockfish-11-linux/stockfish-11-linux/Linux/stockfish_20011801_x64_modern')
        .then(error => {
          if (error === '#ConfigEngineErr'){
            this.errorDialogTitle = this.$i18n.t('modals.settings.failedToSetupEngineTitle');
            this.errorDialogText = this.$i18n.t('modals.settings.failedToSetupEngineText');
            this.$refs['errorDialog'].open();
          }
        })
    },
    methods: {
      newGameRequest: function() {
        this.$refs['newGameConfirmation'].open();
      },
      doStartNewGame: function() {
        this.$refs['gameZone'].newGame();
      },
      stopGame: function() {
        const chessBoard = document.querySelector('loloof64-chessboard');
        chessBoard.stop();
      },
      toggleSide: function() {
        const chessBoard = document.querySelector('loloof64-chessboard');
        chessBoard.toggleSide();
      },
      showSettingsDialog: function() {
        this.$refs['settingsDialog'].open();
      },
      selectEngine: function() {
        window.backend.UciEngine.Load().then(error => {
          if (error === '#ConfigEngineErr'){
            this.errorDialogTitle = this.$i18n.t('modals.settings.failedToSetupEngineTitle');
            this.errorDialogText = this.$i18n.t('modals.settings.failedToSetupEngineText');
            this.$refs['errorDialog'].open();
          }
        });
      },
    },
    components: {
      GamePage,
      ToolbarButton,
      SimpleModalDialog,
    },
    props: {
      source: String
    }
  }
</script>