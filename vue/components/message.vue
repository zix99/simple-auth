<template>
  <article class="message" :class="type">
    <div class="message-body">
      <slot></slot>
    </div>
  </article>
</template>

<script>
export default {
  props: {
    type: null,
  },
};
</script>
