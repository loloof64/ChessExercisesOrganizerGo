<template>
  <div>
    <SimpleModalDialog
      ref="modal"
      :title="$t('modals.settings.title')"
      :confirmAction="update"
      :cancelAction="cancel"
      cancelButton
    >
      <v-row class="ml-6 my-2">
        <v-btn class="mx-6" @click="loadDefaultSettings()">{{$t('modals.settings.loadDefaults')}}</v-btn>
      </v-row>

      <v-row class="ml-6 my-2">
          {{$t('modals.settings.selectedEngine')}} {{ selectedEnginePath }}
      </v-row>

      <v-row class="ml-6 my-2">
        <v-btn class="mx-6" @click="selectEngine()">{{$t('modals.settings.configureEngine')}}</v-btn>
      </v-row>

      <v-row class="ml-6 my-2">
        <v-checkbox v-model="tempSettings.CoordinatesVisible" :label="$t('modals.settings.coordinatesVisible')"></v-checkbox>
      </v-row>

      <v-row class="ml-6 my-2">
        <v-checkbox v-model="tempSettings.LastMoveArrowVisible" :label="$t('modals.settings.lastMoveArrowVisible')"></v-checkbox>
      </v-row>

      <!-- Background color -->
      <v-row class="ml-6 my-2">
        <v-btn cols="3" class="mx-6" @click="manageEditingBoardBackgroundColorVisibility()"
          v-bind:style="{backgroundColor: tempSettings.BoardBackgroundColor}"
        >{{backgroundColorButtonText}}</v-btn>
        <v-btn v-if="isEditingBoardBackgroundColor === true" cols="3" dark class="mx-6 red" @click="cancelEditingBoardBackground()">
          {{$t('modals.global.cancelButton')}}
        </v-btn>
        <div cols="6" class="mx-6 my-2">{{$t('modals.settings.boardBackgroundColor')}}</div>
      </v-row>
      <v-color-picker class="mx-6" v-if="isEditingBoardBackgroundColor" 
        v-model="editingBoardBackgroundColorValue"></v-color-picker>

      <!-- Coordinates color -->
      <v-row class="ml-6 my-2" v-if="mustShowCoordinatesColorOption">
        <v-btn cols="3" class="mx-6" @click="manageEditingBoardCoordinatesColorVisibility()"
          v-bind:style="{backgroundColor: tempSettings.BoardCoordinatesColor}"
        >{{coordinatesColorButtonText}}</v-btn>
        <v-btn v-if="isEditingBoardCoordinatesColor === true" cols="3" dark class="mx-6 red" @click="cancelEditingBoardCoordinates()">
          {{$t('modals.global.cancelButton')}}
        </v-btn>
        <div cols="6" class="mx-6 my-2">{{$t('modals.settings.boardCoordinatesColor')}}</div>
      </v-row>
      <v-color-picker class="mx-6" v-if="isEditingBoardCoordinatesColor" 
        v-model="editingBoardCoordinatesColorValue"></v-color-picker>

      <!-- Last move arrow color -->
      <v-row class="ml-6 my-2" v-if="mustShowLastMoveArrowColorOption">
        <v-btn cols="3" class="mx-6" @click="manageEditingBoardLastMoveArrowColorVisibility()"
          v-bind:style="{backgroundColor: tempSettings.BoardLastMoveArrowColor}"
        >{{lastMoveArrowColorButtonText}}</v-btn>
        <v-btn v-if="isEditingBoardLastMoveArrowColor === true" cols="3" dark class="mx-6 red" @click="cancelEditingBoardLastMoveArrow()">
          {{$t('modals.global.cancelButton')}}
        </v-btn>
        <div cols="6" class="mx-6 my-2">{{$t('modals.settings.boardLastMoveArrowColor')}}</div>
      </v-row>
      <v-color-picker class="mx-6" v-if="isEditingBoardLastMoveArrowColor" 
        v-model="editingBoardLastMoveArrowColorValue"></v-color-picker>

      <!-- White cells color -->
      <v-row class="ml-6 my-2">
        <v-btn cols="3" class="mx-6" @click="manageEditingBoardLastWhiteCellsColorVisibility()"
          v-bind:style="{backgroundColor: tempSettings.BoardWhiteCellsColor}"
        >{{whiteCellsColorButtonText}}</v-btn>
        <v-btn v-if="isEditingBoardWhiteCellsColor === true" cols="3" dark class="mx-6 red" @click="cancelEditingBoardWhiteCells()">
          {{$t('modals.global.cancelButton')}}
        </v-btn>
        <div cols="6" class="mx-6 my-2">{{$t('modals.settings.boardWhiteCellsColor')}}</div>
      </v-row>
      <v-color-picker class="mx-6" v-if="isEditingBoardWhiteCellsColor" 
        v-model="editingBoardWhiteCellsColorValue"></v-color-picker>

      <!-- Black cells color -->
      <v-row class="ml-6 my-2">
        <v-btn cols="3" class="mx-6" @click="manageEditingBoardLastBlackCellsColorVisibility()"
          v-bind:style="{backgroundColor: tempSettings.BoardBlackCellsColor}"
        >{{blackCellsColorButtonText}}</v-btn>
        <v-btn v-if="isEditingBoardBlackCellsColor === true" cols="3" dark class="mx-6 red" @click="cancelEditingBoardBlackCells()">
          {{$t('modals.global.cancelButton')}}
        </v-btn>
        <div cols="6" class="mx-6 my-2">{{$t('modals.settings.boardBlackCellsColor')}}</div>
      </v-row>
      <v-color-picker class="mx-6" v-if="isEditingBoardBlackCellsColor" 
        v-model="editingBoardBlackCellsColorValue"></v-color-picker>

      <!-- Dnd start color -->
      <v-row class="ml-6 my-2">
        <v-btn cols="3" class="mx-6" @click="manageEditingBoardDndStartColorVisibility()"
          v-bind:style="{backgroundColor: tempSettings.BoardDndStartColor}"
        >{{dndStartColorButtonText}}</v-btn>
        <v-btn v-if="isEditingBoardDndStartColor === true" cols="3" dark class="mx-6 red" @click="cancelEditingBoardDndStart()">
          {{$t('modals.global.cancelButton')}}
        </v-btn>
        <div cols="6" class="mx-6 my-2">{{$t('modals.settings.boardDndStartCellColor')}}</div>
      </v-row>
      <v-color-picker class="mx-6" v-if="isEditingBoardDndStartColor" 
        v-model="editingBoardDndStartColorValue"></v-color-picker>
        
      <!-- Dnd end color -->
      <v-row class="ml-6 my-2">
        <v-btn cols="3" class="mx-6" @click="manageEditingBoardDndEndColorVisibility()"
          v-bind:style="{backgroundColor: tempSettings.BoardDndEndColor}"
        >{{dndEndColorButtonText}}</v-btn>
        <v-btn v-if="isEditingBoardDndEndColor === true" cols="3" dark class="mx-6 red" @click="cancelEditingBoardDndEnd()">
          {{$t('modals.global.cancelButton')}}
        </v-btn>
        <div cols="6" class="mx-6 my-2">{{$t('modals.settings.boardDndEndCellColor')}}</div>
      </v-row>
      <v-color-picker class="mx-6" v-if="isEditingBoardDndEndColor" 
        v-model="editingBoardDndEndColorValue"></v-color-picker>

      <!-- Dnd cross color -->
      <v-row class="ml-6 my-2">
        <v-btn cols="3" class="mx-6" @click="manageEditingBoardDndSCrossColorVisibility()"
          v-bind:style="{backgroundColor: tempSettings.BoardDndCrossColor}"
        >{{dndCrossColorButtonText}}</v-btn>
        <v-btn v-if="isEditingBoardDndCrossColor === true" cols="3" dark class="mx-6 red" @click="cancelEditingBoardDndCross()">
          {{$t('modals.global.cancelButton')}}
        </v-btn>
        <div cols="6" class="mx-6 my-2">{{$t('modals.settings.boardDndSCrossColor')}}</div>
      </v-row>
      <v-color-picker class="mx-6" v-if="isEditingBoardDndCrossColor" 
        v-model="editingBoardDndCrossColorValue"></v-color-picker>

      <v-row class="mx-6 my-2">
        <v-text-field
          v-model="tempSettings.PlayerName"
          :label="$t('modals.settings.playerName')"
        >
        </v-text-field>
      </v-row>

      <v-row class="mx-6 my-2">
        <v-text-field
          v-model="tempSettings.ComputerName"
          :label="$t('modals.settings.computerName')"
        >
        </v-text-field>
      </v-row>

    </SimpleModalDialog>

    <SimpleModalDialog ref="errorDialog" :title="errorDialogTitle">
        <v-card-text>{{errorDialogText}}</v-card-text>
    </SimpleModalDialog>
  </div>
