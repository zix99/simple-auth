<template>
  <button class="button" @click="clicked">
    <span class="icon"><i class="fab" :class="icon"></i></span>
    <span><slot></slot></span>
  </button>
</template>

<script>
export default {
  props: {
    id: { required: true },
    icon: { required: true },
    continue: null,
  },
  methods: {
    clicked() {
      let url = `/oidc/${this.id}`;
      if (this.continue) {
        url += `?continue=${encodeURIComponent(this.continue)}`;
      }
      window.location = url;
    },
  },
};
</script>
