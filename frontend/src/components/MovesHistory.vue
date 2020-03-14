<template>
  <v-container class="white indigo--text" id="root" width="400" height="600">
    <v-row class="mx-2">
        <v-col v-for="btn in buttons" :key="btn.tooltipTrKey" cols="3">
          <v-tooltip bottom>
            <template v-slot:activator="{ on }">
                <v-btn v-on="on" @click="btn.action()" class="blue lighten-3 red--text text-darken-1"><v-icon>{{btn.icon}}</v-icon></v-btn>
            </template>
            <span>{{ $t(btn.tooltipTrKey) }}</span>
          </v-tooltip>
        </v-col>
    </v-row>
    <v-row>
      <v-col move cols="4" class="d-flex justify-end">{{$t('history.header.moveNumber')}}</v-col>
      <v-col move cols="4" class="d-flex justify-end">{{$t('history.header.whiteSide')}}</v-col>
      <v-col move cols="4" class="d-flex justify-end">{{$t('history.header.blackSide')}}</v-col>
    </v-row>
    <v-row v-for="(historyLine, index) in history" :key="historyLine.moveNumber">
      <v-col move cols="4" class="d-flex justify-end">{{historyLine.moveNumber}}</v-col>
      <v-col
        move
        cols="4"
        class="d-flex justify-end"
        :class="{ highlight: isSelectedPosition(index, true) }"
        @click="setPosition(index, true)"
      >{{historyLine.white ? historyLine.white.moveFan : ''}}</v-col>
      <v-col
        move
        cols="4"
        class="d-flex justify-end"
        :class="{ highlight: isSelectedPosition(index, false) }"
        @click="setPosition(index, false)"
      >{{historyLine.black? historyLine.black.moveFan : ''}}</v-col>
    </v-row>
  </v-container>
</template>

<script>
import ToolbarButton from './ToolbarButton';
export default {
  name: 'MovesHistory',
  props: {
    history: {
      type: Array,
      default: []
    }
  },
  data() {
    return {
      buttons: [
        {
          action: this.goFirst,
          tooltipTrKey: 'history.buttons.first',
          icon: 'mdi-skip-previous'
        },
        {
          action: this.goPrevious,
          tooltipTrKey: 'history.buttons.previous',
          icon: 'mdi-chevron-left'
        },
        {
          action: this.goNext,
          tooltipTrKey: 'history.buttons.next',
          icon: 'mdi-chevron-right'
        },
        {
          action: this.goLast,
          tooltipTrKey: 'history.buttons.last',
          icon: 'mdi-skip-next'
        },
      ],
      selectedPosition: undefined,
    };
  },
  methods: {
    setPosition: function(historyIndex, whitePlayer) {
      this.$emit("position_requested", { historyIndex, whitePlayer });
    },
    goFirst: function() {

    },
    goPrevious: function() {

    },
    goNext: function() {

    },
    goLast: function() {

    },
    confirmPositionSet: function(evt) {
      this.selectedPosition = evt;
    },
    isSelectedPosition: function(index, whitePlayer) {
      if (this.selectedPosition === undefined) return false;
      return whitePlayer === this.selectedPosition.whitePlayer &&
        index === this.selectedPosition.historyIndex;
    }
  },
  components: {
    ToolbarButton,
  }
};
</script>

<style scoped>
#root {
  color: black;
  display: inline-block;
}

.col[move] {
  border-right: 3px solid black;
  border-bottom: 3px solid black;
  font-size: 24px;
}

.highlight {
  background-color: aqua;
  color: brown;
}
</style>