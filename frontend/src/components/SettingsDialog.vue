<template>
  <div>
    <SimpleModalDialog
      ref="modal"
      :title="$t('modals.settings.title')"
      :confirmAction="update"
      :cancelAction="cancel"
      cancelButton
    >
      <v-btn class="ml-6 my-2" @click="selectEngine()">{{$t('modals.settings.configureEngine')}}</v-btn>

      <v-row class="my-2">
        <v-btn cols="3" class="mx-6" @click="manageEditingBoardBackgroundColorVisibility()"
          v-bind:style="{backgroundColor: tempSettings.BoardBackgroundColor}"
        >{{backgroundColorButtonText}}</v-btn>
        <div cols="9" class="mx-6 my-2">{{$t('modals.settings.boardBackgroundColor')}}</div>
      </v-row>
      <v-color-picker class="mx-6" v-if="isEditingBoardBackgroundColor" 
        v-model="editingBoardBackgroundColorValue"></v-color-picker>

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
      backgroundColorButtonText: this.$i18n.t('modals.settings.modifyColor'),
    };
  },
  mounted() {
    this.$i18n.locale = navigator.language.substring(0, 2);
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
    },
    selectEngine: function() {
      // Production mode : path should be set to the value of
      // window.backend.UciEngine.GetUserEnginePath promise instead
      const path =
        "/home/laurent-bernabe/Programmes/Echecs/stockfish-11-linux/stockfish-11-linux/Linux/stockfish_20011801_x64_modern";
      this.tempSettings.EnginePath = path;
    },
    cancel: function() {
      this.tempSettings = this.settings;
    },
    open: function(currentSettings) {
        const newSettings = JSON.parse(currentSettings);
        this.settings = newSettings;
        this.tempSettings = newSettings;
        this.$refs['modal'].open();
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
    }
  },
  components: {
    SimpleModalDialog
  }
};
</script>