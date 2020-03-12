<template>
  <v-container class="white indigo--text" id="root" width="400" height="600">
    <v-row>
      <v-col cols="4" class="d-flex justify-end">Number</v-col>
      <v-col cols="4" class="d-flex justify-end">White</v-col>
      <v-col cols="4" class="d-flex justify-end">Black</v-col>
    </v-row>
    <v-row v-for="historyLine in history" :key="historyLine.moveNumber">
      <v-col cols="4" class="d-flex justify-end">{{historyLine.moveNumber}}</v-col>
      <v-col
        cols="4"
        class="d-flex justify-end"
      >{{historyLine.white ? historyLine.white.moveFan : ''}}</v-col>
      <v-col
        cols="4"
        class="d-flex justify-end"
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