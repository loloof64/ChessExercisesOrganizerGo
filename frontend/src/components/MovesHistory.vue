<template>
  <v-container class="white indigo--text" id="root" width="400" height="600">
    <v-row>
      <v-col cols="4" class="d-flex justify-end">{{$t('history.header.moveNumber')}}</v-col>
      <v-col cols="4" class="d-flex justify-end">{{$t('history.header.whiteSide')}}</v-col>
      <v-col cols="4" class="d-flex justify-end">{{$t('history.header.blackSide')}}</v-col>
    </v-row>
    <v-row v-for="(historyLine, index) in history" :key="historyLine.moveNumber">
      <v-col cols="4" class="d-flex justify-end">{{historyLine.moveNumber}}</v-col>
      <v-col
        cols="4"
        class="d-flex justify-end"
        @click="setPosition(index, true)"
      >{{historyLine.white ? historyLine.white.moveFan : ''}}</v-col>
      <v-col
        cols="4"
        class="d-flex justify-end"
        @click="setPosition(index, false)"
      >{{historyLine.black? historyLine.black.moveFan : ''}}</v-col>
    </v-row>
  </v-container>
</template>

<script>
export default {
  props: {
    history: {
      type: Array,
      default: []
    }
  },
  methods: {
    setPosition: function(historyIndex, whitePlayer) {
      this.$emit("position_requested", { historyIndex, whitePlayer });
    }
  }
};
</script>

<style scoped>
#root {
  color: black;
  display: inline-block;
}

.col {
  border-right: 3px solid black;
  border-bottom: 3px solid black;
  font-size: 24px;
}
</style>