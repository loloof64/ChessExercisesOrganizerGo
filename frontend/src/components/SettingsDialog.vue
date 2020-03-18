<template>
  <div>
    <SimpleModalDialog
      ref="modal"
      :title="$t('modals.settings.title')"
      :confirmAction="update"
      :cancelAction="cancel"
      cancelButton
    >
      <v-btn class="mx-a" @click="selectEngine()">{{$t('modals.settings.configureEngine')}}</v-btn>
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
    };
  },
  methods: {
    update: function() {
      this.settings = this.tempSettings;
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

      window.backend.UciEngine.LoadEngineWithPathProviden(this.settings.EnginePath).then(error => {
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
    }
  },
  components: {
    SimpleModalDialog
  }
};
</script>