<template>
  <v-app id="inspire">
    <v-navigation-drawer v-model="drawer" clipped fixed app>
      <v-list dense>
        <v-list-item @click="newGame()">
          <v-list-item-icon>
            <v-icon>mdi-restart</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>{{$t('menu.newGame.caption')}}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item @click="toggleSide()">
          <v-list-item-icon>
            <v-icon>mdi-arrow-up-down</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>{{$t('menu.toggleSide.caption')}}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item @click="stopGame()">
          <v-list-item-icon>
            <v-icon>mdi-stop-circle</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>{{$t('menu.stopGame.caption')}}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item @click="showSettingsDialog()">
          <v-list-item-icon>
            <v-icon>mdi-settings</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>{{$t('menu.settings.caption')}}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>
    <v-app-bar app fixed clipped-left>
      <v-app-bar-nav-icon @click.stop="drawer = !drawer"></v-app-bar-nav-icon>
      <v-toolbar-title>Chess Exercises Organizer</v-toolbar-title>
      <v-btn icon @click="newGame()"><v-icon>mdi-restart</v-icon></v-btn>
      <v-btn icon @click="toggleSide()"><v-icon>mdi-arrow-up-down</v-icon></v-btn>
      <v-btn icon @click="stopGame()"><v-icon>mdi-stop-circle</v-icon></v-btn>
    </v-app-bar>

    <v-content>
      <v-container fluid class="px-0">
        <v-layout justify-center align-center class="px-0">
          <game-page ref="gameZone"></game-page>
        </v-layout>

        <v-dialog v-model="settingsDialog" persistent max-width="300">
          <v-card>
            <v-card-title class="headline">{{settingsDialogTitle}}</v-card-title>
            <v-btn @click="selectEngine()">{{configureEngineButtonText}}</v-btn>
            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn color="blue darken-1" text @click="settingsDialog = false">{{okButtonText}}</v-btn>
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
            <v-btn color="green darken-1" text @click="errorDialog = false">{{okButtonText}}</v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>

    <v-footer app fixed>
      <span style="margin-left:1em">&copy; Laurent Bernab&eacute;</span>
    </v-footer>
  </v-app>
</template>

<script>
  import GamePage from "./components/GamePage.vue";

  export default {
    data: () => ({
      drawer: false,
      settingsDialog: false,
      settingsDialogTitle: 'Settings',
      configureEngineButtonText: 'Configure UCI engine',
      errorDialog: false,
      errorDialogTitle: '',
      errorDialogText: '',
      okButtonText: 'Ok',
    }),
    mounted() {
      this.$i18n.locale = navigator.language.substring(0, 2);
      window.backend.UciEngine.SetEnginePathManually('/home/laurent-bernabe/Programmes/Echecs/stockfish-11-linux/stockfish-11-linux/Linux/stockfish_20011801_x64_modern')
        .then(error => {
          if (error === '#ConfigEngineErr'){
            this.errorDialogTitle = 'Could not select this engine';
            this.errorDialogText = 'Failed to select the engine : have you got the right over it ?';
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
            this.errorDialogTitle = 'Could not select this engine';
            this.errorDialogText = 'Failed to select the engine : have you got the right over it ?';
            this.errorDialog = true;
          }
        });
      },
    },
    components: {
      GamePage,
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