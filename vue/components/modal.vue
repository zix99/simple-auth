<template>
  <div class="modal is-active" v-if="active">
    <div class="modal-background"></div>
    <div class="modal-card">
      <header class="modal-card-head">
        <p class="modal-card-title">{{title}}</p>
        <button class="delete" aria-label="close" @click="close()" v-if="showClose"></button>
      </header>
      <section class="modal-card-body">
        <slot></slot>
      </section>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    title: null,
    showClose: { default: true },
  },
  data() {
    return {
      active: false,
    };
  },
  methods: {
    open() {
      this.active = true;
    },
    close() {
      this.active = false;
    },
  },
};
</script>
