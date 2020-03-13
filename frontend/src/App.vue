<template>
  <v-app id="inspire">
    <v-app-bar app fixed clipped-left>
      <v-toolbar-title>{{$t('app.title')}}</v-toolbar-title>
      <ToolbarButton :text="$t('menu.newGame.tooltip')" :action="newGame"><v-icon>mdi-restart</v-icon></ToolbarButton>
      <ToolbarButton :text="$t('menu.toggleSide.tooltip')" :action="toggleSide"><v-icon>mdi-arrow-up-down</v-icon></ToolbarButton>
      <ToolbarButton :text="$t('menu.stopGame.tooltip')" :action="stopGame"><v-icon>mdi-stop-circle</v-icon></ToolbarButton>
      <ToolbarButton :text="$t('menu.settings.tooltip')" :action="showSettingsDialog"><v-icon>mdi-settings</v-icon></ToolbarButton>
    </v-app-bar>

    <v-content>
      <v-container fluid class="px-0">
        <v-layout justify-center align-center class="px-0">
          <game-page ref="gameZone"></game-page>
        </v-layout>

        <v-dialog v-model="settingsDialog" persistent max-width="300">
          <v-card>
            <v-card-title class="headline">{{$t('modals.settings.title')}}</v-card-title>
            <v-btn @click="selectEngine()">{{$t('modals.settings.configureEngine')}}</v-btn>
            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn color="blue darken-1" text @click="settingsDialog = false">{{$t('modals.global.okButton')}}</v-btn>
            </v-card-actions>
          </v-card>
        </v-dialog>
      </v-container>
    </v-content>

    <v-dialog v-model="errorDialog" persistent max-width="300">
        <v-card>
          <v-card-title class="headline">{{errorDialogTitle}}</v-card-title>
          <v-card-text>{{errorDialogText}}</v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="green darken-1" text @click="errorDialog = false">{{$t('modals.global.okButton')}}</v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>

    <v-footer app fixed>
      <span style="margin-left:1em">&copy; Laurent Bernab&eacute; - 2020</span>
    </v-footer>
  </v-app>
</template>

<script>
  import GamePage from "./components/GamePage";
  import ToolbarButton from './components/ToolbarButton';

  export default {
    data: () => ({
      drawer: false,
      settingsDialog: false,
      errorDialog: false,
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
            this.errorDialog = true;
          }
        })
    },
    methods: {
      newGame: function() {
        this.drawer = false;
        this.$refs['gameZone'].newGame();
      },
      stopGame: function() {
        this.drawer = false;
        const chessBoard = document.querySelector('loloof64-chessboard');
        chessBoard.stop();
      },
      toggleSide: function() {
        this.drawer = false;
        const chessBoard = document.querySelector('loloof64-chessboard');
        chessBoard.toggleSide();
      },
      showSettingsDialog: function() {
        this.drawer = false;
        this.settingsDialog = true;
      },
      selectEngine: function() {
        window.backend.UciEngine.Load().then(error => {
          if (error === '#ConfigEngineErr'){
            this.errorDialogTitle = this.$i18n.t('modals.settings.failedToSetupEngineTitle');
            this.errorDialogText = this.$i18n.t('modals.settings.failedToSetupEngineText');
            this.errorDialog = true;
          }
        });
      },
    },
    components: {
      GamePage,
      ToolbarButton,
    },
    props: {
      source: String
    }
  }
</script>

<style>
  .logo {
    width: 16em;
  }
</style>