<template>
  <SimpleModalDialog
    ref="modal"
    :title="$t('modals.settings.title')"
    :confirmAction="update()"
    cancelButton
  >
    <v-btn class="mx-a" @click="selectEngine()">{{$t('modals.settings.configureEngine')}}</v-btn>
  </SimpleModalDialog>
</template>

<script>
import SimpleModalDialog from "./SimpleModalDialog";

export default {
  name: "SettingsDialog",
  methods: {
    update: function() {},
    selectEngine: function() {
      const path =
        "/home/laurent-bernabe/Programmes/Echecs/stockfish-11-linux/stockfish-11-linux/Linux/stockfish_20011801_x64_modern";
      window.backend.UciEngine.LoadEngineWithPathProviden(path).then(error => {
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
    open: function() {
        this.$refs['modal'].open();
    }
  },
  components: {
    SimpleModalDialog
  }
};
</script>