<template>
  <div class="columns is-centered">
    <div class="column is-half">
      <Card :title="title">
        <slot />
      </Card>
    </div>
  </div>
</template>

<script>
import Card from './card.vue';

export default {
  components: { Card },
  props: {
    title: null,
  },
};
</script>