</template>

<script>
import SimpleModalDialog from "./SimpleModalDialog";

export default {
  name: "SettingsDialog",
  data() {
    return {
      settings: {},
      tempSettings: {},
      errorDialogTitle: '',
      errorDialogText: '',

      isEditingBoardBackgroundColor: false,
      editingBoardBackgroundColorValue: undefined,
      backgroundColorButtonText: '',

      isEditingBoardCoordinatesColor: false,
      editingBoardCoordinatesColorValue: undefined,
      coordinatesColorButtonText: '',
      
      isEditingBoardLastMoveArrowColor: false,
      editingBoardLastMoveArrowColorValue: undefined,
      lastMoveArrowColorButtonText: '',
      
      isEditingBoardWhiteCellsColor: false,
      editingBoardWhiteCellsColorValue: undefined,
      whiteCellsColorButtonText: '',
      
      isEditingBoardBlackCellsColor: false,
      editingBoardBlackCellsColorValue: undefined,
      blackCellsColorButtonText: '',
      
      isEditingBoardDndStartColor: false,
      editingBoardDndStartColorValue: undefined,
      dndStartColorButtonText: '',

      isEditingBoardDndEndColor: false,
      editingBoardDndEndColorValue: undefined,
      dndEndColorButtonText: '',

      isEditingBoardDndCrossColor: false,
      editingBoardDndCrossColorValue: undefined,
      dndCrossColorButtonText: '',
    };
  },
  mounted() {
    this.$i18n.locale = navigator.language.substring(0, 2);
    this.backgroundColorButtonText = this.$i18n.t('modals.settings.modifyColor');
    this.coordinatesColorButtonText = this.$i18n.t('modals.settings.modifyColor');
    this.lastMoveArrowColorButtonText = this.$i18n.t('modals.settings.modifyColor');
    this.whiteCellsColorButtonText = this.$i18n.t('modals.settings.modifyColor');
    this.blackCellsColorButtonText = this.$i18n.t('modals.settings.modifyColor');
    this.dndStartColorButtonText = this.$i18n.t('modals.settings.modifyColor');
    this.dndEndColorButtonText = this.$i18n.t('modals.settings.modifyColor');
    this.dndCrossColorButtonText = this.$i18n.t('modals.settings.modifyColor');
  },
  methods: {
    update: function() {
      this.settings = this.tempSettings;

      // Save settings into file
      window.backend.SettingsManager.Save(this.settings).then(error => {
        if (error === '#ConversionError' || error === '#FileSavingError') {
          console.error(error);
          this.errorDialogTitle = this.$i18n.t(
            "modals.settings.failedToSaveTitle"
          );
          this.errorDialogText = this.$i18n.t(
            "modals.settings.failedToSaveText"
          );
          this.$refs["errorDialog"].open();
        }
      });

      this.$emit('configurationUpdated', this.settings);
      this.closeColorChoosers();
    },
    selectEngine: function() {
      // Production mode : path should be set to the value of
      // window.backend.UciEngine.GetUserEnginePath promise
      window.backend.UciEngine.GetUserEnginePath().then(path => {
        if (path.length > 0) {
          this.tempSettings.EnginePath = path;
        }
      });
    },
    cancel: function() {
      this.tempSettings = this.settings;
      this.closeColorChoosers();
    },
    open: function(currentSettings) {
        const newSettings = currentSettings;
        this.settings = newSettings;
        this.tempSettings = newSettings;
        this.$refs['modal'].open();
    },
    loadDefaultSettings: function() {
      window.backend.SettingsManager.GetDefaultSettings().then(result => {
        if (result === '#ConversionToJsonError') {
          console.error(result);
          this.errorDialogTitle = this.$i18n.t('modals.settings.failedLoadingDefaultSettingsTitle');
          this.errorDialogText = this.$i18n.t('modals.settings.failedLoadingDefaultSettingsText');
          this.$refs['errorDialog'].open();
          return;
        }

        this.tempSettings = JSON.parse(result);
      });
    },
    manageEditingBoardBackgroundColorVisibility: function() {
      if (this.isEditingBoardBackgroundColor === true) {
        this.tempSettings.BoardBackgroundColor = this.editingBoardBackgroundColorValue.hex || this.editingBoardBackgroundColorValue;
        this.isEditingBoardBackgroundColor = false;
        this.backgroundColorButtonText = this.$i18n.t('modals.settings.modifyColor');
      }
      else {
        this.editingBoardBackgroundColorValue = this.tempSettings.BoardBackgroundColor;
        this.isEditingBoardBackgroundColor = true;
        this.backgroundColorButtonText = this.$i18n.t('modals.settings.acceptColor');
      }
    },
    manageEditingBoardCoordinatesColorVisibility: function() {
      if (this.isEditingBoardCoordinatesColor === true) {
        this.tempSettings.BoardCoordinatesColor = this.editingBoardCoordinatesColorValue.hex || this.editingBoardCoordinatesColorValue;
        this.isEditingBoardCoordinatesColor = false;
        this.coordinatesColorButtonText = this.$i18n.t('modals.settings.modifyColor');
      }
      else {
        this.editingBoardCoordinatesColorValue = this.tempSettings.BoardCoordinatesColor;
        this.isEditingBoardCoordinatesColor = true;
        this.coordinatesColorButtonText = this.$i18n.t('modals.settings.acceptColor');
      }
    },
    manageEditingBoardLastMoveArrowColorVisibility: function() {
      if (this.isEditingBoardLastMoveArrowColor === true) {
        this.tempSettings.BoardLastMoveArrowColor = this.editingBoardLastMoveArrowColorValue.hex || this.editingBoardLastMoveArrowColorValue;
        this.isEditingBoardLastMoveArrowColor = false;
        this.lastMoveArrowColorButtonText = this.$i18n.t('modals.settings.modifyColor');
      }
      else {
        this.editingBoardLastMoveArrowColorValue = this.tempSettings.BoardLastMoveArrowColor;
        this.isEditingBoardLastMoveArrowColor = true;
        this.lastMoveArrowColorButtonText = this.$i18n.t('modals.settings.acceptColor');
      }
    },
    manageEditingBoardLastWhiteCellsColorVisibility: function() {
      if (this.isEditingBoardWhiteCellsColor === true) {
        this.tempSettings.BoardWhiteCellsColor = this.editingBoardWhiteCellsColorValue.hex || this.editingBoardWhiteCellsColorValue;
        this.isEditingBoardWhiteCellsColor = false;
        this.whiteCellsColorButtonText = this.$i18n.t('modals.settings.modifyColor');
      }
      else {
        this.editingBoardWhiteCellsColorValue = this.tempSettings.BoardWhiteCellsColor;
        this.isEditingBoardWhiteCellsColor = true;
        this.whiteCellsColorButtonText = this.$i18n.t('modals.settings.acceptColor');
      }
    },
    manageEditingBoardLastBlackCellsColorVisibility: function() {
      if (this.isEditingBoardBlackCellsColor === true) {
        this.tempSettings.BoardBlackCellsColor = this.editingBoardBlackCellsColorValue.hex || this.editingBoardBlackCellsColorValue;
        this.isEditingBoardBlackCellsColor = false;
        this.blackCellsColorButtonText = this.$i18n.t('modals.settings.modifyColor');
      }
      else {
        this.editingBoardBlackCellsColorValue = this.tempSettings.BoardBlackCellsColor;
        this.isEditingBoardBlackCellsColor = true;
        this.blackCellsColorButtonText = this.$i18n.t('modals.settings.acceptColor');
      }
    },
    manageEditingBoardDndStartColorVisibility: function() {
      if (this.isEditingBoardDndStartColor === true) {
        this.tempSettings.BoardDndStartColor = this.editingBoardDndStartColorValue.hex || this.editingBoardDndStartColorValue;
        this.isEditingBoardDndStartColor = false;
        this.dndStartColorButtonText = this.$i18n.t('modals.settings.modifyColor');
      }
      else {
        this.editingBoardDndStartColorValue = this.tempSettings.BoardDndStartColor;
        this.isEditingBoardDndStartColor = true;
        this.dndStartColorButtonText = this.$i18n.t('modals.settings.acceptColor');
      }
    },
    manageEditingBoardDndEndColorVisibility: function() {
      if (this.isEditingBoardDndEndColor === true) {
        this.tempSettings.BoardDndEndColor = this.editingBoardDndEndColorValue.hex || this.editingBoardDndEndColorValue;
        this.isEditingBoardDndEndColor = false;
        this.dndEndColorButtonText = this.$i18n.t('modals.settings.modifyColor');
      }
      else {
        this.editingBoardDndEndColorValue = this.tempSettings.BoardDndEndColor;
        this.isEditingBoardDndEndColor = true;
        this.dndEndColorButtonText = this.$i18n.t('modals.settings.acceptColor');
      }
    },
    manageEditingBoardDndSCrossColorVisibility: function() {
      if (this.isEditingBoardDndCrossColor === true) {
        this.tempSettings.BoardDndCrossColor = this.editingBoardDndCrossColorValue.hex || this.editingBoardDndCrossColorValue;
        this.isEditingBoardDndCrossColor = false;
        this.dndCrossColorButtonText = this.$i18n.t('modals.settings.modifyColor');
      }
      else {
        this.editingBoardDndCrossColorValue = this.tempSettings.BoardDndCrossColor;
        this.isEditingBoardDndCrossColor = true;
        this.dndCrossColorButtonText = this.$i18n.t('modals.settings.acceptColor');
      }
    },
    cancelEditingBoardBackground: function() {
      this.isEditingBoardBackgroundColor = false;
      this.backgroundColorButtonText = this.$i18n.t('modals.settings.modifyColor');
    },
    cancelEditingBoardCoordinates: function() {
      this.isEditingBoardCoordinatesColor = false;
      this.coordinatesColorButtonText = this.$i18n.t('modals.settings.modifyColor');
    },
    cancelEditingBoardLastMoveArrow: function() {
      this.isEditingBoardLastMoveArrowColor = false;
      this.lastMoveArrowColorButtonText = this.$i18n.t('modals.settings.modifyColor');
    },
    cancelEditingBoardWhiteCells: function() {
      this.isEditingBoardWhiteCellsColor = false;
      this.whiteCellsColorButtonText = this.$i18n.t('modals.settings.modifyColor');
    },
    cancelEditingBoardBlackCells: function() {
      this.isEditingBoardBlackCellsColor = false;
      this.blackCellsColorButtonText = this.$i18n.t('modals.settings.modifyColor');
    },
    cancelEditingBoardDndStart: function() {
      this.isEditingBoardDndStartColor = false;
      this.dndStartColorButtonText = this.$i18n.t('modals.settings.modifyColor');
    },
    cancelEditingBoardDndEnd: function() {
      this.isEditingBoardDndEndColor = false;
      this.dndSEndColorButtonText = this.$i18n.t('modals.settings.modifyColor');
    },
    cancelEditingBoardDndCross: function() {
      this.isEditingBoardDndCrossColor = false;
      this.dndCrossColorButtonText = this.$i18n.t('modals.settings.modifyColor');
    },
    closeColorChoosers: function() {
      this.isEditingBoardBackgroundColor = false;
      this.isEditingBoardCoordinatesColor = false;
      this.isEditingBoardLastMoveArrowColor = false;
      this.isEditingBoardWhiteCellsColor = false;
      this.isEditingBoardBlackCellsColor = false;
      this.isEditingBoardDndStartColor = false;
      this.isEditingBoardDndEndColor = false;
      this.isEditingBoardDndCrossColor = false;
    },
  },
  computed: {
    selectedEnginePath: function() {
      if (this.tempSettings.EnginePath === undefined) return this.$i18n.t('modals.settings.noEngine');
      if (this.tempSettings.EnginePath.length === 0) return this.$i18n.t('modals.settings.noEngine');
      return this.tempSettings.EnginePath; 
    },
    mustShowCoordinatesColorOption: function() {
      if (this.tempSettings.CoordinatesVisible === undefined) return false;
      return [true, "true"].includes(this.tempSettings.CoordinatesVisible);
    },
    mustShowLastMoveArrowColorOption: function() {
      if (this.tempSettings.LastMoveArrowVisible === undefined) return false;
      return [true, "true"].includes(this.tempSettings.LastMoveArrowVisible);
    },
  },
  components: {
    SimpleModalDialog
  }
};
</script>