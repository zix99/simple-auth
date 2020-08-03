<template>
  <span class="hoverable" :title="longDate">{{shortDate}}</span>
</template>

<script>
import dayjs from 'dayjs';

export default {
  props: {
    date: { required: true },
    shortFormat: { default: 'l' },
    longFormat: { default: 'lll' },
  },
  computed: {
    shortDate() {
      return dayjs(this.date).format(this.shortFormat);
    },
    longDate() {
      return dayjs(this.date).format(this.longFormat);
    },
  },
};
</script>

<style scoped>
span.hoverable {
  border-bottom: 1px dotted #999;
  cursor: help;
}
</style>
